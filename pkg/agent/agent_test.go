// Copyright 2025 Gerry Miller <gerry@gerrymiller.com>
//
// Licensed under the MIT License.
// See LICENSE file in the project root for full license information.

package agent

import (
	"context"
	"errors"
	"strings"
	"testing"

	"deep-thinking-agent/pkg/embedding"
	"deep-thinking-agent/pkg/llm"
	"deep-thinking-agent/pkg/vectorstore"
	"deep-thinking-agent/pkg/workflow"
)

// Mock LLM Provider
type mockLLMProvider struct {
	response string
	err      error
}

func (m *mockLLMProvider) Complete(ctx context.Context, req *llm.CompletionRequest) (*llm.CompletionResponse, error) {
	if m.err != nil {
		return nil, m.err
	}
	return &llm.CompletionResponse{
		Content: m.response,
		Usage:   llm.UsageStats{PromptTokens: 10, CompletionTokens: 20, TotalTokens: 30},
	}, nil
}

func (m *mockLLMProvider) Name() string            { return "mock" }
func (m *mockLLMProvider) ModelName() string       { return "mock-model" }
func (m *mockLLMProvider) SupportsStreaming() bool { return false }

// Mock Embedder
type mockEmbedder struct {
	embeddings [][]float32
	err        error
}

func (m *mockEmbedder) Embed(ctx context.Context, req *embedding.EmbedRequest) (*embedding.EmbedResponse, error) {
	if m.err != nil {
		return nil, m.err
	}

	vectors := make([]embedding.Vector, len(req.Texts))
	for i, text := range req.Texts {
		var emb []float32
		if m.embeddings != nil && i < len(m.embeddings) {
			emb = m.embeddings[i]
		} else {
			// Default embedding
			emb = make([]float32, 128)
			for j := range emb {
				emb[j] = 0.1
			}
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

// Planner Tests
func TestNewPlanner(t *testing.T) {
	provider := &mockLLMProvider{}

	tests := []struct {
		name   string
		config *PlannerConfig
	}{
		{"with nil config", nil},
		{"with custom config", &PlannerConfig{Temperature: 0.5, MaxTokens: 1000}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			planner := NewPlanner(provider, tt.config)
			if planner == nil {
				t.Fatal("NewPlanner returned nil")
			}
			if planner.llm == nil {
				t.Error("LLM provider not set")
			}
		})
	}
}

func TestPlan(t *testing.T) {
	validResponse := `{
		"steps": [
			{
				"index": 0,
				"sub_question": "What are the risks?",
				"tool_type": "doc_search",
				"schema_hint": "focus on risk sections",
				"expected_outputs": ["risk factors"],
				"dependencies": []
			}
		],
		"reasoning": "Test reasoning"
	}`

	tests := []struct {
		name     string
		provider *mockLLMProvider
		question string
		wantErr  bool
	}{
		{
			name:     "successful planning",
			provider: &mockLLMProvider{response: validResponse},
			question: "What are the main risks?",
			wantErr:  false,
		},
		{
			name:     "LLM error",
			provider: &mockLLMProvider{err: errors.New("API error")},
			question: "Test question",
			wantErr:  true,
		},
		{
			name:     "invalid JSON response",
			provider: &mockLLMProvider{response: "not json"},
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			planner := NewPlanner(tt.provider, nil)
			plan, err := planner.Plan(context.Background(), tt.question)

			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if len(plan.Steps) == 0 {
				t.Error("plan has no steps")
			}
		})
	}
}

// Rewriter Tests
func TestNewRewriter(t *testing.T) {
	provider := &mockLLMProvider{}
	rewriter := NewRewriter(provider, nil)
	if rewriter == nil {
		t.Fatal("NewRewriter returned nil")
	}
}

func TestRewrite(t *testing.T) {
	tests := []struct {
		name     string
		provider *mockLLMProvider
		query    string
		state    *workflow.State
		wantErr  bool
	}{
		{
			name:     "successful rewrite",
			provider: &mockLLMProvider{response: "enhanced query with synonyms"},
			query:    "original query",
			state:    nil,
			wantErr:  false,
		},
		{
			name:     "with past steps context",
			provider: &mockLLMProvider{response: "context-aware query"},
			query:    "test query",
			state: &workflow.State{
				PastSteps: []workflow.PastStep{
					{Summary: "Previous finding"},
				},
			},
			wantErr: false,
		},
		{
			name:     "LLM error",
			provider: &mockLLMProvider{err: errors.New("error")},
			query:    "test",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rewriter := NewRewriter(tt.provider, nil)
			result, err := rewriter.Rewrite(context.Background(), tt.query, tt.state)

			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if result == "" {
				t.Error("rewritten query is empty")
			}
		})
	}
}

