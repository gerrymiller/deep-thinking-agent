// Copyright 2025 Gerry Miller <gerry@gerrymiller.com>
//
// Licensed under the MIT License.
// See LICENSE file in the project root for full license information.

package schema

import (
	"testing"
)

func TestNewMetadataBuilder(t *testing.T) {
	schema := &DocumentSchema{
		DocID: "doc1",
	}

	builder := NewMetadataBuilder(schema)
	if builder == nil {
		t.Fatal("NewMetadataBuilder() returned nil")
	}
	if builder.schema != schema {
		t.Error("schema not set correctly")
	}
}

func TestBuildChunkMetadata(t *testing.T) {
	schema := &DocumentSchema{
		DocID:  "doc1",
		Title:  "Test Document",
		Format: "text",
		Sections: []Section{
			{
				ID:       "sec1",
				Title:    "Section 1",
				Level:    1,
				StartPos: 0,
				EndPos:   100,
				Type:     "introduction",
			},
			{
				ID:       "sec2",
				Title:    "Section 2",
				Level:    1,
				StartPos: 100,
				EndPos:   200,
				Type:     "body",
			},
		},
		SemanticRegions: []SemanticRegion{
			{
				ID:         "reg1",
				Type:       "problem_statement",
				Keywords:   []string{"problem", "issue"},
				Boundaries: []Boundary{{StartPos: 0, EndPos: 50}},
				Confidence: 0.9,
			},
		},
		CustomAttributes: map[string]interface{}{
			"doc_type": "test",
			"version":  1,
		},
	}

	builder := NewMetadataBuilder(schema)

	tests := []struct {
		name         string
		chunkIndex   int
		startPos     int
		endPos       int
		chunkMethod  string
		validateFunc func(*testing.T, *ChunkMetadata)
	}{
		{
			name:        "chunk in first section",
			chunkIndex:  0,
			startPos:    10,
			endPos:      50,
			chunkMethod: "section_based",
			validateFunc: func(t *testing.T, m *ChunkMetadata) {
				if m.DocID != "doc1" {
					t.Errorf("DocID = %v, want doc1", m.DocID)
				}
				if m.DocTitle != "Test Document" {
					t.Errorf("DocTitle = %v, want Test Document", m.DocTitle)
				}
				if m.SectionID != "sec1" {
					t.Errorf("SectionID = %v, want sec1", m.SectionID)
				}
				if m.SectionType != "introduction" {
					t.Errorf("SectionType = %v, want introduction", m.SectionType)
				}
				if len(m.SemanticTags) == 0 {
					t.Error("SemanticTags should not be empty")
				}
				if m.CustomAttributes["doc_type"] != "test" {
					t.Error("CustomAttributes not copied correctly")
				}
			},
		},
		{
			name:        "chunk in second section",
			chunkIndex:  1,
			startPos:    120,
			endPos:      180,
			chunkMethod: "sliding_window",
			validateFunc: func(t *testing.T, m *ChunkMetadata) {
				if m.SectionID != "sec2" {
					t.Errorf("SectionID = %v, want sec2", m.SectionID)
				}
				if m.ChunkIndex != 1 {
					t.Errorf("ChunkIndex = %v, want 1", m.ChunkIndex)
				}
				if m.ChunkMethod != "sliding_window" {
					t.Errorf("ChunkMethod = %v, want sliding_window", m.ChunkMethod)
				}
			},
		},
		{
			name:        "chunk outside all sections",
			chunkIndex:  2,
			startPos:    250,
			endPos:      300,
			chunkMethod: "semantic",
			validateFunc: func(t *testing.T, m *ChunkMetadata) {
				if m.SectionID != "" {
					t.Errorf("SectionID = %v, want empty", m.SectionID)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			metadata := builder.BuildChunkMetadata(tt.chunkIndex, tt.startPos, tt.endPos, tt.chunkMethod)

			if metadata == nil {
				t.Fatal("BuildChunkMetadata() returned nil")
			}

			if metadata.StartPos != tt.startPos {
				t.Errorf("StartPos = %v, want %v", metadata.StartPos, tt.startPos)
			}
			if metadata.EndPos != tt.endPos {
				t.Errorf("EndPos = %v, want %v", metadata.EndPos, tt.endPos)
			}

			if tt.validateFunc != nil {
				tt.validateFunc(t, metadata)
			}
		})
	}
}

func TestFindContainingSection(t *testing.T) {
	schema := &DocumentSchema{
		Sections: []Section{
			{ID: "sec1", StartPos: 0, EndPos: 100},
			{ID: "sec2", StartPos: 100, EndPos: 200},
			{ID: "sec3", StartPos: 200, EndPos: 300},
		},
	}

	builder := NewMetadataBuilder(schema)

	tests := []struct {
		name      string
		startPos  int
		endPos    int
		wantSecID string
		wantNil   bool
	}{
		{
			name:      "fully contained in section",
			startPos:  10,
			endPos:    50,
			wantSecID: "sec1",
		},
		{
			name:      "at section boundary",
			startPos:  100,
			endPos:    150,
			wantSecID: "sec2",
		},
		{
			name:      "overlapping section",
			startPos:  90,
			endPos:    110,
			wantSecID: "sec1",
		},
		{
			name:     "outside all sections",
			startPos: 400,
			endPos:   500,
			wantNil:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			section := builder.findContainingSection(tt.startPos, tt.endPos)

			if tt.wantNil {
				if section != nil {
					t.Errorf("expected nil, got section %v", section.ID)
				}
			} else {
				if section == nil {
					t.Fatal("expected section, got nil")
				}
				if section.ID != tt.wantSecID {
					t.Errorf("section ID = %v, want %v", section.ID, tt.wantSecID)
				}
			}
		})
	}
}

