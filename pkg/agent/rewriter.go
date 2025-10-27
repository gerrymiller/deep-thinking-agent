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

// Rewriter enhances queries for better retrieval using a fast LLM.
// It expands queries with synonyms, related terms, and contextual information.
type Rewriter struct {
	llm         llm.Provider
	temperature float32
	maxTokens   int
}

// RewriterConfig contains configuration for the rewriter agent.
type RewriterConfig struct {
	Temperature float32
	MaxTokens   int
}

// NewRewriter creates a new rewriter agent.
func NewRewriter(llmProvider llm.Provider, config *RewriterConfig) *Rewriter {
	if config == nil {
		config = &RewriterConfig{
			Temperature: 0.5,
			MaxTokens:   500,
		}
	}

	return &Rewriter{
		llm:         llmProvider,
		temperature: config.Temperature,
		maxTokens:   config.MaxTokens,
	}
}

// Rewrite enhances a query for better retrieval.
func (r *Rewriter) Rewrite(ctx context.Context, query string, state *workflow.State) (string, error) {
	// Build context from past steps if available
	contextInfo := ""
	if state != nil && len(state.PastSteps) > 0 {
		contextInfo = r.buildContextFromPastSteps(state.PastSteps)
	}

	prompt := r.buildRewritePrompt(query, contextInfo)

	resp, err := r.llm.Complete(ctx, &llm.CompletionRequest{
		Messages: []llm.Message{
			{Role: "system", Content: systemPromptRewriter},
			{Role: "user", Content: prompt},
		},
		Temperature: r.temperature,
		MaxTokens:   r.maxTokens,
	})

	if err != nil {
		return "", fmt.Errorf("LLM rewrite failed: %w", err)
	}

	// Extract rewritten query from response
	rewritten := strings.TrimSpace(resp.Content)
	if rewritten == "" {
		return query, nil // Fallback to original if empty
	}

	return rewritten, nil
}

// buildContextFromPastSteps creates context string from execution history.
func (r *Rewriter) buildContextFromPastSteps(pastSteps []workflow.PastStep) string {
	if len(pastSteps) == 0 {
		return ""
	}

	var builder strings.Builder
	builder.WriteString("Previous findings:\n")

	for i, step := range pastSteps {
		if i >= 3 { // Limit to last 3 steps
			break
		}
		builder.WriteString(fmt.Sprintf("- %s: %s\n", step.Step.SubQuestion, step.Summary))
	}

	return builder.String()
}

// buildRewritePrompt constructs the rewriting prompt.
func (r *Rewriter) buildRewritePrompt(query, context string) string {
	if context == "" {
		return fmt.Sprintf(`Rewrite the following query to be more effective for semantic search.

Original query: %s

Provide an enhanced version that:
- Expands key concepts with synonyms and related terms
- Adds contextual information that would help retrieval
- Maintains the core intent of the original query

Return only the rewritten query, nothing else.`, query)
	}

	return fmt.Sprintf(`Rewrite the following query to be more effective for semantic search, considering the execution context.

Original query: %s

%s

Provide an enhanced version that:
- Incorporates relevant context from previous findings
- Expands key concepts with synonyms and related terms
- Adds specific details that would help retrieval
- Maintains the core intent of the original query

Return only the rewritten query, nothing else.`, query, context)
}

const systemPromptRewriter = `You are a query enhancement specialist for a RAG system.

Your task is to rewrite queries to improve retrieval effectiveness.

Guidelines:
- Expand queries with synonyms, related terms, and domain-specific language
- Add contextual information that helps semantic search
- Keep queries concise but comprehensive
- Preserve the original intent
- Consider execution context from previous steps if provided

Return only the rewritten query without explanations or formatting.`
