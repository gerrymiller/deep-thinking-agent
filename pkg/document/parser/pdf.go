package parser

import (
	"errors"
	"io"
)

// PDFParser handles PDF files.
// TODO: Phase 2 - Implement full PDF parsing using a library like pdfcpu or go-fitz
type PDFParser struct{}

// NewPDFParser creates a new PDF parser instance.
func NewPDFParser() *PDFParser {
	return &PDFParser{}
}

// Parse extracts content from a PDF file.
// Currently returns an error as PDF parsing is not yet implemented.
func (p *PDFParser) Parse(reader io.Reader, sourcePath string) (*Document, error) {
	return nil, errors.New("PDF parsing not yet implemented - will be added in Phase 2")
}

// SupportedFormats returns the file extensions this parser supports.
func (p *PDFParser) SupportedFormats() []string {
	return []string{".pdf", ".PDF"}
}

// Name returns the parser name.
func (p *PDFParser) Name() string {
	return "pdf"
}