// Supervisor Tests
func TestNewSupervisor(t *testing.T) {
	provider := &mockLLMProvider{}
	supervisor := NewSupervisor(provider, nil)
	if supervisor == nil {
		t.Fatal("NewSupervisor returned nil")
	}
}

func TestSelectStrategy(t *testing.T) {
	tests := []struct {
		name             string
		response         string
		expectedStrategy workflow.RetrievalStrategy
	}{
		{"vector strategy", "vector", workflow.StrategyVector},
		{"keyword strategy", "keyword", workflow.StrategyKeyword},
		{"hybrid strategy", "hybrid", workflow.StrategyHybrid},
		{"schema_filtered strategy", "schema_filtered", workflow.StrategySchemaFiltered},
		{"default to hybrid", "unclear", workflow.StrategyHybrid},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			provider := &mockLLMProvider{response: tt.response}
			supervisor := NewSupervisor(provider, nil)

			strategy, err := supervisor.SelectStrategy(context.Background(), "test query", nil)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if strategy != tt.expectedStrategy {
				t.Errorf("got %v, want %v", strategy, tt.expectedStrategy)
			}
		})
	}
}

// Retriever Tests
func TestNewRetriever(t *testing.T) {
	store := &mockVectorStore{}
	embedder := &mockEmbedder{}
	retriever := NewRetriever(store, embedder, nil)
	if retriever == nil {
		t.Fatal("NewRetriever returned nil")
	}
}

func TestRetrieve(t *testing.T) {
	tests := []struct {
		name     string
		store    *mockVectorStore
		embedder *mockEmbedder
		ctx      *workflow.RetrievalContext
		wantErr  bool
	}{
		{
			name: "successful retrieval",
			store: &mockVectorStore{
				searchResults: []vectorstore.Document{
					{ID: "doc1", Content: "content", Score: 0.9},
				},
			},
			embedder: &mockEmbedder{},
			ctx: &workflow.RetrievalContext{
				Query: "test query",
				TopK:  10,
			},
			wantErr: false,
		},
		{
			name:     "nil context",
			store:    &mockVectorStore{},
			embedder: &mockEmbedder{},
			ctx:      nil,
			wantErr:  true,
		},
		{
			name:     "embedder error",
			store:    &mockVectorStore{},
			embedder: &mockEmbedder{err: errors.New("embed error")},
			ctx:      &workflow.RetrievalContext{Query: "test", TopK: 10},
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			retriever := NewRetriever(tt.store, tt.embedder, nil)
			docs, err := retriever.Retrieve(context.Background(), tt.ctx)

			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if len(docs) == 0 {
				t.Error("no documents returned")
			}
		})
	}
}

// Reranker Tests
func TestNewReranker(t *testing.T) {
	reranker := NewReranker(nil)
	if reranker == nil {
		t.Fatal("NewReranker returned nil")
	}
}

func TestRerank(t *testing.T) {
	docs := []vectorstore.Document{
		{ID: "doc1", Score: 0.5},
		{ID: "doc2", Score: 0.9},
		{ID: "doc3", Score: 0.7},
	}

	reranker := NewReranker(&RerankerConfig{TopN: 2})
	reranked := reranker.Rerank(context.Background(), "query", docs)

	if len(reranked) != 2 {
		t.Errorf("expected 2 docs, got %d", len(reranked))
	}

	if reranked[0].ID != "doc2" {
		t.Errorf("first doc should be doc2, got %s", reranked[0].ID)
	}
}

// Distiller Tests
func TestNewDistiller(t *testing.T) {
	provider := &mockLLMProvider{}
	distiller := NewDistiller(provider, nil)
	if distiller == nil {
		t.Fatal("NewDistiller returned nil")
	}
}

