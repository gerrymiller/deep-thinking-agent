// Copyright 2025 Gerry Miller <gerry@gerrymiller.com>
//
// Licensed under the MIT License.
// See LICENSE file in the project root for full license information.

package workflow

import (
	"deep-thinking-agent/pkg/schema"
	"deep-thinking-agent/pkg/vectorstore"
)

// State represents the complete state of the deep thinking RAG workflow.
// This state is passed between nodes in the execution graph and accumulates
// information as the workflow progresses through its steps.
type State struct {
	// Original query context
	OriginalQuestion string
	Plan             *Plan

	// Execution tracking
	CurrentStepIndex int
	PastSteps        []PastStep
	MaxIterations    int // Safety limit to prevent infinite loops

	// Retrieval results (current step)
	RetrievedDocs []vectorstore.Document
	RerankedDocs  []vectorstore.Document

	// Synthesis
	SynthesizedContext string
	FinalAnswer        string

	// Schema context (enables schema-aware retrieval)
	RelevantSchemas map[string]*schema.DocumentSchema // DocID -> Schema
	ActiveFilters   *SchemaFilters

	// Workflow control
	ShouldContinue bool  // Policy agent sets this
	Error          error // Any error encountered during workflow
}

// Plan represents the decomposed query execution plan.
// The planner agent creates this by breaking down the original question
// into sequential steps that can be executed independently.
type Plan struct {
	// Steps to execute in order
	Steps []PlanStep

	// Reasoning explains why this plan was chosen
	Reasoning string
}

// PlanStep represents a single step in the execution plan.
type PlanStep struct {
	// Index in the plan sequence
	Index int

	// SubQuestion is the specific question this step answers
	SubQuestion string

	// ToolType indicates which retrieval tool to use
	// Values: "doc_search", "web_search", "schema_filter"
	ToolType string

	// SchemaHint provides guidance on which document sections to target
	// Example: "focus on risk sections", "financial data only"
	SchemaHint string

	// ExpectedOutputs describes what information this step should find
	ExpectedOutputs []string

	// Dependencies lists step indices that must complete before this step
	Dependencies []int
}

// PastStep records the execution and results of a completed plan step.
// This history enables reflection and informs future steps.
type PastStep struct {
	// Step that was executed
	Step PlanStep

	// Documents retrieved during this step
	RetrievedDocs []vectorstore.Document

	// Summary of findings from this step
	Summary string

	// KeyFindings are the most important pieces of information extracted
	KeyFindings []string

	// ExecutionTime tracks how long this step took (for performance monitoring)
	ExecutionTimeMs int64
}

// SchemaFilters defines constraints for schema-aware retrieval.
// These filters are derived from the current plan step and relevant schemas.
type SchemaFilters struct {
	// DocumentIDs limits search to specific documents
	DocumentIDs []string

	// SectionTypes filters by semantic section types
	// Example: ["risk_factors", "financial_data", "methodology"]
	SectionTypes []string

	// HierarchyPaths filters by document hierarchy
	// Example: ["1.2", "3.1.4"] for specific heading paths
	HierarchyPaths []string

	// SemanticTags filters by LLM-identified semantic tags
	SemanticTags []string

	// CustomAttributes allows filtering by document-specific metadata
	CustomAttributes map[string]interface{}

	// MinRelevanceScore filters results below this threshold
	MinRelevanceScore float32
}

// NodeResult represents the output of executing a workflow node.
// Nodes return this to indicate success/failure and provide updated state.
type NodeResult struct {
	// UpdatedState is the modified state after node execution
	UpdatedState *State

	// NextNode specifies which node to execute next
	// If empty, workflow continues to default next node
	NextNode string

	// Error indicates if the node encountered an error
	Error error
}

// RetrievalStrategy indicates which retrieval approach to use.
type RetrievalStrategy string

const (
	// StrategyVector uses semantic vector similarity search
	StrategyVector RetrievalStrategy = "vector"

	// StrategyKeyword uses BM25/keyword-based search
	StrategyKeyword RetrievalStrategy = "keyword"

	// StrategyHybrid combines vector and keyword with RRF
	StrategyHybrid RetrievalStrategy = "hybrid"

	// StrategySchemaFiltered uses schema metadata for targeted retrieval
	StrategySchemaFiltered RetrievalStrategy = "schema_filtered"
)

// RetrievalContext provides context for retrieval operations.
// This is used by retrieval nodes to understand what to search for and how.
type RetrievalContext struct {
	// Query is the search query (may be rewritten from original)
	Query string

	// Strategy indicates which retrieval approach to use
	Strategy RetrievalStrategy

	// TopK is the number of results to retrieve
	TopK int

	// SchemaFilters constrains the search using schema metadata
	SchemaFilters *SchemaFilters

	// RerankerTopN is the number of results to keep after reranking
	RerankerTopN int

	// IncludeHistory indicates if past step findings should influence retrieval
	IncludeHistory bool
}

// PolicyDecision represents the policy agent's decision on whether to continue.
type PolicyDecision struct {
	// ShouldContinue indicates if the workflow should continue to the next step
	ShouldContinue bool

	// Reasoning explains why this decision was made
	Reasoning string

	// Confidence is the policy agent's confidence in this decision (0.0-1.0)
	Confidence float32

	// SuggestedAction provides guidance if continuing
	// Example: "focus on external sources", "need more specific data"
	SuggestedAction string
}

// NewState creates a new workflow state initialized with defaults.
func NewState(question string) *State {
	return &State{
		OriginalQuestion: question,
		CurrentStepIndex: 0,
		PastSteps:        make([]PastStep, 0),
		MaxIterations:    10, // Default safety limit
		RelevantSchemas:  make(map[string]*schema.DocumentSchema),
		ShouldContinue:   true,
	}
}

// AddPastStep appends a completed step to the history.
func (s *State) AddPastStep(step PastStep) {
	s.PastSteps = append(s.PastSteps, step)
}

// IncrementStep moves to the next step in the plan.
func (s *State) IncrementStep() {
	s.CurrentStepIndex++
}

// CurrentStep returns the current plan step if available.
func (s *State) CurrentStep() *PlanStep {
	if s.Plan == nil || s.CurrentStepIndex >= len(s.Plan.Steps) {
		return nil
	}
	return &s.Plan.Steps[s.CurrentStepIndex]
}

// IsComplete returns true if all plan steps have been executed.
func (s *State) IsComplete() bool {
	if s.Plan == nil {
		return false
	}
	return s.CurrentStepIndex >= len(s.Plan.Steps)
}

// HasReachedMaxIterations returns true if the safety limit has been hit.
func (s *State) HasReachedMaxIterations() bool {
	return len(s.PastSteps) >= s.MaxIterations
}

// GetRetrievalContext builds retrieval context from current state.
func (s *State) GetRetrievalContext() *RetrievalContext {
	currentStep := s.CurrentStep()
	if currentStep == nil {
		return nil
	}

	return &RetrievalContext{
		Query:          currentStep.SubQuestion,
		Strategy:       StrategyHybrid, // Default to hybrid
		TopK:           10,
		RerankerTopN:   3,
		SchemaFilters:  s.ActiveFilters,
		IncludeHistory: len(s.PastSteps) > 0,
	}
}
