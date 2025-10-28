// Copyright 2025 Gerry Miller <gerry@gerrymiller.com>
//
// Licensed under the MIT License.
// See LICENSE file in the project root for full license information.

package parser

import (
	"strings"
	"testing"
)

// TestHTMLParser_Parse tests HTML parsing functionality.
func TestHTMLParser_Parse(t *testing.T) {
	parser := NewHTMLParser()

	t.Run("supported formats", func(t *testing.T) {
		formats := parser.SupportedFormats()
		expectedFormats := []string{".html", ".HTML", ".htm", ".HTM"}

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
		if parser.Name() != "html" {
			t.Errorf("expected format name 'html', got %s", parser.Name())
		}
	})

	t.Run("simple HTML", func(t *testing.T) {
		html := `<!DOCTYPE html>
<html>
<head>
    <title>Test Page</title>
</head>
<body>
    <h1>Hello World</h1>
    <p>This is a test paragraph.</p>
</body>
</html>`

		reader := strings.NewReader(html)
		doc, err := parser.Parse(reader, "test.html")

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if doc == nil {
			t.Fatal("expected document, got nil")
		}

		if doc.Title != "Test Page" {
			t.Errorf("expected title 'Test Page', got %s", doc.Title)
		}

		if !strings.Contains(doc.Content, "Hello World") {
			t.Error("expected content to contain 'Hello World'")
		}

		if !strings.Contains(doc.Content, "test paragraph") {
			t.Error("expected content to contain 'test paragraph'")
		}

		if doc.Format != "html" {
			t.Errorf("expected format 'html', got %s", doc.Format)
		}

		if doc.SourcePath != "test.html" {
			t.Errorf("expected source path 'test.html', got %s", doc.SourcePath)
		}
	})

	t.Run("HTML without title", func(t *testing.T) {
		html := `<html><body><p>Content without title</p></body></html>`

		reader := strings.NewReader(html)
		doc, err := parser.Parse(reader, "notitle.html")

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if doc.Title != "" {
			t.Errorf("expected empty title, got %s", doc.Title)
		}

		if !strings.Contains(doc.Content, "Content without title") {
			t.Error("expected content to be extracted")
		}
	})

	t.Run("HTML with scripts and styles", func(t *testing.T) {
		html := `<!DOCTYPE html>
<html>
<head>
    <title>Script Test</title>
    <style>body { color: red; }</style>
</head>
<body>
    <p>Visible content</p>
    <script>console.log('hidden');</script>
    <p>More visible content</p>
</body>
</html>`

		reader := strings.NewReader(html)
		doc, err := parser.Parse(reader, "scripts.html")

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Should contain visible content
		if !strings.Contains(doc.Content, "Visible content") {
			t.Error("expected to find visible content")
		}

		if !strings.Contains(doc.Content, "More visible content") {
			t.Error("expected to find more visible content")
		}

		// Should NOT contain script or style content
		if strings.Contains(doc.Content, "console.log") {
			t.Error("script content should be excluded")
		}

		if strings.Contains(doc.Content, "color: red") {
			t.Error("style content should be excluded")
		}
	})

	t.Run("HTML with multiple elements", func(t *testing.T) {
		html := `<!DOCTYPE html>
<html>
<head><title>Multi Element</title></head>
<body>
    <h1>Header 1</h1>
    <h2>Header 2</h2>
    <p>Paragraph 1</p>
    <ul>
        <li>Item 1</li>
        <li>Item 2</li>
    </ul>
    <div>
        <span>Nested content</span>
    </div>
</body>
</html>`

		reader := strings.NewReader(html)
		doc, err := parser.Parse(reader, "multi.html")

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Check all content is extracted
		expectedContent := []string{"Header 1", "Header 2", "Paragraph 1", "Item 1", "Item 2", "Nested content"}
		for _, expected := range expectedContent {
			if !strings.Contains(doc.Content, expected) {
				t.Errorf("expected content to contain %q", expected)
			}
		}
	})

	t.Run("empty HTML", func(t *testing.T) {
		reader := strings.NewReader("")
		doc, err := parser.Parse(reader, "empty.html")

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if doc.Title != "" {
			t.Errorf("expected empty title for empty HTML, got %s", doc.Title)
		}
	})

	t.Run("malformed HTML", func(t *testing.T) {
		// HTML parser is lenient and should handle malformed HTML
		html := `<html><body><p>Unclosed paragraph<div>content</body></html>`

		reader := strings.NewReader(html)
		doc, err := parser.Parse(reader, "malformed.html")

		if err != nil {
			t.Fatalf("unexpected error for malformed HTML: %v", err)
		}

		// Should still extract content despite malformed structure
		if !strings.Contains(doc.Content, "Unclosed paragraph") {
			t.Error("expected to extract content from malformed HTML")
		}
	})

	t.Run("HTML with special characters", func(t *testing.T) {
		html := `<!DOCTYPE html>
<html>
<head><title>Special &amp; Characters</title></head>
<body>
    <p>Less than: &lt;</p>
    <p>Greater than: &gt;</p>
    <p>Quote: &quot;</p>
</body>
</html>`

		reader := strings.NewReader(html)
		doc, err := parser.Parse(reader, "special.html")

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if doc.Title != "Special & Characters" {
			t.Errorf("expected title with decoded entities, got %s", doc.Title)
		}

		// Content should have decoded entities
		if !strings.Contains(doc.Content, "Less than") {
			t.Error("expected decoded content")
		}
	})

	t.Run("metadata extraction", func(t *testing.T) {
		html := `<!DOCTYPE html>
<html>
<head>
    <title>Metadata Test</title>
    <meta name="description" content="Test description">
    <meta name="keywords" content="test, html, parser">
</head>
<body><p>Content</p></body>
</html>`

		reader := strings.NewReader(html)
		doc, err := parser.Parse(reader, "metadata.html")

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if doc.Metadata == nil {
			t.Fatal("expected metadata, got nil")
		}

		// Check word count metadata is present
		if _, ok := doc.Metadata["word_count"]; !ok {
			t.Error("expected word_count in metadata")
		}
	})
}
