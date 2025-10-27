package parser

import (
	"errors"
	"io"
)

// HTMLParser handles HTML files.
// TODO: Phase 2 - Implement full HTML parsing using golang.org/x/net/html
type HTMLParser struct{}

// NewHTMLParser creates a new HTML parser instance.
func NewHTMLParser() *HTMLParser {
	return &HTMLParser{}
}

// Parse extracts content from an HTML file.
// Currently returns an error as HTML parsing is not yet implemented.
func (p *HTMLParser) Parse(reader io.Reader, sourcePath string) (*Document, error) {
	return nil, errors.New("HTML parsing not yet implemented - will be added in Phase 2")
}

// SupportedFormats returns the file extensions this parser supports.
func (p *HTMLParser) SupportedFormats() []string {
	return []string{".html", ".HTML", ".htm", ".HTM"}
}

// Name returns the parser name.
func (p *HTMLParser) Name() string {
	return "html"
}
