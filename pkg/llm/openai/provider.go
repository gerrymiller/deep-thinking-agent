// Copyright 2025 Gerry Miller <gerry@gerrymiller.com>
//
// Licensed under the MIT License.
// See LICENSE file in the project root for full license information.

package openai

import (
	"context"
	"errors"
	"fmt"
	"time"

	"deep-thinking-agent/pkg/llm"

	openai "github.com/sashabaranov/go-openai"
)

// Provider implements the llm.Provider interface for OpenAI's API.
type Provider struct {
	client *openai.Client
	model  string
	config *llm.Config
}

// NewProvider creates a new OpenAI provider instance.
// apiKey: OpenAI API key for authentication
// model: Model to use (e.g., "gpt-4", "gpt-4-turbo", "gpt-3.5-turbo")
// config: Optional configuration (can be nil for defaults)
func NewProvider(apiKey, model string, config *llm.Config) (*Provider, error) {
	if apiKey == "" {
		return nil, errors.New("OpenAI API key is required")
	}
	if model == "" {
		return nil, errors.New("model name is required")
	}

	// Apply default config if not provided
	if config == nil {
		config = &llm.Config{
			Provider:           "openai",
			APIKey:             apiKey,
			Model:              model,
			DefaultTemperature: 0.7,
			DefaultMaxTokens:   2048,
			TimeoutSeconds:     60,
		}
	}

	// Create OpenAI client configuration
	clientConfig := openai.DefaultConfig(apiKey)
	if config.BaseURL != "" {
		clientConfig.BaseURL = config.BaseURL
	}

	client := openai.NewClientWithConfig(clientConfig)

	return &Provider{
		client: client,
		model:  model,
		config: config,
	}, nil
}

// Complete generates a completion for the given request.
func (p *Provider) Complete(ctx context.Context, req *llm.CompletionRequest) (*llm.CompletionResponse, error) {
	if req == nil {
		return nil, errors.New("completion request cannot be nil")
	}
	if len(req.Messages) == 0 {
		return nil, errors.New("messages cannot be empty")
	}

	// Apply timeout
	if p.config.TimeoutSeconds > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, time.Duration(p.config.TimeoutSeconds)*time.Second)
		defer cancel()
	}

	// Convert our messages to OpenAI format
	openaiMessages := make([]openai.ChatCompletionMessage, len(req.Messages))
	for i, msg := range req.Messages {
		openaiMessages[i] = openai.ChatCompletionMessage{
			Role:    msg.Role,
			Content: msg.Content,
		}
	}

	// Apply defaults for unspecified parameters
	temperature := req.Temperature
	if temperature == 0 {
		temperature = p.config.DefaultTemperature
	}

	maxTokens := req.MaxTokens
	if maxTokens == 0 {
		maxTokens = p.config.DefaultMaxTokens
	}

	topP := req.TopP
	if topP == 0 {
		topP = 1.0 // OpenAI default
	}

	// Create OpenAI request
	openaiReq := openai.ChatCompletionRequest{
		Model:               p.model,
		Messages:            openaiMessages,
		Temperature:         temperature,
		MaxCompletionTokens: maxTokens,
		TopP:                topP,
		Stop:                req.StopSequences,
	}

	// Execute request
	resp, err := p.client.CreateChatCompletion(ctx, openaiReq)
	if err != nil {
		return nil, fmt.Errorf("OpenAI API error: %w", err)
	}

	// Validate response
	if len(resp.Choices) == 0 {
		return nil, errors.New("OpenAI returned no choices")
	}

	// Convert response to our format
	return &llm.CompletionResponse{
		Content:      resp.Choices[0].Message.Content,
		FinishReason: string(resp.Choices[0].FinishReason),
		Usage: llm.UsageStats{
			PromptTokens:     resp.Usage.PromptTokens,
			CompletionTokens: resp.Usage.CompletionTokens,
			TotalTokens:      resp.Usage.TotalTokens,
		},
		Model: resp.Model,
	}, nil
}

// Name returns the provider name.
func (p *Provider) Name() string {
	return "openai"
}

// ModelName returns the specific model being used.
func (p *Provider) ModelName() string {
	return p.model
}

// SupportsStreaming indicates if this provider supports streaming responses.
func (p *Provider) SupportsStreaming() bool {
	return true // OpenAI supports streaming, but not implemented in Phase 1
}
