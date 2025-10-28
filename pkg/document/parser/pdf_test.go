// Copyright 2025 Gerry Miller <gerry@gerrymiller.com>
//
// Licensed under the MIT License.
// See LICENSE file in the project root for full license information.

package parser

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

// TestPDFParser_Parse tests basic PDF parsing functionality.
func TestPDFParser_Parse(t *testing.T) {
	parser := NewPDFParser()

	t.Run("supported formats", func(t *testing.T) {
		formats := parser.SupportedFormats()
		expectedFormats := []string{".pdf", ".PDF"}

		if len(formats) != len(expectedFormats) {
			t.Errorf("expected %d formats, got %d", len(expectedFormats), len(formats))
		}

		for i, format := range formats {
			if format != expectedFormats[i] {
				t.Errorf("expected format %s at index %d, got %s", expectedFormats[i], i, format)
			}
		}
	})

	t.Run("format name", func(t *testing.T) {
		if parser.Name() != "pdf" {
			t.Errorf("expected format name 'pdf', got %s", parser.Name())
		}
	})

	// Note: Full PDF parsing tests would require creating valid PDF files
	// or using test fixtures. Since PDF parsing uses external library,
	// we test the interface and error handling.

	t.Run("invalid PDF content", func(t *testing.T) {
		invalidPDF := bytes.NewReader([]byte("This is not a valid PDF file"))
		_, err := parser.Parse(invalidPDF, "invalid.pdf")

		if err == nil {
			t.Error("expected error for invalid PDF, got nil")
		}

		if !strings.Contains(err.Error(), "failed to open PDF") {
			t.Errorf("expected 'failed to open PDF' error, got: %v", err)
		}
	})

	t.Run("empty reader", func(t *testing.T) {
		emptyReader := bytes.NewReader([]byte{})
		_, err := parser.Parse(emptyReader, "empty.pdf")

		if err == nil {
			t.Error("expected error for empty PDF, got nil")
		}
	})
}

// TestPDFParser_Integration tests PDF parsing with a real PDF if available.
func TestPDFParser_Integration(t *testing.T) {
	// Skip if no test PDF is available
	testPDFPath := "testdata/sample.pdf"
	if _, err := os.Stat(testPDFPath); os.IsNotExist(err) {
		t.Skip("Skipping integration test: testdata/sample.pdf not found")
	}

	file, err := os.Open(testPDFPath)
	if err != nil {
		t.Fatalf("failed to open test PDF: %v", err)
	}
	defer file.Close()

	parser := NewPDFParser()
	doc, err := parser.Parse(file, testPDFPath)

	if err != nil {
		t.Fatalf("failed to parse PDF: %v", err)
	}

	if doc == nil {
		t.Fatal("expected document, got nil")
	}

	if doc.Content == "" {
		t.Error("expected non-empty content")
	}

	if doc.Format != "pdf" {
		t.Errorf("expected format 'pdf', got %s", doc.Format)
	}

	if doc.SourcePath != testPDFPath {
		t.Errorf("expected source path %s, got %s", testPDFPath, doc.SourcePath)
	}

	if doc.Metadata == nil {
		t.Error("expected metadata, got nil")
	}
}

// TestExtractPDFTitle tests title extraction logic.
func TestExtractPDFTitle(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		expected string
	}{
		{
			name:     "empty text",
			text:     "",
			expected: "",
		},
		{
			name:     "single line",
			text:     "This is a title",
			expected: "This is a title",
		},
		{
			name:     "multiple lines",
			text:     "Document Title\n\nThis is the body text\nwith multiple lines.",
			expected: "Document Title",
		},
		{
			name:     "long first line gets skipped",
			text:     strings.Repeat("a", 250) + "\n\nBody text",
			expected: "Body text", // First line is too long (>200), so second line is used
		},
		{
			name:     "whitespace only",
			text:     "   \n\t\n  \n",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractPDFTitle(tt.text)
			if result != tt.expected {
				t.Errorf("expected title %q, got %q", tt.expected, result)
			}
		})
	}
}
