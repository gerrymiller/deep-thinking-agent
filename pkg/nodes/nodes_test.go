// Copyright 2025 Gerry Miller <gerry@gerrymiller.com>
//
// Licensed under the MIT License.
// See LICENSE file in the project root for full license information.

package nodes

import (
	"context"
	"testing"

	"deep-thinking-agent/pkg/agent"
	"deep-thinking-agent/pkg/llm"
)

// mockLLM implements llm.Provider for testing
type mockLLM struct{}

func (m *mockLLM) Complete(ctx context.Context, req *llm.CompletionRequest) (*llm.CompletionResponse, error) {
	return &llm.CompletionResponse{Content: "mock response", FinishReason: "stop", Model: "mock"}, nil
}

func (m *mockLLM) Name() string            { return "mock" }
func (m *mockLLM) ModelName() string       { return "mock-model" }
func (m *mockLLM) SupportsStreaming() bool { return false }

// TestLLMBasedNodes tests construction and naming of LLM-based nodes
func TestLLMBasedNodes(t *testing.T) {
	ctx := context.Background()
	mockLLMProvider := &mockLLM{}

	tests := []struct {
		name         string
		createNode   func() interface{}
		expectedName string
	}{
		{
			name: "PlannerNode",
			createNode: func() interface{} {
				planner := agent.NewPlanner(mockLLMProvider, nil)
				return NewPlannerNode(ctx, planner)
			},
			expectedName: "planner",
		},
		{
			name: "RewriterNode",
			createNode: func() interface{} {
				rewriter := agent.NewRewriter(mockLLMProvider, nil)
				return NewRewriterNode(ctx, rewriter)
			},
			expectedName: "rewriter",
		},
		{
			name: "SupervisorNode",
			createNode: func() interface{} {
				supervisor := agent.NewSupervisor(mockLLMProvider, nil)
				return NewSupervisorNode(ctx, supervisor)
			},
			expectedName: "supervisor",
		},
		{
			name: "RerankerNode",
			createNode: func() interface{} {
				reranker := agent.NewReranker(nil)
				return NewRerankerNode(ctx, reranker)
			},
			expectedName: "reranker",
		},
		{
			name: "DistillerNode",
			createNode: func() interface{} {
				distiller := agent.NewDistiller(mockLLMProvider, nil)
				return NewDistillerNode(ctx, distiller)
			},
			expectedName: "distiller",
		},
		{
			name: "ReflectorNode",
			createNode: func() interface{} {
				reflector := agent.NewReflector(mockLLMProvider, nil)
				return NewReflectorNode(ctx, reflector)
			},
			expectedName: "reflector",
		},
		{
			name: "PolicyNode",
			createNode: func() interface{} {
				policy := agent.NewPolicy(mockLLMProvider, nil)
				return NewPolicyNode(ctx, policy)
			},
			expectedName: "policy",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node := tt.createNode()
			if node == nil {
				t.Fatal("expected node, got nil")
			}

			// Check if node has Name() method
			type named interface {
				Name() string
			}

			if n, ok := node.(named); ok {
				if n.Name() != tt.expectedName {
					t.Errorf("expected name %s, got %s", tt.expectedName, n.Name())
				}
			} else {
				t.Error("node does not implement Name() method")
			}
		})
	}
}

// TestNodeNamesAreUnique ensures all node names are distinct
func TestNodeNamesAreUnique(t *testing.T) {
	ctx := context.Background()
	mockLLMProvider := &mockLLM{}

	// Create all LLM-based nodes
	planner := NewPlannerNode(ctx, agent.NewPlanner(mockLLMProvider, nil))
	rewriter := NewRewriterNode(ctx, agent.NewRewriter(mockLLMProvider, nil))
	supervisor := NewSupervisorNode(ctx, agent.NewSupervisor(mockLLMProvider, nil))
	reranker := NewRerankerNode(ctx, agent.NewReranker(nil))
	distiller := NewDistillerNode(ctx, agent.NewDistiller(mockLLMProvider, nil))
	reflector := NewReflectorNode(ctx, agent.NewReflector(mockLLMProvider, nil))
	policy := NewPolicyNode(ctx, agent.NewPolicy(mockLLMProvider, nil))

	// Verify all have unique names
	names := make(map[string]bool)
	allNodes := []interface{}{planner, rewriter, supervisor, reranker, distiller, reflector, policy}

	for _, node := range allNodes {
		type named interface {
			Name() string
		}
		if n, ok := node.(named); ok {
			name := n.Name()
			if names[name] {
				t.Errorf("duplicate node name: %s", name)
			}
			names[name] = true
		}
	}

	expectedNames := []string{"planner", "rewriter", "supervisor", "reranker", "distiller", "reflector", "policy"}
	if len(names) != len(expectedNames) {
		t.Errorf("expected %d unique names, got %d", len(expectedNames), len(names))
	}

	for _, expected := range expectedNames {
		if !names[expected] {
			t.Errorf("missing expected node name: %s", expected)
		}
	}
}

// Note: RetrieverNode testing requires complex interface mocking and is deferred to integration tests.
