// Copyright 2025 Gerry Miller <gerry@gerrymiller.com>
//
// Licensed under the MIT License.
// See LICENSE file in the project root for full license information.

package embedding

import (
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
