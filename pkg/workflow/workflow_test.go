// Copyright 2025 Gerry Miller <gerry@gerrymiller.com>
//
// Licensed under the MIT License.
// See LICENSE file in the project root for full license information.

package workflow_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"deep-thinking-agent/pkg/workflow"
)

// ============================================================================
// Mock Nodes for Testing
// ============================================================================

type mockNode struct {
	name        string
	executeFunc func(state *workflow.State) (*workflow.NodeResult, error)
}

func (m *mockNode) Name() string {
	return m.name
}

func (m *mockNode) Execute(state *workflow.State) (*workflow.NodeResult, error) {
	if m.executeFunc != nil {
		return m.executeFunc(state)
	}
	return &workflow.NodeResult{UpdatedState: state}, nil
}

// ============================================================================
// Graph Tests
// ============================================================================

func TestNewGraph(t *testing.T) {
	graph := workflow.NewGraph()
	if graph == nil {
		t.Fatal("NewGraph returned nil")
	}
	// Test behavior rather than internal state
	if graph.GetStartNode() != "" {
		t.Error("start node should be empty initially")
	}
}

func TestGraph_AddNode(t *testing.T) {
	tests := []struct {
		name    string
		node    workflow.Node
		wantErr bool
		errMsg  string
	}{
		{
			name:    "success",
			node:    &mockNode{name: "test"},
			wantErr: false,
		},
		{
			name:    "nil node",
			node:    nil,
			wantErr: true,
			errMsg:  "node is nil",
		},
		{
			name:    "empty name",
			node:    &mockNode{name: ""},
			wantErr: true,
			errMsg:  "node name is empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			graph := workflow.NewGraph()
			err := graph.AddNode(tt.node)
			if (err != nil) != tt.wantErr {
				t.Errorf("AddNode() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil && tt.errMsg != "" && err.Error() != tt.errMsg {
				t.Errorf("AddNode() error = %v, want %v", err.Error(), tt.errMsg)
			}
		})
	}

	// Test duplicate node
	t.Run("duplicate node", func(t *testing.T) {
		graph := workflow.NewGraph()
		node := &mockNode{name: "test"}
		if err := graph.AddNode(node); err != nil {
			t.Fatalf("first AddNode failed: %v", err)
		}
		err := graph.AddNode(node)
		if err == nil {
			t.Error("AddNode should error on duplicate")
		}
		if err != nil && err.Error() != "node test already exists" {
			t.Errorf("unexpected error: %v", err)
		}
	})
}

func TestGraph_AddEdge(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(*workflow.Graph)
		from    string
		to      string
		wantErr bool
		errMsg  string
	}{
		{
			name: "success",
			setup: func(g *workflow.Graph) {
				g.AddNode(&mockNode{name: "node1"})
				g.AddNode(&mockNode{name: "node2"})
			},
			from:    "node1",
			to:      "node2",
			wantErr: false,
		},
		{
			name:    "nonexistent from node",
			setup:   func(g *workflow.Graph) {},
			from:    "node1",
			to:      "node2",
			wantErr: true,
			errMsg:  "from node node1 does not exist",
		},
		{
			name: "nonexistent to node",
			setup: func(g *workflow.Graph) {
				g.AddNode(&mockNode{name: "node1"})
			},
			from:    "node1",
			to:      "node2",
			wantErr: true,
			errMsg:  "to node node2 does not exist",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			graph := workflow.NewGraph()
			tt.setup(graph)
			err := graph.AddEdge(tt.from, tt.to)
			if (err != nil) != tt.wantErr {
				t.Errorf("AddEdge() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil && tt.errMsg != "" && err.Error() != tt.errMsg {
				t.Errorf("AddEdge() error = %v, want %v", err.Error(), tt.errMsg)
			}
		})
	}
}

func TestGraph_SetStart(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		graph := workflow.NewGraph()
		graph.AddNode(&mockNode{name: "start"})
		err := graph.SetStart("start")
		if err != nil {
			t.Errorf("SetStart() error = %v", err)
		}
		if graph.GetStartNode() != "start" {
			t.Errorf("GetStartNode() = %v, want start", graph.GetStartNode())
		}
	})

	t.Run("nonexistent node", func(t *testing.T) {
		graph := workflow.NewGraph()
		err := graph.SetStart("nonexistent")
		if err == nil {
			t.Error("SetStart should error on nonexistent node")
		}
		if err != nil && err.Error() != "start node nonexistent does not exist" {
			t.Errorf("unexpected error: %v", err)
		}
	})
}

