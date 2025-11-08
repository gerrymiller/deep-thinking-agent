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

CRITICAL: Respond with ONLY valid JSON. Do not add markdown, explanations, or extra text.

JSON SCHEMA REQUIREMENTS:
- "dependencies" MUST be an array of integers: [0, 1, 2]
- Use empty array [] if no dependencies (NEVER use null, {}, or empty string)
- Each dependency is a step index (integer) that must complete first
- Example: "dependencies": [0] means this step depends on step 0 completing

Respond with valid JSON in this EXACT format:
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

// planStepJSON handles flexible JSON parsing for dependencies field.
// This supports various formats that different LLM models may generate.
type planStepJSON struct {
	Index           int             `json:"index"`
	SubQuestion     string          `json:"sub_question"`
	ToolType        string          `json:"tool_type"`
	SchemaHint      string          `json:"schema_hint"`
	ExpectedOutputs []string        `json:"expected_outputs"`
	Dependencies    json.RawMessage `json:"dependencies"` // Raw to handle multiple formats
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

	// Parse JSON with flexible dependencies field
	var parsed struct {
		Steps     []planStepJSON `json:"steps"`
		Reasoning string         `json:"reasoning"`
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
		// Parse dependencies with flexible format handling
		deps, err := parseDependencies(s.Dependencies)
		if err != nil {
			// Log warning but continue with empty dependencies
			fmt.Printf("Warning: failed to parse dependencies for step %d: %v (using empty array)\n", s.Index, err)
			deps = []int{}
		}

		plan.Steps[i] = workflow.PlanStep{
			Index:           s.Index,
			SubQuestion:     s.SubQuestion,
			ToolType:        s.ToolType,
			SchemaHint:      s.SchemaHint,
			ExpectedOutputs: s.ExpectedOutputs,
			Dependencies:    deps,
		}
	}

	return plan, nil
}

// parseDependencies converts various JSON formats to []int.
// Supports: arrays, empty objects, null, nested objects, and quoted strings.
// This handles variations across different LLM models (GPT-4o, GPT-5, Claude, etc.)
func parseDependencies(raw json.RawMessage) ([]int, error) {
	if len(raw) == 0 {
		return []int{}, nil
	}

	rawStr := string(raw)

	// Handle null
	if rawStr == "null" {
		return []int{}, nil
	}

	// Try array format first (correct format)
	var arr []int
	if err := json.Unmarshal(raw, &arr); err == nil {
		return arr, nil
	}

	// Try empty object {} (some models like gpt-4o may use this)
	var obj map[string]interface{}
	if err := json.Unmarshal(raw, &obj); err == nil {
		if len(obj) == 0 {
			return []int{}, nil
		}

		// Handle {"indices": [0,1]} format (Claude variations)
		if indices, ok := obj["indices"]; ok {
			if arr, ok := indices.([]interface{}); ok {
				return parseInterfaceArray(arr)
			}
		}

		// Handle numbered keys like {"0": [], "1": [0]}
		// Extract dependency values from map
		return extractDepsFromMap(obj)
	}

	// Try string (edge case - some models quote arrays as strings)
	var str string
	if err := json.Unmarshal(raw, &str); err == nil {
		return parseArrayString(str)
	}

	// Fallback: return empty array (safe default)
	return []int{}, nil
}

// parseInterfaceArray converts []interface{} to []int
func parseInterfaceArray(arr []interface{}) ([]int, error) {
	result := make([]int, 0, len(arr))
	for _, v := range arr {
		switch val := v.(type) {
		case float64:
			result = append(result, int(val))
		case int:
			result = append(result, val)
		case string:
			// Try to parse string as int
			var num int
			if _, err := fmt.Sscanf(val, "%d", &num); err == nil {
				result = append(result, num)
			}
		}
	}
	return result, nil
}

// extractDepsFromMap extracts dependency values from a map structure
func extractDepsFromMap(obj map[string]interface{}) ([]int, error) {
	result := []int{}
	for _, v := range obj {
		if arr, ok := v.([]interface{}); ok {
			if deps, err := parseInterfaceArray(arr); err == nil {
				result = append(result, deps...)
			}
		}
	}
	return result, nil
}

// parseArrayString parses a string like "[0,1,2]" into []int
func parseArrayString(str string) ([]int, error) {
	str = strings.TrimSpace(str)
	if !strings.HasPrefix(str, "[") || !strings.HasSuffix(str, "]") {
		return []int{}, nil
	}

	var arr []int
	if err := json.Unmarshal([]byte(str), &arr); err != nil {
		return []int{}, nil
	}
	return arr, nil
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
