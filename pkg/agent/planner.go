// Copyright 2025 Gerry Miller <gerry@gerrymiller.com>
//
// Licensed under the MIT License.
// See LICENSE file in the project root for full license information.

package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"deep-thinking-agent/pkg/llm"
	"deep-thinking-agent/pkg/workflow"
)

// Planner decomposes complex queries into sequential execution plans.
// It uses a reasoning LLM to break down multi-hop questions into manageable steps.
type Planner struct {
	llm         llm.Provider
	temperature float32
	maxTokens   int
}

// PlannerConfig contains configuration for the planner agent.
type PlannerConfig struct {
	Temperature float32
	MaxTokens   int
}

// NewPlanner creates a new planner agent.
func NewPlanner(llmProvider llm.Provider, config *PlannerConfig) *Planner {
	if config == nil {
		config = &PlannerConfig{
			Temperature: 0.7, // Higher for creative planning
			MaxTokens:   2000,
		}
	}

	return &Planner{
		llm:         llmProvider,
		temperature: config.Temperature,
		maxTokens:   config.MaxTokens,
	}
}

// Plan decomposes a question into an execution plan.
func (p *Planner) Plan(ctx context.Context, question string) (*workflow.Plan, error) {
	prompt := p.buildPlanningPrompt(question)

	resp, err := p.llm.Complete(ctx, &llm.CompletionRequest{
		Messages: []llm.Message{
			{Role: "system", Content: systemPromptPlanner},
			{Role: "user", Content: prompt},
		},
		Temperature: p.temperature,
		MaxTokens:   p.maxTokens,
	})

	if err != nil {
		return nil, fmt.Errorf("LLM planning failed: %w", err)
	}

	plan, err := p.parsePlanResponse(resp.Content)
	if err != nil {
		return nil, fmt.Errorf("failed to parse plan: %w", err)
	}

	return plan, nil
}

// buildPlanningPrompt constructs the planning prompt.
func (p *Planner) buildPlanningPrompt(question string) string {
	return fmt.Sprintf(`Decompose the following question into a sequential execution plan.

Question: %s

Create a plan with 2-5 steps that can be executed independently. Each step should:
1. Answer a specific sub-question
2. Specify which tool to use (doc_search, web_search, or schema_filter)
3. Provide hints for schema-aware retrieval if applicable

Respond with valid JSON in this format:
{
  "steps": [
    {
      "index": 0,
      "sub_question": "What specific information does this step need?",
      "tool_type": "doc_search",
      "schema_hint": "focus on specific document sections",
      "expected_outputs": ["expected finding 1", "expected finding 2"],
      "dependencies": []
    }
  ],
  "reasoning": "Explain why this plan will effectively answer the question"
}`, question)
}

// parsePlanResponse parses the LLM's plan response.
func (p *Planner) parsePlanResponse(response string) (*workflow.Plan, error) {
	// Extract JSON from response
	jsonStart := strings.Index(response, "{")
	if jsonStart == -1 {
		return nil, fmt.Errorf("no JSON found in response")
	}

	jsonEnd := strings.LastIndex(response, "}")
	if jsonEnd == -1 {
		return nil, fmt.Errorf("no JSON found in response")
	}

	jsonStr := response[jsonStart : jsonEnd+1]

	// Parse JSON
	var parsed struct {
		Steps []struct {
			Index           int      `json:"index"`
			SubQuestion     string   `json:"sub_question"`
			ToolType        string   `json:"tool_type"`
			SchemaHint      string   `json:"schema_hint"`
			ExpectedOutputs []string `json:"expected_outputs"`
			Dependencies    []int    `json:"dependencies"`
		} `json:"steps"`
		Reasoning string `json:"reasoning"`
	}

	if err := json.Unmarshal([]byte(jsonStr), &parsed); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	// Convert to workflow.Plan
	plan := &workflow.Plan{
		Steps:     make([]workflow.PlanStep, len(parsed.Steps)),
		Reasoning: parsed.Reasoning,
	}

	for i, s := range parsed.Steps {
		plan.Steps[i] = workflow.PlanStep{
			Index:           s.Index,
			SubQuestion:     s.SubQuestion,
			ToolType:        s.ToolType,
			SchemaHint:      s.SchemaHint,
			ExpectedOutputs: s.ExpectedOutputs,
			Dependencies:    s.Dependencies,
		}
	}

	return plan, nil
}

const systemPromptPlanner = `You are an expert query planner for a deep-thinking RAG system.

Your task is to decompose complex, multi-hop questions into sequential execution plans.

Guidelines:
- Create 2-5 steps that build on each other
- Each step should have a clear sub-question
- Specify the appropriate tool: doc_search (internal documents), web_search (external), or schema_filter (targeted search)
- Provide schema hints to guide retrieval (e.g., "focus on methodology sections")
- List expected outputs to clarify what each step should find
- Indicate dependencies if a step requires information from previous steps

Always respond with valid JSON matching the requested format.`