func TestGraph_GetNode(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		graph := workflow.NewGraph()
		node := &mockNode{name: "test"}
		graph.AddNode(node)
		retrieved, err := graph.GetNode("test")
		if err != nil {
			t.Errorf("GetNode() error = %v", err)
		}
		if retrieved.Name() != "test" {
			t.Errorf("GetNode() name = %v, want test", retrieved.Name())
		}
	})

	t.Run("not found", func(t *testing.T) {
		graph := workflow.NewGraph()
		_, err := graph.GetNode("nonexistent")
		if err == nil {
			t.Error("GetNode should error on nonexistent node")
		}
		if err != nil && err.Error() != "node nonexistent not found" {
			t.Errorf("unexpected error: %v", err)
		}
	})
}

func TestGraph_GetNextNodes(t *testing.T) {
	graph := workflow.NewGraph()
	graph.AddNode(&mockNode{name: "node1"})
	graph.AddNode(&mockNode{name: "node2"})
	graph.AddNode(&mockNode{name: "node3"})
	graph.AddEdge("node1", "node2")
	graph.AddEdge("node1", "node3")

	next := graph.GetNextNodes("node1")
	if len(next) != 2 {
		t.Errorf("GetNextNodes() len = %v, want 2", len(next))
	}

	// Test node with no edges
	next = graph.GetNextNodes("node2")
	if len(next) != 0 {
		t.Errorf("GetNextNodes() len = %v, want 0", len(next))
	}
}

func TestGraph_GetStartNode(t *testing.T) {
	graph := workflow.NewGraph()
	graph.AddNode(&mockNode{name: "start"})
	graph.SetStart("start")

	start := graph.GetStartNode()
	if start != "start" {
		t.Errorf("GetStartNode() = %v, want start", start)
	}
}

// ============================================================================
// Executor Tests
// ============================================================================

func TestNewExecutor(t *testing.T) {
	graph := workflow.NewGraph()

	t.Run("with config", func(t *testing.T) {
		config := &workflow.ExecutorConfig{Timeout: 30 * time.Second}
		executor := workflow.NewExecutor(graph, config)
		if executor == nil {
			t.Fatal("NewExecutor returned nil")
		}
	})

	t.Run("without config", func(t *testing.T) {
		executor := workflow.NewExecutor(graph, nil)
		if executor == nil {
			t.Fatal("NewExecutor returned nil")
		}
	})
}

