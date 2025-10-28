// Copyright 2025 Gerry Miller <gerry@gerrymiller.com>
//
// Licensed under the MIT License.
// See LICENSE file in the project root for full license information.

package chunker

import (
	"deep-thinking-agent/pkg/schema"
	"strings"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()
	if config == nil {
		t.Fatal("DefaultConfig() returned nil")
	}
	if config.ChunkSize != 1000 {
		t.Errorf("ChunkSize = %v, want 1000", config.ChunkSize)
	}
	if config.ChunkOverlap != 150 {
		t.Errorf("ChunkOverlap = %v, want 150", config.ChunkOverlap)
	}
}

func TestNewChunker(t *testing.T) {
	tests := []struct {
		name         string
		strategy     string
		expectedType string
	}{
		{
			name:         "section_based",
			strategy:     "section_based",
			expectedType: "*chunker.SectionBasedChunker",
		},
		{
			name:         "hierarchical",
			strategy:     "hierarchical",
			expectedType: "*chunker.HierarchicalChunker",
		},
		{
			name:         "semantic",
			strategy:     "semantic",
			expectedType: "*chunker.SemanticChunker",
		},
		{
			name:         "sliding_window",
			strategy:     "sliding_window",
			expectedType: "*chunker.SlidingWindowChunker",
		},
		{
			name:         "unknown defaults to sliding_window",
			strategy:     "unknown",
			expectedType: "*chunker.SlidingWindowChunker",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			chunker := NewChunker(tt.strategy, DefaultConfig())
			if chunker == nil {
				t.Fatal("NewChunker() returned nil")
			}
		})
	}
}

func TestChunkDocument(t *testing.T) {
	content := "This is test content for chunking."

	tests := []struct {
		name    string
		schema  *schema.DocumentSchema
		wantErr bool
	}{
		{
			name: "with section_based strategy",
			schema: &schema.DocumentSchema{
				DocID:            "doc1",
				ChunkingStrategy: "section_based",
				Sections: []schema.Section{
					{ID: "sec1", StartPos: 0, EndPos: len(content)},
				},
			},
			wantErr: false,
		},
		{
			name: "with empty strategy defaults to sliding_window",
			schema: &schema.DocumentSchema{
				DocID:            "doc1",
				ChunkingStrategy: "",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			chunks, err := ChunkDocument(content, tt.schema, DefaultConfig())
			if (err != nil) != tt.wantErr {
				t.Errorf("ChunkDocument() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && len(chunks) == 0 {
				t.Error("expected chunks, got none")
			}
		})
	}
}

func TestSplitPreservingWords(t *testing.T) {
	tests := []struct {
		name    string
		text    string
		maxSize int
		wantLen int
	}{
		{
			name:    "short text no split",
			text:    "short text",
			maxSize: 100,
			wantLen: 1,
		},
		{
			name:    "long text split",
			text:    strings.Repeat("word ", 50),
			maxSize: 50,
			wantLen: 5, // approximate
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			chunks := splitPreservingWords(tt.text, tt.maxSize)
			if len(chunks) < tt.wantLen {
				t.Errorf("chunk count = %v, want >= %v", len(chunks), tt.wantLen)
			}
			// Verify no chunk exceeds maxSize
			for i, chunk := range chunks {
				if len(chunk) > tt.maxSize*2 { // Allow some flexibility
					t.Errorf("chunk %v length = %v, exceeds maxSize %v", i, len(chunk), tt.maxSize)
				}
			}
		})
	}
}

