// Copyright 2025 Gerry Miller <gerry@gerrymiller.com>
//
// Licensed under the MIT License.
// See LICENSE file in the project root for full license information.

package retrieval

import (
	"context"
	"fmt"

	"deep-thinking-agent/pkg/embedding"
	"deep-thinking-agent/pkg/vectorstore"
)

// VectorRetriever implements pure semantic vector similarity search.
type VectorRetriever struct {
	store    vectorstore.Store
	embedder embedding.Embedder
}

// NewVectorRetriever creates a new vector retriever.
func NewVectorRetriever(store vectorstore.Store, embedder embedding.Embedder) *VectorRetriever {
	return &VectorRetriever{
		store:    store,
		embedder: embedder,
	}
}

// Search performs semantic vector similarity search.
func (v *VectorRetriever) Search(ctx context.Context, query string, topK int, filters map[string]interface{}) ([]vectorstore.Document, error) {
	// Generate query embedding
	embedResp, err := v.embedder.Embed(ctx, &embedding.EmbedRequest{
		Texts: []string{query},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to embed query: %w", err)
	}

	if len(embedResp.Vectors) == 0 {
		return nil, fmt.Errorf("no embeddings generated")
	}

	// Perform vector search
	searchResp, err := v.store.Search(ctx, &vectorstore.SearchRequest{
		Vector: embedResp.Vectors[0].Embedding,
		TopK:   topK,
		Filter: filters,
	})

	if err != nil {
		return nil, fmt.Errorf("vector search failed: %w", err)
	}

	return searchResp.Documents, nil
}

// Name returns the retriever name.
func (v *VectorRetriever) Name() string {
	return "vector"
}
