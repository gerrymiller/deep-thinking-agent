// Copyright 2025 Gerry Miller <gerry@gerrymiller.com>
//
// Licensed under the MIT License.
// See LICENSE file in the project root for full license information.

package chunker

import (
	"deep-thinking-agent/pkg/schema"
	"strings"
)

// Chunk represents a single chunk of text with its position information.
type Chunk struct {
	Index    int
	Text     string
	StartPos int
	EndPos   int
	Metadata *schema.ChunkMetadata
}

// Chunker defines the interface for document chunking strategies.
type Chunker interface {
	// Chunk splits a document into chunks and attaches metadata.
	Chunk(content string, docSchema *schema.DocumentSchema) ([]Chunk, error)

	// Name returns the chunking strategy name.
	Name() string
}

// ChunkerConfig contains configuration for chunking strategies.
type ChunkerConfig struct {
	ChunkSize    int
	ChunkOverlap int
}

// DefaultConfig returns default chunking configuration.
func DefaultConfig() *ChunkerConfig {
	return &ChunkerConfig{
		ChunkSize:    1000,
		ChunkOverlap: 150,
	}
}

// NewChunker creates a chunker based on the schema's recommended strategy.
func NewChunker(strategy string, config *ChunkerConfig) Chunker {
	if config == nil {
		config = DefaultConfig()
	}

	switch strategy {
	case "section_based":
		return NewSectionBasedChunker()
	case "hierarchical":
		return NewHierarchicalChunker()
	case "semantic":
		return NewSemanticChunker()
	case "sliding_window":
		return NewSlidingWindowChunker(config)
	default:
		// Default to sliding window
		return NewSlidingWindowChunker(config)
	}
}

// Factory function for creating appropriate chunker.
func ChunkDocument(content string, docSchema *schema.DocumentSchema, config *ChunkerConfig) ([]Chunk, error) {
	strategy := docSchema.ChunkingStrategy
	if strategy == "" {
		strategy = "sliding_window"
	}

	chunker := NewChunker(strategy, config)
	return chunker.Chunk(content, docSchema)
}

// Helper functions

// splitPreservingWords splits text at word boundaries.
func splitPreservingWords(text string, maxSize int) []string {
	if len(text) <= maxSize {
		return []string{text}
	}

	var chunks []string
	words := strings.Fields(text)
	current := ""

	for _, word := range words {
		test := current
		if current != "" {
			test += " "
		}
		test += word

		if len(test) > maxSize && current != "" {
			chunks = append(chunks, current)
			current = word
		} else {
			current = test
		}
	}

	if current != "" {
		chunks = append(chunks, current)
	}

	return chunks
}

// attachMetadata attaches schema-based metadata to chunks.
func attachMetadata(chunks []Chunk, docSchema *schema.DocumentSchema, methodName string) []Chunk {
	builder := schema.NewMetadataBuilder(docSchema)

	for i := range chunks {
		chunks[i].Metadata = builder.BuildChunkMetadata(
			chunks[i].Index,
			chunks[i].StartPos,
			chunks[i].EndPos,
			methodName,
		)
	}

	return chunks
}