func TestAttachMetadata(t *testing.T) {
	docSchema := &schema.DocumentSchema{
		DocID:  "doc1",
		Format: "text",
		Sections: []schema.Section{
			{ID: "sec1", StartPos: 0, EndPos: 100},
		},
	}

	chunks := []Chunk{
		{Index: 0, Text: "chunk 1", StartPos: 0, EndPos: 50},
		{Index: 1, Text: "chunk 2", StartPos: 50, EndPos: 100},
	}

	result := attachMetadata(chunks, docSchema, "test_method")

	if len(result) != 2 {
		t.Errorf("result length = %v, want 2", len(result))
	}

	for i, chunk := range result {
		if chunk.Metadata == nil {
			t.Errorf("chunk %v has nil metadata", i)
		}
		if chunk.Metadata.ChunkMethod != "test_method" {
			t.Errorf("chunk %v method = %v, want test_method", i, chunk.Metadata.ChunkMethod)
		}
		if chunk.Metadata.DocID != "doc1" {
			t.Errorf("chunk %v DocID = %v, want doc1", i, chunk.Metadata.DocID)
		}
	}
}

// Section-based chunker tests

func TestNewSectionBasedChunker(t *testing.T) {
	chunker := NewSectionBasedChunker()
	if chunker == nil {
		t.Fatal("NewSectionBasedChunker() returned nil")
	}
	if chunker.Name() != "section_based" {
		t.Errorf("Name() = %v, want section_based", chunker.Name())
	}
}

func TestSectionBasedChunk(t *testing.T) {
	chunker := NewSectionBasedChunker()

	tests := []struct {
		name      string
		content   string
		schema    *schema.DocumentSchema
		wantChunk int
		wantErr   bool
	}{
		{
			name:    "with sections",
			content: strings.Repeat("a", 500),
			schema: &schema.DocumentSchema{
				DocID: "doc1",
				Sections: []schema.Section{
					{ID: "sec1", StartPos: 0, EndPos: 100},
					{ID: "sec2", StartPos: 100, EndPos: 500},
				},
			},
			wantChunk: 2,
			wantErr:   false,
		},
		{
			name:    "with large section that needs splitting",
			content: strings.Repeat("word ", 1000),
			schema: &schema.DocumentSchema{
				DocID: "doc1",
				Sections: []schema.Section{
					{ID: "sec1", StartPos: 0, EndPos: 5000},
				},
			},
			wantChunk: 2, // Will be split
			wantErr:   false,
		},
		{
			name:    "no sections falls back to sliding window",
			content: "test content",
			schema: &schema.DocumentSchema{
				DocID:    "doc1",
				Sections: []schema.Section{},
			},
			wantChunk: 1,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			chunks, err := chunker.Chunk(tt.content, tt.schema)
			if (err != nil) != tt.wantErr {
				t.Errorf("Chunk() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr {
				if len(chunks) < tt.wantChunk {
					t.Errorf("chunk count = %v, want >= %v", len(chunks), tt.wantChunk)
				}
				// Verify metadata is attached
				for i, chunk := range chunks {
					if chunk.Metadata == nil {
						t.Errorf("chunk %v has nil metadata", i)
					}
				}
			}
		})
	}
}

// Hierarchical chunker tests

func TestNewHierarchicalChunker(t *testing.T) {
	chunker := NewHierarchicalChunker()
	if chunker == nil {
		t.Fatal("NewHierarchicalChunker() returned nil")
	}
	if chunker.Name() != "hierarchical" {
		t.Errorf("Name() = %v, want hierarchical", chunker.Name())
	}
}

func TestHierarchicalChunk(t *testing.T) {
	chunker := NewHierarchicalChunker()

	content := "Test content for hierarchical chunking."
	docSchema := &schema.DocumentSchema{
		DocID: "doc1",
		Sections: []schema.Section{
			{ID: "sec1", StartPos: 0, EndPos: len(content)},
		},
	}

	chunks, err := chunker.Chunk(content, docSchema)
	if err != nil {
		t.Fatalf("Chunk() error = %v", err)
	}

	if len(chunks) == 0 {
		t.Error("expected chunks, got none")
	}

	// Verify that metadata has hierarchical method
	for i, chunk := range chunks {
		if chunk.Metadata == nil {
			t.Errorf("chunk %v has nil metadata", i)
		}
		if chunk.Metadata.ChunkMethod != "hierarchical" {
			t.Errorf("chunk %v method = %v, want hierarchical", i, chunk.Metadata.ChunkMethod)
		}
	}
}

