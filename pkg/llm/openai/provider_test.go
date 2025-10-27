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

// Note: Integration tests that make actual API calls should be in a separate
// test file (e.g., provider_integration_test.go) and run with a build tag:
// //go:build integration
// This allows unit tests to run quickly without API dependencies.
