// Copyright 2025 Gerry Miller <gerry@gerrymiller.com>
//
// Licensed under the MIT License.
// See LICENSE file in the project root for full license information.

package retrieval

import (
	"context"
	"sort"

	"deep-thinking-agent/pkg/vectorstore"
)

// HybridRetriever combines vector and keyword search using RRF.
// It uses Reciprocal Rank Fusion to merge results from both strategies.
type HybridRetriever struct {
	vectorRetriever  *VectorRetriever
	keywordRetriever *KeywordRetriever
	rrfK             int // RRF constant
}

// NewHybridRetriever creates a new hybrid retriever.
func NewHybridRetriever(vectorRet *VectorRetriever, keywordRet *KeywordRetriever) *HybridRetriever {
	return &HybridRetriever{
		vectorRetriever:  vectorRet,
		keywordRetriever: keywordRet,
		rrfK:             60, // Standard RRF value
	}
}

// Search performs hybrid search combining vector and keyword results.
func (h *HybridRetriever) Search(ctx context.Context, query string, topK int, filters map[string]interface{}) ([]vectorstore.Document, error) {
	// Retrieve from both strategies
	vectorResults, err := h.vectorRetriever.Search(ctx, query, topK*2, filters)
	if err != nil {
		return nil, err
	}

	keywordResults, err := h.keywordRetriever.Search(ctx, query, topK*2, filters)
	if err != nil {
		return nil, err
	}

	// Apply RRF fusion
	fused := h.fuseRRF(vectorResults, keywordResults)

	// Return top K
	if len(fused) > topK {
		return fused[:topK], nil
	}

	return fused, nil
}

// fuseRRF applies Reciprocal Rank Fusion to merge ranked lists.
func (h *HybridRetriever) fuseRRF(vectorResults, keywordResults []vectorstore.Document) []vectorstore.Document {
	// Build rank maps
	vectorRanks := make(map[string]int)
	for i, doc := range vectorResults {
		vectorRanks[doc.ID] = i + 1
	}

	keywordRanks := make(map[string]int)
	for i, doc := range keywordResults {
		keywordRanks[doc.ID] = i + 1
	}

	// Combine unique documents
	docMap := make(map[string]vectorstore.Document)
	for _, doc := range vectorResults {
		docMap[doc.ID] = doc
	}
	for _, doc := range keywordResults {
		if _, exists := docMap[doc.ID]; !exists {
			docMap[doc.ID] = doc
		}
	}

	// Calculate RRF scores
	type rrfDoc struct {
		doc   vectorstore.Document
		score float64
	}

	rrfResults := make([]rrfDoc, 0, len(docMap))
	for id, doc := range docMap {
		score := 0.0

		// Add vector rank contribution
		if rank, exists := vectorRanks[id]; exists {
			score += 1.0 / float64(rank+h.rrfK)
		}

		// Add keyword rank contribution
		if rank, exists := keywordRanks[id]; exists {
			score += 1.0 / float64(rank+h.rrfK)
		}

		rrfResults = append(rrfResults, rrfDoc{
			doc:   doc,
			score: score,
		})
	}

	// Sort by RRF score
	sort.Slice(rrfResults, func(i, j int) bool {
		return rrfResults[i].score > rrfResults[j].score
	})

	// Convert back to documents
	results := make([]vectorstore.Document, len(rrfResults))
	for i, r := range rrfResults {
		doc := r.doc
		doc.Score = float32(r.score)
		results[i] = doc
	}

	return results
}

// Name returns the retriever name.
func (h *HybridRetriever) Name() string {
	return "hybrid"
}