func TestDistill(t *testing.T) {
	provider := &mockLLMProvider{response: "Synthesized context"}
	distiller := NewDistiller(provider, nil)

	docs := []vectorstore.Document{
		{ID: "doc1", Content: "content 1", Score: 0.9},
		{ID: "doc2", Content: "content 2", Score: 0.8},
	}

	result, err := distiller.Distill(context.Background(), "query", docs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result == "" {
		t.Error("distilled context is empty")
	}

	// Test error cases
	_, err = distiller.Distill(context.Background(), "query", []vectorstore.Document{})
	if err == nil {
		t.Error("expected error for empty docs")
	}
}

// Reflector Tests
func TestNewReflector(t *testing.T) {
	provider := &mockLLMProvider{}
	reflector := NewReflector(provider, nil)
	if reflector == nil {
		t.Fatal("NewReflector returned nil")
	}
}

func TestReflect(t *testing.T) {
	response := `SUMMARY: This step found key risk factors.

KEY FINDINGS:
- Risk factor 1
- Risk factor 2
- Risk factor 3`

	provider := &mockLLMProvider{response: response}
	reflector := NewReflector(provider, nil)

	step := &workflow.PlanStep{
		SubQuestion:     "What are the risks?",
		ExpectedOutputs: []string{"risks"},
	}

	summary, findings, err := reflector.Reflect(context.Background(), step, "synthesized context")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if summary == "" {
		t.Error("summary is empty")
	}

	if len(findings) == 0 {
		t.Error("no key findings extracted")
	}

	// Test error case
	_, _, err = reflector.Reflect(context.Background(), nil, "context")
	if err == nil {
		t.Error("expected error for nil step")
	}
}

// Policy Tests
func TestNewPolicy(t *testing.T) {
	provider := &mockLLMProvider{}
	policy := NewPolicy(provider, nil)
	if policy == nil {
		t.Fatal("NewPolicy returned nil")
	}
}

func TestDecide(t *testing.T) {
	tests := []struct {
		name             string
		state            *workflow.State
		response         string
		expectedContinue bool
		wantErr          bool
	}{
		{
			name: "plan complete",
			state: &workflow.State{
				Plan: &workflow.Plan{
					Steps: []workflow.PlanStep{{}, {}},
				},
				CurrentStepIndex: 2,
			},
			expectedContinue: false,
			wantErr:          false,
		},
		{
			name: "max iterations reached",
			state: &workflow.State{
				MaxIterations: 2,
				PastSteps:     []workflow.PastStep{{}, {}},
			},
			expectedContinue: false,
			wantErr:          false,
		},
		{
			name: "LLM says continue",
			state: &workflow.State{
				OriginalQuestion: "test",
				Plan:             &workflow.Plan{Steps: []workflow.PlanStep{{}, {}, {}}},
				CurrentStepIndex: 1,
				MaxIterations:    10, // Set max iterations
			},
			response:         "DECISION: continue\nREASONING: More steps needed\nCONFIDENCE: 0.8",
			expectedContinue: true,
			wantErr:          false,
		},
		{
			name: "LLM says finish",
			state: &workflow.State{
				OriginalQuestion: "test",
				Plan:             &workflow.Plan{Steps: []workflow.PlanStep{{}, {}}},
				CurrentStepIndex: 1,
				MaxIterations:    10, // Set max iterations
			},
			response:         "DECISION: finish\nREASONING: Sufficient information\nCONFIDENCE: 0.9",
			expectedContinue: false,
			wantErr:          false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var provider *mockLLMProvider
			if tt.response != "" {
				provider = &mockLLMProvider{response: tt.response}
			} else {
				provider = &mockLLMProvider{response: "DECISION: continue"}
			}

			policy := NewPolicy(provider, nil)
			decision, err := policy.Decide(context.Background(), tt.state)

			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if decision.ShouldContinue != tt.expectedContinue {
				t.Errorf("ShouldContinue = %v, want %v", decision.ShouldContinue, tt.expectedContinue)
			}
		})
	}

	// Test nil state error
	provider := &mockLLMProvider{}
	policy := NewPolicy(provider, nil)
	_, err := policy.Decide(context.Background(), nil)
	if err == nil {
		t.Error("expected error for nil state")
	}
}

// Helper function tests
func TestBuildContextFromPastSteps(t *testing.T) {
	rewriter := NewRewriter(&mockLLMProvider{}, nil)

	steps := []workflow.PastStep{
		{Step: workflow.PlanStep{SubQuestion: "Q1"}, Summary: "S1"},
		{Step: workflow.PlanStep{SubQuestion: "Q2"}, Summary: "S2"},
	}

	context := rewriter.buildContextFromPastSteps(steps)
	if context == "" {
		t.Error("context is empty")
	}

	if !strings.Contains(context, "Q1") {
		t.Error("context missing question 1")
	}
}