func TestGetHierarchyPath(t *testing.T) {
	schema := &DocumentSchema{}
	builder := NewMetadataBuilder(schema)

	section := &Section{
		ID:    "sec1",
		Title: "Section 1",
	}

	path := builder.getHierarchyPath(section)
	if path != "sec1" {
		t.Errorf("path = %v, want sec1", path)
	}
}

func TestFindSemanticTags(t *testing.T) {
	schema := &DocumentSchema{
		SemanticRegions: []SemanticRegion{
			{
				ID:       "reg1",
				Type:     "introduction",
				Keywords: []string{"intro", "overview"},
				Boundaries: []Boundary{
					{StartPos: 0, EndPos: 100},
				},
			},
			{
				ID:       "reg2",
				Type:     "conclusion",
				Keywords: []string{"conclusion", "summary"},
				Boundaries: []Boundary{
					{StartPos: 200, EndPos: 300},
				},
			},
		},
	}

	builder := NewMetadataBuilder(schema)

	tests := []struct {
		name         string
		startPos     int
		endPos       int
		wantTags     []string
		wantMinCount int
	}{
		{
			name:         "overlaps with first region",
			startPos:     50,
			endPos:       75,
			wantMinCount: 2,
		},
		{
			name:         "overlaps with second region",
			startPos:     250,
			endPos:       275,
			wantMinCount: 2,
		},
		{
			name:         "no overlap",
			startPos:     400,
			endPos:       500,
			wantMinCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tags := builder.findSemanticTags(tt.startPos, tt.endPos)

			if len(tags) < tt.wantMinCount {
				t.Errorf("tag count = %v, want >= %v", len(tags), tt.wantMinCount)
			}
		})
	}
}

func TestFindSemanticTypes(t *testing.T) {
	schema := &DocumentSchema{
		SemanticRegions: []SemanticRegion{
			{
				ID:   "reg1",
				Type: "introduction",
				Boundaries: []Boundary{
					{StartPos: 0, EndPos: 100},
				},
			},
			{
				ID:   "reg2",
				Type: "methodology",
				Boundaries: []Boundary{
					{StartPos: 100, EndPos: 200},
				},
			},
		},
	}

	builder := NewMetadataBuilder(schema)

	tests := []struct {
		name      string
		startPos  int
		endPos    int
		wantTypes []string
	}{
		{
			name:      "in introduction region",
			startPos:  25,
			endPos:    75,
			wantTypes: []string{"introduction"},
		},
		{
			name:      "spanning both regions",
			startPos:  90,
			endPos:    110,
			wantTypes: []string{"introduction", "methodology"},
		},
		{
			name:      "outside all regions",
			startPos:  300,
			endPos:    400,
			wantTypes: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			types := builder.findSemanticTypes(tt.startPos, tt.endPos)

			if len(types) != len(tt.wantTypes) {
				t.Errorf("types count = %v, want %v", len(types), len(tt.wantTypes))
			}

			// Check that expected types are present
			typeMap := make(map[string]bool)
			for _, typ := range types {
				typeMap[typ] = true
			}
			for _, wantType := range tt.wantTypes {
				if !typeMap[wantType] {
					t.Errorf("expected type %v not found", wantType)
				}
			}
		})
	}
}

