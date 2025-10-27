// Copyright 2025 Gerry Miller <gerry@gerrymiller.com>
//
// Licensed under the MIT License.
// See LICENSE file in the project root for full license information.

package schema

import "fmt"

// MetadataBuilder helps construct chunk metadata from document schemas.
type MetadataBuilder struct {
	schema *DocumentSchema
}

// NewMetadataBuilder creates a new metadata builder for a document schema.
func NewMetadataBuilder(schema *DocumentSchema) *MetadataBuilder {
	return &MetadataBuilder{
		schema: schema,
	}
}

// BuildChunkMetadata creates metadata for a text chunk based on its position.
func (mb *MetadataBuilder) BuildChunkMetadata(chunkIndex, startPos, endPos int, chunkMethod string) *ChunkMetadata {
	metadata := &ChunkMetadata{
		DocID:            mb.schema.DocID,
		DocTitle:         mb.schema.Title,
		Format:           mb.schema.Format,
		StartPos:         startPos,
		EndPos:           endPos,
		ChunkIndex:       chunkIndex,
		ChunkMethod:      chunkMethod,
		CustomAttributes: make(map[string]interface{}),
	}

	// Find section containing this chunk
	if section := mb.findContainingSection(startPos, endPos); section != nil {
		metadata.SectionID = section.ID
		metadata.SectionTitle = section.Title
		metadata.SectionType = section.Type
		metadata.SectionLevel = section.Level
		metadata.HierarchyPath = mb.getHierarchyPath(section)
	}

	// Find semantic tags for this chunk
	metadata.SemanticTags = mb.findSemanticTags(startPos, endPos)
	metadata.SemanticTypes = mb.findSemanticTypes(startPos, endPos)

	// Copy custom attributes from schema
	for k, v := range mb.schema.CustomAttributes {
		metadata.CustomAttributes[k] = v
	}

	return metadata
}

// findContainingSection finds the section that contains the given position range.
func (mb *MetadataBuilder) findContainingSection(startPos, endPos int) *Section {
	for i := range mb.schema.Sections {
		section := &mb.schema.Sections[i]
		// Check if chunk is within section boundaries
		if startPos >= section.StartPos && endPos <= section.EndPos {
			return section
		}
		// Also check if chunk overlaps significantly with section
		if startPos < section.EndPos && endPos > section.StartPos {
			return section
		}
	}
	return nil
}

// getHierarchyPath constructs a hierarchy path for a section.
func (mb *MetadataBuilder) getHierarchyPath(section *Section) string {
	// Simple implementation: use section ID as path
	// TODO: Build proper hierarchical path in future enhancement
	return section.ID
}

// findSemanticTags extracts semantic tags for a position range.
func (mb *MetadataBuilder) findSemanticTags(startPos, endPos int) []string {
	tags := make([]string, 0)
	seen := make(map[string]bool)

	for _, region := range mb.schema.SemanticRegions {
		// Check if position range overlaps with any region boundary
		for _, boundary := range region.Boundaries {
			if overlaps(startPos, endPos, boundary.StartPos, boundary.EndPos) {
				// Add keywords as tags
				for _, keyword := range region.Keywords {
					if !seen[keyword] {
						tags = append(tags, keyword)
						seen[keyword] = true
					}
				}
				break
			}
		}
	}

	return tags
}

// findSemanticTypes extracts semantic region types for a position range.
func (mb *MetadataBuilder) findSemanticTypes(startPos, endPos int) []string {
	types := make([]string, 0)
	seen := make(map[string]bool)

	for _, region := range mb.schema.SemanticRegions {
		for _, boundary := range region.Boundaries {
			if overlaps(startPos, endPos, boundary.StartPos, boundary.EndPos) {
				if !seen[region.Type] {
					types = append(types, region.Type)
					seen[region.Type] = true
				}
				break
			}
		}
	}

	return types
}

// overlaps checks if two ranges overlap.
func overlaps(start1, end1, start2, end2 int) bool {
	return start1 < end2 && end1 > start2
}

// DocumentIndex manages document-level schema indexing.
type DocumentIndex struct {
	entries map[string]*IndexEntry
}

// IndexEntry contains schema and chunk references for a document.
type IndexEntry struct {
	DocID    string
	Schema   *DocumentSchema
	ChunkIDs []string
}

// NewDocumentIndex creates a new document index.
func NewDocumentIndex() *DocumentIndex {
	return &DocumentIndex{
		entries: make(map[string]*IndexEntry),
	}
}

// Add adds or updates a document index entry.
func (di *DocumentIndex) Add(docID string, schema *DocumentSchema, chunkIDs []string) {
	di.entries[docID] = &IndexEntry{
		DocID:    docID,
		Schema:   schema,
		ChunkIDs: chunkIDs,
	}
}

// Get retrieves a document index entry.
func (di *DocumentIndex) Get(docID string) (*IndexEntry, bool) {
	entry, ok := di.entries[docID]
	return entry, ok
}

// Delete removes a document from the index.
func (di *DocumentIndex) Delete(docID string) {
	delete(di.entries, docID)
}

// List returns all document IDs in the index.
func (di *DocumentIndex) List() []string {
	ids := make([]string, 0, len(di.entries))
	for id := range di.entries {
		ids = append(ids, id)
	}
	return ids
}

// Count returns the number of indexed documents.
func (di *DocumentIndex) Count() int {
	return len(di.entries)
}

// GetSchemaForChunk retrieves the schema for a document containing a chunk.
func (di *DocumentIndex) GetSchemaForChunk(chunkID string) (*DocumentSchema, error) {
	for _, entry := range di.entries {
		for _, id := range entry.ChunkIDs {
			if id == chunkID {
				return entry.Schema, nil
			}
		}
	}
	return nil, fmt.Errorf("no schema found for chunk: %s", chunkID)
}
