// Copyright 2025 Gerry Miller <gerry@gerrymiller.com>
//
// Licensed under the MIT License.
// See LICENSE file in the project root for full license information.

package agent

import (
	"context"
	"sort"

	"deep-thinking-agent/pkg/vectorstore"
)

// Reranker applies precision ranking to filter retrieval results.
// In this implementation, it uses a simple scoring approach.
// Future: can be enhanced with cross-encoder models.
type Reranker struct {
	topN int
}

// RerankerConfig contains configuration for the reranker agent.
type RerankerConfig struct {
	TopN int
}

// NewReranker creates a new reranker agent.
func NewReranker(config *RerankerConfig) *Reranker {
	if config == nil {
		config = &RerankerConfig{
			TopN: 3,
		}
	}

	return &Reranker{
		topN: config.TopN,
	}
}

// Rerank selects the top N most relevant documents.
func (r *Reranker) Rerank(ctx context.Context, query string, docs []vectorstore.Document) []vectorstore.Document {
	if len(docs) == 0 {
		return docs
	}

	// Sort by score (already present from vector search)
	sorted := make([]vectorstore.Document, len(docs))
	copy(sorted, docs)

	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Score > sorted[j].Score
	})

	// Return top N
	if len(sorted) > r.topN {
		return sorted[:r.topN]
	}

	return sorted
}

// RerankWithScores reranks and adds additional scoring.
// This is a placeholder for future cross-encoder integration.
func (r *Reranker) RerankWithScores(ctx context.Context, query string, docs []vectorstore.Document) []vectorstore.Document {
	// For now, just use the existing scores
	return r.Rerank(ctx, query, docs)
}