func TestOverlaps(t *testing.T) {
	tests := []struct {
		name   string
		start1 int
		end1   int
		start2 int
		end2   int
		want   bool
	}{
		{
			name:   "fully overlapping",
			start1: 10,
			end1:   50,
			start2: 20,
			end2:   40,
			want:   true,
		},
		{
			name:   "partial overlap",
			start1: 10,
			end1:   50,
			start2: 40,
			end2:   80,
			want:   true,
		},
		{
			name:   "no overlap",
			start1: 10,
			end1:   50,
			start2: 60,
			end2:   100,
			want:   false,
		},
		{
			name:   "touching boundaries",
			start1: 10,
			end1:   50,
			start2: 50,
			end2:   100,
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := overlaps(tt.start1, tt.end1, tt.start2, tt.end2)
			if got != tt.want {
				t.Errorf("overlaps() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewDocumentIndex(t *testing.T) {
	index := NewDocumentIndex()
	if index == nil {
		t.Fatal("NewDocumentIndex() returned nil")
	}
	if index.entries == nil {
		t.Error("entries map not initialized")
	}
	if index.Count() != 0 {
		t.Errorf("Count() = %v, want 0", index.Count())
	}
}

func TestDocumentIndexAdd(t *testing.T) {
	index := NewDocumentIndex()

	schema := &DocumentSchema{
		DocID: "doc1",
	}
	chunkIDs := []string{"chunk1", "chunk2", "chunk3"}

	index.Add("doc1", schema, chunkIDs)

	if index.Count() != 1 {
		t.Errorf("Count() = %v, want 1", index.Count())
	}

	entry, found := index.Get("doc1")
	if !found {
		t.Fatal("expected to find entry")
	}
	if entry.DocID != "doc1" {
		t.Errorf("DocID = %v, want doc1", entry.DocID)
	}
	if len(entry.ChunkIDs) != 3 {
		t.Errorf("ChunkIDs count = %v, want 3", len(entry.ChunkIDs))
	}
}

func TestDocumentIndexGet(t *testing.T) {
	index := NewDocumentIndex()

	schema := &DocumentSchema{DocID: "doc1"}
	index.Add("doc1", schema, []string{"chunk1"})

	tests := []struct {
		name      string
		docID     string
		wantFound bool
	}{
		{
			name:      "existing document",
			docID:     "doc1",
			wantFound: true,
		},
		{
			name:      "non-existent document",
			docID:     "missing",
			wantFound: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			entry, found := index.Get(tt.docID)
			if found != tt.wantFound {
				t.Errorf("Get() found = %v, want %v", found, tt.wantFound)
			}
			if tt.wantFound && entry == nil {
				t.Error("expected entry, got nil")
			}
		})
	}
}

func TestDocumentIndexDelete(t *testing.T) {
	index := NewDocumentIndex()

	schema := &DocumentSchema{DocID: "doc1"}
	index.Add("doc1", schema, []string{"chunk1"})

	index.Delete("doc1")

	if index.Count() != 0 {
		t.Errorf("Count() = %v, want 0 after delete", index.Count())
	}

	_, found := index.Get("doc1")
	if found {
		t.Error("document should not be found after delete")
	}
}

func TestDocumentIndexList(t *testing.T) {
	index := NewDocumentIndex()

	index.Add("doc1", &DocumentSchema{DocID: "doc1"}, []string{"chunk1"})
	index.Add("doc2", &DocumentSchema{DocID: "doc2"}, []string{"chunk2"})
	index.Add("doc3", &DocumentSchema{DocID: "doc3"}, []string{"chunk3"})

	list := index.List()

	if len(list) != 3 {
		t.Errorf("List() length = %v, want 3", len(list))
	}

	// Check that all doc IDs are present
	idMap := make(map[string]bool)
	for _, id := range list {
		idMap[id] = true
	}

	for _, expectedID := range []string{"doc1", "doc2", "doc3"} {
		if !idMap[expectedID] {
			t.Errorf("expected ID %v not in list", expectedID)
		}
	}
}

func TestDocumentIndexCount(t *testing.T) {
	index := NewDocumentIndex()

	if index.Count() != 0 {
		t.Errorf("Count() = %v, want 0", index.Count())
	}

	index.Add("doc1", &DocumentSchema{}, []string{})
	if index.Count() != 1 {
		t.Errorf("Count() = %v, want 1", index.Count())
	}

	index.Add("doc2", &DocumentSchema{}, []string{})
	if index.Count() != 2 {
		t.Errorf("Count() = %v, want 2", index.Count())
	}
}

func TestDocumentIndexGetSchemaForChunk(t *testing.T) {
	index := NewDocumentIndex()

	schema1 := &DocumentSchema{DocID: "doc1"}
	schema2 := &DocumentSchema{DocID: "doc2"}

	index.Add("doc1", schema1, []string{"chunk1", "chunk2"})
	index.Add("doc2", schema2, []string{"chunk3", "chunk4"})

	tests := []struct {
		name      string
		chunkID   string
		wantDocID string
		wantErr   bool
	}{
		{
			name:      "chunk from doc1",
			chunkID:   "chunk1",
			wantDocID: "doc1",
			wantErr:   false,
		},
		{
			name:      "chunk from doc2",
			chunkID:   "chunk4",
			wantDocID: "doc2",
			wantErr:   false,
		},
		{
			name:    "non-existent chunk",
			chunkID: "missing",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			schema, err := index.GetSchemaForChunk(tt.chunkID)

			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if schema.DocID != tt.wantDocID {
				t.Errorf("DocID = %v, want %v", schema.DocID, tt.wantDocID)
			}
		})
	}
}
