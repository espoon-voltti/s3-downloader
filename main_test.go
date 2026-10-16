// SPDX-FileCopyrightText: 2018-2020 City of Espoo
//
// SPDX-License-Identifier: MIT

package main

import (
	"testing"
)

func TestKeyPrefix(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "add missing slash",
			input:    "documents",
			expected: "documents/",
		},
		{
			name:     "slash already exists",
			input:    "documents/",
			expected: "documents/",
		},
		{
			name:     "nested prefix",
			input:    "documents/subfolder",
			expected: "documents/subfolder/",
		},
		{
			name:     "empty prefix stays empty to match the whole bucket",
			input:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := keyPrefix(tt.input)
			if result != tt.expected {
				t.Errorf("keyPrefix(%q) = %q; want %q", tt.input, result, tt.expected)
			}
		})
	}
}
