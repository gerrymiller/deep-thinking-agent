// Copyright 2025 Gerry Miller <gerry@gerrymiller.com>
//
// Licensed under the MIT License.
// See LICENSE file in the project root for full license information.

package retrieval

import (
	"context"

	"deep-thinking-agent/pkg/vectorstore"
	"deep-thinking-agent/pkg/workflow"
)

// SchemaRetriever implements schema-aware targeted retrieval.
// It applies metadata filters based on document schemas.
type SchemaRetriever struct {
	vectorRetriever *VectorRetriever
}

// NewSchemaRetriever creates a new schema-filtered retriever.
func NewSchemaRetriever(vectorRet *VectorRetriever) *SchemaRetriever {
	return &SchemaRetriever{
		vectorRetriever: vectorRet,
	}
}

// Search performs schema-filtered retrieval.
func (s *SchemaRetriever) Search(ctx context.Context, query string, topK int, schemaFilters *workflow.SchemaFilters) ([]vectorstore.Document, error) {
	// Convert schema filters to metadata filters
	filters := s.buildMetadataFilters(schemaFilters)

	// Use vector search with enhanced filters
	return s.vectorRetriever.Search(ctx, query, topK, filters)
}

// SearchWithFilters performs retrieval with explicit metadata filters.
func (s *SchemaRetriever) SearchWithFilters(ctx context.Context, query string, topK int, filters map[string]interface{}) ([]vectorstore.Document, error) {
	return s.vectorRetriever.Search(ctx, query, topK, filters)
}

// buildMetadataFilters converts schema filters to vector store filters.
func (s *SchemaRetriever) buildMetadataFilters(schemaFilters *workflow.SchemaFilters) map[string]interface{} {
	if schemaFilters == nil {
		return nil
	}

	filters := make(map[string]interface{})

	// Document IDs
	if len(schemaFilters.DocumentIDs) > 0 {
		filters["doc_id"] = schemaFilters.DocumentIDs
	}

	// Section types
	if len(schemaFilters.SectionTypes) > 0 {
		filters["section_type"] = schemaFilters.SectionTypes
	}

	// Hierarchy paths
	if len(schemaFilters.HierarchyPaths) > 0 {
		filters["hierarchy_path"] = schemaFilters.HierarchyPaths
	}

	// Semantic tags
	if len(schemaFilters.SemanticTags) > 0 {
		filters["semantic_tags"] = schemaFilters.SemanticTags
	}

	// Minimum relevance score
	if schemaFilters.MinRelevanceScore > 0 {
		filters["min_score"] = schemaFilters.MinRelevanceScore
	}

	// Custom attributes
	for key, value := range schemaFilters.CustomAttributes {
		filters[key] = value
	}

	return filters
}

// Name returns the retriever name.
func (s *SchemaRetriever) Name() string {
	return "schema_filtered"
}
