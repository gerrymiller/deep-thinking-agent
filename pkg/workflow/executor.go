// Copyright 2025 Gerry Miller <gerry@gerrymiller.com>
//
// Licensed under the MIT License.
// See LICENSE file in the project root for full license information.

package workflow

import (
	"context"
	"fmt"
	"time"
)

// Executor runs the workflow graph with state management.
type Executor struct {
	graph   *Graph
	timeout time.Duration
}

// ExecutorConfig contains configuration for the executor.
type ExecutorConfig struct {
	Timeout time.Duration
}

// NewExecutor creates a new workflow executor.
func NewExecutor(graph *Graph, config *ExecutorConfig) *Executor {
	if config == nil {
		config = &ExecutorConfig{
			Timeout: 5 * time.Minute,
		}
	}

	return &Executor{
		graph:   graph,
		timeout: config.Timeout,
	}
}

// Execute runs the workflow graph starting from the initial state.
func (e *Executor) Execute(ctx context.Context, initialState *State) (*State, error) {
	if e.graph == nil {
		return nil, fmt.Errorf("graph is nil")
	}

	if initialState == nil {
		return nil, fmt.Errorf("initial state is nil")
	}

	// Apply timeout
	if e.timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, e.timeout)
		defer cancel()
	}

	// Get starting node
	currentNodeName := e.graph.GetStartNode()
	if currentNodeName == "" {
		return nil, fmt.Errorf("no start node defined")
	}

	state := initialState
	iterationCount := 0

	// Execute nodes in sequence
	for {
		// Check context cancellation
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("execution timeout or cancelled: %w", ctx.Err())
		default:
		}

		// Safety check for infinite loops
		iterationCount++
		if iterationCount > 100 {
			return nil, fmt.Errorf("exceeded maximum iteration count (100)")
		}

		// Get current node
		node, err := e.graph.GetNode(currentNodeName)
		if err != nil {
			return nil, fmt.Errorf("failed to get node %s: %w", currentNodeName, err)
		}

		// Execute node
		result, err := node.Execute(state)
		if err != nil {
			return nil, fmt.Errorf("node %s execution failed: %w", currentNodeName, err)
		}

		if result == nil {
			return nil, fmt.Errorf("node %s returned nil result", currentNodeName)
		}

		// Update state
		state = result.UpdatedState
		if state == nil {
			return nil, fmt.Errorf("node %s returned nil state", currentNodeName)
		}

		// Check for errors in state
		if state.Error != nil {
			return state, fmt.Errorf("workflow error: %w", state.Error)
		}

		// Determine next node
		if result.NextNode != "" {
			// Explicit next node specified
			currentNodeName = result.NextNode
		} else {
			// Use default routing from graph
			nextNodes := e.graph.GetNextNodes(currentNodeName)

			if len(nextNodes) == 0 {
				// No more nodes, workflow complete
				break
			}

			if len(nextNodes) == 1 {
				// Single next node
				currentNodeName = nextNodes[0]
			} else {
				// Multiple possible next nodes - use routing logic
				currentNodeName = e.routeNext(state, nextNodes)
			}
		}

		// Check if policy says to finish
		if currentNodeName == "finish" || !state.ShouldContinue {
			break
		}

		// Check if plan is complete
		if currentNodeName == "rewriter" && state.IsComplete() {
			// All steps done, exit loop
			break
		}

		// Check max iterations safety
		if state.HasReachedMaxIterations() {
			break
		}
	}

	return state, nil
}

// routeNext determines the next node based on state and available options.
func (e *Executor) routeNext(state *State, options []string) string {
	// For policy node, check ShouldContinue
	if !state.ShouldContinue {
		return "finish"
	}

	// Default: return first option
	if len(options) > 0 {
		return options[0]
	}

	return "finish"
}

// ExecuteStep runs a single step of the workflow (for debugging/testing).
func (e *Executor) ExecuteStep(ctx context.Context, state *State, nodeName string) (*State, error) {
	node, err := e.graph.GetNode(nodeName)
	if err != nil {
		return nil, err
	}

	result, err := node.Execute(state)
	if err != nil {
		return nil, err
	}

	return result.UpdatedState, nil
}
