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
	"deep-thinking-agent/pkg/vectorstore"
	"deep-thinking-agent/pkg/workflow"
)

// mockLLM implements llm.Provider for testing
type mockLLM struct{}

func (m *mockLLM) Complete(ctx context.Context, req *llm.CompletionRequest) (*llm.CompletionResponse, error) {
	// Return valid JSON for planner agent
	content := `{
		"steps": [
			{"sub_question": "Step 1", "rationale": "First step"},
			{"sub_question": "Step 2", "rationale": "Second step"}
		]
	}`
	return &llm.CompletionResponse{Content: content, FinishReason: "stop", Model: "mock"}, nil
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

func TestPlannerNode_Execute(t *testing.T) {
	ctx := context.Background()
	mockLLMProvider := &mockLLM{}
	planner := agent.NewPlanner(mockLLMProvider, nil)
	node := NewPlannerNode(ctx, planner)

	t.Run("successful planning", func(t *testing.T) {
		state := &workflow.State{
			OriginalQuestion: "What is the capital of France?",
		}

		result, err := node.Execute(state)
		if err != nil {
			t.Fatalf("Execute() failed: %v", err)
		}

		if result == nil {
			t.Fatal("expected result, got nil")
		}

		if result.UpdatedState == nil {
			t.Fatal("expected updated state, got nil")
		}

		if result.UpdatedState.Plan == nil {
			t.Error("expected plan to be set in state")
		}
	})

	t.Run("node name", func(t *testing.T) {
		if node.Name() != "planner" {
			t.Errorf("expected name 'planner', got %s", node.Name())
		}
	})
}

func TestRewriterNode_Execute(t *testing.T) {
	ctx := context.Background()
	mockLLMProvider := &mockLLM{}
	rewriter := agent.NewRewriter(mockLLMProvider, nil)
	node := NewRewriterNode(ctx, rewriter)

	t.Run("no current step", func(t *testing.T) {
		state := &workflow.State{
			OriginalQuestion: "Test question",
		}

		_, err := node.Execute(state)
		if err == nil {
			t.Error("expected error for no current step, got nil")
		}
	})

	t.Run("node name", func(t *testing.T) {
		if node.Name() != "rewriter" {
			t.Errorf("expected name 'rewriter', got %s", node.Name())
		}
	})
}

func TestSupervisorNode_Execute(t *testing.T) {
	ctx := context.Background()
	mockLLMProvider := &mockLLM{}
	supervisor := agent.NewSupervisor(mockLLMProvider, nil)
	node := NewSupervisorNode(ctx, supervisor)

	t.Run("no current step", func(t *testing.T) {
		state := &workflow.State{
			OriginalQuestion: "Test question",
		}

		_, err := node.Execute(state)
		if err == nil {
			t.Error("expected error for no current step, got nil")
		}
	})

	t.Run("node name", func(t *testing.T) {
		if node.Name() != "supervisor" {
			t.Errorf("expected name 'supervisor', got %s", node.Name())
		}
	})
}

func TestRerankerNode_Execute(t *testing.T) {
	ctx := context.Background()
	reranker := agent.NewReranker(nil)
	node := NewRerankerNode(ctx, reranker)

	t.Run("no documents to rerank", func(t *testing.T) {
		state := &workflow.State{
			OriginalQuestion: "Test question",
			RetrievedDocs:    []vectorstore.Document{},
		}

		result, err := node.Execute(state)
		if err != nil {
			t.Fatalf("Execute() failed: %v", err)
		}

		if result == nil {
			t.Fatal("expected result, got nil")
		}

		if len(result.UpdatedState.RerankedDocs) != 0 {
			t.Errorf("expected empty reranked docs, got %d", len(result.UpdatedState.RerankedDocs))
		}
	})

	t.Run("no current step with documents", func(t *testing.T) {
		state := &workflow.State{
			OriginalQuestion: "Test question",
			RetrievedDocs: []vectorstore.Document{
				{Content: "Test document", Metadata: make(map[string]interface{})},
			},
		}

		_, err := node.Execute(state)
		if err == nil {
			t.Error("expected error for no current step, got nil")
		}
	})

	t.Run("node name", func(t *testing.T) {
		if node.Name() != "reranker" {
			t.Errorf("expected name 'reranker', got %s", node.Name())
		}
	})
}

func TestDistillerNode_Execute(t *testing.T) {
	ctx := context.Background()
	mockLLMProvider := &mockLLM{}
	distiller := agent.NewDistiller(mockLLMProvider, nil)
	node := NewDistillerNode(ctx, distiller)

	t.Run("no documents to distill", func(t *testing.T) {
		state := &workflow.State{
			OriginalQuestion: "Test question",
			RerankedDocs:     []vectorstore.Document{},
		}

		result, err := node.Execute(state)
		if err != nil {
			t.Fatalf("Execute() failed: %v", err)
		}

		if result == nil {
			t.Fatal("expected result, got nil")
		}

		if result.UpdatedState.SynthesizedContext != "" {
			t.Errorf("expected empty synthesized context, got %s", result.UpdatedState.SynthesizedContext)
		}
	})

	t.Run("no current step with documents", func(t *testing.T) {
		state := &workflow.State{
			OriginalQuestion: "Test question",
			RerankedDocs: []vectorstore.Document{
				{Content: "Test document", Metadata: make(map[string]interface{})},
			},
		}

		_, err := node.Execute(state)
		if err == nil {
			t.Error("expected error for no current step, got nil")
		}
	})

	t.Run("node name", func(t *testing.T) {
		if node.Name() != "distiller" {
			t.Errorf("expected name 'distiller', got %s", node.Name())
		}
	})
}

func TestReflectorNode_Execute(t *testing.T) {
	ctx := context.Background()
	mockLLMProvider := &mockLLM{}
	reflector := agent.NewReflector(mockLLMProvider, nil)
	node := NewReflectorNode(ctx, reflector)

	t.Run("no current step", func(t *testing.T) {
		state := &workflow.State{
			OriginalQuestion: "Test question",
		}

		_, err := node.Execute(state)
		if err == nil {
			t.Error("expected error for no current step, got nil")
		}
	})

	t.Run("node name", func(t *testing.T) {
		if node.Name() != "reflector" {
			t.Errorf("expected name 'reflector', got %s", node.Name())
		}
	})
}

func TestPolicyNode_Execute(t *testing.T) {
	ctx := context.Background()
	mockLLMProvider := &mockLLM{}
	policy := agent.NewPolicy(mockLLMProvider, nil)
	node := NewPolicyNode(ctx, policy)

	t.Run("continue decision", func(t *testing.T) {
		state := &workflow.State{
			OriginalQuestion: "Test question",
		}
		// Initialize the plan to have some steps
		state.Plan = &workflow.Plan{
			Steps: []workflow.PlanStep{
				{SubQuestion: "Step 1"},
				{SubQuestion: "Step 2"},
			},
		}

		result, err := node.Execute(state)
		if err != nil {
			t.Fatalf("Execute() failed: %v", err)
		}

		if result == nil {
			t.Fatal("expected result, got nil")
		}

		// Policy will decide based on state - test that result is returned
		if result.UpdatedState == nil {
			t.Error("expected updated state, got nil")
		}

		// Verify next node is set
		if result.NextNode == "" {
			t.Error("expected next node to be set")
		}
	})

	t.Run("node name", func(t *testing.T) {
		if node.Name() != "policy" {
			t.Errorf("expected name 'policy', got %s", node.Name())
		}
	})
}
