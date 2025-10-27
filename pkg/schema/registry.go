// Copyright 2025 Gerry Miller <gerry@gerrymiller.com>
//
// Licensed under the MIT License.
// See LICENSE file in the project root for full license information.

package schema

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sort"
	"sync"
)

// Registry stores and manages predefined schema patterns.
// Patterns can be registered programmatically or loaded from files.
type Registry struct {
	patterns map[string]SchemaPattern
	mu       sync.RWMutex
}

// NewRegistry creates a new schema pattern registry.
func NewRegistry() *Registry {
	return &Registry{
		patterns: make(map[string]SchemaPattern),
	}
}

// Register adds a schema pattern to the registry.
func (r *Registry) Register(pattern SchemaPattern) error {
	if pattern.Name == "" {
		return errors.New("pattern name is required")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	r.patterns[pattern.Name] = pattern
	return nil
}

// Get retrieves a schema pattern by name.
func (r *Registry) Get(name string) (*SchemaPattern, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	pattern, ok := r.patterns[name]
	if !ok {
		return nil, false
	}

	return &pattern, true
}

// List returns all registered patterns sorted by priority (highest first).
func (r *Registry) List() []SchemaPattern {
	r.mu.RLock()
	defer r.mu.RUnlock()

	patterns := make([]SchemaPattern, 0, len(r.patterns))
	for _, pattern := range r.patterns {
		patterns = append(patterns, pattern)
	}

	// Sort by priority (descending)
	sort.Slice(patterns, func(i, j int) bool {
		return patterns[i].Priority > patterns[j].Priority
	})

	return patterns
}

// Delete removes a schema pattern from the registry.
func (r *Registry) Delete(name string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.patterns[name]; !ok {
		return fmt.Errorf("pattern not found: %s", name)
	}

	delete(r.patterns, name)
	return nil
}

// Count returns the number of registered patterns.
func (r *Registry) Count() int {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return len(r.patterns)
}

// Clear removes all patterns from the registry.
func (r *Registry) Clear() {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.patterns = make(map[string]SchemaPattern)
}

// LoadFromFile loads a schema pattern from a JSON file.
func (r *Registry) LoadFromFile(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	var pattern SchemaPattern
	if err := json.Unmarshal(data, &pattern); err != nil {
		return fmt.Errorf("failed to parse JSON: %w", err)
	}

	return r.Register(pattern)
}

// SaveToFile saves a schema pattern to a JSON file.
func (r *Registry) SaveToFile(name, path string) error {
	r.mu.RLock()
	pattern, ok := r.patterns[name]
	r.mu.RUnlock()

	if !ok {
		return fmt.Errorf("pattern not found: %s", name)
	}

	data, err := json.MarshalIndent(pattern, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

// LoadDirectory loads all JSON schema patterns from a directory.
func (r *Registry) LoadDirectory(dirPath string) error {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return fmt.Errorf("failed to read directory: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		// Only process .json files
		name := entry.Name()
		if len(name) > 5 && name[len(name)-5:] == ".json" {
			path := fmt.Sprintf("%s/%s", dirPath, name)
			if err := r.LoadFromFile(path); err != nil {
				// Log error but continue loading other files
				fmt.Printf("Warning: failed to load %s: %v\n", path, err)
			}
		}
	}

	return nil
}

// RegisterBuiltInPatterns registers common predefined patterns.
// This provides examples and can be extended with more patterns.
func (r *Registry) RegisterBuiltInPatterns() {
	// Example: SEC 10-K filing pattern
	sec10K := SchemaPattern{
		Name:        "sec_10k",
		Description: "SEC 10-K Annual Report filing",
		Indicators: []string{
			"ITEM 1A. Risk Factors",
			"ITEM 1. Business",
			"ITEM 7. Management's Discussion",
			"Form 10-K",
		},
		Priority:               100,
		RequiresLLMEnhancement: false,
		Template: DocumentSchema{
			Format:           "pdf",
			ChunkingStrategy: "section_based",
			Sections: []Section{
				{ID: "item1", Title: "Item 1. Business", Level: 1, Type: "business_overview"},
				{ID: "item1a", Title: "Item 1A. Risk Factors", Level: 1, Type: "risk_factors"},
				{ID: "item7", Title: "Item 7. MD&A", Level: 1, Type: "financial_analysis"},
			},
			CustomAttributes: map[string]interface{}{
				"document_type": "sec_filing",
				"filing_type":   "10-K",
			},
		},
	}
	r.Register(sec10K)

	// Example: Research paper pattern
	researchPaper := SchemaPattern{
		Name:        "research_paper",
		Description: "Academic research paper",
		Indicators: []string{
			"Abstract",
			"Introduction",
			"Methodology",
			"Results",
			"Conclusion",
			"References",
		},
		Priority:               90,
		RequiresLLMEnhancement: true,
		Template: DocumentSchema{
			Format:           "pdf",
			ChunkingStrategy: "hierarchical",
			Sections: []Section{
				{ID: "abstract", Title: "Abstract", Level: 1, Type: "abstract"},
				{ID: "intro", Title: "Introduction", Level: 1, Type: "introduction"},
				{ID: "method", Title: "Methodology", Level: 1, Type: "methodology"},
				{ID: "results", Title: "Results", Level: 1, Type: "results"},
				{ID: "conclusion", Title: "Conclusion", Level: 1, Type: "conclusion"},
			},
			CustomAttributes: map[string]interface{}{
				"document_type": "research_paper",
			},
		},
	}
	r.Register(researchPaper)
}