func TestExecutor_Execute(t *testing.T) {
	ctx := context.Background()

	t.Run("nil graph", func(t *testing.T) {
		executor := &workflow.Executor{}
		state := workflow.NewState("test question")
		_, err := executor.Execute(ctx, state)
		if err == nil {
			t.Error("Execute should error on nil graph")
		}
	})

	t.Run("nil state", func(t *testing.T) {
		graph := workflow.NewGraph()
		executor := workflow.NewExecutor(graph, nil)
		_, err := executor.Execute(ctx, nil)
		if err == nil {
			t.Error("Execute should error on nil state")
		}
	})

	t.Run("no start node", func(t *testing.T) {
		graph := workflow.NewGraph()
		executor := workflow.NewExecutor(graph, nil)
		state := workflow.NewState("test question")
		_, err := executor.Execute(ctx, state)
		if err == nil {
			t.Error("Execute should error on no start node")
		}
	})

	t.Run("single node success", func(t *testing.T) {
		graph := workflow.NewGraph()
		executed := false
		node := &mockNode{
			name: "test",
			executeFunc: func(state *workflow.State) (*workflow.NodeResult, error) {
				executed = true
				state.FinalAnswer = "done"
				return &workflow.NodeResult{UpdatedState: state}, nil
			},
		}
		graph.AddNode(node)
		graph.SetStart("test")

		executor := workflow.NewExecutor(graph, nil)
		state := workflow.NewState("test question")
		result, err := executor.Execute(ctx, state)
		if err != nil {
			t.Errorf("Execute() error = %v", err)
		}
		if !executed {
			t.Error("node was not executed")
		}
		if result.FinalAnswer != "done" {
			t.Errorf("FinalAnswer = %v, want done", result.FinalAnswer)
		}
	})

	t.Run("node error", func(t *testing.T) {
		graph := workflow.NewGraph()
		node := &mockNode{
			name: "test",
			executeFunc: func(state *workflow.State) (*workflow.NodeResult, error) {
				return nil, errors.New("node error")
			},
		}
		graph.AddNode(node)
		graph.SetStart("test")

		executor := workflow.NewExecutor(graph, nil)
		state := workflow.NewState("test question")
		_, err := executor.Execute(ctx, state)
		if err == nil {
			t.Error("Execute should propagate node error")
		}
	})

	t.Run("max iterations", func(t *testing.T) {
		graph := workflow.NewGraph()
		callCount := 0
		node := &mockNode{
			name: "loop",
			executeFunc: func(state *workflow.State) (*workflow.NodeResult, error) {
				callCount++
				// Always loop back to self
				return &workflow.NodeResult{
					UpdatedState: state,
					NextNode:     "loop",
				}, nil
			},
		}
		graph.AddNode(node)
		graph.SetStart("loop")

		executor := workflow.NewExecutor(graph, &workflow.ExecutorConfig{Timeout: 10 * time.Second})
		state := workflow.NewState("test question")
		_, err := executor.Execute(ctx, state)
		if err == nil {
			t.Error("Execute should error on max iterations")
		}
		if callCount != 100 { // Exactly 100 iterations before hitting limit
			t.Errorf("callCount = %v, want 100", callCount)
		}
	})

	t.Run("timeout", func(t *testing.T) {
		graph := workflow.NewGraph()
		callCount := 0
		node := &mockNode{
			name: "loop",
			executeFunc: func(state *workflow.State) (*workflow.NodeResult, error) {
				callCount++
				time.Sleep(20 * time.Millisecond)
				// Loop back to self
				return &workflow.NodeResult{
					UpdatedState: state,
					NextNode:     "loop",
				}, nil
			},
		}
		graph.AddNode(node)
		graph.SetStart("loop")

		start := time.Now()
		executor := workflow.NewExecutor(graph, &workflow.ExecutorConfig{Timeout: 50 * time.Millisecond})
		state := workflow.NewState("test question")
		_, err := executor.Execute(ctx, state)
		elapsed := time.Since(start)

		if err == nil {
			t.Error("Execute should error on timeout")
		}
		// Should timeout after ~50ms, having completed 2-3 iterations
		if elapsed > 100*time.Millisecond {
			t.Errorf("execution took too long: %v", elapsed)
		}
		if callCount < 2 {
			t.Errorf("callCount = %v, should have completed at least 2 iterations", callCount)
		}
	})

	t.Run("multi-node workflow", func(t *testing.T) {
		graph := workflow.NewGraph()

		node1 := &mockNode{
			name: "node1",
			executeFunc: func(state *workflow.State) (*workflow.NodeResult, error) {
				state.FinalAnswer = "step1"
				return &workflow.NodeResult{UpdatedState: state}, nil
			},
		}
		node2 := &mockNode{
			name: "node2",
			executeFunc: func(state *workflow.State) (*workflow.NodeResult, error) {
				state.FinalAnswer += " step2"
				return &workflow.NodeResult{UpdatedState: state}, nil
			},
		}

		graph.AddNode(node1)
		graph.AddNode(node2)
		graph.AddEdge("node1", "node2")
		graph.SetStart("node1")

		executor := workflow.NewExecutor(graph, nil)
		state := workflow.NewState("test question")
		result, err := executor.Execute(ctx, state)
		if err != nil {
			t.Errorf("Execute() error = %v", err)
		}
		if result.FinalAnswer != "step1 step2" {
			t.Errorf("FinalAnswer = %v, want 'step1 step2'", result.FinalAnswer)
		}
	})
}

func TestExecutor_ExecuteStep(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		graph := workflow.NewGraph()
		executed := false
		node := &mockNode{
			name: "test",
			executeFunc: func(state *workflow.State) (*workflow.NodeResult, error) {
				executed = true
				state.FinalAnswer = "done"
				return &workflow.NodeResult{UpdatedState: state}, nil
			},
		}
		graph.AddNode(node)

		executor := workflow.NewExecutor(graph, nil)
		state := workflow.NewState("test question")
		result, err := executor.ExecuteStep(ctx, state, "test")
		if err != nil {
			t.Errorf("ExecuteStep() error = %v", err)
		}
		if !executed {
			t.Error("node was not executed")
		}
		if result.FinalAnswer != "done" {
			t.Errorf("FinalAnswer = %v, want done", result.FinalAnswer)
		}
	})

	t.Run("nonexistent node", func(t *testing.T) {
		graph := workflow.NewGraph()
		executor := workflow.NewExecutor(graph, nil)
		state := workflow.NewState("test question")
		_, err := executor.ExecuteStep(ctx, state, "nonexistent")
		if err == nil {
			t.Error("ExecuteStep should error on nonexistent node")
		}
	})
}

