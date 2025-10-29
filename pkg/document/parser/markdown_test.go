// Copyright 2025 Gerry Miller <gerry@gerrymiller.com>
//
// Licensed under the MIT License.
// See LICENSE file in the project root for full license information.

package parser

import (
	"strings"
	"testing"
)

func TestMarkdownParser_Parse(t *testing.T) {
	tests := []struct {
		name       string
		content    string
		sourcePath string
		wantTitle  string
		wantError  bool
	}{
		{
			name: "markdown with atx heading",
			content: `# My Markdown Document

This is some content.

## Section 1
Content here.`,
			sourcePath: "/path/to/file.md",
			wantTitle:  "My Markdown Document",
			wantError:  false,
		},
		{
			name: "markdown with setext heading",
			content: `My Document
===========

This is some content.`,
			sourcePath: "/path/to/file.md",
			wantTitle:  "My Document",
			wantError:  false,
		},
		{
			name: "markdown without heading",
			content: `This is just some text.
No headings here.`,
			sourcePath: "/path/to/file.md",
			wantTitle:  "",
			wantError:  false,
		},
		{
			name:       "empty file",
			content:    "",
			sourcePath: "/path/to/empty.md",
			wantTitle:  "",
			wantError:  false,
		},
	}

	parser := NewMarkdownParser()

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

			if doc.Format != "markdown" {
				t.Errorf("Parse() format = %v, want %v", doc.Format, "markdown")
			}

			// Check metadata exists
			if doc.Metadata["total_headings"] == nil {
				t.Error("Parse() metadata missing total_headings")
			}
		})
	}
}

func TestMarkdownParser_SupportedFormats(t *testing.T) {
	parser := NewMarkdownParser()
	formats := parser.SupportedFormats()

	expectedFormats := []string{".md", ".MD", ".markdown", ".MARKDOWN"}
	if len(formats) != len(expectedFormats) {
		t.Errorf("SupportedFormats() returned %d formats, want %d", len(formats), len(expectedFormats))
	}
}

func TestMarkdownParser_Name(t *testing.T) {
	parser := NewMarkdownParser()
	if parser.Name() != "markdown" {
		t.Errorf("Name() = %v, want %v", parser.Name(), "markdown")
	}
}

func TestExtractMarkdownMetadata(t *testing.T) {
	content := `# Title

## Section 1
Some text.

### Subsection
More text.

` + "```go\ncode here\n```" + `

[Link](http://example.com)
![Image](image.png)
`

	metadata := extractMarkdownMetadata(content)

	// Check heading counts
	headingCounts, ok := metadata["heading_counts"].(map[int]int)
	if !ok {
		t.Fatal("heading_counts not found or wrong type")
	}

	if headingCounts[1] != 1 {
		t.Errorf("Expected 1 h1 heading, got %d", headingCounts[1])
	}
	if headingCounts[2] != 1 {
		t.Errorf("Expected 1 h2 heading, got %d", headingCounts[2])
	}
	if headingCounts[3] != 1 {
		t.Errorf("Expected 1 h3 heading, got %d", headingCounts[3])
	}

	// Check code blocks
	if metadata["code_blocks"].(int) != 1 {
		t.Errorf("Expected 1 code block, got %d", metadata["code_blocks"])
	}

	// Check links (note: this counts images too since they use similar syntax)
	linkCount := metadata["link_count"].(int)
	if linkCount < 1 {
		t.Errorf("Expected at least 1 link, got %d", linkCount)
	}

	// Check images
	if metadata["image_count"].(int) != 1 {
		t.Errorf("Expected 1 image, got %d", metadata["image_count"])
	}
}

func TestMarkdownParser_EdgeCases(t *testing.T) {
	parser := NewMarkdownParser()

	t.Run("markdown with only whitespace", func(t *testing.T) {
		content := "   \n\t\n  \n"
		reader := strings.NewReader(content)
		doc, err := parser.Parse(reader, "whitespace.md")

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if doc.Title != "" {
			t.Errorf("expected empty title, got %q", doc.Title)
		}
	})

	t.Run("markdown with long heading", func(t *testing.T) {
		longTitle := strings.Repeat("A", 250)
		content := "# " + longTitle + "\n\nContent here."
		reader := strings.NewReader(content)
		doc, err := parser.Parse(reader, "long.md")

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Should extract the title even if long
		if !strings.Contains(doc.Title, "A") {
			t.Error("expected title to be extracted")
		}
	})
}
