// Copyright 2025 Gerry Miller <gerry@gerrymiller.com>
//
// Licensed under the MIT License.
// See LICENSE file in the project root for full license information.

package embedding

import "context"

// Vector represents an embedding vector with its associated metadata.
type Vector struct {
	// Embedding is the dense vector representation
	Embedding []float32

	// Text is the original text that was embedded
	Text string

	// Metadata contains additional information about this embedding
	Metadata map[string]interface{}
}

// EmbedRequest contains parameters for an embedding request.
type EmbedRequest struct {
	// Texts are the strings to embed
	Texts []string

	// Model specifies which embedding model to use
	Model string

	// Metadata will be attached to each resulting vector
	Metadata map[string]interface{}
}

// EmbedResponse contains the results of an embedding request.
type EmbedResponse struct {
	// Vectors are the generated embeddings with metadata
	Vectors []Vector

	// Usage contains token/request usage statistics
	Usage UsageStats

	// Model is the actual model used
	Model string
}

// UsageStats tracks usage for embedding requests.
type UsageStats struct {
	// PromptTokens is the number of tokens in the input
	PromptTokens int

	// TotalTokens includes any additional tokens used
	TotalTokens int
}

// Embedder defines the interface for generating embeddings from text.
// This abstraction allows using different embedding models (OpenAI, local models, etc.)
type Embedder interface {
	// Embed generates embeddings for the given texts.
	// Returns vectors with the same ordering as input texts.
	Embed(ctx context.Context, req *EmbedRequest) (*EmbedResponse, error)

	// Dimensions returns the dimensionality of the embeddings produced by this embedder.
	Dimensions() int

	// ModelName returns the name of the embedding model being used.
	ModelName() string
}

// Config contains configuration for embedding generation.
type Config struct {
	// Provider specifies which embedding provider to use (e.g., "openai", "local")
	Provider string

	// APIKey for authentication (if required)
	APIKey string

	// BaseURL allows overriding the default API endpoint
	BaseURL string

	// Model specifies which embedding model to use
	Model string

	// BatchSize controls how many texts to embed in a single request
	BatchSize int

	// TimeoutSeconds for API requests
	TimeoutSeconds int
}
