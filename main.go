// SPDX-FileCopyrightText: 2018-2020 City of Espoo
//
// SPDX-License-Identifier: MIT

package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/transfermanager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/sirupsen/logrus"
)

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
	ctx := context.Background()
	bucket, prefix, targetDir := readArgs(os.Args)

	// Load AWS configuration (automatically reads AWS_ENDPOINT_URL from environment)
	cfg, err := config.LoadDefaultConfig(ctx)
	check(err, "Unable to load AWS configuration")

	svc := s3.NewFromConfig(cfg, func(o *s3.Options) {
		// Enable path-style addressing for S3-compatible services like LocalStack
		if os.Getenv("AWS_S3_FORCE_PATH_STYLE") == "true" {
			o.UsePathStyle = true
		}
	})

	out, err := transfermanager.New(svc).DownloadDirectory(ctx, &transfermanager.DownloadDirectoryInput{
		Bucket:      aws.String(bucket),
		KeyPrefix:   aws.String(prefix),
		Destination: aws.String(targetDir),
		Callback:    logDownload{},
	})
	check(err, "Unable to download objects from bucket %v with prefix %q to %q", bucket, prefix, targetDir)

	if out.ObjectsDownloaded == 0 {
		log.Warn(fmt.Sprintf("No objects found in bucket %v with prefix %q", bucket, prefix))
		return
	}
	log.Info(fmt.Sprintf("Downloaded %d objects from bucket %v with prefix %q to %q", out.ObjectsDownloaded, bucket, prefix, targetDir))
}

// logDownload logs the key of each object just before DownloadDirectory requests it
type logDownload struct{}

func (logDownload) UpdateRequest(in *transfermanager.GetObjectInput) {
	log.Info(fmt.Sprintf("Downloading object key: %v", aws.ToString(in.Key)))
}

func readArgs(argv []string) (string, string, string) {
	argc := len(argv)
	if argc != 4 {
		log.Error(fmt.Sprintf("Wrong number of arguments provided. Usage: %s bucket prefix targetDir", argv[0]))
		os.Exit(1)
	}
	return argv[1], keyPrefix(argv[2]), argv[3]
}

// keyPrefix makes a prefix match whole path segments only, so that "config"
// does not also match "config-old/file". An empty prefix matches the whole bucket.
func keyPrefix(prefix string) string {
	if prefix == "" || strings.HasSuffix(prefix, "/") {
		return prefix
	}
	return prefix + "/"
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

// Get environment variable value or default if not defined
func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}
