// Copyright 2025 Gerry Miller <gerry@gerrymiller.com>
//
// Licensed under the MIT License.
// See LICENSE file in the project root for full license information.

package retrieval

import (
	"context"
	"errors"
	"testing"

	"deep-thinking-agent/pkg/embedding"
	"deep-thinking-agent/pkg/vectorstore"
	"deep-thinking-agent/pkg/workflow"
)

// Mock Embedder
type mockEmbedder struct {
	embeddings [][]float32
	err        error
}

func (m *mockEmbedder) Embed(ctx context.Context, req *embedding.EmbedRequest) (*embedding.EmbedResponse, error) {
	if m.err != nil {
		return nil, m.err
	}

	// If embeddings is explicitly set (even to empty), use that
	if m.embeddings != nil {
		vectors := make([]embedding.Vector, len(m.embeddings))
		for i, emb := range m.embeddings {
			vectors[i] = embedding.Vector{
				Embedding: emb,
				Text:      req.Texts[i],
			}
		}
		return &embedding.EmbedResponse{
			Vectors: vectors,
			Usage:   embedding.UsageStats{PromptTokens: 10, TotalTokens: 10},
		}, nil
	}

	// Default behavior: generate embeddings
	vectors := make([]embedding.Vector, len(req.Texts))
	for i, text := range req.Texts {
		emb := make([]float32, 128)
		for j := range emb {
			emb[j] = 0.1
		}
		vectors[i] = embedding.Vector{
			Embedding: emb,
			Text:      text,
		}
	}

	return &embedding.EmbedResponse{
		Vectors: vectors,
		Usage:   embedding.UsageStats{PromptTokens: 10, TotalTokens: 10},
	}, nil
}

func (m *mockEmbedder) Dimensions() int   { return 128 }
func (m *mockEmbedder) ModelName() string { return "mock-embed" }

// Mock VectorStore
type mockVectorStore struct {
	searchResults []vectorstore.Document
	err           error
}

func (m *mockVectorStore) Search(ctx context.Context, req *vectorstore.SearchRequest) (*vectorstore.SearchResponse, error) {
	if m.err != nil {
		return nil, m.err
	}
	return &vectorstore.SearchResponse{
		Documents:    m.searchResults,
		TotalResults: len(m.searchResults),
	}, nil
}

func (m *mockVectorStore) Insert(ctx context.Context, req *vectorstore.InsertRequest) (*vectorstore.InsertResponse, error) {
	if m.err != nil {
		return nil, m.err
	}
	ids := make([]string, len(req.Documents))
	for i, doc := range req.Documents {
		ids[i] = doc.ID
	}
	return &vectorstore.InsertResponse{InsertedIDs: ids}, nil
}

func (m *mockVectorStore) Delete(ctx context.Context, req *vectorstore.DeleteRequest) (*vectorstore.DeleteResponse, error) {
	if m.err != nil {
		return nil, m.err
	}
	return &vectorstore.DeleteResponse{DeletedCount: len(req.IDs)}, nil
}

func (m *mockVectorStore) Get(ctx context.Context, collectionName string, ids []string) ([]vectorstore.Document, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.searchResults, nil
}

func (m *mockVectorStore) CreateCollection(ctx context.Context, name string, dimension int, metadata map[string]interface{}) error {
	return m.err
}

func (m *mockVectorStore) DeleteCollection(ctx context.Context, name string) error {
	return m.err
}

func (m *mockVectorStore) ListCollections(ctx context.Context) ([]vectorstore.CollectionInfo, error) {
	if m.err != nil {
		return nil, m.err
	}
	return []vectorstore.CollectionInfo{}, nil
}

func (m *mockVectorStore) GetCollection(ctx context.Context, name string) (*vectorstore.CollectionInfo, error) {
	if m.err != nil {
		return nil, m.err
	}
	return &vectorstore.CollectionInfo{Name: name}, nil
}

func (m *mockVectorStore) Close() error { return m.err }
func (m *mockVectorStore) Name() string { return "mock-store" }

// Vector Retriever Tests
func TestNewVectorRetriever(t *testing.T) {
	store := &mockVectorStore{}
	embedder := &mockEmbedder{}

	retriever := NewVectorRetriever(store, embedder)
	if retriever == nil {
		t.Fatal("NewVectorRetriever returned nil")
	}

	if retriever.Name() != "vector" {
		t.Errorf("Name() = %v, want vector", retriever.Name())
	}
}

