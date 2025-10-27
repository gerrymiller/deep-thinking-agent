// Copyright 2025 Gerry Miller <gerry@gerrymiller.com>
//
// Licensed under the MIT License.
// See LICENSE file in the project root for full license information.

package parser

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/ledongthuc/pdf"
)

// PDFParser handles PDF files.
type PDFParser struct{}

// NewPDFParser creates a new PDF parser instance.
func NewPDFParser() *PDFParser {
	return &PDFParser{}
}

// Parse extracts content from a PDF file.
func (p *PDFParser) Parse(reader io.Reader, sourcePath string) (*Document, error) {
	// Read all content to a temporary file (ledongthuc/pdf requires a file path)
	tempFile, err := os.CreateTemp("", "pdf-*.pdf")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp file: %w", err)
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	// Copy reader content to temp file
	content, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to read PDF content: %w", err)
	}

	if _, err := tempFile.Write(content); err != nil {
		return nil, fmt.Errorf("failed to write temp file: %w", err)
	}
	tempFile.Close()

	// Open PDF file
	f, pdfReader, err := pdf.Open(tempFile.Name())
	if err != nil {
		return nil, fmt.Errorf("failed to open PDF: %w", err)
	}
	defer f.Close()

	// Extract text from all pages
	var textBuilder strings.Builder
	numPages := pdfReader.NumPage()

	for pageNum := 1; pageNum <= numPages; pageNum++ {
		page := pdfReader.Page(pageNum)
		if page.V.IsNull() {
			continue
		}

		// Extract text from page
		pageText, err := page.GetPlainText(nil)
		if err != nil {
			// Continue on error, some pages might be problematic
			continue
		}

		textBuilder.WriteString(pageText)
		textBuilder.WriteString("\n\n")
	}

	text := textBuilder.String()

	// Try to extract title (often in first page or metadata)
	title := extractPDFTitle(text)

	// Build metadata
	metadata := make(map[string]interface{})
	metadata["format"] = "pdf"
	metadata["page_count"] = numPages
	metadata["char_count"] = len(text)
	metadata["word_count"] = len(strings.Fields(text))

	return &Document{
		Content:    text,
		Format:     "pdf",
		SourcePath: sourcePath,
		Title:      title,
		Metadata:   metadata,
		RawContent: content,
	}, nil
}

// extractPDFTitle attempts to extract a title from PDF content.
func extractPDFTitle(text string) string {
	// Simple heuristic: first non-empty line that's not too long
	lines := strings.Split(text, "\n")
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if len(trimmed) > 0 && len(trimmed) < 200 {
			return trimmed
		}
	}
	return ""
}

// SupportedFormats returns the file extensions this parser supports.
func (p *PDFParser) SupportedFormats() []string {
	return []string{".pdf", ".PDF"}
}

// Name returns the parser name.
func (p *PDFParser) Name() string {
	return "pdf"
}
