// Copyright 2025 Gerry Miller <gerry@gerrymiller.com>
//
// Licensed under the MIT License.
// See LICENSE file in the project root for full license information.

package workflow

import "fmt"

// Graph represents the workflow execution graph.
// It defines nodes and their connections for the deep thinking loop.
type Graph struct {
	nodes map[string]Node
	edges map[string][]string // node name -> list of possible next nodes
	start string              // starting node name
}

// Node represents a single node in the workflow graph.
type Node interface {
	// Execute runs this node with the given state and returns updated state
	Execute(state *State) (*NodeResult, error)

	// Name returns the node's unique identifier
	Name() string
}

// NewGraph creates a new workflow graph.
func NewGraph() *Graph {
	return &Graph{
		nodes: make(map[string]Node),
		edges: make(map[string][]string),
	}
}

// AddNode adds a node to the graph.
func (g *Graph) AddNode(node Node) error {
	if node == nil {
		return fmt.Errorf("node is nil")
	}

	name := node.Name()
	if name == "" {
		return fmt.Errorf("node name is empty")
	}

	if _, exists := g.nodes[name]; exists {
		return fmt.Errorf("node %s already exists", name)
	}

	g.nodes[name] = node
	return nil
}

// AddEdge adds a directed edge from one node to another.
func (g *Graph) AddEdge(from, to string) error {
	if _, exists := g.nodes[from]; !exists {
		return fmt.Errorf("from node %s does not exist", from)
	}
	if _, exists := g.nodes[to]; !exists {
		return fmt.Errorf("to node %s does not exist", to)
	}

	g.edges[from] = append(g.edges[from], to)
	return nil
}

// SetStart sets the starting node for execution.
func (g *Graph) SetStart(nodeName string) error {
	if _, exists := g.nodes[nodeName]; !exists {
		return fmt.Errorf("start node %s does not exist", nodeName)
	}

	g.start = nodeName
	return nil
}

// GetNode retrieves a node by name.
func (g *Graph) GetNode(name string) (Node, error) {
	node, exists := g.nodes[name]
	if !exists {
		return nil, fmt.Errorf("node %s not found", name)
	}
	return node, nil
}

// GetNextNodes returns the possible next nodes from a given node.
func (g *Graph) GetNextNodes(nodeName string) []string {
	return g.edges[nodeName]
}

// GetStartNode returns the starting node name.
func (g *Graph) GetStartNode() string {
	return g.start
}

// BuildDeepThinkingGraph constructs the standard deep thinking workflow graph.
// Flow: Plan → Rewrite → Supervise → Retrieve → Rerank → Distill → Reflect → Policy
// Policy decides: continue (loop back) or finish
func BuildDeepThinkingGraph(nodes map[string]Node) (*Graph, error) {
	graph := NewGraph()

	// Expected node names
	expectedNodes := []string{
		"planner",
		"rewriter",
		"supervisor",
		"retriever",
		"reranker",
		"distiller",
		"reflector",
		"policy",
	}

	// Add all nodes
	for _, name := range expectedNodes {
		node, exists := nodes[name]
		if !exists {
			return nil, fmt.Errorf("required node %s not provided", name)
		}
		if err := graph.AddNode(node); err != nil {
			return nil, fmt.Errorf("failed to add node %s: %w", name, err)
		}
	}

	// Build the workflow pipeline
	// Planner runs once at start
	if err := graph.AddEdge("planner", "rewriter"); err != nil {
		return nil, err
	}

	// Main loop: Rewrite → Supervise → Retrieve → Rerank → Distill → Reflect → Policy
	if err := graph.AddEdge("rewriter", "supervisor"); err != nil {
		return nil, err
	}
	if err := graph.AddEdge("supervisor", "retriever"); err != nil {
		return nil, err
	}
	if err := graph.AddEdge("retriever", "reranker"); err != nil {
		return nil, err
	}
	if err := graph.AddEdge("reranker", "distiller"); err != nil {
		return nil, err
	}
	if err := graph.AddEdge("distiller", "reflector"); err != nil {
		return nil, err
	}
	if err := graph.AddEdge("reflector", "policy"); err != nil {
		return nil, err
	}

	// Policy decides: continue back to rewriter, or finish
	if err := graph.AddEdge("policy", "rewriter"); err != nil {
		return nil, err
	}
	// Policy can also go to "finish" (handled by executor)

	// Set start node
	if err := graph.SetStart("planner"); err != nil {
		return nil, err
	}

	return graph, nil
}
