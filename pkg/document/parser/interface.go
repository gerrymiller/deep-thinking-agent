package parser

import "io"

// Document represents a parsed document with its content and metadata.
type Document struct {
	// Content is the extracted text content
	Content string

	// Format indicates the original format (pdf, html, markdown, txt)
	Format string

	// Metadata contains format-specific information
	Metadata map[string]interface{}

	// SourcePath is the original file path or URL
	SourcePath string

	// Title of the document (if extractable)
	Title string

	// RawContent stores the original content before parsing (optional)
	RawContent []byte
}

// Parser defines the interface for document parsers.
// Each format (PDF, HTML, MD, TXT) implements this interface.
type Parser interface {
	// Parse extracts content from the input reader.
	// Returns a Document with extracted text and metadata.
	Parse(reader io.Reader, sourcePath string) (*Document, error)

	// SupportedFormats returns the file extensions this parser supports.
	// Example: []string{".pdf", ".PDF"}
	SupportedFormats() []string

	// Name returns the parser name for identification.
	Name() string
}

// ParserRegistry maintains a collection of parsers for different formats.
// This allows the system to automatically select the appropriate parser
// based on file extension.
type ParserRegistry struct {
	parsers map[string]Parser // extension -> parser
}

// NewParserRegistry creates a new registry with default parsers.
func NewParserRegistry() *ParserRegistry {
	registry := &ParserRegistry{
		parsers: make(map[string]Parser),
	}

	// Register default parsers
	registry.Register(NewTextParser())
	registry.Register(NewMarkdownParser())
	// PDF and HTML parsers will be registered when implemented

	return registry
}

// Register adds a parser to the registry.
func (r *ParserRegistry) Register(parser Parser) {
	for _, ext := range parser.SupportedFormats() {
		r.parsers[ext] = parser
	}
}

// GetParser returns the appropriate parser for the given file extension.
// extension should include the dot (e.g., ".pdf", ".txt")
func (r *ParserRegistry) GetParser(extension string) (Parser, bool) {
	parser, ok := r.parsers[extension]
	return parser, ok
}

// ParseFile is a convenience method that selects the appropriate parser
// and parses the content.
func (r *ParserRegistry) ParseFile(reader io.Reader, sourcePath string, extension string) (*Document, error) {
	parser, ok := r.GetParser(extension)
	if !ok {
		// Fall back to text parser for unknown formats
		parser = NewTextParser()
	}

	return parser.Parse(reader, sourcePath)
}
