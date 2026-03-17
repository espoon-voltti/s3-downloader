// SPDX-FileCopyrightText: 2018-2020 City of Espoo
//
// SPDX-License-Identifier: MIT

package main

import (
	"path/filepath"
	"testing"
)

func TestTargetFilePath(t *testing.T) {
	tests := []struct {
		name      string
		key       string
		prefix    string
		targetDir string
		expected  string
	}{
		{
			name:      "basic file path",
			key:       "documents/file.txt",
			prefix:    "documents/",
			targetDir: "/tmp/download/",
			expected:  filepath.Join("/tmp/download", "file.txt"),
		},
		{
			name:      "nested file path",
			key:       "documents/subfolder/file.pdf",
			prefix:    "documents/",
			targetDir: "/tmp/download/",
			expected:  filepath.Join("/tmp/download", "subfolder/file.pdf"),
		},
		{
			name:      "no prefix match",
			key:       "other/file.txt",
			prefix:    "documents/",
			targetDir: "/tmp/download/",
			expected:  filepath.Join("/tmp/download", "other/file.txt"),
		},
		{
			name:      "empty prefix",
			key:       "file.txt",
			prefix:    "",
			targetDir: "/tmp/download/",
			expected:  filepath.Join("/tmp/download", "file.txt"),
		},
		{
			name:      "exact prefix match",
			key:       "documents/",
			prefix:    "documents/",
			targetDir: "/tmp/download/",
			expected:  filepath.Join("/tmp/download", ""),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := targetFilePath(tt.key, tt.prefix, tt.targetDir)
			if result != tt.expected {
				t.Errorf("targetFilePath(%q, %q, %q) = %q; want %q",
					tt.key, tt.prefix, tt.targetDir, result, tt.expected)
			}
		})
	}
}

func TestEnsureSuffix(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		suffix   string
		expected string
	}{
		{
			name:     "add missing suffix",
			input:    "path",
			suffix:   "/",
			expected: "path/",
		},
		{
			name:     "suffix already exists",
			input:    "path/",
			suffix:   "/",
			expected: "path/",
		},
		{
			name:     "empty string",
			input:    "",
			suffix:   "/",
			expected: "/",
		},
		{
			name:     "multiple character suffix",
			input:    "file.txt",
			suffix:   ".bak",
			expected: "file.txt.bak",
		},
		{
			name:     "multiple character suffix already exists",
			input:    "file.txt.bak",
			suffix:   ".bak",
			expected: "file.txt.bak",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ensureSuffix(tt.input, tt.suffix)
			if result != tt.expected {
				t.Errorf("ensureSuffix(%q, %q) = %q; want %q",
					tt.input, tt.suffix, result, tt.expected)
			}
		})
	}
}
