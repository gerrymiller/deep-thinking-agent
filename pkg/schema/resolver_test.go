// Copyright 2025 Gerry Miller <gerry@gerrymiller.com>
//
// Licensed under the MIT License.
// See LICENSE file in the project root for full license information.

package schema

import (
	"context"
	"testing"
	"time"
)

func TestNewResolver(t *testing.T) {
	tests := []struct {
		name   string
		config *ResolverConfig
	}{
		{
			name:   "with nil config uses defaults",
			config: nil,
		},
		{
			name: "with custom config",
			config: &ResolverConfig{
				EnablePatternMatching: false,
				EnableLLMAnalysis:     false,
				EnableCaching:         false,
				CacheTTL:              1 * time.Hour,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			provider := &mockLLMProvider{response: "{}"}
			resolver := NewResolver(provider, tt.config)

			if resolver.analyzer == nil {
				t.Error("analyzer not initialized")
			}
			if resolver.registry == nil {
				t.Error("registry not initialized")
			}
			if tt.config == nil || tt.config.EnableCaching {
				if resolver.cache == nil {
					t.Error("cache should be initialized when caching enabled")
				}
			}
		})
	}
}

func TestResolveExplicitSchema(t *testing.T) {
	provider := &mockLLMProvider{response: "{}"}
	resolver := NewResolver(provider, nil)

	explicitSchema := &DocumentSchema{
		DocID:            "doc1",
		Format:           "text",
		ChunkingStrategy: "section_based",
	}

	result, err := resolver.Resolve(context.Background(), "doc1", "content", "text", explicitSchema)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Strategy != StrategyExplicit {
		t.Errorf("Strategy = %v, want %v", result.Strategy, StrategyExplicit)
	}
	if result.Confidence != 1.0 {
		t.Errorf("Confidence = %v, want 1.0", result.Confidence)
	}
	if result.Schema != explicitSchema {
		t.Error("Schema should be the same as provided")
	}
}

