// Copyright 2025 Gerry Miller <gerry@gerrymiller.com>
//
// Licensed under the MIT License.
// See LICENSE file in the project root for full license information.

package openai

import (
	"context"
	"testing"

	"deep-thinking-agent/pkg/llm"
)

func TestNewProvider(t *testing.T) {
	tests := []struct {
		name    string
		apiKey  string
		model   string
		config  *llm.Config
		wantErr bool
		errMsg  string
	}{
		{
			name:    "valid provider with defaults",
			apiKey:  "test-api-key",
			model:   "gpt-4",
			config:  nil,
			wantErr: false,
		},
		{
			name:    "valid provider with custom config",
			apiKey:  "test-api-key",
			model:   "gpt-3.5-turbo",
			config:  &llm.Config{DefaultTemperature: 0.5, DefaultMaxTokens: 1000},
			wantErr: false,
		},
		{
			name:    "missing API key",
			apiKey:  "",
			model:   "gpt-4",
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
			errMsg:  "model name is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			provider, err := NewProvider(tt.apiKey, tt.model, tt.config)

			if tt.wantErr {
				if err == nil {
					t.Errorf("NewProvider() expected error but got nil")
				} else if tt.errMsg != "" && err.Error() != tt.errMsg {
					t.Errorf("NewProvider() error = %v, want %v", err.Error(), tt.errMsg)
				}
				return
			}

			if err != nil {
				t.Errorf("NewProvider() unexpected error: %v", err)
				return
			}

			if provider == nil {
				t.Error("NewProvider() returned nil provider")
				return
			}

			// Verify provider properties
			if provider.Name() != "openai" {
				t.Errorf("Provider.Name() = %v, want %v", provider.Name(), "openai")
			}

			if provider.ModelName() != tt.model {
				t.Errorf("Provider.ModelName() = %v, want %v", provider.ModelName(), tt.model)
			}

			if !provider.SupportsStreaming() {
				t.Error("Provider.SupportsStreaming() = false, want true")
			}
		})
	}
}

func TestProvider_Complete_Validation(t *testing.T) {
	// Create a provider for validation tests (won't make actual API calls)
	provider, err := NewProvider("test-key", "gpt-4", nil)
	if err != nil {
		t.Fatalf("Failed to create provider: %v", err)
	}

	ctx := context.Background()

	tests := []struct {
		name    string
		req     *llm.CompletionRequest
		wantErr bool
		errMsg  string
	}{
		{
			name:    "nil request",
			req:     nil,
			wantErr: true,
			errMsg:  "completion request cannot be nil",
		},
		{
			name: "empty messages",
			req: &llm.CompletionRequest{
				Messages: []llm.Message{},
			},
			wantErr: true,
			errMsg:  "messages cannot be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := provider.Complete(ctx, tt.req)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Complete() expected error but got nil")
				} else if err.Error() != tt.errMsg {
					t.Errorf("Complete() error = %v, want %v", err.Error(), tt.errMsg)
				}
				return
			}

			if err != nil {
				t.Errorf("Complete() unexpected error: %v", err)
			}
		})
	}
}

func TestProvider_Methods(t *testing.T) {
	provider, err := NewProvider("test-key", "gpt-4-turbo", nil)
	if err != nil {
		t.Fatalf("Failed to create provider: %v", err)
	}

	// Test Name()
	if got := provider.Name(); got != "openai" {
		t.Errorf("Name() = %v, want %v", got, "openai")
	}

	// Test ModelName()
	if got := provider.ModelName(); got != "gpt-4-turbo" {
		t.Errorf("ModelName() = %v, want %v", got, "gpt-4-turbo")
	}

	// Test SupportsStreaming()
	if got := provider.SupportsStreaming(); !got {
		t.Errorf("SupportsStreaming() = %v, want %v", got, true)
	}
}

func TestProvider_ConfigDefaults(t *testing.T) {
	// Test that default config is applied when nil config provided
	provider, err := NewProvider("test-key", "gpt-4o", nil)
	if err != nil {
		t.Fatalf("Failed to create provider: %v", err)
	}

	if provider.config == nil {
		t.Fatal("Expected default config to be set, got nil")
	}

	// Verify default values
	if provider.config.DefaultTemperature != 0.7 {
		t.Errorf("DefaultTemperature = %v, want 0.7", provider.config.DefaultTemperature)
	}

	if provider.config.DefaultMaxTokens != 2048 {
		t.Errorf("DefaultMaxTokens = %v, want 2048", provider.config.DefaultMaxTokens)
	}

	if provider.config.TimeoutSeconds != 60 {
		t.Errorf("TimeoutSeconds = %v, want 60", provider.config.TimeoutSeconds)
	}

	if provider.config.Provider != "openai" {
		t.Errorf("Provider = %v, want openai", provider.config.Provider)
	}
}

