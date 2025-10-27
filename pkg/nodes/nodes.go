// Copyright 2025 Gerry Miller <gerry@gerrymiller.com>
//
// Licensed under the MIT License.
// See LICENSE file in the project root for full license information.

package nodes

import (
	"context"
	"fmt"

	"deep-thinking-agent/pkg/agent"
	"deep-thinking-agent/pkg/vectorstore"
	"deep-thinking-agent/pkg/workflow"
)

// PlannerNode wraps the planner agent as a workflow node.
type PlannerNode struct {
	planner *agent.Planner
	ctx     context.Context
}

// NewPlannerNode creates a new planner node.
func NewPlannerNode(ctx context.Context, planner *agent.Planner) *PlannerNode {
	return &PlannerNode{
		planner: planner,
		ctx:     ctx,
	}
}

// Execute runs the planner to create a query execution plan.
func (n *PlannerNode) Execute(state *workflow.State) (*workflow.NodeResult, error) {
	plan, err := n.planner.Plan(n.ctx, state.OriginalQuestion)
	if err != nil {
		return nil, fmt.Errorf("planning failed: %w", err)
	}

	state.Plan = plan
	return &workflow.NodeResult{UpdatedState: state}, nil
}

// Name returns the node name.
func (n *PlannerNode) Name() string {
	return "planner"
}

// RewriterNode wraps the rewriter agent as a workflow node.
type RewriterNode struct {
	rewriter *agent.Rewriter
	ctx      context.Context
}

// NewRewriterNode creates a new rewriter node.
func NewRewriterNode(ctx context.Context, rewriter *agent.Rewriter) *RewriterNode {
	return &RewriterNode{
		rewriter: rewriter,
		ctx:      ctx,
	}
}

// Execute enhances the current query for better retrieval.
func (n *RewriterNode) Execute(state *workflow.State) (*workflow.NodeResult, error) {
	currentStep := state.CurrentStep()
	if currentStep == nil {
		return nil, fmt.Errorf("no current step available")
	}

	rewritten, err := n.rewriter.Rewrite(n.ctx, currentStep.SubQuestion, state)
	if err != nil {
		return nil, fmt.Errorf("rewriting failed: %w", err)
	}

	// Update the query in retrieval context
	if state.GetRetrievalContext() != nil {
		state.GetRetrievalContext().Query = rewritten
	}

	return &workflow.NodeResult{UpdatedState: state}, nil
}

// Name returns the node name.
func (n *RewriterNode) Name() string {
	return "rewriter"
}

// SupervisorNode wraps the supervisor agent as a workflow node.
type SupervisorNode struct {
	supervisor *agent.Supervisor
	ctx        context.Context
}

// NewSupervisorNode creates a new supervisor node.
func NewSupervisorNode(ctx context.Context, supervisor *agent.Supervisor) *SupervisorNode {
	return &SupervisorNode{
		supervisor: supervisor,
		ctx:        ctx,
	}
}

// Execute selects the optimal retrieval strategy.
func (n *SupervisorNode) Execute(state *workflow.State) (*workflow.NodeResult, error) {
	currentStep := state.CurrentStep()
	if currentStep == nil {
		return nil, fmt.Errorf("no current step available")
	}

	strategy, err := n.supervisor.SelectStrategy(n.ctx, currentStep.SubQuestion, state)
	if err != nil {
		return nil, fmt.Errorf("strategy selection failed: %w", err)
	}

	// Update retrieval context with selected strategy
	retrievalCtx := state.GetRetrievalContext()
	if retrievalCtx != nil {
		retrievalCtx.Strategy = strategy
	}

	return &workflow.NodeResult{UpdatedState: state}, nil
}

// Name returns the node name.
func (n *SupervisorNode) Name() string {
	return "supervisor"
}

// RetrieverNode wraps the retriever agent as a workflow node.
type RetrieverNode struct {
	retriever *agent.Retriever
	ctx       context.Context
}

// NewRetrieverNode creates a new retriever node.
func NewRetrieverNode(ctx context.Context, retriever *agent.Retriever) *RetrieverNode {
	return &RetrieverNode{
		retriever: retriever,
		ctx:       ctx,
	}
}

// Execute retrieves relevant documents.
func (n *RetrieverNode) Execute(state *workflow.State) (*workflow.NodeResult, error) {
	retrievalCtx := state.GetRetrievalContext()
	if retrievalCtx == nil {
		return nil, fmt.Errorf("no retrieval context available")
	}

	docs, err := n.retriever.Retrieve(n.ctx, retrievalCtx)
	if err != nil {
		return nil, fmt.Errorf("retrieval failed: %w", err)
	}

	state.RetrievedDocs = docs
	return &workflow.NodeResult{UpdatedState: state}, nil
}

// Name returns the node name.
func (n *RetrieverNode) Name() string {
	return "retriever"
}

