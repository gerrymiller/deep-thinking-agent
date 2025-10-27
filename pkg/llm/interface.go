package llm

import "context"

// Message represents a single message in a conversation between user and assistant.
// Role can be "system", "user", or "assistant".
type Message struct {
	Role    string // "system", "user", or "assistant"
	Content string
}

// CompletionRequest contains all parameters needed for an LLM completion request.
type CompletionRequest struct {
	// Messages is the conversation history including system prompts
	Messages []Message

	// Temperature controls randomness (0.0 = deterministic, 1.0 = creative)
	Temperature float32

	// MaxTokens is the maximum number of tokens to generate
	MaxTokens int

	// TopP controls nucleus sampling (0.0-1.0)
	TopP float32

	// Stop sequences that will halt generation
	StopSequences []string

	// Stream enables streaming responses (not implemented in Phase 1)
	Stream bool
}

// CompletionResponse contains the LLM's response to a completion request.
type CompletionResponse struct {
	// Content is the generated text
	Content string

	// FinishReason indicates why generation stopped ("stop", "length", "error")
	FinishReason string

	// Usage contains token usage statistics
	Usage UsageStats

	// Model is the actual model used (may differ from requested model)
	Model string
}

// UsageStats tracks token usage for a completion request.
type UsageStats struct {
	PromptTokens     int
	CompletionTokens int
	TotalTokens      int
}

// Provider defines the interface that all LLM providers must implement.
// This abstraction allows swapping between OpenAI, Anthropic, Ollama, etc.
type Provider interface {
	// Complete generates a completion for the given request.
	// Returns the response or an error if the request fails.
	Complete(ctx context.Context, req *CompletionRequest) (*CompletionResponse, error)

	// Name returns the provider name (e.g., "openai", "anthropic", "ollama")
	Name() string

	// ModelName returns the specific model being used
	ModelName() string

	// SupportsStreaming indicates if this provider supports streaming responses
	SupportsStreaming() bool
}

// Config contains common configuration options for LLM providers.
type Config struct {
	// Provider specifies which LLM provider to use
	Provider string

	// APIKey for authentication (if required)
	APIKey string

	// BaseURL allows overriding the default API endpoint (useful for proxies/local deployments)
	BaseURL string

	// Model specifies which model to use (e.g., "gpt-4", "claude-3-sonnet")
	Model string

	// DefaultTemperature is used when requests don't specify temperature
	DefaultTemperature float32

	// DefaultMaxTokens is used when requests don't specify max tokens
	DefaultMaxTokens int

	// Timeout in seconds for API requests
	TimeoutSeconds int
}
