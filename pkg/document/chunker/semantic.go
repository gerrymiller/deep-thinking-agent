// Copyright 2025 Gerry Miller <gerry@gerrymiller.com>
//
// Licensed under the MIT License.
// See LICENSE file in the project root for full license information.

package chunker

import "deep-thinking-agent/pkg/schema"

// SemanticChunker chunks documents based on semantic regions.
// Uses LLM-identified semantic boundaries rather than structural markers.
type SemanticChunker struct{}

// NewSemanticChunker creates a new semantic chunker.
func NewSemanticChunker() *SemanticChunker {
	return &SemanticChunker{}
}

// Chunk splits a document based on semantic regions.
func (c *SemanticChunker) Chunk(content string, docSchema *schema.DocumentSchema) ([]Chunk, error) {
	if len(docSchema.SemanticRegions) == 0 {
		// Fall back to section-based if no semantic regions
		fallback := NewSectionBasedChunker()
		return fallback.Chunk(content, docSchema)
	}

	chunks := make([]Chunk, 0)
	chunkIndex := 0

	for _, region := range docSchema.SemanticRegions {
		for _, boundary := range region.Boundaries {
			start := boundary.StartPos
			end := boundary.EndPos

			if end > len(content) {
				end = len(content)
			}
			if start >= end || start >= len(content) {
				continue
			}

			regionContent := content[start:end]

			// If region is too large, split it
			const maxRegionSize = 2000
			if len(regionContent) > maxRegionSize {
				subChunks := splitPreservingWords(regionContent, maxRegionSize)
				currentPos := start

				for _, subChunk := range subChunks {
					chunks = append(chunks, Chunk{
						Index:    chunkIndex,
						Text:     subChunk,
						StartPos: currentPos,
						EndPos:   currentPos + len(subChunk),
					})
					chunkIndex++
					currentPos += len(subChunk) + 1
				}
			} else {
				chunks = append(chunks, Chunk{
					Index:    chunkIndex,
					Text:     regionContent,
					StartPos: start,
					EndPos:   end,
				})
				chunkIndex++
			}
		}
	}

	// Attach metadata
	chunks = attachMetadata(chunks, docSchema, "semantic")

	return chunks, nil
}

// Name returns the chunking strategy name.
func (c *SemanticChunker) Name() string {
	return "semantic"
}