func TestResolvePatternMatching(t *testing.T) {
	provider := &mockLLMProvider{response: "{}"}
	resolver := NewResolver(provider, nil)

	// Register a pattern
	pattern := SchemaPattern{
		Name:        "test_pattern",
		Description: "Test pattern",
		Indicators:  []string{"ITEM 1A", "ITEM 1", "Form 10-K"},
		Priority:    100,
		Template: DocumentSchema{
			Format:           "text",
			ChunkingStrategy: "section_based",
		},
	}
	resolver.registry.Register(pattern)

	// Content that matches pattern (has 2 out of 3 indicators)
	content := "This document contains ITEM 1A and ITEM 1 sections."

	result, err := resolver.Resolve(context.Background(), "doc1", content, "text", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Strategy != StrategyPattern {
		t.Errorf("Strategy = %v, want %v", result.Strategy, StrategyPattern)
	}
	if result.PatternUsed != "test_pattern" {
		t.Errorf("PatternUsed = %v, want test_pattern", result.PatternUsed)
	}
	if result.Schema.DocID != "doc1" {
		t.Errorf("DocID = %v, want doc1", result.Schema.DocID)
	}
}

func TestResolveLLMAnalysis(t *testing.T) {
	validResponse := `{
		"title": "Test Document",
		"sections": [],
		"semantic_regions": [],
		"custom_attributes": {},
		"chunking_strategy": "sliding_window",
		"confidence": 0.85
	}`

	provider := &mockLLMProvider{response: validResponse}
	resolver := NewResolver(provider, nil)

	result, err := resolver.Resolve(context.Background(), "doc1", "content", "text", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Strategy != StrategyLLM {
		t.Errorf("Strategy = %v, want %v", result.Strategy, StrategyLLM)
	}
	if result.Schema.DocID != "doc1" {
		t.Errorf("DocID = %v, want doc1", result.Schema.DocID)
	}
}

func TestResolveWithCache(t *testing.T) {
	validResponse := `{
		"title": "Test Document",
		"sections": [],
		"semantic_regions": [],
		"custom_attributes": {},
		"chunking_strategy": "sliding_window",
		"confidence": 0.85
	}`

	provider := &mockLLMProvider{response: validResponse}
	resolver := NewResolver(provider, nil)

	// First call - should hit LLM
	result1, err := resolver.Resolve(context.Background(), "doc1", "content", "text", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if provider.callCount != 1 {
		t.Errorf("callCount = %v, want 1", provider.callCount)
	}

	// Second call with same docID - should hit cache
	result2, err := resolver.Resolve(context.Background(), "doc1", "different content", "text", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if provider.callCount != 1 {
		t.Errorf("callCount = %v, want 1 (should use cache)", provider.callCount)
	}
	if result2.Schema.DocID != result1.Schema.DocID {
		t.Error("cached result should have same DocID")
	}
}

func TestMatchPattern(t *testing.T) {
	provider := &mockLLMProvider{response: "{}"}
	resolver := NewResolver(provider, nil)

	pattern := &SchemaPattern{
		Name:        "test",
		Indicators:  []string{"keyword1", "keyword2", "keyword3"},
		Priority:    100,
	}

	tests := []struct {
		name            string
		content         string
		expectMatch     bool
		minConfidence   float32
	}{
		{
			name:          "no match",
			content:       "This content has none of the keywords",
			expectMatch:   false,
			minConfidence: 0.0,
		},
		{
			name:          "partial match below threshold",
			content:       "This has keyword1 only",
			expectMatch:   false,
			minConfidence: 0.0,
		},
		{
			name:          "match above threshold",
			content:       "This has keyword1 and keyword2",
			expectMatch:   true,
			minConfidence: 0.5,
		},
		{
			name:          "full match",
			content:       "This has keyword1, keyword2, and keyword3",
			expectMatch:   true,
			minConfidence: 0.9,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			matches, confidence := resolver.matchPattern(tt.content, "text", pattern)

			if matches != tt.expectMatch {
				t.Errorf("matches = %v, want %v", matches, tt.expectMatch)
			}
			if matches && confidence < tt.minConfidence {
				t.Errorf("confidence = %v, want >= %v", confidence, tt.minConfidence)
			}
		})
	}
}

func TestEnhanceWithLLM(t *testing.T) {
	provider := &mockLLMProvider{response: "{}"}
	resolver := NewResolver(provider, nil)

	schema := &DocumentSchema{
		DocID:  "doc1",
		Format: "text",
	}

	// Currently just returns schema as-is
	enhanced, err := resolver.enhanceWithLLM(context.Background(), schema, "content")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if enhanced != schema {
		t.Error("enhanced should be same as input (not yet implemented)")
	}
}

func TestContains(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		substr   string
		expected bool
	}{
		{
			name:     "empty text",
			text:     "",
			substr:   "test",
			expected: false,
		},
		{
			name:     "empty substr",
			text:     "test",
			substr:   "",
			expected: false,
		},
		{
			name:     "both non-empty",
			text:     "hello world",
			substr:   "world",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := contains(tt.text, tt.substr)
			if result != tt.expected {
				t.Errorf("contains() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestSchemaCache(t *testing.T) {
	cache := NewSchemaCache(100 * time.Millisecond)

	result := &ResolutionResult{
		Schema: &DocumentSchema{
			DocID: "doc1",
		},
		Strategy:   StrategyLLM,
		Confidence: 0.9,
	}

	// Test Set and Get
	cache.Set("doc1", result)
	cached := cache.Get("doc1")
	if cached == nil {
		t.Fatal("expected cached result, got nil")
	}
	if cached.Schema.DocID != "doc1" {
		t.Errorf("cached DocID = %v, want doc1", cached.Schema.DocID)
	}

	// Test Get non-existent
	notCached := cache.Get("doc2")
	if notCached != nil {
		t.Error("expected nil for non-existent key")
	}

	// Test expiration
	time.Sleep(150 * time.Millisecond)
	expired := cache.Get("doc1")
	if expired != nil {
		t.Error("expected nil for expired cache entry")
	}

	// Test Clear
	cache.Set("doc3", result)
	cache.Clear()
	cleared := cache.Get("doc3")
	if cleared != nil {
		t.Error("expected nil after clear")
	}
}
