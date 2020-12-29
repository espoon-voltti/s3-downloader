// SPDX-FileCopyrightText: 2018-2020 City of Espoo
//
// SPDX-License-Identifier: MIT

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/sirupsen/logrus"
)

// DownloadDetails of an S3 object to be downloaded
type DownloadDetails struct {
	Bucket     string
	Key        string
	TargetFile string
}

// Shared logger to always include static fields
var log *logrus.Entry

func init() {
	logrus.SetFormatter(&logrus.JSONFormatter{
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime:  "@timestamp",
			logrus.FieldKeyLevel: "logLevel",
			logrus.FieldKeyMsg:   "message",
		},
	})

	log = logrus.WithFields(logrus.Fields{
		"appBuild":   getEnv("APP_BUILD", "local"),
		"appCommit":  getEnv("APP_COMMIT", "HEAD"),
		"appName":    os.Getenv("APP_NAME"),
		"env":        getEnv("VOLTTI_ENV", getEnv("ENV", "local")),
		"hostIp":     os.Getenv("HOST_IP"),
		"type":       "app-misc", // Voltti log type
		"userIdHash": "",         // Required field, should be empty here
		"version":    "1",        // Version of the log format
	})
}

func main() {
	bucket, prefix, targetDir := readArgs(os.Args)

	sess, err := session.NewSession()
	svc := s3.New(sess)

	resp, err := svc.ListObjects(&s3.ListObjectsInput{
		Bucket: aws.String(bucket),
		Prefix: aws.String(prefix),
	})
	check(err, "Unable to list objects in bucket: %v", bucket)

	downloader := s3manager.NewDownloader(sess)
	wg := sync.WaitGroup{}

	for _, obj := range resp.Contents {
		if strings.HasSuffix(*obj.Key, "/") {
			// Skip folders
			continue
		}
		wg.Add(1)
		go downloadObject(downloader, &wg, &DownloadDetails{
			Bucket:     bucket,
			Key:        *obj.Key,
			TargetFile: targetFilePath(*obj.Key, prefix, targetDir),
		})
	}
	// Wait for all download to finish before exiting
	wg.Wait()
}

func targetFilePath(key string, prefix string, targetDir string) string {
	pathInDir := strings.TrimPrefix(key, prefix)
	return filepath.Join(targetDir, pathInDir)
}

func readArgs(argv []string) (string, string, string) {
	argc := len(argv)
	if argc != 4 {
		log.Error(fmt.Sprintf("Wrong number of arguments provided. Usage: %s bucket prefix targetDir", argv[0]))
		os.Exit(1)
	}
	bucket := argv[1]
	prefix := ensureSuffix(argv[2], "/")
	targetDir := ensureSuffix(argv[3], "/")
	return bucket, prefix, targetDir
}

func downloadObject(downloader *s3manager.Downloader, wg *sync.WaitGroup, details *DownloadDetails) {
	defer wg.Done()

	log.Info(fmt.Sprintf("Downloading object key: %v, file: %q", details.Key, details.TargetFile))
	file := createFile(details.TargetFile)
	defer file.Close()

	_, err := downloader.Download(file, &s3.GetObjectInput{
		Bucket: aws.String(details.Bucket),
		Key:    aws.String(details.Key),
	})
	check(err, "Unable to download item %q", details.TargetFile)
}

func createFile(targetFile string) *os.File {
	dir := filepath.Dir(targetFile)

	err := os.MkdirAll(dir, os.ModePerm)
	check(err, "Unable to create directory: %v", dir)

	file, err := os.Create(targetFile)
	check(err, "Unable to open file: %q", targetFile)

	return file
}

func check(err error, msg string, args ...interface{}) {
	if err != nil {
		log.WithFields(logrus.Fields{
			"exception":  "Error",
			"stackTrace": err.Error(),
		}).Error(fmt.Sprintf(msg, args...))
		os.Exit(2)
	}
}

func ensureSuffix(s string, suffix string) string {
	if strings.HasSuffix(s, suffix) {
		return s
	}
	return s + suffix
}

// Get environment variable value or default if not defined
func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}
