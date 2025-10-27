// Copyright 2025 Gerry Miller <gerry@gerrymiller.com>
//
// Licensed under the MIT License.
// See LICENSE file in the project root for full license information.

package parser

import (
	"strings"
	"testing"
)

func TestTextParser_Parse(t *testing.T) {
	tests := []struct {
		name       string
		content    string
		sourcePath string
		wantTitle  string
		wantError  bool
	}{
		{
			name: "simple text file",
			content: `This is a test document.
It has multiple lines.
And some content.`,
			sourcePath: "/path/to/file.txt",
			wantTitle:  "",
			wantError:  false,
		},
		{
			name: "text with title",
			content: `My Document Title

This is the content of the document.
It starts after a blank line.`,
			sourcePath: "/path/to/file.txt",
			wantTitle:  "My Document Title",
			wantError:  false,
		},
		{
			name:       "empty file",
			content:    "",
			sourcePath: "/path/to/empty.txt",
			wantTitle:  "",
			wantError:  false,
		},
	}

	parser := NewTextParser()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := strings.NewReader(tt.content)
			doc, err := parser.Parse(reader, tt.sourcePath)

			if tt.wantError {
				if err == nil {
					t.Error("Parse() expected error but got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("Parse() unexpected error: %v", err)
				return
			}

			if doc.Content != tt.content {
				t.Errorf("Parse() content = %v, want %v", doc.Content, tt.content)
			}

			if doc.Title != tt.wantTitle {
				t.Errorf("Parse() title = %v, want %v", doc.Title, tt.wantTitle)
			}

			if doc.Format != "txt" {
				t.Errorf("Parse() format = %v, want %v", doc.Format, "txt")
			}

			if doc.SourcePath != tt.sourcePath {
				t.Errorf("Parse() sourcePath = %v, want %v", doc.SourcePath, tt.sourcePath)
			}

			// Check metadata
			if doc.Metadata["line_count"] == nil {
				t.Error("Parse() metadata missing line_count")
			}
			if doc.Metadata["char_count"] == nil {
				t.Error("Parse() metadata missing char_count")
			}
		})
	}
}

func TestTextParser_SupportedFormats(t *testing.T) {
	parser := NewTextParser()
	formats := parser.SupportedFormats()

	expectedFormats := []string{".txt", ".TXT", ".text", ".TEXT"}
	if len(formats) != len(expectedFormats) {
		t.Errorf("SupportedFormats() returned %d formats, want %d", len(formats), len(expectedFormats))
	}

	for _, expected := range expectedFormats {
		found := false
		for _, format := range formats {
			if format == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("SupportedFormats() missing expected format: %s", expected)
		}
	}
}

func TestTextParser_Name(t *testing.T) {
	parser := NewTextParser()
	if parser.Name() != "text" {
		t.Errorf("Name() = %v, want %v", parser.Name(), "text")
	}
}