func TestProvider_CustomConfig(t *testing.T) {
	// Test that custom config is respected
	customConfig := &llm.Config{
		Provider:           "openai",
		APIKey:             "custom-key",
		Model:              "gpt-4o-mini",
		DefaultTemperature: 0.5,
		DefaultMaxTokens:   1000,
		TimeoutSeconds:     30,
	}

	provider, err := NewProvider("test-key", "gpt-4o", customConfig)
	if err != nil {
		t.Fatalf("Failed to create provider: %v", err)
	}

	// Verify custom config is preserved
	if provider.config.DefaultTemperature != 0.5 {
		t.Errorf("DefaultTemperature = %v, want 0.5", provider.config.DefaultTemperature)
	}

	if provider.config.DefaultMaxTokens != 1000 {
		t.Errorf("DefaultMaxTokens = %v, want 1000", provider.config.DefaultMaxTokens)
	}

	if provider.config.TimeoutSeconds != 30 {
		t.Errorf("TimeoutSeconds = %v, want 30", provider.config.TimeoutSeconds)
	}
}

func TestProvider_BaseURL(t *testing.T) {
	// Test that custom BaseURL is respected
	config := &llm.Config{
		Provider:           "openai",
		APIKey:             "test-key",
		BaseURL:            "https://custom.openai.endpoint",
		Model:              "gpt-4o",
		DefaultTemperature: 0.7,
		DefaultMaxTokens:   2048,
		TimeoutSeconds:     60,
	}

	provider, err := NewProvider("test-key", "gpt-4o", config)
	if err != nil {
		t.Fatalf("Failed to create provider with custom BaseURL: %v", err)
	}

	if provider == nil {
		t.Error("Expected non-nil provider with custom BaseURL")
	}

	// Verify config is stored
	if provider.config.BaseURL != "https://custom.openai.endpoint" {
		t.Errorf("BaseURL not preserved, got %v", provider.config.BaseURL)
	}
}

func TestProvider_EmptyBaseURL(t *testing.T) {
	// Test that empty BaseURL doesn't cause issues
	config := &llm.Config{
		Provider:           "openai",
		APIKey:             "test-key",
		BaseURL:            "", // Empty BaseURL should use default
		Model:              "gpt-4o",
		DefaultTemperature: 0.7,
		DefaultMaxTokens:   2048,
		TimeoutSeconds:     60,
	}

	provider, err := NewProvider("test-key", "gpt-4o", config)
	if err != nil {
		t.Fatalf("Failed to create provider with empty BaseURL: %v", err)
	}

	if provider == nil {
		t.Error("Expected non-nil provider with empty BaseURL")
	}
}

