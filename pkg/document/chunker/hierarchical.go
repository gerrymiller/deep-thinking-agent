// Copyright 2025 Gerry Miller <gerry@gerrymiller.com>
//
// Licensed under the MIT License.
// See LICENSE file in the project root for full license information.

package chunker

import "deep-thinking-agent/pkg/schema"

// HierarchicalChunker chunks documents based on hierarchical structure.
// Preserves parent-child relationships in chunk boundaries.
type HierarchicalChunker struct{}

// NewHierarchicalChunker creates a new hierarchical chunker.
func NewHierarchicalChunker() *HierarchicalChunker {
	return &HierarchicalChunker{}
}

// Chunk splits a document based on hierarchical structure.
func (c *HierarchicalChunker) Chunk(content string, docSchema *schema.DocumentSchema) ([]Chunk, error) {
	// For now, use section-based chunking as hierarchical implementation
	// TODO: Implement true hierarchical chunking that preserves parent-child relationships
	sectionChunker := NewSectionBasedChunker()
	chunks, err := sectionChunker.Chunk(content, docSchema)
	if err != nil {
		return nil, err
	}

	// Update method name in metadata
	for i := range chunks {
		if chunks[i].Metadata != nil {
			chunks[i].Metadata.ChunkMethod = "hierarchical"
		}
	}

	return chunks, nil
}

// Name returns the chunking strategy name.
func (c *HierarchicalChunker) Name() string {
	return "hierarchical"
}
