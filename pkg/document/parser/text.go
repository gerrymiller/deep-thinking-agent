package parser

import (
	"bufio"
	"io"
	"strings"
)

// TextParser handles plain text files.
type TextParser struct{}

// NewTextParser creates a new text parser instance.
func NewTextParser() *TextParser {
	return &TextParser{}
}

// Parse extracts content from a text file.
func (p *TextParser) Parse(reader io.Reader, sourcePath string) (*Document, error) {
	// Read all content
	content, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	text := string(content)

	// Try to extract title from first line if it looks like a title
	title := extractTitleFromText(text)

	doc := &Document{
		Content:    text,
		Format:     "txt",
		SourcePath: sourcePath,
		Title:      title,
		Metadata: map[string]interface{}{
			"line_count": countLines(text),
			"char_count": len(text),
		},
		RawContent: content,
	}

	return doc, nil
}

// SupportedFormats returns the file extensions this parser supports.
func (p *TextParser) SupportedFormats() []string {
	return []string{".txt", ".TXT", ".text", ".TEXT"}
}

// Name returns the parser name.
func (p *TextParser) Name() string {
	return "text"
}

// extractTitleFromText attempts to extract a title from the first line.
// If the first line is short and followed by a blank line, treat it as a title.
func extractTitleFromText(text string) string {
	lines := strings.Split(text, "\n")
	if len(lines) < 2 {
		return ""
	}

	firstLine := strings.TrimSpace(lines[0])
	secondLine := strings.TrimSpace(lines[1])

	// If first line is short (< 100 chars) and second line is empty, likely a title
	if len(firstLine) > 0 && len(firstLine) < 100 && secondLine == "" {
		return firstLine
	}

	return ""
}

// countLines counts the number of lines in the text.
func countLines(text string) int {
	scanner := bufio.NewScanner(strings.NewReader(text))
	count := 0
	for scanner.Scan() {
		count++
	}
	return count
}
