// Copyright 2025 Gerry Miller <gerry@gerrymiller.com>
//
// Licensed under the MIT License.
// See LICENSE file in the project root for full license information.

package common

import (
	"os"
	"path/filepath"
	"testing"
)

// TestLoadConfig_EnvFiles ensures that .env files are loaded and supply API keys when missing in the JSON config.
func TestLoadConfig_EnvFiles(t *testing.T) {
	tmpDir := t.TempDir()

	// Ensure the OPENAI_API_KEY environment variable is cleared for the test duration.
	if original, ok := os.LookupEnv("OPENAI_API_KEY"); ok {
		defer func() {
			_ = os.Setenv("OPENAI_API_KEY", original)
		}()
	} else {
		defer os.Unsetenv("OPENAI_API_KEY")
	}
	os.Unsetenv("OPENAI_API_KEY")

	// Create config file without API keys.
	configContent := `{
        "llm": {
            "reasoning_llm": {"provider": "openai", "model": "gpt-4o"},
            "fast_llm": {"provider": "openai", "model": "gpt-4o-mini"}
        },
        "embedding": {"provider": "openai", "model": "text-embedding-3-small"},
        "vector_store": {"type": "qdrant", "address": "localhost:6334"},
        "workflow": {"max_iterations": 10, "top_k_retrieval": 10, "top_n_reranking": 3, "default_strategy": "hybrid"}
    }`

	configPath := filepath.Join(tmpDir, "config.json")
	if err := os.WriteFile(configPath, []byte(configContent), 0o600); err != nil {
		t.Fatalf("failed to write config file: %v", err)
	}

	// Write .env and .env.local files. .env.local should override .env.
	if err := os.WriteFile(filepath.Join(tmpDir, ".env"), []byte("OPENAI_API_KEY=base-key\n"), 0o600); err != nil {
		t.Fatalf("failed to write .env file: %v", err)
	}

	if err := os.WriteFile(filepath.Join(tmpDir, ".env.local"), []byte("OPENAI_API_KEY=local-key\n"), 0o600); err != nil {
		t.Fatalf("failed to write .env.local file: %v", err)
	}

	// Run the load from within the temporary directory so the helper sees the test .env files.
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working directory: %v", err)
	}
	defer func() {
		_ = os.Chdir(wd)
	}()

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to change directory: %v", err)
	}

	cfg, err := LoadConfig(configPath)
	if err != nil {
		t.Fatalf("LoadConfig returned error: %v", err)
	}

	if cfg.LLM.ReasoningLLM.APIKey != "local-key" {
		t.Fatalf("expected reasoning API key from .env.local, got %q", cfg.LLM.ReasoningLLM.APIKey)
	}

	if cfg.LLM.FastLLM.APIKey != "local-key" {
		t.Fatalf("expected fast API key from .env.local, got %q", cfg.LLM.FastLLM.APIKey)
	}

	if cfg.Embedding.APIKey != "local-key" {
		t.Fatalf("expected embedding API key from .env.local, got %q", cfg.Embedding.APIKey)
	}
}
