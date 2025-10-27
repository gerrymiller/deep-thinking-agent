// Copyright 2025 Gerry Miller <gerry@gerrymiller.com>
//
// Licensed under the MIT License.
// See LICENSE file in the project root for full license information.

package schema

// DocumentSchema represents the LLM-derived structural analysis of a document.
// This schema is used to enable intelligent chunking and schema-aware retrieval.
type DocumentSchema struct {
	// DocID uniquely identifies this document
	DocID string

	// Format indicates the document type (PDF, HTML, MD, TXT)
	Format string

	// Title of the document (if available)
	Title string

	// Sections are the logical divisions identified by the LLM
	Sections []Section

	// Hierarchy represents the document's structural hierarchy
	Hierarchy *HierarchyTree

	// SemanticRegions are topic-based regions identified by the LLM
	SemanticRegions []SemanticRegion

	// CustomAttributes are document-specific metadata fields suggested by the LLM
	CustomAttributes map[string]interface{}

	// ChunkingStrategy indicates how this document should be chunked
	// Values: "section_based", "sliding_window", "semantic", "hierarchical"
	ChunkingStrategy string

	// ChunkMetadata contains strategy-specific parameters
	ChunkMetadata map[string]interface{}

	// ParsingMethod describes how this document was parsed
	// Example: "regex_patterns", "llm_analysis", "predefined_schema"
	ParsingMethod string

	// Confidence is the LLM's confidence in this schema (0.0-1.0)
	Confidence float32

	// CreatedAt timestamp
	CreatedAt int64
}

// Section represents a logical section of a document.
type Section struct {
	// ID uniquely identifies this section within the document
	ID string

	// Title is the section heading/title
	Title string

	// Level indicates the hierarchy level (1 = top level, 2 = subsection, etc.)
	Level int

	// StartPos is the character position where this section begins
	StartPos int

	// EndPos is the character position where this section ends
	EndPos int

	// Type is the semantic type of this section
	// Example: "risk_factors", "financial_data", "methodology", "introduction"
	Type string

	// Summary is a brief description of this section's content
	Summary string

	// Keywords are important terms associated with this section
	Keywords []string

	// ParentID references the parent section (empty for top-level sections)
	ParentID string

	// ChildIDs reference child sections
	ChildIDs []string
}

// HierarchyTree represents the document's hierarchical structure.
type HierarchyTree struct {
	// Root is the top-level node
	Root *HierarchyNode

	// MaxDepth is the maximum nesting level
	MaxDepth int
}

// HierarchyNode represents a single node in the document hierarchy.
type HierarchyNode struct {
	// ID uniquely identifies this node
	ID string

	// Path is the hierarchical path (e.g., "1.2.3")
	Path string

	// Title of this node
	Title string

	// Level in the hierarchy (1 = top level)
	Level int

	// Content range
	StartPos int
	EndPos   int

	// Children nodes
	Children []*HierarchyNode
}

// SemanticRegion represents a topic-based region identified by the LLM.
// Unlike sections which follow document structure, semantic regions are
// based on content meaning and may span multiple structural sections.
type SemanticRegion struct {
	// ID uniquely identifies this region
	ID string

	// Type is the LLM-identified semantic type
	// Example: "problem_statement", "solution_approach", "results_analysis"
	Type string

	// Description explains what this region contains
	Description string

	// Keywords are important terms for this region
	Keywords []string

	// Boundaries define the character positions for this region
	Boundaries []Boundary

	// Confidence is the LLM's confidence in this identification (0.0-1.0)
	Confidence float32

	// RelatedSections lists section IDs that overlap with this region
	RelatedSections []string
}

// Boundary defines a character range for content.
type Boundary struct {
	StartPos int
	EndPos   int
}

// SchemaPattern represents a reusable schema template.
// These patterns can be registered for common document types (like 10-Ks)
// to avoid re-analysis of similar documents.
type SchemaPattern struct {
	// Name of this pattern
	Name string

	// Description of when to use this pattern
	Description string

	// Indicators are features that suggest this pattern applies
	// Example: ["ITEM 1A. Risk Factors", "Form 10-K", "SEC filing"]
	Indicators []string

	// Template is the schema structure to apply
	Template DocumentSchema

	// Matchers are functions that determine if this pattern applies
	// (Defined separately in analyzer.go)

	// Priority affects pattern matching order (higher = checked first)
	Priority int

	// RequiresLLMEnhancement indicates if the pattern should be refined by LLM
	RequiresLLMEnhancement bool
}

// ChunkMetadata represents metadata attached to a chunk in the vector store.
// This metadata enables schema-aware retrieval.
type ChunkMetadata struct {
	// Document identification
	DocID    string
	DocTitle string
	Format   string

	// Section information
	SectionID    string
	SectionTitle string
	SectionType  string
	SectionLevel int

	// Hierarchy information
	HierarchyPath string // Example: "1.2.3"

	// Semantic information
	SemanticTags  []string // Tags from semantic regions
	SemanticTypes []string // Types of semantic regions this chunk belongs to

	// Position information
	StartPos int
	EndPos   int

	// Custom attributes from document schema
	CustomAttributes map[string]interface{}

	// Chunking information
	ChunkIndex  int    // Position in document
	ChunkMethod string // How this chunk was created
}

// ResolverStrategy indicates how to resolve a document schema.
type ResolverStrategy string

const (
	// StrategyExplicit means user provided explicit schema
	StrategyExplicit ResolverStrategy = "explicit"

	// StrategyPattern means schema was matched from a predefined pattern
	StrategyPattern ResolverStrategy = "pattern"

	// StrategyLLM means schema was derived via LLM analysis
	StrategyLLM ResolverStrategy = "llm"

	// StrategyHybrid means pattern was matched then enhanced by LLM
	StrategyHybrid ResolverStrategy = "hybrid"
)

// ResolutionResult contains the outcome of schema resolution.
type ResolutionResult struct {
	// Schema is the resolved document schema
	Schema *DocumentSchema

	// Strategy indicates how the schema was resolved
	Strategy ResolverStrategy

	// PatternUsed is the name of the pattern if StrategyPattern or StrategyHybrid
	PatternUsed string

	// Confidence in the resolved schema (0.0-1.0)
	Confidence float32

	// ProcessingTimeMs tracks how long resolution took
	ProcessingTimeMs int64
}
