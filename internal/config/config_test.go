// Copyright 2025 Gerry Miller <gerry@gerrymiller.com>
//
// Licensed under the MIT License.
// See LICENSE file in the project root for full license information.

package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

// TestLoadFromFile tests loading configuration from a JSON file.
func TestLoadFromFile(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		wantErr  bool
		validate func(*testing.T, *Config)
	}{
		{
			name: "valid minimal config",
			content: `{
				"llm": {
					"reasoning_llm": {
						"provider": "openai",
						"model": "gpt-5"
					},
					"fast_llm": {
						"provider": "openai",
						"model": "gpt-5-mini"
					}
				},
				"embedding": {
					"provider": "openai",
					"model": "text-embedding-3-small"
				},
				"vector_store": {
					"type": "qdrant",
					"address": "localhost:6334"
				},
				"workflow": {},
				"schema": {}
			}`,
			wantErr: false,
			validate: func(t *testing.T, c *Config) {
				if c.LLM.ReasoningLLM.Provider != "openai" {
					t.Errorf("expected provider openai, got %s", c.LLM.ReasoningLLM.Provider)
				}
				// Check defaults were applied
				if c.LLM.ReasoningLLM.DefaultTemperature != 0.7 {
					t.Errorf("expected default temperature 0.7, got %f", c.LLM.ReasoningLLM.DefaultTemperature)
				}
				if c.Workflow.MaxIterations != 10 {
					t.Errorf("expected max iterations 10, got %d", c.Workflow.MaxIterations)
				}
			},
		},
		{
			name: "valid complete config",
			content: `{
				"llm": {
					"reasoning_llm": {
						"provider": "anthropic",
						"api_key": "test-key",
						"model": "claude-3-5-sonnet-20241022",
						"default_temperature": 0.8,
						"default_max_tokens": 4096,
						"timeout_seconds": 90
					},
					"fast_llm": {
						"provider": "openai",
						"model": "gpt-5-mini",
						"default_temperature": 0.3,
						"default_max_tokens": 512,
						"timeout_seconds": 20
					}
				},
				"embedding": {
					"provider": "openai",
					"api_key": "embed-key",
					"model": "text-embedding-3-large",
					"batch_size": 50,
					"timeout_seconds": 45
				},
				"vector_store": {
					"type": "qdrant",
					"address": "qdrant:6334",
					"api_key": "qdrant-key",
					"timeout_seconds": 60,
					"default_collection": "my_docs"
				},
				"web_search": {
					"enabled": true,
					"provider": "serper",
					"api_key": "search-key",
					"max_results": 5
				},
				"workflow": {
					"max_iterations": 15,
					"top_k_retrieval": 20,
					"top_n_reranking": 5,
					"min_relevance_score": 0.6,
					"default_strategy": "vector"
				},
				"schema": {
					"enable_llm_analysis": false,
					"enable_pattern_matching": true,
					"predefined_schemas": ["/path/to/schema.json"],
					"cache_schemas": false,
					"analysis_timeout_seconds": 120
				}
			}`,
			wantErr: false,
			validate: func(t *testing.T, c *Config) {
				// Verify custom values weren't overridden by defaults
				if c.LLM.ReasoningLLM.DefaultTemperature != 0.8 {
					t.Errorf("expected temperature 0.8, got %f", c.LLM.ReasoningLLM.DefaultTemperature)
				}
				if c.Embedding.BatchSize != 50 {
					t.Errorf("expected batch size 50, got %d", c.Embedding.BatchSize)
				}
				if c.Workflow.MaxIterations != 15 {
					t.Errorf("expected max iterations 15, got %d", c.Workflow.MaxIterations)
				}
				if c.WebSearch == nil {
					t.Error("expected web search config, got nil")
				} else if !c.WebSearch.Enabled {
					t.Error("expected web search enabled")
				}
			},
		},
		{
			name:    "invalid JSON",
			content: `{invalid json}`,
			wantErr: true,
		},
		{
			name:    "empty file",
			content: "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temporary file
			tmpDir := t.TempDir()
			tmpFile := filepath.Join(tmpDir, "config.json")

			if err := os.WriteFile(tmpFile, []byte(tt.content), 0644); err != nil {
				t.Fatalf("failed to write test file: %v", err)
			}

			// Test loading
			config, err := LoadFromFile(tmpFile)

			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if config == nil {
				t.Fatal("expected config, got nil")
			}

			if tt.validate != nil {
				tt.validate(t, config)
			}
		})
	}
}

