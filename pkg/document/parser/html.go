// Copyright 2025 Gerry Miller <gerry@gerrymiller.com>
//
// Licensed under the MIT License.
// See LICENSE file in the project root for full license information.

package parser

import (
	"bytes"
	"io"
	"strings"

	"golang.org/x/net/html"
)

// HTMLParser handles HTML files.
type HTMLParser struct{}

// NewHTMLParser creates a new HTML parser instance.
func NewHTMLParser() *HTMLParser {
	return &HTMLParser{}
}

// Parse extracts content from an HTML file.
func (p *HTMLParser) Parse(reader io.Reader, sourcePath string) (*Document, error) {
	// Read all content
	content, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	// Parse HTML
	doc, err := html.Parse(bytes.NewReader(content))
	if err != nil {
		return nil, err
	}

	// Extract text content
	var textBuilder strings.Builder
	var title string
	metadata := make(map[string]interface{})

	// Traverse HTML tree and extract text
	var extractText func(*html.Node)
	var findTitle func(*html.Node)

	findTitle = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "title" {
			if n.FirstChild != nil {
				title = n.FirstChild.Data
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			findTitle(c)
		}
	}

	extractText = func(n *html.Node) {
		if n.Type == html.TextNode {
			text := strings.TrimSpace(n.Data)
			if text != "" {
				textBuilder.WriteString(text)
				textBuilder.WriteString(" ")
			}
		}
		// Skip script and style content
		if n.Type == html.ElementNode && (n.Data == "script" || n.Data == "style") {
			return
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			extractText(c)
		}
	}

	// Extract title
	findTitle(doc)

	// Extract text
	extractText(doc)
	text := textBuilder.String()

	// Extract metadata
	metadata["format"] = "html"
	metadata["char_count"] = len(text)
	metadata["word_count"] = len(strings.Fields(text))

	// Count specific elements
	metadata["heading_count"] = countElements(doc, []string{"h1", "h2", "h3", "h4", "h5", "h6"})
	metadata["paragraph_count"] = countElements(doc, []string{"p"})
	metadata["link_count"] = countElements(doc, []string{"a"})

	return &Document{
		Content:    text,
		Format:     "html",
		SourcePath: sourcePath,
		Title:      title,
		Metadata:   metadata,
		RawContent: content,
	}, nil
}

// countElements counts specific HTML elements in the tree.
func countElements(n *html.Node, tags []string) int {
	count := 0
	var walk func(*html.Node)
	walk = func(node *html.Node) {
		if node.Type == html.ElementNode {
			for _, tag := range tags {
				if node.Data == tag {
					count++
					break
				}
			}
		}
		for c := node.FirstChild; c != nil; c = c.NextSibling {
			walk(c)
		}
	}
	walk(n)
	return count
}

// SupportedFormats returns the file extensions this parser supports.
func (p *HTMLParser) SupportedFormats() []string {
	return []string{".html", ".HTML", ".htm", ".HTM"}
}

// Name returns the parser name.
func (p *HTMLParser) Name() string {
	return "html"
}
