// Copyright 2025 Gerry Miller <gerry@gerrymiller.com>
//
// Licensed under the MIT License.
// See LICENSE file in the project root for full license information.

package chunker

import "deep-thinking-agent/pkg/schema"

// SectionBasedChunker chunks documents based on identified sections.
// Each section becomes one or more chunks depending on size.
type SectionBasedChunker struct{}

// NewSectionBasedChunker creates a new section-based chunker.
func NewSectionBasedChunker() *SectionBasedChunker {
	return &SectionBasedChunker{}
}

// Chunk splits a document into section-based chunks.
func (c *SectionBasedChunker) Chunk(content string, docSchema *schema.DocumentSchema) ([]Chunk, error) {
	if len(docSchema.Sections) == 0 {
		// Fall back to sliding window if no sections
		fallback := NewSlidingWindowChunker(DefaultConfig())
		return fallback.Chunk(content, docSchema)
	}

	chunks := make([]Chunk, 0)
	chunkIndex := 0

	for _, section := range docSchema.Sections {
		// Extract section content
		start := section.StartPos
		end := section.EndPos
		if end > len(content) {
			end = len(content)
		}
		if start >= end {
			continue
		}

		sectionContent := content[start:end]

		// If section is too large, split it
		const maxSectionSize = 2000
		if len(sectionContent) > maxSectionSize {
			subChunks := splitPreservingWords(sectionContent, maxSectionSize)
			currentPos := start

			for _, subChunk := range subChunks {
				chunks = append(chunks, Chunk{
					Index:    chunkIndex,
					Text:     subChunk,
					StartPos: currentPos,
					EndPos:   currentPos + len(subChunk),
				})
				chunkIndex++
				currentPos += len(subChunk) + 1 // +1 for space
			}
		} else {
			chunks = append(chunks, Chunk{
				Index:    chunkIndex,
				Text:     sectionContent,
				StartPos: start,
				EndPos:   end,
			})
			chunkIndex++
		}
	}

	// Attach metadata
	chunks = attachMetadata(chunks, docSchema, "section_based")

	return chunks, nil
}

// Name returns the chunking strategy name.
func (c *SectionBasedChunker) Name() string {
	return "section_based"
}