// TestLoadFromFile_MissingFile tests loading from non-existent file.
func TestLoadFromFile_MissingFile(t *testing.T) {
	_, err := LoadFromFile("/nonexistent/path/config.json")
	if err == nil {
		t.Error("expected error for missing file, got nil")
	}
}

// TestLoadFromEnv tests loading configuration from environment variables.
func TestLoadFromEnv(t *testing.T) {
	tests := []struct {
		name     string
		envVars  map[string]string
		validate func(*testing.T, *Config)
	}{
		{
			name:    "default values with no env vars",
			envVars: map[string]string{},
			validate: func(t *testing.T, c *Config) {
				if c.LLM.ReasoningLLM.Provider != "openai" {
					t.Errorf("expected default provider openai, got %s", c.LLM.ReasoningLLM.Provider)
				}
				if c.LLM.ReasoningLLM.Model != "gpt-4o" {
					t.Errorf("expected default model gpt-4o, got %s", c.LLM.ReasoningLLM.Model)
				}
				if c.LLM.FastLLM.Model != "gpt-4o-mini" {
					t.Errorf("expected default fast model gpt-4o-mini, got %s", c.LLM.FastLLM.Model)
				}
				if c.Embedding.Model != "text-embedding-3-small" {
					t.Errorf("expected default embedding model, got %s", c.Embedding.Model)
				}
				if c.VectorStore.Type != "qdrant" {
					t.Errorf("expected default vector store qdrant, got %s", c.VectorStore.Type)
				}
				if c.VectorStore.Address != "localhost:6334" {
					t.Errorf("expected default address localhost:6334, got %s", c.VectorStore.Address)
				}
			},
		},
		{
			name: "custom env vars",
			envVars: map[string]string{
				"REASONING_LLM_PROVIDER":  "anthropic",
				"REASONING_LLM_API_KEY":   "test-key-reasoning",
				"REASONING_LLM_MODEL":     "claude-3-5-sonnet-20241022",
				"FAST_LLM_PROVIDER":       "openai",
				"FAST_LLM_API_KEY":        "test-key-fast",
				"FAST_LLM_MODEL":          "gpt-5-mini",
				"EMBEDDING_PROVIDER":      "openai",
				"EMBEDDING_API_KEY":       "test-key-embed",
				"EMBEDDING_MODEL":         "text-embedding-3-large",
				"VECTOR_STORE_TYPE":       "weaviate",
				"VECTOR_STORE_ADDRESS":    "weaviate:8080",
				"VECTOR_STORE_COLLECTION": "custom_docs",
			},
			validate: func(t *testing.T, c *Config) {
				if c.LLM.ReasoningLLM.Provider != "anthropic" {
					t.Errorf("expected provider anthropic, got %s", c.LLM.ReasoningLLM.Provider)
				}
				if c.LLM.ReasoningLLM.APIKey != "test-key-reasoning" {
					t.Errorf("expected reasoning API key, got %s", c.LLM.ReasoningLLM.APIKey)
				}
				if c.LLM.FastLLM.APIKey != "test-key-fast" {
					t.Errorf("expected fast API key, got %s", c.LLM.FastLLM.APIKey)
				}
				if c.VectorStore.Type != "weaviate" {
					t.Errorf("expected vector store weaviate, got %s", c.VectorStore.Type)
				}
				if c.VectorStore.DefaultCollection != "custom_docs" {
					t.Errorf("expected collection custom_docs, got %s", c.VectorStore.DefaultCollection)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Save and clear environment
			oldEnv := make(map[string]string)
			envKeys := []string{
				"REASONING_LLM_PROVIDER", "REASONING_LLM_API_KEY", "REASONING_LLM_MODEL",
				"FAST_LLM_PROVIDER", "FAST_LLM_API_KEY", "FAST_LLM_MODEL",
				"EMBEDDING_PROVIDER", "EMBEDDING_API_KEY", "EMBEDDING_MODEL",
				"VECTOR_STORE_TYPE", "VECTOR_STORE_ADDRESS", "VECTOR_STORE_COLLECTION",
			}
			for _, key := range envKeys {
				oldEnv[key] = os.Getenv(key)
				os.Unsetenv(key)
			}
			defer func() {
				for key, val := range oldEnv {
					if val != "" {
						os.Setenv(key, val)
					} else {
						os.Unsetenv(key)
					}
				}
			}()

			// Set test environment variables
			for key, val := range tt.envVars {
				os.Setenv(key, val)
			}

			// Load config
			config := LoadFromEnv()

			if config == nil {
				t.Fatal("expected config, got nil")
			}

			if tt.validate != nil {
				tt.validate(t, config)
			}
		})
	}
}

// TestLoadFromEnv_EnvFiles verifies that .env files populate configuration values when environment variables are otherwise unset.
func TestLoadFromEnv_EnvFiles(t *testing.T) {
	tmpDir := t.TempDir()

	// Ensure relevant environment variables are cleared (set to empty so helper can populate them, and they are restored after test).
	envKeys := []string{
		"REASONING_LLM_PROVIDER",
		"REASONING_LLM_API_KEY",
		"REASONING_LLM_MODEL",
		"FAST_LLM_PROVIDER",
		"FAST_LLM_API_KEY",
		"FAST_LLM_MODEL",
		"EMBEDDING_PROVIDER",
		"EMBEDDING_API_KEY",
		"EMBEDDING_MODEL",
		"VECTOR_STORE_TYPE",
		"VECTOR_STORE_ADDRESS",
		"VECTOR_STORE_COLLECTION",
	}

	for _, key := range envKeys {
		t.Setenv(key, "")
	}

	// Create .env and .env.local files. The local file should override shared defaults.
	envContent := "REASONING_LLM_PROVIDER=openai\nREASONING_LLM_API_KEY=base-key\nFAST_LLM_PROVIDER=openai\nFAST_LLM_API_KEY=base-key\n"
	if err := os.WriteFile(filepath.Join(tmpDir, ".env"), []byte(envContent), 0o600); err != nil {
		t.Fatalf("failed to write .env: %v", err)
	}

	localContent := "REASONING_LLM_PROVIDER=anthropic\nREASONING_LLM_API_KEY=local-key\nFAST_LLM_PROVIDER=anthropic\nFAST_LLM_API_KEY=local-key\nEMBEDDING_PROVIDER=openai\nEMBEDDING_API_KEY=embed-key\nEMBEDDING_MODEL=text-embedding-3-large\nVECTOR_STORE_TYPE=weaviate\nVECTOR_STORE_ADDRESS=weaviate:8080\nVECTOR_STORE_COLLECTION=custom_docs\n"
	if err := os.WriteFile(filepath.Join(tmpDir, ".env.local"), []byte(localContent), 0o600); err != nil {
		t.Fatalf("failed to write .env.local: %v", err)
	}

	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working directory: %v", err)
	}
	defer func() {
		_ = os.Chdir(wd)
	}()

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to chdir: %v", err)
	}

	cfg := LoadFromEnv()

	if cfg.LLM.ReasoningLLM.Provider != "anthropic" {
		t.Fatalf("expected reasoning provider from .env.local, got %s", cfg.LLM.ReasoningLLM.Provider)
	}
	if cfg.LLM.ReasoningLLM.APIKey != "local-key" {
		t.Fatalf("expected reasoning API key from .env.local, got %s", cfg.LLM.ReasoningLLM.APIKey)
	}
	if cfg.LLM.FastLLM.Provider != "anthropic" {
		t.Fatalf("expected fast provider from .env.local, got %s", cfg.LLM.FastLLM.Provider)
	}
	if cfg.Embedding.APIKey != "embed-key" {
		t.Fatalf("expected embedding API key from .env.local, got %s", cfg.Embedding.APIKey)
	}
	if cfg.VectorStore.Type != "weaviate" {
		t.Fatalf("expected vector store type from .env.local, got %s", cfg.VectorStore.Type)
	}
	if cfg.VectorStore.Address != "weaviate:8080" {
		t.Fatalf("expected vector store address from .env.local, got %s", cfg.VectorStore.Address)
	}
	if cfg.VectorStore.DefaultCollection != "custom_docs" {
		t.Fatalf("expected vector store collection from .env.local, got %s", cfg.VectorStore.DefaultCollection)
	}
}

