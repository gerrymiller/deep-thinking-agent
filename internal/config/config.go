// Copyright 2025 Gerry Miller <gerry@gerrymiller.com>
//
// Licensed under the MIT License.
// See LICENSE file in the project root for full license information.

package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/joho/godotenv"

	"deep-thinking-agent/pkg/embedding"
	"deep-thinking-agent/pkg/llm"
	"deep-thinking-agent/pkg/vectorstore"
)

// Config represents the complete configuration for the deep thinking agent.
type Config struct {
	// LLM configuration
	LLM LLMConfig `json:"llm"`

	// Embedding configuration
	Embedding EmbeddingConfig `json:"embedding"`

	// VectorStore configuration
	VectorStore VectorStoreConfig `json:"vector_store"`

	// WebSearch configuration (optional)
	WebSearch *WebSearchConfig `json:"web_search,omitempty"`

	// Workflow configuration
	Workflow WorkflowConfig `json:"workflow"`

	// Schema configuration
	Schema SchemaConfig `json:"schema"`
}

// LLMConfig contains settings for LLM providers.
type LLMConfig struct {
	// ReasoningLLM is used for complex reasoning tasks (planning, reflection, policy)
	ReasoningLLM LLMProviderConfig `json:"reasoning_llm"`

	// FastLLM is used for quick tasks (query rewriting, summarization)
	FastLLM LLMProviderConfig `json:"fast_llm"`
}

// LLMProviderConfig contains settings for a specific LLM provider.
type LLMProviderConfig struct {
	Provider           string  `json:"provider"` // "openai", "anthropic", "ollama"
	APIKey             string  `json:"api_key,omitempty"`
	BaseURL            string  `json:"base_url,omitempty"`
	Model              string  `json:"model"`
	DefaultTemperature float32 `json:"default_temperature"`
	DefaultMaxTokens   int     `json:"default_max_tokens"`
	TimeoutSeconds     int     `json:"timeout_seconds"`
}

// EmbeddingConfig contains settings for embedding generation.
type EmbeddingConfig struct {
	Provider       string `json:"provider"` // "openai", "local"
	APIKey         string `json:"api_key,omitempty"`
	BaseURL        string `json:"base_url,omitempty"`
	Model          string `json:"model"`
	BatchSize      int    `json:"batch_size"`
	TimeoutSeconds int    `json:"timeout_seconds"`
}

// VectorStoreConfig contains settings for the vector store.
type VectorStoreConfig struct {
	Type              string                 `json:"type"` // "qdrant", "weaviate", "milvus"
	Address           string                 `json:"address"`
	APIKey            string                 `json:"api_key,omitempty"`
	TimeoutSeconds    int                    `json:"timeout_seconds"`
	DefaultCollection string                 `json:"default_collection"`
	Extra             map[string]interface{} `json:"extra,omitempty"`
}

// WebSearchConfig contains settings for web search (optional feature).
type WebSearchConfig struct {
	Enabled    bool   `json:"enabled"`
	Provider   string `json:"provider"` // "serper", "brave", etc.
	APIKey     string `json:"api_key,omitempty"`
	MaxResults int    `json:"max_results"`
}

// WorkflowConfig contains settings for the workflow execution.
type WorkflowConfig struct {
	MaxIterations     int     `json:"max_iterations"`
	TopKRetrieval     int     `json:"top_k_retrieval"`
	TopNReranking     int     `json:"top_n_reranking"`
	MinRelevanceScore float32 `json:"min_relevance_score"`
	DefaultStrategy   string  `json:"default_strategy"` // "vector", "keyword", "hybrid"
}

// SchemaConfig contains settings for schema analysis.
type SchemaConfig struct {
	EnableLLMAnalysis     bool     `json:"enable_llm_analysis"`
	EnablePatternMatching bool     `json:"enable_pattern_matching"`
	PredefinedSchemas     []string `json:"predefined_schemas"` // Paths to schema files
	CacheSchemas          bool     `json:"cache_schemas"`
	AnalysisTimeout       int      `json:"analysis_timeout_seconds"`
}

// LoadFromFile loads configuration from a JSON file.
func LoadFromFile(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Apply defaults
	applyDefaults(&config)

	return &config, nil
}

// LoadFromEnv loads configuration from environment variables.
// This is useful for containerized deployments.
func LoadFromEnv() *Config {
	loadEnvFiles()

	config := &Config{
		LLM: LLMConfig{
			ReasoningLLM: LLMProviderConfig{
				Provider:           getEnv("REASONING_LLM_PROVIDER", "openai"),
				APIKey:             getEnv("REASONING_LLM_API_KEY", ""),
				Model:              getEnv("REASONING_LLM_MODEL", "gpt-4o"), // Latest: gpt-4o (May 2024), supports vision. Alternative: gpt-5 when available
				DefaultTemperature: 0.7,
				DefaultMaxTokens:   2048,
				TimeoutSeconds:     60,
			},
			FastLLM: LLMProviderConfig{
				Provider:           getEnv("FAST_LLM_PROVIDER", "openai"),
				APIKey:             getEnv("FAST_LLM_API_KEY", ""),
				Model:              getEnv("FAST_LLM_MODEL", "gpt-4o-mini"), // Replaces deprecated gpt-3.5-turbo (July 2024)
				DefaultTemperature: 0.5,
				DefaultMaxTokens:   1024,
				TimeoutSeconds:     30,
			},
		},
		Embedding: EmbeddingConfig{
			Provider:       getEnv("EMBEDDING_PROVIDER", "openai"),
			APIKey:         getEnv("EMBEDDING_API_KEY", ""),
			Model:          getEnv("EMBEDDING_MODEL", "text-embedding-3-small"),
			BatchSize:      100,
			TimeoutSeconds: 30,
		},
		VectorStore: VectorStoreConfig{
			Type:              getEnv("VECTOR_STORE_TYPE", "qdrant"),
			Address:           getEnv("VECTOR_STORE_ADDRESS", "localhost:6334"),
			DefaultCollection: getEnv("VECTOR_STORE_COLLECTION", "documents"),
			TimeoutSeconds:    30,
		},
		Workflow: WorkflowConfig{
			MaxIterations:     10,
			TopKRetrieval:     10,
			TopNReranking:     3,
			MinRelevanceScore: 0.0,
			DefaultStrategy:   "hybrid",
		},
		Schema: SchemaConfig{
			EnableLLMAnalysis:     true,
			EnablePatternMatching: true,
			CacheSchemas:          true,
			AnalysisTimeout:       60,
		},
	}

	return config
}

