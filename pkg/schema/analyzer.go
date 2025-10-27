// Copyright 2025 Gerry Miller <gerry@gerrymiller.com>
//
// Licensed under the MIT License.
// See LICENSE file in the project root for full license information.

package schema

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"deep-thinking-agent/pkg/llm"
)

// Analyzer uses an LLM to derive document schemas.
// It analyzes document structure to identify sections, hierarchy,
// semantic regions, and recommend chunking strategies.
type Analyzer struct {
	llmProvider llm.Provider
	temperature float32
	maxTokens   int
	timeout     time.Duration
}

// AnalyzerConfig contains configuration for the schema analyzer.
type AnalyzerConfig struct {
	Temperature float32
	MaxTokens   int
	Timeout     time.Duration
}

// NewAnalyzer creates a new schema analyzer instance.
func NewAnalyzer(provider llm.Provider, config *AnalyzerConfig) *Analyzer {
	if config == nil {
		config = &AnalyzerConfig{
			Temperature: 0.3, // Lower temperature for more structured output
			MaxTokens:   3000,
			Timeout:     60 * time.Second,
		}
	}

	return &Analyzer{
		llmProvider: provider,
		temperature: config.Temperature,
		maxTokens:   config.MaxTokens,
		timeout:     config.Timeout,
	}
}

// AnalyzeDocument performs LLM-based analysis of a document to derive its schema.
// Returns a DocumentSchema with identified structure and metadata.
func (a *Analyzer) AnalyzeDocument(ctx context.Context, docID, content, format string) (*DocumentSchema, error) {
	startTime := time.Now()

	// Apply timeout
	if a.timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, a.timeout)
		defer cancel()
	}

	// Build analysis prompt
	prompt := a.buildAnalysisPrompt(content, format)

	// Call LLM
	resp, err := a.llmProvider.Complete(ctx, &llm.CompletionRequest{
		Messages: []llm.Message{
			{
				Role:    "system",
				Content: systemPrompt,
			},
			{
				Role:    "user",
				Content: prompt,
			},
		},
		Temperature: a.temperature,
		MaxTokens:   a.maxTokens,
	})

	if err != nil {
		return nil, fmt.Errorf("LLM analysis failed: %w", err)
	}

	// Parse LLM response into schema
	schema, err := a.parseAnalysisResponse(resp.Content, docID, format)
	if err != nil {
		return nil, fmt.Errorf("failed to parse analysis: %w", err)
	}

	// Set metadata
	schema.ParsingMethod = "llm_analysis"
	schema.CreatedAt = time.Now().Unix()

	// Calculate processing time
	processingTime := time.Since(startTime).Milliseconds()
	_ = processingTime // Store in schema metadata if needed

	return schema, nil
}

// buildAnalysisPrompt constructs the prompt for LLM analysis.
func (a *Analyzer) buildAnalysisPrompt(content, format string) string {
	// Truncate content if too long (leave room for response)
	maxContentLength := 8000 // Rough estimate to stay within context limits
	truncatedContent := content
	if len(content) > maxContentLength {
		truncatedContent = content[:maxContentLength] + "\n\n[Content truncated for analysis...]"
	}

	return fmt.Sprintf(`Analyze the following %s document and provide a detailed structural schema.

Document content:
---
%s
---

Provide your analysis as a JSON object with the following structure:
{
  "title": "document title if identifiable",
  "sections": [
    {
      "id": "unique_section_id",
      "title": "section title",
      "level": 1,
      "start_pos": 0,
      "end_pos": 100,
      "type": "semantic_type (e.g., introduction, methodology, results)",
      "summary": "brief section summary",
      "keywords": ["key", "terms"]
    }
  ],
  "semantic_regions": [
    {
      "id": "region_id",
      "type": "region type (e.g., problem_statement, solution_approach)",
      "description": "what this region contains",
      "keywords": ["relevant", "terms"],
      "boundaries": [{"start_pos": 0, "end_pos": 100}],
      "confidence": 0.9
    }
  ],
  "custom_attributes": {
    "key": "value pairs of document-specific metadata"
  },
  "chunking_strategy": "recommended strategy: section_based, hierarchical, semantic, or sliding_window",
  "confidence": 0.9
}

Focus on:
1. Identifying logical sections with clear boundaries
2. Building hierarchical structure (headings, subheadings)
3. Recognizing semantic regions that span multiple structural sections
4. Extracting meaningful custom attributes
5. Recommending an appropriate chunking strategy`, format, truncatedContent)
}