// TestSaveToFile tests saving configuration to a JSON file.
func TestSaveToFile(t *testing.T) {
	config := &Config{
		LLM: LLMConfig{
			ReasoningLLM: LLMProviderConfig{
				Provider:           "openai",
				Model:              "gpt-5",
				DefaultTemperature: 0.7,
				DefaultMaxTokens:   2048,
				TimeoutSeconds:     60,
			},
			FastLLM: LLMProviderConfig{
				Provider:           "openai",
				Model:              "gpt-5-mini",
				DefaultTemperature: 0.5,
				DefaultMaxTokens:   1024,
				TimeoutSeconds:     30,
			},
		},
		Embedding: EmbeddingConfig{
			Provider:       "openai",
			Model:          "text-embedding-3-small",
			BatchSize:      100,
			TimeoutSeconds: 30,
		},
		VectorStore: VectorStoreConfig{
			Type:              "qdrant",
			Address:           "localhost:6334",
			DefaultCollection: "documents",
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

	t.Run("successful save", func(t *testing.T) {
		tmpDir := t.TempDir()
		tmpFile := filepath.Join(tmpDir, "config.json")

		if err := config.SaveToFile(tmpFile); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Verify file exists and can be read back
		data, err := os.ReadFile(tmpFile)
		if err != nil {
			t.Fatalf("failed to read saved file: %v", err)
		}

		var loaded Config
		if err := json.Unmarshal(data, &loaded); err != nil {
			t.Fatalf("failed to unmarshal saved config: %v", err)
		}

		// Verify a few key fields
		if loaded.LLM.ReasoningLLM.Provider != "openai" {
			t.Errorf("expected provider openai, got %s", loaded.LLM.ReasoningLLM.Provider)
		}
		if loaded.Workflow.MaxIterations != 10 {
			t.Errorf("expected max iterations 10, got %d", loaded.Workflow.MaxIterations)
		}
	})

	t.Run("invalid path", func(t *testing.T) {
		err := config.SaveToFile("/nonexistent/dir/config.json")
		if err == nil {
			t.Error("expected error for invalid path, got nil")
		}
	})
}

// TestToLLMConfig tests conversion to LLM config.
func TestToLLMConfig(t *testing.T) {
	config := &Config{
		LLM: LLMConfig{
			ReasoningLLM: LLMProviderConfig{
				Provider:           "openai",
				APIKey:             "test-key",
				BaseURL:            "https://api.openai.com",
				Model:              "gpt-5",
				DefaultTemperature: 0.8,
				DefaultMaxTokens:   3000,
				TimeoutSeconds:     90,
			},
		},
	}

	llmConfig := config.ToLLMConfig()

	if llmConfig.Provider != "openai" {
		t.Errorf("expected provider openai, got %s", llmConfig.Provider)
	}
	if llmConfig.APIKey != "test-key" {
		t.Errorf("expected API key test-key, got %s", llmConfig.APIKey)
	}
	if llmConfig.Model != "gpt-5" {
		t.Errorf("expected model gpt-4, got %s", llmConfig.Model)
	}
	if llmConfig.DefaultTemperature != 0.8 {
		t.Errorf("expected temperature 0.8, got %f", llmConfig.DefaultTemperature)
	}
	if llmConfig.DefaultMaxTokens != 3000 {
		t.Errorf("expected max tokens 3000, got %d", llmConfig.DefaultMaxTokens)
	}
}

// TestToFastLLMConfig tests conversion to fast LLM config.
func TestToFastLLMConfig(t *testing.T) {
	config := &Config{
		LLM: LLMConfig{
			FastLLM: LLMProviderConfig{
				Provider:           "anthropic",
				APIKey:             "fast-key",
				Model:              "claude-3-5-haiku-20241022",
				DefaultTemperature: 0.3,
				DefaultMaxTokens:   1000,
				TimeoutSeconds:     20,
			},
		},
	}

	llmConfig := config.ToFastLLMConfig()

	if llmConfig.Provider != "anthropic" {
		t.Errorf("expected provider anthropic, got %s", llmConfig.Provider)
	}
	if llmConfig.Model != "claude-3-5-haiku-20241022" {
		t.Errorf("expected model claude-3-5-haiku-20241022, got %s", llmConfig.Model)
	}
	if llmConfig.DefaultTemperature != 0.3 {
		t.Errorf("expected temperature 0.3, got %f", llmConfig.DefaultTemperature)
	}
}

// TestToEmbeddingConfig tests conversion to embedding config.
func TestToEmbeddingConfig(t *testing.T) {
	config := &Config{
		Embedding: EmbeddingConfig{
			Provider:       "openai",
			APIKey:         "embed-key",
			BaseURL:        "https://api.openai.com",
			Model:          "text-embedding-3-large",
			BatchSize:      50,
			TimeoutSeconds: 45,
		},
	}

	embedConfig := config.ToEmbeddingConfig()

	if embedConfig.Provider != "openai" {
		t.Errorf("expected provider openai, got %s", embedConfig.Provider)
	}
	if embedConfig.Model != "text-embedding-3-large" {
		t.Errorf("expected model text-embedding-3-large, got %s", embedConfig.Model)
	}
	if embedConfig.BatchSize != 50 {
		t.Errorf("expected batch size 50, got %d", embedConfig.BatchSize)
	}
}

// TestToVectorStoreConfig tests conversion to vector store config.
func TestToVectorStoreConfig(t *testing.T) {
	extra := map[string]interface{}{"key": "value"}
	config := &Config{
		VectorStore: VectorStoreConfig{
			Type:              "qdrant",
			Address:           "qdrant:6334",
			APIKey:            "qdrant-key",
			TimeoutSeconds:    60,
			DefaultCollection: "my_collection",
			Extra:             extra,
		},
	}

	vsConfig := config.ToVectorStoreConfig()

	if vsConfig.Type != "qdrant" {
		t.Errorf("expected type qdrant, got %s", vsConfig.Type)
	}
	if vsConfig.Address != "qdrant:6334" {
		t.Errorf("expected address qdrant:6334, got %s", vsConfig.Address)
	}
	if vsConfig.DefaultCollection != "my_collection" {
		t.Errorf("expected collection my_collection, got %s", vsConfig.DefaultCollection)
	}
	if vsConfig.Extra == nil {
		t.Error("expected extra config, got nil")
	}
}

// TestApplyDefaults tests the default value application logic.
func TestApplyDefaults(t *testing.T) {
	tests := []struct {
		name     string
		config   *Config
		validate func(*testing.T, *Config)
	}{
		{
			name: "empty config gets all defaults",
			config: &Config{
				LLM: LLMConfig{
					ReasoningLLM: LLMProviderConfig{Provider: "openai"},
					FastLLM:      LLMProviderConfig{Provider: "openai"},
				},
				Embedding:   EmbeddingConfig{},
				VectorStore: VectorStoreConfig{},
				Workflow:    WorkflowConfig{},
				Schema:      SchemaConfig{},
			},
			validate: func(t *testing.T, c *Config) {
				// LLM defaults
				if c.LLM.ReasoningLLM.DefaultTemperature != 0.7 {
					t.Errorf("expected default temperature 0.7, got %f", c.LLM.ReasoningLLM.DefaultTemperature)
				}
				if c.LLM.ReasoningLLM.DefaultMaxTokens != 2048 {
					t.Errorf("expected default max tokens 2048, got %d", c.LLM.ReasoningLLM.DefaultMaxTokens)
				}
				if c.LLM.FastLLM.DefaultTemperature != 0.5 {
					t.Errorf("expected fast default temperature 0.5, got %f", c.LLM.FastLLM.DefaultTemperature)
				}

				// Embedding defaults
				if c.Embedding.BatchSize != 100 {
					t.Errorf("expected batch size 100, got %d", c.Embedding.BatchSize)
				}
				if c.Embedding.TimeoutSeconds != 30 {
					t.Errorf("expected timeout 30, got %d", c.Embedding.TimeoutSeconds)
				}

				// VectorStore defaults
				if c.VectorStore.DefaultCollection != "documents" {
					t.Errorf("expected collection documents, got %s", c.VectorStore.DefaultCollection)
				}

				// Workflow defaults
				if c.Workflow.MaxIterations != 10 {
					t.Errorf("expected max iterations 10, got %d", c.Workflow.MaxIterations)
				}
				if c.Workflow.DefaultStrategy != "hybrid" {
					t.Errorf("expected strategy hybrid, got %s", c.Workflow.DefaultStrategy)
				}

				// Schema defaults
				if c.Schema.AnalysisTimeout != 60 {
					t.Errorf("expected analysis timeout 60, got %d", c.Schema.AnalysisTimeout)
				}
			},
		},
		{
			name: "custom values not overridden",
			config: &Config{
				LLM: LLMConfig{
					ReasoningLLM: LLMProviderConfig{
						DefaultTemperature: 0.9,
						DefaultMaxTokens:   4000,
						TimeoutSeconds:     120,
					},
					FastLLM: LLMProviderConfig{
						DefaultTemperature: 0.2,
						DefaultMaxTokens:   500,
						TimeoutSeconds:     15,
					},
				},
				Embedding: EmbeddingConfig{
					BatchSize:      200,
					TimeoutSeconds: 60,
				},
				VectorStore: VectorStoreConfig{
					DefaultCollection: "custom",
					TimeoutSeconds:    90,
				},
				Workflow: WorkflowConfig{
					MaxIterations:   20,
					TopKRetrieval:   30,
					DefaultStrategy: "vector",
				},
				Schema: SchemaConfig{
					AnalysisTimeout: 180,
				},
			},
			validate: func(t *testing.T, c *Config) {
				// Verify custom values weren't changed
				if c.LLM.ReasoningLLM.DefaultTemperature != 0.9 {
					t.Errorf("custom temperature was overridden")
				}
				if c.Embedding.BatchSize != 200 {
					t.Errorf("custom batch size was overridden")
				}
				if c.Workflow.MaxIterations != 20 {
					t.Errorf("custom max iterations was overridden")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			applyDefaults(tt.config)
			if tt.validate != nil {
				tt.validate(t, tt.config)
			}
		})
	}
}

// TestGetEnv tests the environment variable retrieval helper.
func TestGetEnv(t *testing.T) {
	tests := []struct {
		name         string
		key          string
		defaultValue string
		envValue     string
		expected     string
	}{
		{
			name:         "env var set",
			key:          "TEST_VAR",
			defaultValue: "default",
			envValue:     "custom",
			expected:     "custom",
		},
		{
			name:         "env var not set",
			key:          "UNSET_VAR",
			defaultValue: "default",
			envValue:     "",
			expected:     "default",
		},
		{
			name:         "empty default",
			key:          "ANOTHER_UNSET",
			defaultValue: "",
			envValue:     "",
			expected:     "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Save original value
			orig := os.Getenv(tt.key)
			defer func() {
				if orig != "" {
					os.Setenv(tt.key, orig)
				} else {
					os.Unsetenv(tt.key)
				}
			}()

			// Set test value
			if tt.envValue != "" {
				os.Setenv(tt.key, tt.envValue)
			} else {
				os.Unsetenv(tt.key)
			}

			// Test
			result := getEnv(tt.key, tt.defaultValue)
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}
