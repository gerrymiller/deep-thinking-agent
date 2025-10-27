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
	"deep-thinking-agent/pkg/vectorstore"
)

// Distiller synthesizes retrieved chunks into coherent context.
// It compresses and distills information while preserving key insights.
type Distiller struct {
	llm         llm.Provider
	temperature float32
	maxTokens   int
}

// DistillerConfig contains configuration for the distiller agent.
type DistillerConfig struct {
	Temperature float32
	MaxTokens   int
}

// NewDistiller creates a new distiller agent.
func NewDistiller(llmProvider llm.Provider, config *DistillerConfig) *Distiller {
	if config == nil {
		config = &DistillerConfig{
			Temperature: 0.5,
			MaxTokens:   1500,
		}
	}

	return &Distiller{
		llm:         llmProvider,
		temperature: config.Temperature,
		maxTokens:   config.MaxTokens,
	}
}

// Distill synthesizes document chunks into coherent context.
func (d *Distiller) Distill(ctx context.Context, query string, docs []vectorstore.Document) (string, error) {
	if len(docs) == 0 {
		return "", fmt.Errorf("no documents to distill")
	}

	prompt := d.buildDistillationPrompt(query, docs)

	resp, err := d.llm.Complete(ctx, &llm.CompletionRequest{
		Messages: []llm.Message{
			{Role: "system", Content: systemPromptDistiller},
			{Role: "user", Content: prompt},
		},
		Temperature: d.temperature,
		MaxTokens:   d.maxTokens,
	})

	if err != nil {
		return "", fmt.Errorf("LLM distillation failed: %w", err)
	}

	return strings.TrimSpace(resp.Content), nil
}

// buildDistillationPrompt constructs the distillation prompt.
func (d *Distiller) buildDistillationPrompt(query string, docs []vectorstore.Document) string {
	var builder strings.Builder

	builder.WriteString(fmt.Sprintf("Query: %s\n\n", query))
	builder.WriteString("Retrieved documents:\n\n")

	for i, doc := range docs {
		builder.WriteString(fmt.Sprintf("--- Document %d (Score: %.3f) ---\n", i+1, doc.Score))
		builder.WriteString(doc.Content)
		builder.WriteString("\n\n")
	}

	builder.WriteString("Synthesize the above documents into a coherent, comprehensive summary that addresses the query. ")
	builder.WriteString("Include all relevant information while removing redundancy.")

	return builder.String()
}

const systemPromptDistiller = `You are an information synthesis expert for a RAG system.

Your task is to distill retrieved document chunks into coherent, comprehensive context.

Guidelines:
- Synthesize information from all provided documents
- Preserve key facts, findings, and insights
- Remove redundancy and irrelevant details
- Maintain accuracy - do not add information not present in the documents
- Organize information logically
- Be concise but comprehensive

Provide only the synthesized context without meta-commentary.`
