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

// Policy makes decisions about whether to continue or finish the workflow.
// It evaluates progress, completeness, and iteration limits.
type Policy struct {
	llm         llm.Provider
	temperature float32
	maxTokens   int
}

// PolicyConfig contains configuration for the policy agent.
type PolicyConfig struct {
	Temperature float32
	MaxTokens   int
}

// NewPolicy creates a new policy agent.
func NewPolicy(llmProvider llm.Provider, config *PolicyConfig) *Policy {
	if config == nil {
		config = &PolicyConfig{
			Temperature: 0.3, // Low for consistent decisions
			MaxTokens:   500,
		}
	}

	return &Policy{
		llm:         llmProvider,
		temperature: config.Temperature,
		maxTokens:   config.MaxTokens,
	}
}

// Decide determines whether the workflow should continue or finish.
func (p *Policy) Decide(ctx context.Context, state *workflow.State) (*workflow.PolicyDecision, error) {
	if state == nil {
		return nil, fmt.Errorf("state is nil")
	}

	// Check hard limits first
	if state.IsComplete() {
		return &workflow.PolicyDecision{
			ShouldContinue: false,
			Reasoning:      "All plan steps completed",
			Confidence:     1.0,
		}, nil
	}

	if state.HasReachedMaxIterations() {
		return &workflow.PolicyDecision{
			ShouldContinue: false,
			Reasoning:      "Maximum iteration limit reached",
			Confidence:     1.0,
		}, nil
	}

	// Use LLM to evaluate progress
	prompt := p.buildPolicyPrompt(state)

	resp, err := p.llm.Complete(ctx, &llm.CompletionRequest{
		Messages: []llm.Message{
			{Role: "system", Content: systemPromptPolicy},
			{Role: "user", Content: prompt},
		},
		Temperature: p.temperature,
		MaxTokens:   p.maxTokens,
	})

	if err != nil {
		return nil, fmt.Errorf("LLM policy decision failed: %w", err)
	}

	decision := p.parsePolicyResponse(resp.Content)
	return decision, nil
}

// buildPolicyPrompt constructs the policy decision prompt.
func (p *Policy) buildPolicyPrompt(state *workflow.State) string {
	var builder strings.Builder

	builder.WriteString(fmt.Sprintf("Original question: %s\n\n", state.OriginalQuestion))

	if state.Plan != nil {
		builder.WriteString(fmt.Sprintf("Plan: %d steps total\n", len(state.Plan.Steps)))
		builder.WriteString(fmt.Sprintf("Completed: %d steps\n\n", len(state.PastSteps)))
	}

	builder.WriteString("Progress summary:\n")
	for i, step := range state.PastSteps {
		builder.WriteString(fmt.Sprintf("Step %d: %s\n", i+1, step.Summary))
	}

	builder.WriteString("\nDecide: Should the workflow continue to the next step, or is there sufficient information to answer the original question?")
	builder.WriteString("\n\nRespond in format:\nDECISION: continue OR finish\nREASONING: [explanation]\nCONFIDENCE: [0.0-1.0]")

	return builder.String()
}

// parsePolicyResponse extracts the policy decision.
func (p *Policy) parsePolicyResponse(response string) *workflow.PolicyDecision {
	lines := strings.Split(response, "\n")

	decision := &workflow.PolicyDecision{
		ShouldContinue: true, // Default to continue
		Confidence:     0.5,
	}

	for _, line := range lines {
		line = strings.TrimSpace(line)
		upper := strings.ToUpper(line)

		if strings.HasPrefix(upper, "DECISION:") {
			decisionText := strings.TrimSpace(strings.TrimPrefix(upper, "DECISION:"))
			if strings.Contains(decisionText, "FINISH") || strings.Contains(decisionText, "STOP") {
				decision.ShouldContinue = false
			}
		}

		if strings.HasPrefix(upper, "REASONING:") {
			decision.Reasoning = strings.TrimSpace(strings.TrimPrefix(line, "REASONING:"))
			decision.Reasoning = strings.TrimSpace(strings.TrimPrefix(decision.Reasoning, "Reasoning:"))
		}

		if strings.HasPrefix(upper, "CONFIDENCE:") {
			confidenceStr := strings.TrimSpace(strings.TrimPrefix(upper, "CONFIDENCE:"))
			var conf float32
			fmt.Sscanf(confidenceStr, "%f", &conf)
			if conf >= 0.0 && conf <= 1.0 {
				decision.Confidence = conf
			}
		}
	}

	return decision
}

const systemPromptPolicy = `You are a workflow control expert for a RAG system.

Your task is to decide whether the workflow should continue to the next step or finish.

Decision criteria:
- Continue if: More steps remain and would add valuable information
- Finish if: The original question can be adequately answered with current findings
- Finish if: Additional steps would be redundant or provide diminishing returns

Guidelines:
- Evaluate completeness of findings relative to the original question
- Consider the quality and relevance of information gathered
- Balance thoroughness with efficiency
- Be decisive - avoid unnecessary iterations

Respond in format:
DECISION: continue OR finish
REASONING: [clear explanation]
CONFIDENCE: [0.0-1.0]`
