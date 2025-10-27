// Copyright 2025 Gerry Miller <gerry@gerrymiller.com>
//
// Licensed under the MIT License.
// See LICENSE file in the project root for full license information.

package schema

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewRegistry(t *testing.T) {
	registry := NewRegistry()
	if registry == nil {
		t.Fatal("NewRegistry() returned nil")
	}
	if registry.patterns == nil {
		t.Error("patterns map not initialized")
	}
	if registry.Count() != 0 {
		t.Errorf("Count() = %v, want 0", registry.Count())
	}
}

func TestRegister(t *testing.T) {
	registry := NewRegistry()

	tests := []struct {
		name    string
		pattern SchemaPattern
		wantErr bool
	}{
		{
			name: "valid pattern",
			pattern: SchemaPattern{
				Name:        "test_pattern",
				Description: "Test",
				Indicators:  []string{"test"},
				Priority:    100,
			},
			wantErr: false,
		},
		{
			name: "pattern without name",
			pattern: SchemaPattern{
				Description: "Test",
				Indicators:  []string{"test"},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := registry.Register(tt.pattern)
			if (err != nil) != tt.wantErr {
				t.Errorf("Register() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr {
				if registry.Count() == 0 {
					t.Error("pattern not registered")
				}
			}
		})
	}
}

func TestGet(t *testing.T) {
	registry := NewRegistry()

	pattern := SchemaPattern{
		Name:        "test_pattern",
		Description: "Test",
		Indicators:  []string{"test"},
		Priority:    100,
	}
	registry.Register(pattern)

	tests := []struct {
		name      string
		key       string
		wantFound bool
	}{
		{
			name:      "existing pattern",
			key:       "test_pattern",
			wantFound: true,
		},
		{
			name:      "non-existent pattern",
			key:       "missing",
			wantFound: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, found := registry.Get(tt.key)
			if found != tt.wantFound {
				t.Errorf("Get() found = %v, want %v", found, tt.wantFound)
			}
			if tt.wantFound && got.Name != tt.key {
				t.Errorf("Get() Name = %v, want %v", got.Name, tt.key)
			}
		})
	}
}

func TestList(t *testing.T) {
	registry := NewRegistry()

	// Add patterns with different priorities
	patterns := []SchemaPattern{
		{Name: "low", Priority: 10},
		{Name: "high", Priority: 100},
		{Name: "medium", Priority: 50},
	}

	for _, p := range patterns {
		registry.Register(p)
	}

	list := registry.List()

	if len(list) != 3 {
		t.Errorf("List() length = %v, want 3", len(list))
	}

	// Check that list is sorted by priority (descending)
	if len(list) == 3 {
		if list[0].Name != "high" {
			t.Errorf("first pattern = %v, want high", list[0].Name)
		}
		if list[1].Name != "medium" {
			t.Errorf("second pattern = %v, want medium", list[1].Name)
		}
		if list[2].Name != "low" {
			t.Errorf("third pattern = %v, want low", list[2].Name)
		}
	}
}

func TestDelete(t *testing.T) {
	registry := NewRegistry()

	pattern := SchemaPattern{
		Name:     "test",
		Priority: 100,
	}
	registry.Register(pattern)

	// Delete existing pattern
	err := registry.Delete("test")
	if err != nil {
		t.Errorf("Delete() error = %v, want nil", err)
	}

	if registry.Count() != 0 {
		t.Errorf("Count() = %v, want 0", registry.Count())
	}

	// Delete non-existent pattern
	err = registry.Delete("missing")
	if err == nil {
		t.Error("Delete() error = nil, want error")
	}
}

func TestCount(t *testing.T) {
	registry := NewRegistry()

	if registry.Count() != 0 {
		t.Errorf("Count() = %v, want 0", registry.Count())
	}

	registry.Register(SchemaPattern{Name: "p1", Priority: 100})
	if registry.Count() != 1 {
		t.Errorf("Count() = %v, want 1", registry.Count())
	}

	registry.Register(SchemaPattern{Name: "p2", Priority: 100})
	if registry.Count() != 2 {
		t.Errorf("Count() = %v, want 2", registry.Count())
	}

	registry.Delete("p1")
	if registry.Count() != 1 {
		t.Errorf("Count() = %v, want 1", registry.Count())
	}
}

func TestClear(t *testing.T) {
	registry := NewRegistry()

	registry.Register(SchemaPattern{Name: "p1", Priority: 100})
	registry.Register(SchemaPattern{Name: "p2", Priority: 100})

	registry.Clear()

	if registry.Count() != 0 {
		t.Errorf("Count() = %v, want 0 after Clear()", registry.Count())
	}
}