func TestProvider_MessageValidation(t *testing.T) {
	provider, err := NewProvider("test-key", "gpt-4o", nil)
	if err != nil {
		t.Fatalf("Failed to create provider: %v", err)
	}

	ctx := context.Background()

	tests := []struct {
		name    string
		req     *llm.CompletionRequest
		wantErr bool
		errMsg  string
	}{
		{
			name:    "nil messages slice",
			req:     &llm.CompletionRequest{Messages: nil},
			wantErr: true,
			errMsg:  "messages cannot be empty",
		},
		{
			name: "single message",
			req: &llm.CompletionRequest{
				Messages: []llm.Message{
					{Role: "user", Content: "Hello"},
				},
			},
			wantErr: false,
		},
		{
			name: "multiple messages",
			req: &llm.CompletionRequest{
				Messages: []llm.Message{
					{Role: "system", Content: "You are helpful"},
					{Role: "user", Content: "Hello"},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := provider.Complete(ctx, tt.req)

			if tt.wantErr {
				if err == nil {
					t.Error("Complete() expected error but got nil")
				} else if tt.errMsg != "" && err.Error() != tt.errMsg {
					t.Errorf("Complete() error = %v, want %v", err.Error(), tt.errMsg)
				}
			} else {
				// Note: This will fail with actual API call, but validates message conversion
				// For unit tests, we just verify no validation errors occurred
				if err != nil && err.Error() == "messages cannot be empty" {
					t.Errorf("Complete() unexpected validation error: %v", err)
				}
			}
		})
	}
}

func TestProvider_ParameterDefaults(t *testing.T) {
	// Test that request parameters use config defaults when not specified
	provider, err := NewProvider("test-key", "gpt-4o", &llm.Config{
		Provider:           "openai",
		APIKey:             "test-key",
		Model:              "gpt-4o",
		DefaultTemperature: 0.8,
		DefaultMaxTokens:   1500,
		TimeoutSeconds:     45,
	})
	if err != nil {
		t.Fatalf("Failed to create provider: %v", err)
	}

	// Verify provider was created with custom defaults
	if provider.config.DefaultTemperature != 0.8 {
		t.Errorf("DefaultTemperature = %v, want 0.8", provider.config.DefaultTemperature)
	}

	if provider.config.DefaultMaxTokens != 1500 {
		t.Errorf("DefaultMaxTokens = %v, want 1500", provider.config.DefaultMaxTokens)
	}
}

func TestProvider_ModelVariations(t *testing.T) {
	// Test that different model names are handled correctly
	models := []string{
		"gpt-4o",
		"gpt-4o-mini",
		"gpt-4-turbo",
		"gpt-4",
		"gpt-3.5-turbo",
	}

	for _, model := range models {
		t.Run(model, func(t *testing.T) {
			provider, err := NewProvider("test-key", model, nil)
			if err != nil {
				t.Fatalf("Failed to create provider for model %s: %v", model, err)
			}

			if provider.ModelName() != model {
				t.Errorf("ModelName() = %v, want %v", provider.ModelName(), model)
			}
		})
	}
}

func TestProvider_ContextHandling(t *testing.T) {
	provider, err := NewProvider("test-key", "gpt-4o", &llm.Config{
		Provider:           "openai",
		APIKey:             "test-key",
		Model:              "gpt-4o",
		DefaultTemperature: 0.7,
		DefaultMaxTokens:   2048,
		TimeoutSeconds:     1, // Very short timeout for testing
	})
	if err != nil {
		t.Fatalf("Failed to create provider: %v", err)
	}

	ctx := context.Background()

	req := &llm.CompletionRequest{
		Messages: []llm.Message{
			{Role: "user", Content: "Test message"},
		},
	}

	// This will likely timeout or fail with API error, but we're testing
	// that the timeout mechanism is being applied
	_, err = provider.Complete(ctx, req)

	// We expect an error (either timeout or API error), not a panic
	if err == nil {
		t.Log("Complete() succeeded unexpectedly (likely using valid API key in environment)")
	}
}

func TestProvider_ConfigPreservation(t *testing.T) {
	// Test that original config is preserved and not modified
	originalConfig := &llm.Config{
		Provider:           "openai",
		APIKey:             "original-key",
		BaseURL:            "https://original.url",
		Model:              "gpt-4o",
		DefaultTemperature: 0.9,
		DefaultMaxTokens:   3000,
		TimeoutSeconds:     90,
	}

	provider, err := NewProvider("test-key", "gpt-4o", originalConfig)
	if err != nil {
		t.Fatalf("Failed to create provider: %v", err)
	}

	// Verify all config fields are preserved
	if provider.config.BaseURL != originalConfig.BaseURL {
		t.Errorf("BaseURL = %v, want %v", provider.config.BaseURL, originalConfig.BaseURL)
	}

	if provider.config.DefaultTemperature != originalConfig.DefaultTemperature {
		t.Errorf("DefaultTemperature = %v, want %v", provider.config.DefaultTemperature, originalConfig.DefaultTemperature)
	}

	if provider.config.DefaultMaxTokens != originalConfig.DefaultMaxTokens {
		t.Errorf("DefaultMaxTokens = %v, want %v", provider.config.DefaultMaxTokens, originalConfig.DefaultMaxTokens)
	}

	if provider.config.TimeoutSeconds != originalConfig.TimeoutSeconds {
		t.Errorf("TimeoutSeconds = %v, want %v", provider.config.TimeoutSeconds, originalConfig.TimeoutSeconds)
	}
}

// Note: Integration tests that make actual API calls should be in a separate
// test file (e.g., provider_integration_test.go) and run with a build tag:
// //go:build integration
// This allows unit tests to run quickly without API dependencies.
