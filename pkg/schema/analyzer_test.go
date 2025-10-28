// Copyright 2025 Gerry Miller <gerry@gerrymiller.com>
//
// Licensed under the MIT License.
// See LICENSE file in the project root for full license information.

package schema

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"deep-thinking-agent/pkg/llm"
)

// MockLLMProvider for testing
type mockLLMProvider struct {
	response    string
	err         error
	callCount   int
	lastRequest *llm.CompletionRequest
	shouldDelay time.Duration
}

func (m *mockLLMProvider) Complete(ctx context.Context, req *llm.CompletionRequest) (*llm.CompletionResponse, error) {
	m.callCount++
	m.lastRequest = req

	if m.shouldDelay > 0 {
		select {
		case <-time.After(m.shouldDelay):
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}

	if m.err != nil {
		return nil, m.err
	}

	return &llm.CompletionResponse{
		Content:      m.response,
		FinishReason: "stop",
		Usage: llm.UsageStats{
			PromptTokens:     100,
			CompletionTokens: 200,
			TotalTokens:      300,
		},
	}, nil
}

func (m *mockLLMProvider) Name() string {
	return "mock"
}

func (m *mockLLMProvider) ModelName() string {
	return "mock-model"
}

func (m *mockLLMProvider) SupportsStreaming() bool {
	return false
}

func TestNewAnalyzer(t *testing.T) {
	tests := []struct {
		name           string
		config         *AnalyzerConfig
		expectedTemp   float32
		expectedTokens int
		expectedTO     time.Duration
	}{
		{
			name:           "with nil config uses defaults",
			config:         nil,
			expectedTemp:   0.3,
			expectedTokens: 3000,
			expectedTO:     60 * time.Second,
		},
		{
			name: "with custom config",
			config: &AnalyzerConfig{
				Temperature: 0.5,
				MaxTokens:   5000,
				Timeout:     30 * time.Second,
			},
			expectedTemp:   0.5,
			expectedTokens: 5000,
			expectedTO:     30 * time.Second,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			provider := &mockLLMProvider{}
			analyzer := NewAnalyzer(provider, tt.config)

			if analyzer.llmProvider == nil {
				t.Error("LLM provider not set")
			}
			if analyzer.temperature != tt.expectedTemp {
				t.Errorf("temperature = %v, want %v", analyzer.temperature, tt.expectedTemp)
			}
			if analyzer.maxTokens != tt.expectedTokens {
				t.Errorf("maxTokens = %v, want %v", analyzer.maxTokens, tt.expectedTokens)
			}
			if analyzer.timeout != tt.expectedTO {
				t.Errorf("timeout = %v, want %v", analyzer.timeout, tt.expectedTO)
			}
		})
	}
}

func TestAnalyzeDocument(t *testing.T) {
	validResponse := `{
		"title": "Test Document",
		"sections": [
			{
				"id": "sec1",
				"title": "Section 1",
				"level": 1,
				"start_pos": 0,
				"end_pos": 100,
				"type": "introduction",
				"summary": "First section",
				"keywords": ["test", "intro"]
			}
		],
		"semantic_regions": [
			{
				"id": "reg1",
				"type": "problem_statement",
				"description": "Problem description",
				"keywords": ["problem"],
				"boundaries": [{"start_pos": 0, "end_pos": 50}],
				"confidence": 0.9
			}
		],
		"custom_attributes": {
			"doc_type": "test"
		},
		"chunking_strategy": "section_based",
		"confidence": 0.95
	}`

	tests := []struct {
		name        string
		provider    *mockLLMProvider
		content     string
		format      string
		wantErr     bool
		errContains string
	}{
		{
			name: "successful analysis",
			provider: &mockLLMProvider{
				response: validResponse,
			},
			content: "Test content",
			format:  "text",
			wantErr: false,
		},
		{
			name: "LLM error",
			provider: &mockLLMProvider{
				err: errors.New("API error"),
			},
			content:     "Test content",
			format:      "text",
			wantErr:     true,
			errContains: "LLM analysis failed",
		},
		{
			name: "invalid JSON response",
			provider: &mockLLMProvider{
				response: "Not valid JSON at all",
			},
			content:     "Test content",
			format:      "text",
			wantErr:     true,
			errContains: "failed to parse analysis",
		},
		{
			name: "JSON with explanation text",
			provider: &mockLLMProvider{
				response: "Here is the analysis:\n" + validResponse + "\nThat's the result.",
			},
			content: "Test content",
			format:  "text",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			analyzer := NewAnalyzer(tt.provider, nil)
			schema, err := analyzer.AnalyzeDocument(context.Background(), "doc1", tt.content, tt.format)

			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				if tt.errContains != "" && !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("error = %v, should contain %v", err, tt.errContains)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if schema.DocID != "doc1" {
				t.Errorf("DocID = %v, want doc1", schema.DocID)
			}
			if schema.Format != tt.format {
				t.Errorf("Format = %v, want %v", schema.Format, tt.format)
			}
			if schema.ParsingMethod != "llm_analysis" {
				t.Errorf("ParsingMethod = %v, want llm_analysis", schema.ParsingMethod)
			}
			if schema.CreatedAt == 0 {
				t.Error("CreatedAt not set")
			}
		})
	}
}

func TestAnalyzeDocumentTimeout(t *testing.T) {
	provider := &mockLLMProvider{
		shouldDelay: 200 * time.Millisecond,
		response:    "{}",
	}

	config := &AnalyzerConfig{
		Temperature: 0.3,
		MaxTokens:   3000,
		Timeout:     50 * time.Millisecond,
	}

	analyzer := NewAnalyzer(provider, config)
	ctx := context.Background()

	_, err := analyzer.AnalyzeDocument(ctx, "doc1", "content", "text")
	if err == nil {
		t.Fatal("expected timeout error, got nil")
	}
}