func TestVectorSearch(t *testing.T) {
	tests := []struct {
		name      string
		store     *mockVectorStore
		embedder  *mockEmbedder
		query     string
		topK      int
		wantErr   bool
		wantCount int
	}{
		{
			name: "successful search",
			store: &mockVectorStore{
				searchResults: []vectorstore.Document{
					{ID: "doc1", Content: "content", Score: 0.9},
				},
			},
			embedder:  &mockEmbedder{},
			query:     "test query",
			topK:      10,
			wantErr:   false,
			wantCount: 1,
		},
		{
			name:     "embedder error",
			store:    &mockVectorStore{},
			embedder: &mockEmbedder{err: errors.New("embed error")},
			query:    "test",
			topK:     10,
			wantErr:  true,
		},
		{
			name:     "no embeddings generated",
			store:    &mockVectorStore{},
			embedder: &mockEmbedder{embeddings: [][]float32{}},
			query:    "test",
			topK:     10,
			wantErr:  true,
		},
		{
			name:     "vector store error",
			store:    &mockVectorStore{err: errors.New("search error")},
			embedder: &mockEmbedder{},
			query:    "test",
			topK:     10,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			retriever := NewVectorRetriever(tt.store, tt.embedder)
			results, err := retriever.Search(context.Background(), tt.query, tt.topK, nil)

			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if len(results) != tt.wantCount {
				t.Errorf("got %d results, want %d", len(results), tt.wantCount)
			}
		})
	}
}

// Keyword Retriever Tests
func TestNewKeywordRetriever(t *testing.T) {
	store := &mockVectorStore{}
	retriever := NewKeywordRetriever(store)

	if retriever == nil {
		t.Fatal("NewKeywordRetriever returned nil")
	}

	if retriever.Name() != "keyword" {
		t.Errorf("Name() = %v, want keyword", retriever.Name())
	}
}

