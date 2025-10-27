// Copyright 2025 Gerry Miller <gerry@gerrymiller.com>
//
// Licensed under the MIT License.
// See LICENSE file in the project root for full license information.

package chunker

import "deep-thinking-agent/pkg/schema"

// SlidingWindowChunker implements traditional sliding window chunking.
// Chunks have fixed size with overlap for context preservation.
type SlidingWindowChunker struct {
	chunkSize    int
	chunkOverlap int
}

// NewSlidingWindowChunker creates a new sliding window chunker.
func NewSlidingWindowChunker(config *ChunkerConfig) *SlidingWindowChunker {
	return &SlidingWindowChunker{
		chunkSize:    config.ChunkSize,
		chunkOverlap: config.ChunkOverlap,
	}
}

// Chunk splits a document using sliding window with overlap.
func (c *SlidingWindowChunker) Chunk(content string, docSchema *schema.DocumentSchema) ([]Chunk, error) {
	chunks := make([]Chunk, 0)
	chunkIndex := 0
	currentPos := 0
	contentLen := len(content)

	for currentPos < contentLen {
		// Calculate chunk end position
		endPos := currentPos + c.chunkSize
		if endPos > contentLen {
			endPos = contentLen
		}

		// Extract chunk text
		chunkText := content[currentPos:endPos]

		// Try to break at sentence or word boundary if not at document end
		if endPos < contentLen {
			// Look for sentence boundary (period, question mark, exclamation)
			lastSentence := -1
			for i := len(chunkText) - 1; i >= len(chunkText)-50 && i >= 0; i-- {
				if chunkText[i] == '.' || chunkText[i] == '?' || chunkText[i] == '!' {
					if i+1 < len(chunkText) && chunkText[i+1] == ' ' {
						lastSentence = i + 1
						break
					}
				}
			}

			// If found sentence boundary, use it
			if lastSentence > 0 {
				chunkText = chunkText[:lastSentence]
				endPos = currentPos + lastSentence
			} else {
				// Fall back to word boundary
				lastSpace := -1
				for i := len(chunkText) - 1; i >= len(chunkText)-50 && i >= 0; i-- {
					if chunkText[i] == ' ' {
						lastSpace = i
						break
					}
				}

				if lastSpace > 0 {
					chunkText = chunkText[:lastSpace]
					endPos = currentPos + lastSpace
				}
			}
		}

		chunks = append(chunks, Chunk{
			Index:    chunkIndex,
			Text:     chunkText,
			StartPos: currentPos,
			EndPos:   endPos,
		})

		chunkIndex++

		// Move window forward (with overlap)
		step := c.chunkSize - c.chunkOverlap
		if step <= 0 {
			step = c.chunkSize / 2 // Safety: ensure we make progress
		}
		currentPos += step

		// If we're at the end, break
		if endPos >= contentLen {
			break
		}
	}

	// Attach metadata
	chunks = attachMetadata(chunks, docSchema, "sliding_window")

	return chunks, nil
}

// Name returns the chunking strategy name.
func (c *SlidingWindowChunker) Name() string {
	return "sliding_window"
}
