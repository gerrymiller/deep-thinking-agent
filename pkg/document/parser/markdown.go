// Copyright 2025 Gerry Miller <gerry@gerrymiller.com>
//
// Licensed under the MIT License.
// See LICENSE file in the project root for full license information.

package parser

import (
	"io"
	"regexp"
	"strings"
)

// MarkdownParser handles Markdown files.
type MarkdownParser struct{}

// NewMarkdownParser creates a new markdown parser instance.
func NewMarkdownParser() *MarkdownParser {
	return &MarkdownParser{}
}

// Parse extracts content from a Markdown file.
// For Phase 1, this is a simple implementation that preserves Markdown structure.
// Future phases could render to HTML or extract more sophisticated metadata.
func (p *MarkdownParser) Parse(reader io.Reader, sourcePath string) (*Document, error) {
	// Read all content
	content, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	text := string(content)

	// Extract title from first heading
	title := extractTitleFromMarkdown(text)

	// Extract metadata
	metadata := extractMarkdownMetadata(text)
	metadata["format"] = "markdown"

	doc := &Document{
		Content:    text,
		Format:     "markdown",
		SourcePath: sourcePath,
		Title:      title,
		Metadata:   metadata,
		RawContent: content,
	}

	return doc, nil
}

// SupportedFormats returns the file extensions this parser supports.
func (p *MarkdownParser) SupportedFormats() []string {
	return []string{".md", ".MD", ".markdown", ".MARKDOWN"}
}

// Name returns the parser name.
func (p *MarkdownParser) Name() string {
	return "markdown"
}

// extractTitleFromMarkdown extracts the first H1 heading as the title.
func extractTitleFromMarkdown(text string) string {
	// Try to find # heading (atx style)
	atxPattern := regexp.MustCompile(`(?m)^#\s+(.+)$`)
	if matches := atxPattern.FindStringSubmatch(text); len(matches) > 1 {
		return strings.TrimSpace(matches[1])
	}

	// Try to find underlined heading (setext style)
	setextPattern := regexp.MustCompile(`(?m)^(.+)\n=+\s*$`)
	if matches := setextPattern.FindStringSubmatch(text); len(matches) > 1 {
		return strings.TrimSpace(matches[1])
	}

	return ""
}

// extractMarkdownMetadata extracts basic metadata from Markdown content.
func extractMarkdownMetadata(text string) map[string]interface{} {
	metadata := make(map[string]interface{})

	// Count headings at each level
	headingCounts := make(map[int]int)
	headingPattern := regexp.MustCompile(`(?m)^(#{1,6})\s+.+$`)
	matches := headingPattern.FindAllStringSubmatch(text, -1)
	for _, match := range matches {
		level := len(match[1])
		headingCounts[level]++
	}

	metadata["heading_counts"] = headingCounts
	metadata["total_headings"] = len(matches)

	// Count code blocks
	codeBlockPattern := regexp.MustCompile("```")
	codeBlockMatches := codeBlockPattern.FindAllString(text, -1)
	metadata["code_blocks"] = len(codeBlockMatches) / 2 // Each block has opening and closing

	// Count links
	linkPattern := regexp.MustCompile(`\[.+?\]\(.+?\)`)
	linkMatches := linkPattern.FindAllString(text, -1)
	metadata["link_count"] = len(linkMatches)

	// Count images
	imagePattern := regexp.MustCompile(`!\[.+?\]\(.+?\)`)
	imageMatches := imagePattern.FindAllString(text, -1)
	metadata["image_count"] = len(imageMatches)

	metadata["char_count"] = len(text)
	metadata["line_count"] = strings.Count(text, "\n") + 1

	return metadata
}
