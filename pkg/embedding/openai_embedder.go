// Copyright 2025 Gerry Miller <gerry@gerrymiller.com>
//
// Licensed under the MIT License.
// See LICENSE file in the project root for full license information.

package embedding

import (
	"context"
	"errors"
	"fmt"
	"time"

	openai "github.com/sashabaranov/go-openai"
)

// OpenAIEmbedder implements the Embedder interface using OpenAI's embedding models.
type OpenAIEmbedder struct {
	client     *openai.Client
	model      string
	dimensions int
	config     *Config
}

// Model dimensions for common OpenAI embedding models
const (
	DimensionsTextEmbedding3Small = 1536
	DimensionsTextEmbedding3Large = 3072
	DimensionsTextEmbeddingAda002 = 1536
)

// NewOpenAIEmbedder creates a new OpenAI embedder instance.
// apiKey: OpenAI API key for authentication
// model: Embedding model to use (e.g., "text-embedding-3-small", "text-embedding-ada-002")
// config: Optional configuration (can be nil for defaults)
func NewOpenAIEmbedder(apiKey, model string, config *Config) (*OpenAIEmbedder, error) {
	if apiKey == "" {
		return nil, errors.New("OpenAI API key is required")
	}
	if model == "" {
		return nil, errors.New("embedding model name is required")
	}

	// Apply default config if not provided
	if config == nil {
		config = &Config{
			Provider:       "openai",
			APIKey:         apiKey,
			Model:          model,
			BatchSize:      100,
			TimeoutSeconds: 30,
		}
	}

	// Determine dimensions based on model
	dimensions := getDimensionsForModel(model)

	// Create OpenAI client configuration
	clientConfig := openai.DefaultConfig(apiKey)
	if config.BaseURL != "" {
		clientConfig.BaseURL = config.BaseURL
	}

	client := openai.NewClientWithConfig(clientConfig)

	return &OpenAIEmbedder{
		client:     client,
		model:      model,
		dimensions: dimensions,
		config:     config,
	}, nil
}

// getDimensionsForModel returns the embedding dimensions for a given model.
func getDimensionsForModel(model string) int {
	switch model {
	case "text-embedding-3-small":
		return DimensionsTextEmbedding3Small
	case "text-embedding-3-large":
		return DimensionsTextEmbedding3Large
	case "text-embedding-ada-002":
		return DimensionsTextEmbeddingAda002
	default:
		// Default to ada-002 dimensions for unknown models
		return DimensionsTextEmbeddingAda002
	}
}

// Embed generates embeddings for the given texts.
func (e *OpenAIEmbedder) Embed(ctx context.Context, req *EmbedRequest) (*EmbedResponse, error) {
	if req == nil {
		return nil, errors.New("embed request cannot be nil")
	}
	if len(req.Texts) == 0 {
		return nil, errors.New("texts cannot be empty")
	}

	// Apply timeout
	if e.config.TimeoutSeconds > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, time.Duration(e.config.TimeoutSeconds)*time.Second)
		defer cancel()
	}

	// Process texts in batches to avoid API limits
	batchSize := e.config.BatchSize
	if batchSize <= 0 {
		batchSize = 100
	}

	var allVectors []Vector
	totalPromptTokens := 0
	totalTokens := 0

	for i := 0; i < len(req.Texts); i += batchSize {
		end := i + batchSize
		if end > len(req.Texts) {
			end = len(req.Texts)
		}

		batch := req.Texts[i:end]

		// Create OpenAI embedding request
		openaiReq := openai.EmbeddingRequest{
			Input: batch,
			Model: openai.EmbeddingModel(e.model),
		}

		// Execute request
		resp, err := e.client.CreateEmbeddings(ctx, openaiReq)
		if err != nil {
			return nil, fmt.Errorf("OpenAI embedding API error: %w", err)
		}

		// Convert response to our format
		for j, data := range resp.Data {
			vector := Vector{
				Embedding: data.Embedding,
				Text:      batch[j],
				Metadata:  make(map[string]interface{}),
			}

			// Copy metadata from request if provided
			if req.Metadata != nil {
				for k, v := range req.Metadata {
					vector.Metadata[k] = v
				}
			}

			allVectors = append(allVectors, vector)
		}

		// Accumulate usage stats
		totalPromptTokens += resp.Usage.PromptTokens
		totalTokens += resp.Usage.TotalTokens
	}

	return &EmbedResponse{
		Vectors: allVectors,
		Usage: UsageStats{
			PromptTokens: totalPromptTokens,
			TotalTokens:  totalTokens,
		},
		Model: e.model,
	}, nil
}

// Dimensions returns the dimensionality of the embeddings produced by this embedder.
func (e *OpenAIEmbedder) Dimensions() int {
	return e.dimensions
}

// ModelName returns the name of the embedding model being used.
func (e *OpenAIEmbedder) ModelName() string {
	return e.model
}