// SaveToFile saves the configuration to a JSON file.
func (c *Config) SaveToFile(path string) error {
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// ToLLMConfig converts to llm.Config for the reasoning LLM.
func (c *Config) ToLLMConfig() *llm.Config {
	return &llm.Config{
		Provider:           c.LLM.ReasoningLLM.Provider,
		APIKey:             c.LLM.ReasoningLLM.APIKey,
		BaseURL:            c.LLM.ReasoningLLM.BaseURL,
		Model:              c.LLM.ReasoningLLM.Model,
		DefaultTemperature: c.LLM.ReasoningLLM.DefaultTemperature,
		DefaultMaxTokens:   c.LLM.ReasoningLLM.DefaultMaxTokens,
		TimeoutSeconds:     c.LLM.ReasoningLLM.TimeoutSeconds,
	}
}

// ToFastLLMConfig converts to llm.Config for the fast LLM.
func (c *Config) ToFastLLMConfig() *llm.Config {
	return &llm.Config{
		Provider:           c.LLM.FastLLM.Provider,
		APIKey:             c.LLM.FastLLM.APIKey,
		BaseURL:            c.LLM.FastLLM.BaseURL,
		Model:              c.LLM.FastLLM.Model,
		DefaultTemperature: c.LLM.FastLLM.DefaultTemperature,
		DefaultMaxTokens:   c.LLM.FastLLM.DefaultMaxTokens,
		TimeoutSeconds:     c.LLM.FastLLM.TimeoutSeconds,
	}
}

// ToEmbeddingConfig converts to embedding.Config.
func (c *Config) ToEmbeddingConfig() *embedding.Config {
	return &embedding.Config{
		Provider:       c.Embedding.Provider,
		APIKey:         c.Embedding.APIKey,
		BaseURL:        c.Embedding.BaseURL,
		Model:          c.Embedding.Model,
		BatchSize:      c.Embedding.BatchSize,
		TimeoutSeconds: c.Embedding.TimeoutSeconds,
	}
}

// ToVectorStoreConfig converts to vectorstore.Config.
func (c *Config) ToVectorStoreConfig() *vectorstore.Config {
	return &vectorstore.Config{
		Type:              c.VectorStore.Type,
		Address:           c.VectorStore.Address,
		APIKey:            c.VectorStore.APIKey,
		TimeoutSeconds:    c.VectorStore.TimeoutSeconds,
		DefaultCollection: c.VectorStore.DefaultCollection,
		Extra:             c.VectorStore.Extra,
	}
}

// applyDefaults fills in default values for unspecified config fields.
func applyDefaults(config *Config) {
	// LLM defaults
	if config.LLM.ReasoningLLM.DefaultTemperature == 0 {
		config.LLM.ReasoningLLM.DefaultTemperature = 0.7
	}
	if config.LLM.ReasoningLLM.DefaultMaxTokens == 0 {
		config.LLM.ReasoningLLM.DefaultMaxTokens = 2048
	}
	if config.LLM.ReasoningLLM.TimeoutSeconds == 0 {
		config.LLM.ReasoningLLM.TimeoutSeconds = 60
	}

	if config.LLM.FastLLM.DefaultTemperature == 0 {
		config.LLM.FastLLM.DefaultTemperature = 0.5
	}
	if config.LLM.FastLLM.DefaultMaxTokens == 0 {
		config.LLM.FastLLM.DefaultMaxTokens = 1024
	}
	if config.LLM.FastLLM.TimeoutSeconds == 0 {
		config.LLM.FastLLM.TimeoutSeconds = 30
	}

	// Embedding defaults
	if config.Embedding.BatchSize == 0 {
		config.Embedding.BatchSize = 100
	}
	if config.Embedding.TimeoutSeconds == 0 {
		config.Embedding.TimeoutSeconds = 30
	}

	// VectorStore defaults
	if config.VectorStore.TimeoutSeconds == 0 {
		config.VectorStore.TimeoutSeconds = 30
	}
	if config.VectorStore.DefaultCollection == "" {
		config.VectorStore.DefaultCollection = "documents"
	}

	// Workflow defaults
	if config.Workflow.MaxIterations == 0 {
		config.Workflow.MaxIterations = 10
	}
	if config.Workflow.TopKRetrieval == 0 {
		config.Workflow.TopKRetrieval = 10
	}
	if config.Workflow.TopNReranking == 0 {
		config.Workflow.TopNReranking = 3
	}
	if config.Workflow.DefaultStrategy == "" {
		config.Workflow.DefaultStrategy = "hybrid"
	}

	// Schema defaults
	if config.Schema.AnalysisTimeout == 0 {
		config.Schema.AnalysisTimeout = 60
	}
}

// getEnv retrieves an environment variable or returns a default value.
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
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
