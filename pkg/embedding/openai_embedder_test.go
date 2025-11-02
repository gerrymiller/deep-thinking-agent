// Copyright 2025 Gerry Miller <gerry@gerrymiller.com>
//
// Licensed under the MIT License.
// See LICENSE file in the project root for full license information.

package embedding

import (
	"context"
	"testing"
)

func TestNewOpenAIEmbedder(t *testing.T) {
	tests := []struct {
		name    string
		apiKey  string
		model   string
		config  *Config
		wantErr bool
		errMsg  string
	}{
		{
			name:    "valid embedder with defaults",
			apiKey:  "test-api-key",
			model:   "text-embedding-3-small",
			config:  nil,
			wantErr: false,
		},
		{
			name:    "valid embedder with custom config",
			apiKey:  "test-api-key",
			model:   "text-embedding-ada-002",
			config:  &Config{BatchSize: 50, TimeoutSeconds: 60},
			wantErr: false,
		},
		{
			name:    "missing API key",
			apiKey:  "",
			model:   "text-embedding-3-small",
			config:  nil,
			wantErr: true,
			errMsg:  "OpenAI API key is required",
		},
		{
			name:    "missing model",
			apiKey:  "test-api-key",
			model:   "",
			config:  nil,
			wantErr: true,
			errMsg:  "embedding model name is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			embedder, err := NewOpenAIEmbedder(tt.apiKey, tt.model, tt.config)

			if tt.wantErr {
				if err == nil {
					t.Errorf("NewOpenAIEmbedder() expected error but got nil")
				} else if tt.errMsg != "" && err.Error() != tt.errMsg {
					t.Errorf("NewOpenAIEmbedder() error = %v, want %v", err.Error(), tt.errMsg)
				}
				return
			}

			if err != nil {
				t.Errorf("NewOpenAIEmbedder() unexpected error: %v", err)
				return
			}

			if embedder == nil {
				t.Error("NewOpenAIEmbedder() returned nil embedder")
				return
			}

			// Verify embedder properties
			if embedder.ModelName() != tt.model {
				t.Errorf("Embedder.ModelName() = %v, want %v", embedder.ModelName(), tt.model)
			}

			// Verify dimensions are set correctly
			expectedDims := getDimensionsForModel(tt.model)
			if embedder.Dimensions() != expectedDims {
				t.Errorf("Embedder.Dimensions() = %v, want %v", embedder.Dimensions(), expectedDims)
			}
		})
	}
}

func TestGetDimensionsForModel(t *testing.T) {
	tests := []struct {
		model      string
		dimensions int
	}{
		{"text-embedding-3-small", 1536},
		{"text-embedding-3-large", 3072},
		{"text-embedding-ada-002", 1536},
		{"unknown-model", 1536}, // Default
	}

	for _, tt := range tests {
		t.Run(tt.model, func(t *testing.T) {
			got := getDimensionsForModel(tt.model)
			if got != tt.dimensions {
				t.Errorf("getDimensionsForModel(%s) = %v, want %v", tt.model, got, tt.dimensions)
			}
		})
	}
}

func TestOpenAIEmbedder_Methods(t *testing.T) {
	embedder, err := NewOpenAIEmbedder("test-key", "text-embedding-3-small", nil)
	if err != nil {
		t.Fatalf("Failed to create embedder: %v", err)
	}

	// Test ModelName()
	if got := embedder.ModelName(); got != "text-embedding-3-small" {
		t.Errorf("ModelName() = %v, want %v", got, "text-embedding-3-small")
	}

	// Test Dimensions()
	if got := embedder.Dimensions(); got != 1536 {
		t.Errorf("Dimensions() = %v, want %v", got, 1536)
	}
}

func TestOpenAIEmbedder_Embed_ErrorCases(t *testing.T) {
	embedder, err := NewOpenAIEmbedder("test-key", "text-embedding-3-small", nil)
	if err != nil {
		t.Fatalf("Failed to create embedder: %v", err)
	}

	ctx := context.Background()

	tests := []struct {
		name    string
		req     *EmbedRequest
		wantErr bool
		errMsg  string
	}{
		{
			name:    "nil request",
			req:     nil,
			wantErr: true,
			errMsg:  "embed request cannot be nil",
		},
		{
			name:    "empty texts",
			req:     &EmbedRequest{Texts: []string{}},
			wantErr: true,
			errMsg:  "texts cannot be empty",
		},
		{
			name:    "nil texts slice",
			req:     &EmbedRequest{Texts: nil},
			wantErr: true,
			errMsg:  "texts cannot be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := embedder.Embed(ctx, tt.req)

			if tt.wantErr {
				if err == nil {
					t.Error("Embed() expected error but got nil")
				} else if tt.errMsg != "" && err.Error() != tt.errMsg {
					t.Errorf("Embed() error = %v, want %v", err.Error(), tt.errMsg)
				}
			} else if err != nil {
				t.Errorf("Embed() unexpected error: %v", err)
			}
		})
	}
}

func TestOpenAIEmbedder_BaseURL(t *testing.T) {
	// Test that custom BaseURL is respected
	config := &Config{
		Provider:       "openai",
		APIKey:         "test-key",
		BaseURL:        "https://custom.api.endpoint",
		Model:          "text-embedding-3-small",
		BatchSize:      50,
		TimeoutSeconds: 30,
	}

	embedder, err := NewOpenAIEmbedder("test-key", "text-embedding-3-small", config)
	if err != nil {
		t.Fatalf("Failed to create embedder with custom BaseURL: %v", err)
	}

	if embedder == nil {
		t.Error("Expected non-nil embedder with custom BaseURL")
	}

	// Verify config is stored
	if embedder.config.BaseURL != "https://custom.api.endpoint" {
		t.Errorf("BaseURL not preserved, got %v", embedder.config.BaseURL)
	}
}

func TestOpenAIEmbedder_BatchProcessing(t *testing.T) {
	// Test that batch size configuration is respected
	testCases := []struct {
		name      string
		batchSize int
		expected  int
	}{
		{"default batch size", 0, 100},
		{"custom batch size", 50, 50},
		{"large batch size", 200, 200},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			config := &Config{
				Provider:       "openai",
				APIKey:         "test-key",
				Model:          "text-embedding-3-small",
				BatchSize:      tc.batchSize,
				TimeoutSeconds: 30,
			}

			embedder, err := NewOpenAIEmbedder("test-key", "text-embedding-3-small", config)
			if err != nil {
				t.Fatalf("Failed to create embedder: %v", err)
			}

			if embedder.config.BatchSize != tc.batchSize {
				t.Errorf("BatchSize not set correctly, got %d, want %d", embedder.config.BatchSize, tc.batchSize)
			}
		})
	}
}

func TestOpenAIEmbedder_ModelDimensions(t *testing.T) {
	tests := []struct {
		model      string
		dimensions int
	}{
		{"text-embedding-3-small", DimensionsTextEmbedding3Small},
		{"text-embedding-3-large", DimensionsTextEmbedding3Large},
		{"text-embedding-ada-002", DimensionsTextEmbeddingAda002},
		{"future-unknown-model", DimensionsTextEmbeddingAda002}, // Should default
	}

	for _, tt := range tests {
		t.Run(tt.model, func(t *testing.T) {
			embedder, err := NewOpenAIEmbedder("test-key", tt.model, nil)
			if err != nil {
				t.Fatalf("Failed to create embedder: %v", err)
			}

			if embedder.Dimensions() != tt.dimensions {
				t.Errorf("Dimensions() = %d, want %d", embedder.Dimensions(), tt.dimensions)
			}
		})
	}
}
