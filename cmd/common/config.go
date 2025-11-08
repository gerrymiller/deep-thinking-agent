// Copyright 2025 Gerry Miller <gerry@gerrymiller.com>
//
// Licensed under the MIT License.
// See LICENSE file in the project root for full license information.

package common

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// Config represents the complete application configuration.
type Config struct {
	LLM         LLMConfig         `json:"llm"`
	Embedding   EmbeddingConfig   `json:"embedding"`
	VectorStore VectorStoreConfig `json:"vector_store"`
	Workflow    WorkflowConfig    `json:"workflow"`
}

// LLMConfig contains configuration for LLM providers.
type LLMConfig struct {
	ReasoningLLM LLMProviderConfig `json:"reasoning_llm"`
	FastLLM      LLMProviderConfig `json:"fast_llm"`
}

// LLMProviderConfig contains configuration for a specific LLM provider.
type LLMProviderConfig struct {
	Provider           string  `json:"provider"`
	Model              string  `json:"model"`
	APIKey             string  `json:"api_key,omitempty"`
	DefaultTemperature float32 `json:"default_temperature"`
}

// EmbeddingConfig contains configuration for embedding generation.
type EmbeddingConfig struct {
	Provider string `json:"provider"`
	Model    string `json:"model"`
	APIKey   string `json:"api_key,omitempty"`
}

// VectorStoreConfig contains configuration for the vector database.
type VectorStoreConfig struct {
	Type              string `json:"type"`
	Address           string `json:"address"`
	DefaultCollection string `json:"default_collection"`
}

// WorkflowConfig contains configuration for workflow execution.
type WorkflowConfig struct {
	MaxIterations   int    `json:"max_iterations"`
	TopKRetrieval   int    `json:"top_k_retrieval"`
	TopNReranking   int    `json:"top_n_reranking"`
	DefaultStrategy string `json:"default_strategy"`
}

// LoadConfig loads configuration from a JSON file.
func LoadConfig(path string) (*Config, error) {
	loadEnvFiles()

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	// Load API keys from environment if not in config
	if config.LLM.ReasoningLLM.APIKey == "" {
		config.LLM.ReasoningLLM.APIKey = os.Getenv("OPENAI_API_KEY")
	}
	if config.LLM.FastLLM.APIKey == "" {
		config.LLM.FastLLM.APIKey = os.Getenv("OPENAI_API_KEY")
	}
	if config.Embedding.APIKey == "" {
		config.Embedding.APIKey = os.Getenv("OPENAI_API_KEY")
	}

	return &config, nil
}

// DefaultConfig returns a default configuration suitable for initial setup.
func DefaultConfig() *Config {
	return &Config{
		LLM: LLMConfig{
			ReasoningLLM: LLMProviderConfig{
				Provider:           "openai",
				Model:              "gpt-4o", // Fast and capable model
				DefaultTemperature: 0.7,
			},
			FastLLM: LLMProviderConfig{
				Provider:           "openai",
				Model:              "gpt-4o-mini", // Fast model for simple tasks
				DefaultTemperature: 0.5,
			},
		},
		Embedding: EmbeddingConfig{
			Provider: "openai",
			Model:    "text-embedding-3-small",
		},
		VectorStore: VectorStoreConfig{
			Type:              "qdrant",
			Address:           "localhost:6334",
			DefaultCollection: "documents",
		},
		Workflow: WorkflowConfig{
			MaxIterations:   10,
			TopKRetrieval:   10,
			TopNReranking:   3,
			DefaultStrategy: "hybrid",
		},
	}
}

func loadEnvFiles() {
	envFiles := []string{".env", ".env.local"}
	merged := make(map[string]string)

	for _, file := range envFiles {
		envMap, err := godotenv.Read(file)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				continue
			}
			// Ignore other read errors to avoid blocking config loading.
			continue
		}
		for key, value := range envMap {
			merged[key] = value
		}
	}

	for key, value := range merged {
		current, exists := os.LookupEnv(key)
		if !exists || current == "" {
			_ = os.Setenv(key, value)
		}
	}
}