func TestBuildAnalysisPrompt(t *testing.T) {
	analyzer := NewAnalyzer(&mockLLMProvider{}, nil)

	tests := []struct {
		name           string
		content        string
		format         string
		expectTruncate bool
	}{
		{
			name:           "short content not truncated",
			content:        "Short content",
			format:         "text",
			expectTruncate: false,
		},
		{
			name:           "long content truncated",
			content:        strings.Repeat("a", 9000),
			format:         "pdf",
			expectTruncate: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			prompt := analyzer.buildAnalysisPrompt(tt.content, tt.format)

			if !strings.Contains(prompt, tt.format) {
				t.Errorf("prompt should contain format %v", tt.format)
			}

			if tt.expectTruncate {
				if strings.Contains(prompt, strings.Repeat("a", 9000)) {
					t.Error("content should be truncated")
				}
				if !strings.Contains(prompt, "[Content truncated for analysis...]") {
					t.Error("should contain truncation message")
				}
			}
		})
	}
}

func TestParseAnalysisResponse(t *testing.T) {
	analyzer := NewAnalyzer(&mockLLMProvider{}, nil)

	tests := []struct {
		name        string
		response    string
		wantErr     bool
		errContains string
		validate    func(*testing.T, *DocumentSchema)
	}{
		{
			name: "valid JSON",
			response: `{
				"title": "Test Doc",
				"sections": [
					{
						"id": "s1",
						"title": "Sec 1",
						"level": 1,
						"start_pos": 0,
						"end_pos": 100,
						"type": "intro",
						"summary": "Summary",
						"keywords": ["key"]
					}
				],
				"semantic_regions": [],
				"custom_attributes": {"type": "test"},
				"chunking_strategy": "section_based",
				"confidence": 0.9
			}`,
			wantErr: false,
			validate: func(t *testing.T, s *DocumentSchema) {
				if s.Title != "Test Doc" {
					t.Errorf("Title = %v, want Test Doc", s.Title)
				}
				if len(s.Sections) != 1 {
					t.Errorf("Sections count = %v, want 1", len(s.Sections))
				}
				if s.Sections[0].ID != "s1" {
					t.Errorf("Section ID = %v, want s1", s.Sections[0].ID)
				}
				if s.ChunkingStrategy != "section_based" {
					t.Errorf("ChunkingStrategy = %v, want section_based", s.ChunkingStrategy)
				}
			},
		},
		{
			name:        "no JSON found",
			response:    "This is plain text with no JSON",
			wantErr:     true,
			errContains: "no valid JSON found",
		},
		{
			name:        "malformed JSON",
			response:    `{"title": "Test", "sections": [}`,
			wantErr:     true,
			errContains: "failed to parse JSON",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			schema, err := analyzer.parseAnalysisResponse(tt.response, "doc1", "text")

			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				if tt.errContains != "" && !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("error = %v, should contain %v", err, tt.errContains)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if tt.validate != nil {
				tt.validate(t, schema)
			}
		})
	}
}

func TestBuildHierarchy(t *testing.T) {
	analyzer := NewAnalyzer(&mockLLMProvider{}, nil)

	tests := []struct {
		name        string
		sections    []Section
		expectDepth int
		expectRoots int
	}{
		{
			name:        "empty sections",
			sections:    []Section{},
			expectDepth: 0,
			expectRoots: 0,
		},
		{
			name: "single level sections",
			sections: []Section{
				{ID: "s1", Title: "Section 1", Level: 1},
				{ID: "s2", Title: "Section 2", Level: 1},
			},
			expectDepth: 1,
			expectRoots: 2,
		},
		{
			name: "multi-level sections",
			sections: []Section{
				{ID: "s1", Title: "Section 1", Level: 1},
				{ID: "s2", Title: "Section 1.1", Level: 2},
				{ID: "s3", Title: "Section 1.1.1", Level: 3},
			},
			expectDepth: 3,
			expectRoots: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hierarchy := analyzer.buildHierarchy(tt.sections)

			if hierarchy.MaxDepth != tt.expectDepth {
				t.Errorf("MaxDepth = %v, want %v", hierarchy.MaxDepth, tt.expectDepth)
			}

			if tt.expectRoots > 0 && hierarchy.Root != nil {
				if len(hierarchy.Root.Children) != tt.expectRoots {
					t.Errorf("Root children = %v, want %v", len(hierarchy.Root.Children), tt.expectRoots)
				}
			}
		})
	}
}

func TestFindJSONStart(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		expected int
	}{
		{
			name:     "JSON at start",
			text:     `{"key": "value"}`,
			expected: 0,
		},
		{
			name:     "JSON with prefix text",
			text:     `Here is the result: {"key": "value"}`,
			expected: 20,
		},
		{
			name:     "no JSON",
			text:     "No JSON here",
			expected: -1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := findJSONStart(tt.text)
			if result != tt.expected {
				t.Errorf("findJSONStart() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestFindJSONEnd(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		expected int
	}{
		{
			name:     "simple JSON",
			text:     `{"key": "value"}`,
			expected: 15,
		},
		{
			name:     "nested JSON",
			text:     `{"outer": {"inner": "value"}}`,
			expected: 28,
		},
		{
			name:     "JSON with suffix",
			text:     `{"key": "value"} and more text`,
			expected: 15,
		},
		{
			name:     "unclosed JSON",
			text:     `{"key": "value"`,
			expected: -1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := findJSONEnd(tt.text)
			if result != tt.expected {
				t.Errorf("findJSONEnd() = %v, want %v", result, tt.expected)
			}
		})
	}
}
