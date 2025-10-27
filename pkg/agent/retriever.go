// Copyright 2025 Gerry Miller <gerry@gerrymiller.com>
//
// Licensed under the MIT License.
// See LICENSE file in the project root for full license information.

package agent

import (
	"context"
	"fmt"

	"deep-thinking-agent/pkg/embedding"
	"deep-thinking-agent/pkg/vectorstore"
	"deep-thinking-agent/pkg/workflow"
)

// Retriever performs schema-aware document retrieval.
// It integrates with various retrieval strategies and applies schema filters.
type Retriever struct {
	vectorStore vectorstore.Store
	embedder    embedding.Embedder
}

// RetrieverConfig contains configuration for the retriever agent.
type RetrieverConfig struct {
	DefaultTopK int
}

// NewRetriever creates a new retriever agent.
func NewRetriever(store vectorstore.Store, embedder embedding.Embedder, config *RetrieverConfig) *Retriever {
	return &Retriever{
		vectorStore: store,
		embedder:    embedder,
	}
}

// Retrieve fetches relevant documents using the specified strategy.
func (r *Retriever) Retrieve(ctx context.Context, retrivalCtx *workflow.RetrievalContext) ([]vectorstore.Document, error) {
	if retrivalCtx == nil {
		return nil, fmt.Errorf("retrieval context is nil")
	}

	// Generate query embedding
	embedResp, err := r.embedder.Embed(ctx, &embedding.EmbedRequest{
		Texts: []string{retrivalCtx.Query},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to embed query: %w", err)
	}

	if len(embedResp.Vectors) == 0 {
		return nil, fmt.Errorf("no embeddings generated")
	}

	queryEmbedding := embedResp.Vectors[0].Embedding

	// Build metadata filters from schema filters
	metadataFilters := r.buildMetadataFilters(retrivalCtx.SchemaFilters)

	// Perform search
	searchResp, err := r.vectorStore.Search(ctx, &vectorstore.SearchRequest{
		Vector: queryEmbedding,
		TopK:   retrivalCtx.TopK,
		Filter: metadataFilters,
	})

	if err != nil {
		return nil, fmt.Errorf("vector search failed: %w", err)
	}

	return searchResp.Documents, nil
}

// buildMetadataFilters converts schema filters to vector store filters.
func (r *Retriever) buildMetadataFilters(schemaFilters *workflow.SchemaFilters) map[string]interface{} {
	if schemaFilters == nil {
		return nil
	}

	filters := make(map[string]interface{})

	if len(schemaFilters.DocumentIDs) > 0 {
		filters["doc_id"] = schemaFilters.DocumentIDs
	}

	if len(schemaFilters.SectionTypes) > 0 {
		filters["section_type"] = schemaFilters.SectionTypes
	}

	if len(schemaFilters.SemanticTags) > 0 {
		filters["semantic_tags"] = schemaFilters.SemanticTags
	}

	// Add custom attributes
	for key, value := range schemaFilters.CustomAttributes {
		filters[key] = value
	}

	return filters
}