// RerankerNode wraps the reranker agent as a workflow node.
type RerankerNode struct {
	reranker *agent.Reranker
	ctx      context.Context
}

// NewRerankerNode creates a new reranker node.
func NewRerankerNode(ctx context.Context, reranker *agent.Reranker) *RerankerNode {
	return &RerankerNode{
		reranker: reranker,
		ctx:      ctx,
	}
}

// Execute reranks retrieved documents for precision.
func (n *RerankerNode) Execute(state *workflow.State) (*workflow.NodeResult, error) {
	if len(state.RetrievedDocs) == 0 {
		// No documents to rerank, continue
		state.RerankedDocs = []vectorstore.Document{}
		return &workflow.NodeResult{UpdatedState: state}, nil
	}

	currentStep := state.CurrentStep()
	if currentStep == nil {
		return nil, fmt.Errorf("no current step available")
	}

	reranked := n.reranker.Rerank(n.ctx, currentStep.SubQuestion, state.RetrievedDocs)
	state.RerankedDocs = reranked

	return &workflow.NodeResult{UpdatedState: state}, nil
}

// Name returns the node name.
func (n *RerankerNode) Name() string {
	return "reranker"
}

// DistillerNode wraps the distiller agent as a workflow node.
type DistillerNode struct {
	distiller *agent.Distiller
	ctx       context.Context
}

// NewDistillerNode creates a new distiller node.
func NewDistillerNode(ctx context.Context, distiller *agent.Distiller) *DistillerNode {
	return &DistillerNode{
		distiller: distiller,
		ctx:       ctx,
	}
}

// Execute synthesizes documents into coherent context.
func (n *DistillerNode) Execute(state *workflow.State) (*workflow.NodeResult, error) {
	if len(state.RerankedDocs) == 0 {
		// No documents to distill
		state.SynthesizedContext = ""
		return &workflow.NodeResult{UpdatedState: state}, nil
	}

	currentStep := state.CurrentStep()
	if currentStep == nil {
		return nil, fmt.Errorf("no current step available")
	}

	synthesized, err := n.distiller.Distill(n.ctx, currentStep.SubQuestion, state.RerankedDocs)
	if err != nil {
		return nil, fmt.Errorf("distillation failed: %w", err)
	}

	state.SynthesizedContext = synthesized
	return &workflow.NodeResult{UpdatedState: state}, nil
}

// Name returns the node name.
func (n *DistillerNode) Name() string {
	return "distiller"
}

// ReflectorNode wraps the reflector agent as a workflow node.
type ReflectorNode struct {
	reflector *agent.Reflector
	ctx       context.Context
}

// NewReflectorNode creates a new reflector node.
func NewReflectorNode(ctx context.Context, reflector *agent.Reflector) *ReflectorNode {
	return &ReflectorNode{
		reflector: reflector,
		ctx:       ctx,
	}
}

// Execute reflects on the completed step and extracts key findings.
func (n *ReflectorNode) Execute(state *workflow.State) (*workflow.NodeResult, error) {
	currentStep := state.CurrentStep()
	if currentStep == nil {
		return nil, fmt.Errorf("no current step available")
	}

	summary, keyFindings, err := n.reflector.Reflect(n.ctx, currentStep, state.SynthesizedContext)
	if err != nil {
		return nil, fmt.Errorf("reflection failed: %w", err)
	}

	// Create past step record
	pastStep := workflow.PastStep{
		Step:          *currentStep,
		RetrievedDocs: state.RerankedDocs,
		Summary:       summary,
		KeyFindings:   keyFindings,
	}

	state.AddPastStep(pastStep)
	state.IncrementStep()

	return &workflow.NodeResult{UpdatedState: state}, nil
}

// Name returns the node name.
func (n *ReflectorNode) Name() string {
	return "reflector"
}

// PolicyNode wraps the policy agent as a workflow node.
type PolicyNode struct {
	policy *agent.Policy
	ctx    context.Context
}

// NewPolicyNode creates a new policy node.
func NewPolicyNode(ctx context.Context, policy *agent.Policy) *PolicyNode {
	return &PolicyNode{
		policy: policy,
		ctx:    ctx,
	}
}

// Execute decides whether to continue or finish the workflow.
func (n *PolicyNode) Execute(state *workflow.State) (*workflow.NodeResult, error) {
	decision, err := n.policy.Decide(n.ctx, state)
	if err != nil {
		return nil, fmt.Errorf("policy decision failed: %w", err)
	}

	state.ShouldContinue = decision.ShouldContinue

	// Determine next node
	nextNode := ""
	if decision.ShouldContinue {
		nextNode = "rewriter" // Continue to next iteration
	} else {
		nextNode = "finish" // End workflow
	}

	return &workflow.NodeResult{
		UpdatedState: state,
		NextNode:     nextNode,
	}, nil
}

// Name returns the node name.
func (n *PolicyNode) Name() string {
	return "policy"
}