func TestKeywordSearch(t *testing.T) {
	// Create mock store with sample documents
	store := &mockVectorStore{
		searchResults: []vectorstore.Document{
			{
				ID:      "doc1",
				Content: "Machine learning algorithms are used for pattern recognition and data analysis",
			},
			{
				ID:      "doc2",
				Content: "Deep learning networks use neural architectures for complex tasks",
			},
			{
				ID:      "doc3",
				Content: "Natural language processing enables computers to understand human language",
			},
		},
	}
	retriever := NewKeywordRetriever(store)

	t.Run("finds relevant documents", func(t *testing.T) {
		results, err := retriever.Search(context.Background(), "machine learning algorithms", 2, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Should find documents containing query terms
		if len(results) == 0 {
			t.Error("expected at least some results for keyword search")
		}

		// First result should have highest BM25 score
		if len(results) > 1 && results[0].Score < results[1].Score {
			t.Error("results should be sorted by score descending")
		}
	})

	t.Run("empty query returns empty results", func(t *testing.T) {
		results, err := retriever.Search(context.Background(), "", 10, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(results) != 0 {
			t.Errorf("expected 0 results for empty query, got %d", len(results))
		}
	})

	t.Run("no matching documents", func(t *testing.T) {
		results, err := retriever.Search(context.Background(), "quantum physics", 10, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Should return empty results when no documents match
		if len(results) != 0 {
			t.Errorf("expected 0 results for non-matching query, got %d", len(results))
		}
	})

	t.Run("respects topK limit", func(t *testing.T) {
		results, err := retriever.Search(context.Background(), "learning", 1, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(results) > 1 {
			t.Errorf("expected at most 1 result, got %d", len(results))
		}
	})
}

func TestTokenize(t *testing.T) {
	retriever := NewKeywordRetriever(&mockVectorStore{})

	tests := []struct {
		name    string
		text    string
		wantMin int
	}{
		{
			name:    "basic tokenization",
			text:    "This is a test",
			wantMin: 1, // "test" after stopword removal
		},
		{
			name:    "with stopwords",
			text:    "the and or but in on at",
			wantMin: 0, // all stopwords
		},
		{
			name:    "complex text",
			text:    "Risk factors include market volatility",
			wantMin: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokens := retriever.tokenize(tt.text)
			if len(tokens) < tt.wantMin {
				t.Errorf("got %d tokens, want at least %d", len(tokens), tt.wantMin)
			}
		})
	}
}

// Hybrid Retriever Tests
func TestNewHybridRetriever(t *testing.T) {
	store := &mockVectorStore{}
	embedder := &mockEmbedder{}

	vectorRet := NewVectorRetriever(store, embedder)
	keywordRet := NewKeywordRetriever(store)

	retriever := NewHybridRetriever(vectorRet, keywordRet)

	if retriever == nil {
		t.Fatal("NewHybridRetriever returned nil")
	}

	if retriever.Name() != "hybrid" {
		t.Errorf("Name() = %v, want hybrid", retriever.Name())
	}
}

func TestHybridSearch(t *testing.T) {
	// NOTE: Since keyword search is simplified in Phase 3,
	// hybrid search will only return vector results

	docs := []vectorstore.Document{
		{ID: "doc1", Content: "content 1", Score: 0.9},
		{ID: "doc2", Content: "content 2", Score: 0.8},
		{ID: "doc3", Content: "content 3", Score: 0.7},
	}

	store := &mockVectorStore{
		searchResults: docs,
	}
	embedder := &mockEmbedder{}

	vectorRet := NewVectorRetriever(store, embedder)
	keywordRet := NewKeywordRetriever(store)
	retriever := NewHybridRetriever(vectorRet, keywordRet)

	results, err := retriever.Search(context.Background(), "test query", 2, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(results) == 0 {
		t.Error("expected some results")
	}

	// Results should be scored with RRF
	if len(results) > 0 && results[0].Score == 0 {
		t.Error("documents should have RRF scores")
	}
}

func TestFuseRRF(t *testing.T) {
	vectorResults := []vectorstore.Document{
		{ID: "doc1", Score: 0.9},
		{ID: "doc2", Score: 0.8},
		{ID: "doc3", Score: 0.7},
	}

	keywordResults := []vectorstore.Document{
		{ID: "doc2", Score: 0.95},
		{ID: "doc1", Score: 0.85},
		{ID: "doc4", Score: 0.75},
	}

	store := &mockVectorStore{}
	embedder := &mockEmbedder{}
	vectorRet := NewVectorRetriever(store, embedder)
	keywordRet := NewKeywordRetriever(store)
	retriever := NewHybridRetriever(vectorRet, keywordRet)

	fused := retriever.fuseRRF(vectorResults, keywordResults)

	if len(fused) == 0 {
		t.Fatal("fused results are empty")
	}

	// doc1 and doc2 should rank highest (both appear in both lists with good ranks)
	// They have nearly identical RRF scores, so check both are in top 2
	topTwo := make(map[string]bool)
	if len(fused) >= 2 {
		topTwo[fused[0].ID] = true
		topTwo[fused[1].ID] = true

		if !topTwo["doc1"] || !topTwo["doc2"] {
			t.Errorf("expected doc1 and doc2 in top 2, got %s and %s", fused[0].ID, fused[1].ID)
		}
	}
}

// Schema Retriever Tests
func TestNewSchemaRetriever(t *testing.T) {
	store := &mockVectorStore{}
	embedder := &mockEmbedder{}
	vectorRet := NewVectorRetriever(store, embedder)

	retriever := NewSchemaRetriever(vectorRet)

	if retriever == nil {
		t.Fatal("NewSchemaRetriever returned nil")
	}

	if retriever.Name() != "schema_filtered" {
		t.Errorf("Name() = %v, want schema_filtered", retriever.Name())
	}
}

func TestSchemaSearch(t *testing.T) {
	docs := []vectorstore.Document{
		{ID: "doc1", Content: "content", Score: 0.9},
	}

	store := &mockVectorStore{searchResults: docs}
	embedder := &mockEmbedder{}
	vectorRet := NewVectorRetriever(store, embedder)
	retriever := NewSchemaRetriever(vectorRet)

	schemaFilters := &workflow.SchemaFilters{
		DocumentIDs:  []string{"doc1"},
		SectionTypes: []string{"risk_factors"},
		SemanticTags: []string{"risk"},
	}

	results, err := retriever.Search(context.Background(), "test query", 10, schemaFilters)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(results) == 0 {
		t.Error("no results returned")
	}
}

func TestBuildMetadataFilters(t *testing.T) {
	store := &mockVectorStore{}
	embedder := &mockEmbedder{}
	vectorRet := NewVectorRetriever(store, embedder)
	retriever := NewSchemaRetriever(vectorRet)

	tests := []struct {
		name          string
		schemaFilters *workflow.SchemaFilters
		wantNil       bool
		checkKeys     []string
	}{
		{
			name:          "nil filters",
			schemaFilters: nil,
			wantNil:       true,
		},
		{
			name: "with document IDs",
			schemaFilters: &workflow.SchemaFilters{
				DocumentIDs: []string{"doc1", "doc2"},
			},
			wantNil:   false,
			checkKeys: []string{"doc_id"},
		},
		{
			name: "with all filters",
			schemaFilters: &workflow.SchemaFilters{
				DocumentIDs:       []string{"doc1"},
				SectionTypes:      []string{"risk"},
				SemanticTags:      []string{"tag1"},
				MinRelevanceScore: 0.8,
				CustomAttributes: map[string]interface{}{
					"custom": "value",
				},
			},
			wantNil:   false,
			checkKeys: []string{"doc_id", "section_type", "semantic_tags", "min_score", "custom"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filters := retriever.buildMetadataFilters(tt.schemaFilters)

			if tt.wantNil {
				if filters != nil {
					t.Error("expected nil filters")
				}
				return
			}

			if filters == nil {
				t.Fatal("filters should not be nil")
			}

			for _, key := range tt.checkKeys {
				if _, exists := filters[key]; !exists {
					t.Errorf("filter key %s not found", key)
				}
			}
		})
	}
}

func TestSearchWithFilters(t *testing.T) {
	docs := []vectorstore.Document{
		{ID: "doc1", Content: "content", Score: 0.9},
	}

	store := &mockVectorStore{searchResults: docs}
	embedder := &mockEmbedder{}
	vectorRet := NewVectorRetriever(store, embedder)
	retriever := NewSchemaRetriever(vectorRet)

	filters := map[string]interface{}{
		"section_type": "risk_factors",
	}

	results, err := retriever.SearchWithFilters(context.Background(), "test", 10, filters)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(results) == 0 {
		t.Error("no results returned")
	}
}

func TestCalculateAvgDocLength(t *testing.T) {
	store := &mockVectorStore{}
	retriever := NewKeywordRetriever(store)

	t.Run("empty documents", func(t *testing.T) {
		docs := []vectorstore.Document{}
		avg := retriever.calculateAvgDocLength(docs)

		if avg != 0 {
			t.Errorf("expected 0 for empty docs, got %f", avg)
		}
	})

	t.Run("single document", func(t *testing.T) {
		docs := []vectorstore.Document{
			{Content: "hello world test"},
		}
		avg := retriever.calculateAvgDocLength(docs)

		if avg != 3.0 {
			t.Errorf("expected 3.0, got %f", avg)
		}
	})

	t.Run("multiple documents", func(t *testing.T) {
		docs := []vectorstore.Document{
			{Content: "hello world"},        // 2 tokens
			{Content: "test document here"}, // 3 tokens
			{Content: "one"},                // 1 token
		}
		avg := retriever.calculateAvgDocLength(docs)

		expected := 2.0 // (2+3+1)/3
		if avg != expected {
			t.Errorf("expected %f, got %f", expected, avg)
		}
	})
}

func TestScoreBM25(t *testing.T) {
	store := &mockVectorStore{}
	retriever := NewKeywordRetriever(store)

	t.Run("document with matching terms", func(t *testing.T) {
		queryTerms := []string{"test", "document"}
		doc := vectorstore.Document{
			Content: "this is a test document with some test content",
		}
		avgDocLen := 5.0
		corpusSize := 10

		score := retriever.scoreBM25(queryTerms, doc, avgDocLen, corpusSize)

		// Score should be positive since there are matching terms
		if score <= 0 {
			t.Errorf("expected positive score, got %f", score)
		}
	})

	t.Run("document with no matching terms", func(t *testing.T) {
		queryTerms := []string{"nonexistent", "terms"}
		doc := vectorstore.Document{
			Content: "this is a test document",
		}
		avgDocLen := 5.0
		corpusSize := 10

		score := retriever.scoreBM25(queryTerms, doc, avgDocLen, corpusSize)

		// Score should be 0 since no matching terms
		if score != 0 {
			t.Errorf("expected 0 score, got %f", score)
		}
	})

	t.Run("document with partial matches", func(t *testing.T) {
		queryTerms := []string{"test", "nonexistent"}
		doc := vectorstore.Document{
			Content: "this is a test document",
		}
		avgDocLen := 5.0
		corpusSize := 10

		score := retriever.scoreBM25(queryTerms, doc, avgDocLen, corpusSize)

		// Score should be positive but less than full match
		if score <= 0 {
			t.Errorf("expected positive score for partial match, got %f", score)
		}
	})
}
