// Copyright 2025 Gerry Miller <gerry@gerrymiller.com>
//
// Licensed under the MIT License.
// See LICENSE file in the project root for full license information.

package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"os"
	"strings"

	"deep-thinking-agent/cmd/common"
	"deep-thinking-agent/pkg/workflow"
)

func runQuery(args []string) error {
	fs := flag.NewFlagSet("query", flag.ExitOnError)
	configPath := fs.String("config", "config.json", "Path to configuration file")
	interactive := fs.Bool("interactive", false, "Run in interactive mode")
	verbose := fs.Bool("verbose", false, "Show detailed execution information")
	maxIterations := fs.Int("max-iterations", 10, "Maximum number of reasoning iterations")

	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, `Usage: deep-thinking-agent query [options] <question>

Execute a deep thinking query using multi-hop reasoning.

Options:
  -config string
        Path to configuration file (default "config.json")
  -interactive
        Run in interactive mode for multiple queries
  -verbose
        Show detailed execution information
  -max-iterations int
        Maximum number of reasoning iterations (default 10)

Examples:
  # Single query
  deep-thinking-agent query "What are the main risk factors mentioned in the document?"

  # Interactive mode
  deep-thinking-agent query -interactive

  # With custom config
  deep-thinking-agent query -config prod.json "Analyze the financial trends"
`)
	}

	if err := fs.Parse(args); err != nil {
		return err
	}

	// Load configuration
	config, err := common.LoadConfig(*configPath)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Initialize system
	system, err := common.InitializeSystem(config)
	if err != nil {
		return fmt.Errorf("failed to initialize system: %w", err)
	}
	defer system.Close()

	if *interactive {
		return runInteractiveQuery(system, *verbose, *maxIterations)
	}

	// Single query mode
	if fs.NArg() < 1 {
		return fmt.Errorf("question is required")
	}

	question := strings.Join(fs.Args(), " ")
	return executeQuery(system, question, *verbose, *maxIterations)
}

func runInteractiveQuery(system *common.System, verbose bool, maxIterations int) error {
	fmt.Println("Deep Thinking Agent - Interactive Mode")
	fmt.Println("Type 'exit' or 'quit' to exit")
	fmt.Println()

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Query> ")
		if !scanner.Scan() {
			break
		}

		question := strings.TrimSpace(scanner.Text())
		if question == "" {
			continue
		}

		if question == "exit" || question == "quit" {
			fmt.Println("Goodbye!")
			break
		}

		if err := executeQuery(system, question, verbose, maxIterations); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		}
		fmt.Println()
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("scanner error: %w", err)
	}

	return nil
}

func executeQuery(system *common.System, question string, verbose bool, maxIterations int) error {
	ctx := context.Background()

	fmt.Printf("Question: %s\n\n", question)

	if verbose {
		fmt.Println("Executing deep thinking workflow...")
		fmt.Println()
	}

	// Create initial state
	state := workflow.NewState(question)
	state.MaxIterations = maxIterations

	// Execute workflow
	result, err := system.Executor.Execute(ctx, state)
	if err != nil {
		return fmt.Errorf("execution failed: %w", err)
	}

	// Display results
	if verbose {
		displayVerboseResults(result)
	} else {
		displayCompactResults(result)
	}

	return nil
}

func displayVerboseResults(state *workflow.State) {
	fmt.Println("=== Execution Plan ===")
	if state.Plan != nil {
		fmt.Printf("Reasoning: %s\n", state.Plan.Reasoning)
		fmt.Printf("Steps: %d\n\n", len(state.Plan.Steps))

		for i, step := range state.Plan.Steps {
			fmt.Printf("%d. %s\n", i+1, step.SubQuestion)
			if step.SchemaHint != "" {
				fmt.Printf("   Schema hint: %s\n", step.SchemaHint)
			}
			if step.ToolType != "" {
				fmt.Printf("   Tool: %s\n", step.ToolType)
			}
		}
		fmt.Println()
	}

	fmt.Println("=== Execution History ===")
	for i, pastStep := range state.PastSteps {
		fmt.Printf("Step %d: %s\n", i+1, pastStep.Step.SubQuestion)
		fmt.Printf("Summary: %s\n", pastStep.Summary)
		if len(pastStep.KeyFindings) > 0 {
			fmt.Println("Key Findings:")
			for _, finding := range pastStep.KeyFindings {
				fmt.Printf("  - %s\n", finding)
			}
		}
		fmt.Printf("Documents: %d\n", len(pastStep.RetrievedDocs))
		fmt.Println()
	}

	fmt.Println("=== Final Answer ===")
	if state.FinalAnswer != "" {
		fmt.Println(state.FinalAnswer)
	} else {
		fmt.Println("No final answer generated.")
	}
}

func displayCompactResults(state *workflow.State) {
	if state.Plan != nil && len(state.Plan.Steps) > 0 {
		fmt.Printf("Executed %d reasoning steps\n", len(state.PastSteps))
		fmt.Println()
	}

	fmt.Println("Answer:")
	if state.FinalAnswer != "" {
		fmt.Println(state.FinalAnswer)
	} else {
		// If no final answer, show key findings from all steps
		fmt.Println("\nKey Findings:")
		for _, pastStep := range state.PastSteps {
			for _, finding := range pastStep.KeyFindings {
				fmt.Printf("- %s\n", finding)
			}
		}
	}
}