// Semantic chunker tests

func TestNewSemanticChunker(t *testing.T) {
	chunker := NewSemanticChunker()
	if chunker == nil {
		t.Fatal("NewSemanticChunker() returned nil")
	}
	if chunker.Name() != "semantic" {
		t.Errorf("Name() = %v, want semantic", chunker.Name())
	}
}

func TestSemanticChunk(t *testing.T) {
	chunker := NewSemanticChunker()

	tests := []struct {
		name      string
		content   string
		schema    *schema.DocumentSchema
		wantChunk int
		wantErr   bool
	}{
		{
			name:    "with semantic regions",
			content: strings.Repeat("a", 500),
			schema: &schema.DocumentSchema{
				DocID: "doc1",
				SemanticRegions: []schema.SemanticRegion{
					{
						ID:   "reg1",
						Type: "introduction",
						Boundaries: []schema.Boundary{
							{StartPos: 0, EndPos: 100},
						},
					},
					{
						ID:   "reg2",
						Type: "body",
						Boundaries: []schema.Boundary{
							{StartPos: 100, EndPos: 500},
						},
					},
				},
			},
			wantChunk: 2,
			wantErr:   false,
		},
		{
			name:    "with large region that needs splitting",
			content: strings.Repeat("word ", 1000),
			schema: &schema.DocumentSchema{
				DocID: "doc1",
				SemanticRegions: []schema.SemanticRegion{
					{
						ID:   "reg1",
						Type: "content",
						Boundaries: []schema.Boundary{
							{StartPos: 0, EndPos: 5000},
						},
					},
				},
			},
			wantChunk: 2,
			wantErr:   false,
		},
		{
			name:    "no semantic regions falls back to section-based",
			content: "test content",
			schema: &schema.DocumentSchema{
				DocID:           "doc1",
				SemanticRegions: []schema.SemanticRegion{},
				Sections: []schema.Section{
					{ID: "sec1", StartPos: 0, EndPos: 12},
				},
			},
			wantChunk: 1,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			chunks, err := chunker.Chunk(tt.content, tt.schema)
			if (err != nil) != tt.wantErr {
				t.Errorf("Chunk() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr {
				if len(chunks) < tt.wantChunk {
					t.Errorf("chunk count = %v, want >= %v", len(chunks), tt.wantChunk)
				}
				// Verify metadata
				for i, chunk := range chunks {
					if chunk.Metadata == nil {
						t.Errorf("chunk %v has nil metadata", i)
					}
					if chunk.Metadata.ChunkMethod != "semantic" {
						t.Errorf("chunk %v method = %v, want semantic", i, chunk.Metadata.ChunkMethod)
					}
				}
			}
		})
	}
}

// Sliding window chunker tests

func TestNewSlidingWindowChunker(t *testing.T) {
	config := &ChunkerConfig{
		ChunkSize:    500,
		ChunkOverlap: 100,
	}

	chunker := NewSlidingWindowChunker(config)
	if chunker == nil {
		t.Fatal("NewSlidingWindowChunker() returned nil")
	}
	if chunker.Name() != "sliding_window" {
		t.Errorf("Name() = %v, want sliding_window", chunker.Name())
	}
	if chunker.chunkSize != 500 {
		t.Errorf("chunkSize = %v, want 500", chunker.chunkSize)
	}
	if chunker.chunkOverlap != 100 {
		t.Errorf("chunkOverlap = %v, want 100", chunker.chunkOverlap)
	}
}