// parseAnalysisResponse parses the LLM's JSON response into a DocumentSchema.
func (a *Analyzer) parseAnalysisResponse(response, docID, format string) (*DocumentSchema, error) {
	// Try to extract JSON from response (LLM might include explanation text)
	jsonStart := findJSONStart(response)
	jsonEnd := findJSONEnd(response[jsonStart:])
	if jsonStart == -1 || jsonEnd == -1 {
		return nil, fmt.Errorf("no valid JSON found in LLM response")
	}

	jsonStr := response[jsonStart : jsonStart+jsonEnd+1]

	// Parse JSON
	var parsed struct {
		Title    string `json:"title"`
		Sections []struct {
			ID       string   `json:"id"`
			Title    string   `json:"title"`
			Level    int      `json:"level"`
			StartPos int      `json:"start_pos"`
			EndPos   int      `json:"end_pos"`
			Type     string   `json:"type"`
			Summary  string   `json:"summary"`
			Keywords []string `json:"keywords"`
		} `json:"sections"`
		SemanticRegions []struct {
			ID          string   `json:"id"`
			Type        string   `json:"type"`
			Description string   `json:"description"`
			Keywords    []string `json:"keywords"`
			Boundaries  []struct {
				StartPos int `json:"start_pos"`
				EndPos   int `json:"end_pos"`
			} `json:"boundaries"`
			Confidence float32 `json:"confidence"`
		} `json:"semantic_regions"`
		CustomAttributes map[string]interface{} `json:"custom_attributes"`
		ChunkingStrategy string                 `json:"chunking_strategy"`
		Confidence       float32                `json:"confidence"`
	}

	if err := json.Unmarshal([]byte(jsonStr), &parsed); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	// Convert to DocumentSchema
	schema := &DocumentSchema{
		DocID:            docID,
		Format:           format,
		Title:            parsed.Title,
		Sections:         make([]Section, len(parsed.Sections)),
		SemanticRegions:  make([]SemanticRegion, len(parsed.SemanticRegions)),
		CustomAttributes: parsed.CustomAttributes,
		ChunkingStrategy: parsed.ChunkingStrategy,
		Confidence:       parsed.Confidence,
		ChunkMetadata:    make(map[string]interface{}),
	}

	// Convert sections
	for i, s := range parsed.Sections {
		schema.Sections[i] = Section{
			ID:       s.ID,
			Title:    s.Title,
			Level:    s.Level,
			StartPos: s.StartPos,
			EndPos:   s.EndPos,
			Type:     s.Type,
			Summary:  s.Summary,
			Keywords: s.Keywords,
			ChildIDs: []string{},
		}
	}

	// Convert semantic regions
	for i, r := range parsed.SemanticRegions {
		boundaries := make([]Boundary, len(r.Boundaries))
		for j, b := range r.Boundaries {
			boundaries[j] = Boundary{
				StartPos: b.StartPos,
				EndPos:   b.EndPos,
			}
		}

		schema.SemanticRegions[i] = SemanticRegion{
			ID:          r.ID,
			Type:        r.Type,
			Description: r.Description,
			Keywords:    r.Keywords,
			Boundaries:  boundaries,
			Confidence:  r.Confidence,
		}
	}

	// Build hierarchy from sections
	schema.Hierarchy = a.buildHierarchy(schema.Sections)

	return schema, nil
}

// buildHierarchy constructs a hierarchical tree from flat sections.
func (a *Analyzer) buildHierarchy(sections []Section) *HierarchyTree {
	if len(sections) == 0 {
		return &HierarchyTree{MaxDepth: 0}
	}

	// Find root sections (level 1)
	var rootNodes []*HierarchyNode
	maxDepth := 0

	for _, section := range sections {
		if section.Level > maxDepth {
			maxDepth = section.Level
		}

		if section.Level == 1 {
			node := &HierarchyNode{
				ID:       section.ID,
				Path:     fmt.Sprintf("%d", len(rootNodes)+1),
				Title:    section.Title,
				Level:    section.Level,
				StartPos: section.StartPos,
				EndPos:   section.EndPos,
				Children: []*HierarchyNode{},
			}
			rootNodes = append(rootNodes, node)
		}
	}

	// For simplicity, create flat structure
	// TODO: Build proper parent-child relationships in future enhancement
	root := &HierarchyNode{
		ID:       "root",
		Path:     "0",
		Title:    "Document Root",
		Level:    0,
		Children: rootNodes,
	}

	return &HierarchyTree{
		Root:     root,
		MaxDepth: maxDepth,
	}
}

// findJSONStart finds the start of JSON in text (looks for opening brace).
func findJSONStart(text string) int {
	for i, ch := range text {
		if ch == '{' {
			return i
		}
	}
	return -1
}

// findJSONEnd finds the matching closing brace for JSON.
func findJSONEnd(text string) int {
	depth := 0
	for i, ch := range text {
		if ch == '{' {
			depth++
		} else if ch == '}' {
			depth--
			if depth == 0 {
				return i
			}
		}
	}
	return -1
}

// systemPrompt is the system message for the schema analyzer.
const systemPrompt = `You are a document structure analysis expert. Your task is to analyze documents and extract their structural schema.

Your analysis should identify:
1. Logical sections with clear boundaries and semantic types
2. Hierarchical structure (headings, subheadings, nesting)
3. Semantic regions (topic-based areas that may span multiple sections)
4. Custom attributes specific to the document type
5. An appropriate chunking strategy for RAG systems

Always respond with valid JSON matching the requested structure. Be precise with position markers (start_pos, end_pos) and provide high-quality semantic types and keywords.`
