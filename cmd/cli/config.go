// Copyright 2025 Gerry Miller <gerry@gerrymiller.com>
//
// Licensed under the MIT License.
// See LICENSE file in the project root for full license information.

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"deep-thinking-agent/cmd/common"
)

func runConfig(args []string) error {
	fs := flag.NewFlagSet("config", flag.ExitOnError)

	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, `Usage: deep-thinking-agent config <subcommand> [options]

Manage configuration for the deep thinking agent.

Subcommands:
  show      Display current configuration
  init      Create a default configuration file
  validate  Validate a configuration file

Examples:
  # Show current config
  deep-thinking-agent config show

  # Create default config
  deep-thinking-agent config init

  # Validate config
  deep-thinking-agent config validate config.json
`)
	}

	if err := fs.Parse(args); err != nil {
		return err
	}

	if fs.NArg() < 1 {
		fs.Usage()
		return fmt.Errorf("subcommand is required")
	}

	subcommand := fs.Arg(0)

	switch subcommand {
	case "show":
		return showConfig(fs.Args()[1:])
	case "init":
		return initConfig(fs.Args()[1:])
	case "validate":
		return validateConfig(fs.Args()[1:])
	default:
		return fmt.Errorf("unknown subcommand: %s", subcommand)
	}
}

func showConfig(args []string) error {
	configPath := "config.json"
	if len(args) > 0 {
		configPath = args[0]
	}

	config, err := common.LoadConfig(configPath)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Pretty print config
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	fmt.Println(string(data))
	return nil
}

func initConfig(args []string) error {
	outputPath := "config.json"
	if len(args) > 0 {
		outputPath = args[0]
	}

	// Check if file exists
	if _, err := os.Stat(outputPath); err == nil {
		return fmt.Errorf("config file already exists: %s (delete it first or specify a different path)", outputPath)
	}

	// Create default config
	config := common.DefaultConfig()

	// Write to file
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(outputPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	fmt.Printf("Created default configuration: %s\n", outputPath)
	fmt.Println("\nNext steps:")
	fmt.Println("1. Edit the config file to add your API keys")
	fmt.Println("2. Configure your vector store connection")
	fmt.Printf("3. Run 'deep-thinking-agent config validate %s' to verify\n", outputPath)

	return nil
}

func validateConfig(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("config file path is required")
	}

	configPath := args[0]

	config, err := common.LoadConfig(configPath)
	if err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	// Perform validation checks
	var errors []string

	// Check LLM config
	if config.LLM.ReasoningLLM.Provider == "" {
		errors = append(errors, "reasoning_llm.provider is required")
	}
	if config.LLM.ReasoningLLM.Model == "" {
		errors = append(errors, "reasoning_llm.model is required")
	}
	if config.LLM.FastLLM.Provider == "" {
		errors = append(errors, "fast_llm.provider is required")
	}
	if config.LLM.FastLLM.Model == "" {
		errors = append(errors, "fast_llm.model is required")
	}

	// Check embedding config
	if config.Embedding.Provider == "" {
		errors = append(errors, "embedding.provider is required")
	}
	if config.Embedding.Model == "" {
		errors = append(errors, "embedding.model is required")
	}

	// Check vector store config
	if config.VectorStore.Type == "" {
		errors = append(errors, "vector_store.type is required")
	}
	if config.VectorStore.Address == "" {
		errors = append(errors, "vector_store.address is required")
	}

	if len(errors) > 0 {
		fmt.Println("Validation errors:")
		for _, err := range errors {
			fmt.Printf("  - %s\n", err)
		}
		return fmt.Errorf("configuration is invalid")
	}

	fmt.Printf("Configuration is valid: %s\n", configPath)
	return nil
}