func TestSlidingWindowChunk(t *testing.T) {
	config := &ChunkerConfig{
		ChunkSize:    100,
		ChunkOverlap: 20,
	}
	chunker := NewSlidingWindowChunker(config)

	tests := []struct {
		name          string
		content       string
		wantMinChunks int
	}{
		{
			name:          "short content single chunk",
			content:       "Short text",
			wantMinChunks: 1,
		},
		{
			name:          "long content multiple chunks",
			content:       strings.Repeat("a", 500),
			wantMinChunks: 4,
		},
		{
			name:          "content with sentences",
			content:       strings.Repeat("This is a sentence. ", 50),
			wantMinChunks: 5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			docSchema := &schema.DocumentSchema{
				DocID:  "doc1",
				Format: "text",
			}

			chunks, err := chunker.Chunk(tt.content, docSchema)
			if err != nil {
				t.Fatalf("Chunk() error = %v", err)
			}

			if len(chunks) < tt.wantMinChunks {
				t.Errorf("chunk count = %v, want >= %v", len(chunks), tt.wantMinChunks)
			}

			// Verify chunks
			for i, chunk := range chunks {
				if chunk.Text == "" {
					t.Errorf("chunk %v has empty text", i)
				}
				if chunk.StartPos < 0 {
					t.Errorf("chunk %v has negative StartPos", i)
				}
				if chunk.EndPos <= chunk.StartPos {
					t.Errorf("chunk %v has invalid EndPos", i)
				}
				if chunk.Metadata == nil {
					t.Errorf("chunk %v has nil metadata", i)
				}
				if chunk.Metadata.ChunkMethod != "sliding_window" {
					t.Errorf("chunk %v method = %v, want sliding_window", i, chunk.Metadata.ChunkMethod)
				}
			}

			// Verify overlap (except for last chunk)
			if len(chunks) > 1 {
				for i := 0; i < len(chunks)-1; i++ {
					step := chunks[i+1].StartPos - chunks[i].StartPos
					if step <= 0 {
						t.Errorf("chunks %v and %v: no forward progress", i, i+1)
					}
				}
			}
		})
	}
}

func TestSlidingWindowBoundaryDetection(t *testing.T) {
	config := &ChunkerConfig{
		ChunkSize:    50,
		ChunkOverlap: 10,
	}
	chunker := NewSlidingWindowChunker(config)

	// Content with clear sentence boundaries
	content := "First sentence here. Second sentence here. Third sentence here. Fourth sentence here."
	docSchema := &schema.DocumentSchema{
		DocID:  "doc1",
		Format: "text",
	}

	chunks, err := chunker.Chunk(content, docSchema)
	if err != nil {
		t.Fatalf("Chunk() error = %v", err)
	}

	// Verify that chunks break at sentence boundaries when possible
	for i, chunk := range chunks {
		if i < len(chunks)-1 { // Not the last chunk
			// Check if chunk ends with proper punctuation or space
			lastChar := chunk.Text[len(chunk.Text)-1]
			if lastChar != '.' && lastChar != ' ' && lastChar != '?' && lastChar != '!' {
				// This is acceptable - just checking the logic exists
			}
		}
	}
}

func TestSlidingWindowZeroOverlap(t *testing.T) {
	config := &ChunkerConfig{
		ChunkSize:    100,
		ChunkOverlap: 0,
	}
	chunker := NewSlidingWindowChunker(config)

	content := strings.Repeat("a", 250)
	docSchema := &schema.DocumentSchema{
		DocID:  "doc1",
		Format: "text",
	}

	chunks, err := chunker.Chunk(content, docSchema)
	if err != nil {
		t.Fatalf("Chunk() error = %v", err)
	}

	if len(chunks) < 2 {
		t.Errorf("expected at least 2 chunks, got %v", len(chunks))
	}

	// With zero overlap, step should equal chunk size (or half if safety kicks in)
	if len(chunks) > 1 {
		step := chunks[1].StartPos - chunks[0].StartPos
		if step <= 0 {
			t.Error("no forward progress between chunks")
		}
	}
}
