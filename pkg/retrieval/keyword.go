// Copyright 2025 Gerry Miller <gerry@gerrymiller.com>
//
// Licensed under the MIT License.
// See LICENSE file in the project root for full license information.

package retrieval

import (
	"context"
	"math"
	"strings"

	"deep-thinking-agent/pkg/vectorstore"
)

// KeywordRetriever implements BM25-style keyword search.
// This is a simplified implementation for Phase 3.
type KeywordRetriever struct {
	store vectorstore.Store
	k1    float64 // BM25 parameter
	b     float64 // BM25 parameter
}

// NewKeywordRetriever creates a new keyword retriever.
func NewKeywordRetriever(store vectorstore.Store) *KeywordRetriever {
	return &KeywordRetriever{
		store: store,
		k1:    1.5, // Standard BM25 values
		b:     0.75,
	}
}

// Search performs keyword-based search using BM25 scoring.
// NOTE: This is a simplified implementation for Phase 3.
// Production systems would use an inverted index (Elasticsearch, etc.)
func (k *KeywordRetriever) Search(ctx context.Context, query string, topK int, filters map[string]interface{}) ([]vectorstore.Document, error) {
	// For Phase 3, keyword search is not fully implemented as it requires
	// an inverted index which is not part of the basic vector store interface.
	// This is a placeholder that returns an empty result.
	// In production, integrate with Elasticsearch or similar.

	// Return empty results for now
	// TODO: Implement proper keyword search with inverted index
	return []vectorstore.Document{}, nil
}

// tokenize splits text into terms.
func (k *KeywordRetriever) tokenize(text string) []string {
	text = strings.ToLower(text)
	words := strings.Fields(text)

	// Remove common stopwords
	stopwords := map[string]bool{
		"a": true, "an": true, "the": true, "and": true, "or": true,
		"but": true, "in": true, "on": true, "at": true, "to": true,
		"for": true, "of": true, "with": true, "by": true, "from": true,
		"is": true, "was": true, "are": true, "were": true, "be": true,
	}

	filtered := make([]string, 0, len(words))
	for _, word := range words {
		if !stopwords[word] && len(word) > 2 {
			filtered = append(filtered, word)
		}
	}

	return filtered
}

// calculateAvgDocLength computes average document length.
func (k *KeywordRetriever) calculateAvgDocLength(docs []vectorstore.Document) float64 {
	if len(docs) == 0 {
		return 0
	}

	total := 0
	for _, doc := range docs {
		total += len(k.tokenize(doc.Content))
	}

	return float64(total) / float64(len(docs))
}

// scoreBM25 calculates BM25 score for a document.
func (k *KeywordRetriever) scoreBM25(queryTerms []string, doc vectorstore.Document, avgDocLen float64, corpusSize int) float64 {
	docTerms := k.tokenize(doc.Content)
	docLength := float64(len(docTerms))

	// Build term frequency map
	termFreq := make(map[string]int)
	for _, term := range docTerms {
		termFreq[term]++
	}

	score := 0.0
	for _, term := range queryTerms {
		tf := float64(termFreq[term])
		if tf == 0 {
			continue
		}

		// Simplified IDF (in production, would track document frequencies)
		idf := math.Log(float64(corpusSize) / (tf + 1.0))

		// BM25 formula
		numerator := tf * (k.k1 + 1.0)
		denominator := tf + k.k1*(1.0-k.b+k.b*(docLength/avgDocLen))

		score += idf * (numerator / denominator)
	}

	return score
}

// Name returns the retriever name.
func (k *KeywordRetriever) Name() string {
	return "keyword"
}

type scoredDoc struct {
	doc   vectorstore.Document
	score float64
}