// ============================================================================
// State Tests
// ============================================================================

func TestNewState(t *testing.T) {
	state := workflow.NewState("test question")
	if state == nil {
		t.Fatal("NewState returned nil")
	}
	if state.OriginalQuestion != "test question" {
		t.Errorf("OriginalQuestion = %v, want 'test question'", state.OriginalQuestion)
	}
	if state.CurrentStepIndex != 0 {
		t.Error("CurrentStepIndex should be 0")
	}
	if !state.ShouldContinue {
		t.Error("ShouldContinue should be true")
	}
	if state.MaxIterations != 10 {
		t.Errorf("MaxIterations = %v, want 10", state.MaxIterations)
	}
}

func TestState_CurrentStep(t *testing.T) {
	state := workflow.NewState("test")

	// No plan
	if state.CurrentStep() != nil {
		t.Error("CurrentStep should return nil when no plan")
	}

	// Plan with steps
	state.Plan = &workflow.Plan{
		Steps: []workflow.PlanStep{
			{Index: 0, SubQuestion: "step 0"},
			{Index: 1, SubQuestion: "step 1"},
		},
	}
	step := state.CurrentStep()
	if step == nil {
		t.Fatal("CurrentStep returned nil")
	}
	if step.SubQuestion != "step 0" {
		t.Errorf("SubQuestion = %v, want 'step 0'", step.SubQuestion)
	}

	// Move to next step
	state.IncrementStep()
	step = state.CurrentStep()
	if step == nil {
		t.Fatal("CurrentStep returned nil after increment")
	}
	if step.SubQuestion != "step 1" {
		t.Errorf("SubQuestion = %v, want 'step 1'", step.SubQuestion)
	}

	// Beyond plan
	state.IncrementStep()
	if state.CurrentStep() != nil {
		t.Error("CurrentStep should return nil when beyond plan")
	}
}

func TestState_IsComplete(t *testing.T) {
	state := workflow.NewState("test")

	// No plan
	if state.IsComplete() {
		t.Error("IsComplete should be false when no plan")
	}

	// Plan with steps
	state.Plan = &workflow.Plan{
		Steps: []workflow.PlanStep{
			{Index: 0},
			{Index: 1},
		},
	}
	if state.IsComplete() {
		t.Error("IsComplete should be false at start")
	}

	state.IncrementStep()
	if state.IsComplete() {
		t.Error("IsComplete should be false at step 1")
	}

	state.IncrementStep()
	if !state.IsComplete() {
		t.Error("IsComplete should be true after all steps")
	}
}

func TestState_HasReachedMaxIterations(t *testing.T) {
	state := workflow.NewState("test")
	state.MaxIterations = 3

	if state.HasReachedMaxIterations() {
		t.Error("Should not have reached max iterations initially")
	}

	// Add past steps
	for i := 0; i < 3; i++ {
		state.AddPastStep(workflow.PastStep{
			Summary: "test",
		})
	}

	if !state.HasReachedMaxIterations() {
		t.Error("Should have reached max iterations")
	}
}

func TestState_GetRetrievalContext(t *testing.T) {
	state := workflow.NewState("test")

	// No plan
	if state.GetRetrievalContext() != nil {
		t.Error("GetRetrievalContext should return nil when no plan")
	}

	// Plan with steps
	state.Plan = &workflow.Plan{
		Steps: []workflow.PlanStep{
			{SubQuestion: "test question"},
		},
	}

	ctx := state.GetRetrievalContext()
	if ctx == nil {
		t.Fatal("GetRetrievalContext returned nil")
	}
	if ctx.Query != "test question" {
		t.Errorf("Query = %v, want 'test question'", ctx.Query)
	}
	if ctx.Strategy != workflow.StrategyHybrid {
		t.Errorf("Strategy = %v, want hybrid", ctx.Strategy)
	}
	if ctx.TopK != 10 {
		t.Errorf("TopK = %v, want 10", ctx.TopK)
	}
	if ctx.RerankerTopN != 3 {
		t.Errorf("RerankerTopN = %v, want 3", ctx.RerankerTopN)
	}
}
