// Copyright 2025 Gerry Miller <gerry@gerrymiller.com>
//
// Licensed under the MIT License.
// See LICENSE file in the project root for full license information.

package agent

import (
	"context"
	"fmt"
	"strings"

	"deep-thinking-agent/pkg/llm"
	"deep-thinking-agent/pkg/workflow"
)

// Supervisor selects the optimal retrieval strategy for each query.
// It analyzes the query characteristics and context to choose between
// vector, keyword, hybrid, or schema-filtered approaches.
type Supervisor struct {
	llm         llm.Provider
	temperature float32
	maxTokens   int
}

// SupervisorConfig contains configuration for the supervisor agent.
type SupervisorConfig struct {
	Temperature float32
	MaxTokens   int
}

// NewSupervisor creates a new supervisor agent.
func NewSupervisor(llmProvider llm.Provider, config *SupervisorConfig) *Supervisor {
	if config == nil {
		config = &SupervisorConfig{
			Temperature: 0.3, // Low for consistent decisions
			MaxTokens:   300,
		}
	}

	return &Supervisor{
		llm:         llmProvider,
		temperature: config.Temperature,
		maxTokens:   config.MaxTokens,
	}
}

// SelectStrategy determines the best retrieval strategy for a query.
func (s *Supervisor) SelectStrategy(ctx context.Context, query string, state *workflow.State) (workflow.RetrievalStrategy, error) {
	prompt := s.buildStrategyPrompt(query, state)

	resp, err := s.llm.Complete(ctx, &llm.CompletionRequest{
		Messages: []llm.Message{
			{Role: "system", Content: systemPromptSupervisor},
			{Role: "user", Content: prompt},
		},
		Temperature: s.temperature,
		MaxTokens:   s.maxTokens,
	})

	if err != nil {
		return workflow.StrategyHybrid, fmt.Errorf("LLM strategy selection failed: %w", err)
	}

	strategy := s.parseStrategyResponse(resp.Content)
	return strategy, nil
}

// buildStrategyPrompt constructs the strategy selection prompt.
func (s *Supervisor) buildStrategyPrompt(query string, state *workflow.State) string {
	contextInfo := ""
	if state != nil && state.CurrentStep() != nil {
		step := state.CurrentStep()
		contextInfo = fmt.Sprintf("\nTool type: %s\nSchema hint: %s", step.ToolType, step.SchemaHint)
	}

	return fmt.Sprintf(`Select the optimal retrieval strategy for this query.

Query: %s
%s

Available strategies:
- vector: Semantic similarity search (best for conceptual queries)
- keyword: BM25 keyword search (best for exact terms, names, specific facts)
- hybrid: Combination of vector and keyword (best for balanced queries)
- schema_filtered: Schema-aware targeted search (best when specific document sections are needed)

Return only the strategy name: vector, keyword, hybrid, or schema_filtered`, query, contextInfo)
}

// parseStrategyResponse extracts the strategy from the LLM response.
func (s *Supervisor) parseStrategyResponse(response string) workflow.RetrievalStrategy {
	response = strings.ToLower(strings.TrimSpace(response))

	// Check for strategy keywords
	if strings.Contains(response, "vector") && !strings.Contains(response, "hybrid") {
		return workflow.StrategyVector
	}
	if strings.Contains(response, "keyword") && !strings.Contains(response, "hybrid") {
		return workflow.StrategyKeyword
	}
	if strings.Contains(response, "schema") {
		return workflow.StrategySchemaFiltered
	}
	if strings.Contains(response, "hybrid") {
		return workflow.StrategyHybrid
	}

	// Default to hybrid if unclear
	return workflow.StrategyHybrid
}

const systemPromptSupervisor = `You are a retrieval strategy expert for a RAG system.

Your task is to select the most effective retrieval strategy based on query characteristics.

Strategy selection guidelines:
- vector: Use for conceptual, semantic, or exploratory queries
- keyword: Use for exact matches, specific names, identifiers, or factual lookups
- hybrid: Use for balanced queries that benefit from both semantic and keyword matching
- schema_filtered: Use when the query targets specific document sections or types

Consider:
- Query specificity (exact terms vs. concepts)
- Tool type hints from the execution plan
- Schema hints that suggest targeted retrieval

Return only the strategy name without explanation.`
