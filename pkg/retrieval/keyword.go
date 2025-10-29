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
// This implementation builds an in-memory inverted index from the corpus
// for BM25 scoring. For large-scale production systems, integrate with
// Elasticsearch or similar for better performance.
func (k *KeywordRetriever) Search(ctx context.Context, query string, topK int, filters map[string]interface{}) ([]vectorstore.Document, error) {
	// Tokenize query
	queryTerms := k.tokenize(query)
	if len(queryTerms) == 0 {
		return []vectorstore.Document{}, nil
	}

	// Fetch documents from store for BM25 scoring
	// Note: In production, use an inverted index to avoid full corpus scan
	// For now, we fetch a larger set and score them
	fetchLimit := topK * 10 // Fetch more docs than needed for better scoring
	if fetchLimit < 100 {
		fetchLimit = 100 // Minimum corpus size for meaningful BM25
	}

	// Create a dummy query vector for store.Search (not used for scoring)
	// We're using the vector store's Search to fetch documents with filters
	dummyVector := make([]float32, 1536) // Standard embedding dimension
	searchReq := &vectorstore.SearchRequest{
		Vector:   dummyVector,
		TopK:     fetchLimit,
		Filter:   filters,
		MinScore: 0.0, // Get all docs regardless of vector similarity
	}

	searchResp, err := k.store.Search(ctx, searchReq)
	if err != nil {
		return nil, err
	}

	allDocs := searchResp.Documents
	if len(allDocs) == 0 {
		return []vectorstore.Document{}, nil
	}

	// Build document frequency map for IDF calculation
	termDocFreq := make(map[string]int)
	for _, doc := range allDocs {
		docTerms := k.tokenize(doc.Content)
		seenTerms := make(map[string]bool)
		for _, term := range docTerms {
			if !seenTerms[term] {
				termDocFreq[term]++
				seenTerms[term] = true
			}
		}
	}

	// Calculate average document length
	avgDocLen := k.calculateAvgDocLength(allDocs)
	corpusSize := len(allDocs)

	// Score all documents using BM25
	scored := make([]scoredDoc, 0, len(allDocs))
	for _, doc := range allDocs {
		score := k.scoreBM25WithIDF(queryTerms, doc, avgDocLen, corpusSize, termDocFreq)
		if score > 0 {
			scored = append(scored, scoredDoc{doc: doc, score: score})
		}
	}

	// Sort by score descending
	for i := 0; i < len(scored); i++ {
		for j := i + 1; j < len(scored); j++ {
			if scored[j].score > scored[i].score {
				scored[i], scored[j] = scored[j], scored[i]
			}
		}
	}

	// Return top K results
	results := make([]vectorstore.Document, 0, topK)
	for i := 0; i < topK && i < len(scored); i++ {
		doc := scored[i].doc
		doc.Score = float32(scored[i].score) // Convert float64 to float32
		results = append(results, doc)
	}

	return results, nil
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

// scoreBM25 calculates BM25 score for a document (simplified IDF version).
// This is kept for backward compatibility with tests.
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

// scoreBM25WithIDF calculates BM25 score using actual document frequencies for IDF.
func (k *KeywordRetriever) scoreBM25WithIDF(queryTerms []string, doc vectorstore.Document, avgDocLen float64, corpusSize int, termDocFreq map[string]int) float64 {
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

		// Calculate proper IDF using document frequency
		df := termDocFreq[term]
		if df == 0 {
			df = 1 // Smoothing for unseen terms
		}
		// IDF = log((N - df + 0.5) / (df + 0.5) + 1)
		idf := math.Log((float64(corpusSize)-float64(df)+0.5)/(float64(df)+0.5) + 1.0)

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