func TestLoadFromFile(t *testing.T) {
	registry := NewRegistry()

	// Create temp directory for test files
	tmpDir := t.TempDir()

	// Valid JSON file
	validJSON := `{
		"name": "test_pattern",
		"description": "Test pattern",
		"indicators": ["test"],
		"priority": 100,
		"requires_llm_enhancement": false,
		"template": {
			"doc_id": "",
			"format": "text",
			"title": "",
			"sections": [],
			"hierarchy": null,
			"semantic_regions": [],
			"custom_attributes": {},
			"chunking_strategy": "section_based",
			"chunk_metadata": {},
			"parsing_method": "",
			"confidence": 0.0,
			"created_at": 0
		}
	}`

	validPath := filepath.Join(tmpDir, "valid.json")
	if err := os.WriteFile(validPath, []byte(validJSON), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	// Invalid JSON file
	invalidJSON := `{"name": "test", invalid json}`
	invalidPath := filepath.Join(tmpDir, "invalid.json")
	if err := os.WriteFile(invalidPath, []byte(invalidJSON), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	tests := []struct {
		name    string
		path    string
		wantErr bool
	}{
		{
			name:    "valid JSON file",
			path:    validPath,
			wantErr: false,
		},
		{
			name:    "invalid JSON file",
			path:    invalidPath,
			wantErr: true,
		},
		{
			name:    "non-existent file",
			path:    filepath.Join(tmpDir, "missing.json"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := registry.LoadFromFile(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadFromFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSaveToFile(t *testing.T) {
	registry := NewRegistry()
	tmpDir := t.TempDir()

	pattern := SchemaPattern{
		Name:        "test_pattern",
		Description: "Test",
		Indicators:  []string{"test"},
		Priority:    100,
	}
	registry.Register(pattern)

	tests := []struct {
		name    string
		key     string
		path    string
		wantErr bool
	}{
		{
			name:    "save existing pattern",
			key:     "test_pattern",
			path:    filepath.Join(tmpDir, "output.json"),
			wantErr: false,
		},
		{
			name:    "save non-existent pattern",
			key:     "missing",
			path:    filepath.Join(tmpDir, "missing.json"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := registry.SaveToFile(tt.key, tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("SaveToFile() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr {
				// Verify file was created
				if _, err := os.Stat(tt.path); os.IsNotExist(err) {
					t.Error("file was not created")
				}
			}
		})
	}
}

func TestLoadDirectory(t *testing.T) {
	registry := NewRegistry()
	tmpDir := t.TempDir()

	// Create valid JSON file
	validJSON := `{
		"name": "pattern1",
		"description": "Test",
		"indicators": ["test"],
		"priority": 100,
		"requires_llm_enhancement": false,
		"template": {
			"doc_id": "",
			"format": "text",
			"title": "",
			"sections": [],
			"hierarchy": null,
			"semantic_regions": [],
			"custom_attributes": {},
			"chunking_strategy": "section_based",
			"chunk_metadata": {},
			"parsing_method": "",
			"confidence": 0.0,
			"created_at": 0
		}
	}`
	os.WriteFile(filepath.Join(tmpDir, "pattern1.json"), []byte(validJSON), 0644)

	// Create non-JSON file (should be skipped)
	os.WriteFile(filepath.Join(tmpDir, "readme.txt"), []byte("not json"), 0644)

	// Create subdirectory (should be skipped)
	subDir := filepath.Join(tmpDir, "subdir")
	os.Mkdir(subDir, 0755)

	err := registry.LoadDirectory(tmpDir)
	if err != nil {
		t.Errorf("LoadDirectory() error = %v, want nil", err)
	}

	if registry.Count() != 1 {
		t.Errorf("Count() = %v, want 1", registry.Count())
	}

	// Test with non-existent directory
	err = registry.LoadDirectory("/nonexistent/path")
	if err == nil {
		t.Error("LoadDirectory() with non-existent path should return error")
	}
}

func TestRegisterBuiltInPatterns(t *testing.T) {
	registry := NewRegistry()
	registry.RegisterBuiltInPatterns()

	if registry.Count() == 0 {
		t.Error("RegisterBuiltInPatterns() did not register any patterns")
	}

	// Check for specific built-in patterns
	sec10k, found := registry.Get("sec_10k")
	if !found {
		t.Error("sec_10k pattern not registered")
	} else {
		if sec10k.Priority != 100 {
			t.Errorf("sec_10k priority = %v, want 100", sec10k.Priority)
		}
	}

	researchPaper, found := registry.Get("research_paper")
	if !found {
		t.Error("research_paper pattern not registered")
	} else {
		if researchPaper.Priority != 90 {
			t.Errorf("research_paper priority = %v, want 90", researchPaper.Priority)
		}
	}
}