func TestRerankWithScores(t *testing.T) {
	reranker := NewReranker(&RerankerConfig{TopN: 2})

	docs := []vectorstore.Document{
		{ID: "doc1", Content: "content 1", Score: 0.5},
		{ID: "doc2", Content: "content 2", Score: 0.9},
		{ID: "doc3", Content: "content 3", Score: 0.7},
	}

	reranked := reranker.RerankWithScores(context.Background(), "test query", docs)

	if len(reranked) != 2 {
		t.Errorf("expected 2 documents, got %d", len(reranked))
	}

	// Should return top 2 by score
	if reranked[0].ID != "doc2" {
		t.Errorf("expected doc2 first, got %s", reranked[0].ID)
	}
}

func TestBuildMetadataFilters(t *testing.T) {
	retriever := NewRetriever(&mockVectorStore{}, &mockEmbedder{}, nil)

	tests := []struct {
		name     string
		filters  *workflow.SchemaFilters
		expected map[string]interface{}
	}{
		{
			name:     "nil filters",
			filters:  nil,
			expected: nil,
		},
		{
			name: "document IDs only",
			filters: &workflow.SchemaFilters{
				DocumentIDs: []string{"doc1", "doc2"},
			},
			expected: map[string]interface{}{
				"doc_id": []string{"doc1", "doc2"},
			},
		},
		{
			name: "section types only",
			filters: &workflow.SchemaFilters{
				SectionTypes: []string{"methodology", "results"},
			},
			expected: map[string]interface{}{
				"section_type": []string{"methodology", "results"},
			},
		},
		{
			name: "semantic tags only",
			filters: &workflow.SchemaFilters{
				SemanticTags: []string{"important", "key-finding"},
			},
			expected: map[string]interface{}{
				"semantic_tags": []string{"important", "key-finding"},
			},
		},
		{
			name: "custom attributes only",
			filters: &workflow.SchemaFilters{
				CustomAttributes: map[string]interface{}{
					"year":   2025,
					"author": "Smith",
				},
			},
			expected: map[string]interface{}{
				"year":   2025,
				"author": "Smith",
			},
		},
		{
			name: "all filters combined",
			filters: &workflow.SchemaFilters{
				DocumentIDs:  []string{"doc1"},
				SectionTypes: []string{"abstract"},
				SemanticTags: []string{"summary"},
				CustomAttributes: map[string]interface{}{
					"lang": "en",
				},
			},
			expected: map[string]interface{}{
				"doc_id":        []string{"doc1"},
				"section_type":  []string{"abstract"},
				"semantic_tags": []string{"summary"},
				"lang":          "en",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := retriever.buildMetadataFilters(tt.filters)

			if tt.expected == nil {
				if result != nil {
					t.Errorf("expected nil, got %v", result)
				}
				return
			}

			if result == nil {
				t.Fatal("expected non-nil result")
			}

			if len(result) != len(tt.expected) {
				t.Errorf("expected %d filters, got %d", len(tt.expected), len(result))
			}

			for key := range tt.expected {
				if result[key] == nil {
					t.Errorf("missing key %s", key)
				}
			}

			// Verify all keys from result exist in expected
			for key := range result {
				if tt.expected[key] == nil {
					t.Errorf("unexpected key %s in result", key)
				}
			}
		})
	}
}

func TestBuildStrategyPrompt(t *testing.T) {
	supervisor := NewSupervisor(&mockLLMProvider{}, nil)

	state := &workflow.State{
		OriginalQuestion: "What is the capital of France?",
		PastSteps: []workflow.PastStep{
			{Step: workflow.PlanStep{SubQuestion: "Q1"}, Summary: "S1"},
		},
	}

	prompt := supervisor.buildStrategyPrompt("test query", state)

	if prompt == "" {
		t.Error("prompt should not be empty")
	}

	if !strings.Contains(prompt, "test query") {
		t.Error("prompt should contain the query")
	}

	// Verify prompt contains strategy options
	if !strings.Contains(prompt, "vector") {
		t.Error("prompt should contain 'vector' strategy")
	}
	if !strings.Contains(prompt, "keyword") {
		t.Error("prompt should contain 'keyword' strategy")
	}
	if !strings.Contains(prompt, "hybrid") {
		t.Error("prompt should contain 'hybrid' strategy")
	}
}
