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

// Reflector summarizes completed steps and extracts key findings.
// It creates concise summaries that inform future steps in the workflow.
type Reflector struct {
	llm         llm.Provider
	temperature float32
	maxTokens   int
}

// ReflectorConfig contains configuration for the reflector agent.
type ReflectorConfig struct {
	Temperature float32
	MaxTokens   int
}

// NewReflector creates a new reflector agent.
func NewReflector(llmProvider llm.Provider, config *ReflectorConfig) *Reflector {
	if config == nil {
		config = &ReflectorConfig{
			Temperature: 0.4,
			MaxTokens:   800,
		}
	}

	return &Reflector{
		llm:         llmProvider,
		temperature: config.Temperature,
		maxTokens:   config.MaxTokens,
	}
}

// Reflect generates a summary and extracts key findings from a completed step.
func (r *Reflector) Reflect(ctx context.Context, step *workflow.PlanStep, synthesizedContext string) (string, []string, error) {
	if step == nil {
		return "", nil, fmt.Errorf("step is nil")
	}

	prompt := r.buildReflectionPrompt(step, synthesizedContext)

	resp, err := r.llm.Complete(ctx, &llm.CompletionRequest{
		Messages: []llm.Message{
			{Role: "system", Content: systemPromptReflector},
			{Role: "user", Content: prompt},
		},
		Temperature: r.temperature,
		MaxTokens:   r.maxTokens,
	})

	if err != nil {
		return "", nil, fmt.Errorf("LLM reflection failed: %w", err)
	}

	summary, keyFindings := r.parseReflectionResponse(resp.Content)
	return summary, keyFindings, nil
}

// buildReflectionPrompt constructs the reflection prompt.
func (r *Reflector) buildReflectionPrompt(step *workflow.PlanStep, synthesizedContext string) string {
	return fmt.Sprintf(`Reflect on the completed execution step and synthesized findings.

Step question: %s
Expected outputs: %v

Synthesized context:
%s

Provide:
1. A concise summary (2-3 sentences) of what was found
2. A bulleted list of 3-5 key findings

Format your response as:
SUMMARY: [your summary here]

KEY FINDINGS:
- [finding 1]
- [finding 2]
- [finding 3]`, step.SubQuestion, step.ExpectedOutputs, synthesizedContext)
}

// parseReflectionResponse extracts summary and key findings.
func (r *Reflector) parseReflectionResponse(response string) (string, []string) {
	lines := strings.Split(response, "\n")

	var summary string
	var keyFindings []string
	inFindings := false

	for _, line := range lines {
		line = strings.TrimSpace(line)

		if strings.HasPrefix(strings.ToUpper(line), "SUMMARY:") {
			summary = strings.TrimSpace(strings.TrimPrefix(strings.ToUpper(line), "SUMMARY:"))
			summary = strings.TrimSpace(strings.TrimPrefix(line, "SUMMARY:"))
			summary = strings.TrimSpace(strings.TrimPrefix(summary, "Summary:"))
			continue
		}

		if strings.Contains(strings.ToUpper(line), "KEY FINDINGS") {
			inFindings = true
			continue
		}

		if inFindings && strings.HasPrefix(line, "-") {
			finding := strings.TrimSpace(strings.TrimPrefix(line, "-"))
			if finding != "" {
				keyFindings = append(keyFindings, finding)
			}
		} else if inFindings && line != "" && !strings.HasPrefix(line, "-") {
			// Also consider lines as part of summary if not bullets
			if summary == "" {
				summary = line
			}
		}
	}

	// Fallback: if no structured response, use entire response as summary
	if summary == "" {
		summary = strings.TrimSpace(response)
	}

	return summary, keyFindings
}

const systemPromptReflector = `You are a reflection and summarization expert for a RAG system.

Your task is to reflect on completed execution steps and extract key insights.

Guidelines:
- Provide a concise summary of what was found in this step
- Extract 3-5 specific key findings that answer the step's question
- Focus on actionable information that informs future steps
- Be precise and factual
- Follow the requested format

Always structure your response with:
SUMMARY: [2-3 sentence summary]

KEY FINDINGS:
- [specific finding 1]
- [specific finding 2]
- [specific finding 3]`
