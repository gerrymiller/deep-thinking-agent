// Copyright 2025 Gerry Miller <gerry@gerrymiller.com>
//
// Licensed under the MIT License.
// See LICENSE file in the project root for full license information.

package parser

import (
	"strings"
	"testing"
)

func TestNewParserRegistry(t *testing.T) {
	registry := NewParserRegistry()

	if registry == nil {
		t.Fatal("NewParserRegistry() returned nil")
	}

	if registry.parsers == nil {
		t.Error("parsers map should be initialized")
	}
}

func TestParserRegistry_Register(t *testing.T) {
	registry := NewParserRegistry()

	// Register text parser
	textParser := NewTextParser()
	registry.Register(textParser)

	// Verify it's registered for all its supported formats
	for _, format := range textParser.SupportedFormats() {
		if registry.parsers[format] == nil {
			t.Errorf("parser not registered for format %s", format)
		}
	}
}

func TestParserRegistry_GetParser(t *testing.T) {
	registry := NewParserRegistry()

	// Register text parser
	textParser := NewTextParser()
	registry.Register(textParser)

	tests := []struct {
		name      string
		extension string
		wantFound bool
	}{
		{
			name:      "text file extension",
			extension: ".txt",
			wantFound: true,
		},
		{
			name:      "TXT extension uppercase",
			extension: ".TXT",
			wantFound: true,
		},
		{
			name:      "unsupported format",
			extension: ".xyz",
			wantFound: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser, found := registry.GetParser(tt.extension)

			if tt.wantFound {
				if !found {
					t.Errorf("expected to find parser for %s", tt.extension)
				}
				if parser == nil {
					t.Errorf("expected parser for %s, got nil", tt.extension)
				}
			} else {
				if found {
					t.Errorf("expected not to find parser for %s", tt.extension)
				}
			}
		})
	}
}

func TestParserRegistry_ParseFile(t *testing.T) {
	registry := NewParserRegistry()

	// Register text parser
	textParser := NewTextParser()
	registry.Register(textParser)

	t.Run("successful parse", func(t *testing.T) {
		content := "This is test content"
		reader := strings.NewReader(content)

		doc, err := registry.ParseFile(reader, "test.txt", ".txt")

		if err != nil {
			t.Fatalf("ParseFile() failed: %v", err)
		}

		if doc == nil {
			t.Fatal("expected document, got nil")
		}

		if !strings.Contains(doc.Content, content) {
			t.Errorf("expected content to contain %q", content)
		}
	})

	t.Run("unsupported format falls back to text parser", func(t *testing.T) {
		content := "Plain text content"
		reader := strings.NewReader(content)

		// ParseFile falls back to text parser for unknown formats
		doc, err := registry.ParseFile(reader, "test.xyz", ".xyz")

		if err != nil {
			t.Fatalf("ParseFile() failed: %v", err)
		}

		if doc == nil {
			t.Fatal("expected document, got nil")
		}

		// Should parse as text
		if !strings.Contains(doc.Content, content) {
			t.Errorf("expected content to contain %q", content)
		}
	})
}

func TestParserRegistry_MultipleFormats(t *testing.T) {
	registry := NewParserRegistry()

	// Register multiple parsers
	registry.Register(NewTextParser())
	registry.Register(NewMarkdownParser())

	// Verify both are registered
	if _, found := registry.GetParser(".txt"); !found {
		t.Error("text parser not found")
	}

	if _, found := registry.GetParser(".md"); !found {
		t.Error("markdown parser not found")
	}
}
