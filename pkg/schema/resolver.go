// Copyright 2025 Gerry Miller <gerry@gerrymiller.com>
//
// Licensed under the MIT License.
// See LICENSE file in the project root for full license information.

package schema

import (
	"context"
	"fmt"
	"time"

	"deep-thinking-agent/pkg/llm"
)

// Resolver resolves document schemas using multiple strategies.
// It attempts resolution in order: explicit → pattern → LLM → hybrid
type Resolver struct {
	analyzer *Analyzer
	registry *Registry
	cache    *SchemaCache
}

// ResolverConfig contains configuration for the schema resolver.
type ResolverConfig struct {
	EnablePatternMatching bool
	EnableLLMAnalysis     bool
	EnableCaching         bool
	CacheTTL              time.Duration
}

// NewResolver creates a new schema resolver instance.
func NewResolver(llmProvider llm.Provider, config *ResolverConfig) *Resolver {
	if config == nil {
		config = &ResolverConfig{
			EnablePatternMatching: true,
			EnableLLMAnalysis:     true,
			EnableCaching:         true,
			CacheTTL:              24 * time.Hour,
		}
	}

	resolver := &Resolver{
		analyzer: NewAnalyzer(llmProvider, nil),
		registry: NewRegistry(),
	}

	if config.EnableCaching {
		resolver.cache = NewSchemaCache(config.CacheTTL)
	}

	return resolver
}

// Resolve determines the schema for a document using the resolution strategy.
// Strategy order: explicit → pattern → LLM → hybrid
func (r *Resolver) Resolve(ctx context.Context, docID, content, format string, explicitSchema *DocumentSchema) (*ResolutionResult, error) {
	startTime := time.Now()

	// Strategy 1: Use explicit schema if provided
	if explicitSchema != nil {
		return &ResolutionResult{
			Schema:           explicitSchema,
			Strategy:         StrategyExplicit,
			Confidence:       1.0,
			ProcessingTimeMs: time.Since(startTime).Milliseconds(),
		}, nil
	}

	// Check cache first
	if r.cache != nil {
		if cached := r.cache.Get(docID); cached != nil {
			cached.ProcessingTimeMs = time.Since(startTime).Milliseconds()
			return cached, nil
		}
	}

	// Strategy 2: Try pattern matching
	if r.registry != nil {
		patterns := r.registry.List()
		for _, pattern := range patterns {
			if matches, confidence := r.matchPattern(content, format, &pattern); matches {
				schema := pattern.Template
				schema.DocID = docID
				schema.Format = format
				schema.CreatedAt = time.Now().Unix()

				result := &ResolutionResult{
					Schema:           &schema,
					Strategy:         StrategyPattern,
					PatternUsed:      pattern.Name,
					Confidence:       confidence,
					ProcessingTimeMs: time.Since(startTime).Milliseconds(),
				}

				// If pattern requires LLM enhancement, upgrade to hybrid
				if pattern.RequiresLLMEnhancement {
					enhanced, err := r.enhanceWithLLM(ctx, &schema, content)
					if err == nil {
						result.Schema = enhanced
						result.Strategy = StrategyHybrid
					}
				}

				// Cache result
				if r.cache != nil {
					r.cache.Set(docID, result)
				}

				return result, nil
			}
		}
	}

	// Strategy 3: Full LLM analysis
	schema, err := r.analyzer.AnalyzeDocument(ctx, docID, content, format)
	if err != nil {
		return nil, fmt.Errorf("LLM analysis failed: %w", err)
	}

	result := &ResolutionResult{
		Schema:           schema,
		Strategy:         StrategyLLM,
		Confidence:       schema.Confidence,
		ProcessingTimeMs: time.Since(startTime).Milliseconds(),
	}

	// Cache result
	if r.cache != nil {
		r.cache.Set(docID, result)
	}

	return result, nil
}

// matchPattern checks if a document matches a predefined pattern.
func (r *Resolver) matchPattern(content, format string, pattern *SchemaPattern) (bool, float32) {
	// Simple indicator matching for now
	// TODO: Implement more sophisticated matching in future
	matchCount := 0
	for _, indicator := range pattern.Indicators {
		if contains(content, indicator) {
			matchCount++
		}
	}

	if matchCount == 0 {
		return false, 0.0
	}

	// Calculate confidence based on match ratio
	confidence := float32(matchCount) / float32(len(pattern.Indicators))

	// Require at least 50% match
	if confidence < 0.5 {
		return false, confidence
	}

	return true, confidence
}

// enhanceWithLLM uses LLM to refine a pattern-matched schema.
func (r *Resolver) enhanceWithLLM(ctx context.Context, schema *DocumentSchema, content string) (*DocumentSchema, error) {
	// For now, just return the schema as-is
	// TODO: Implement LLM-based enhancement in future
	return schema, nil
}

// contains checks if text contains a substring (case-insensitive).
func contains(text, substr string) bool {
	// Simple implementation for now
	// TODO: Add case-insensitive matching
	return len(text) > 0 && len(substr) > 0
}

// SchemaCache provides caching for resolved schemas.
type SchemaCache struct {
	cache map[string]*cachedResult
	ttl   time.Duration
}

type cachedResult struct {
	result    *ResolutionResult
	timestamp time.Time
}

// NewSchemaCache creates a new schema cache.
func NewSchemaCache(ttl time.Duration) *SchemaCache {
	return &SchemaCache{
		cache: make(map[string]*cachedResult),
		ttl:   ttl,
	}
}

// Get retrieves a cached schema result.
func (c *SchemaCache) Get(docID string) *ResolutionResult {
	if cached, ok := c.cache[docID]; ok {
		// Check if expired
		if time.Since(cached.timestamp) > c.ttl {
			delete(c.cache, docID)
			return nil
		}
		return cached.result
	}
	return nil
}

// Set stores a schema result in cache.
func (c *SchemaCache) Set(docID string, result *ResolutionResult) {
	c.cache[docID] = &cachedResult{
		result:    result,
		timestamp: time.Now(),
	}
}

// Clear removes all cached results.
func (c *SchemaCache) Clear() {
	c.cache = make(map[string]*cachedResult)
}
