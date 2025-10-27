// Copyright 2025 Gerry Miller <gerry@gerrymiller.com>
//
// Licensed under the MIT License.
// See LICENSE file in the project root for full license information.

package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "query":
		if err := runQuery(os.Args[2:]); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	case "ingest":
		if err := runIngest(os.Args[2:]); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	case "config":
		if err := runConfig(os.Args[2:]); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	case "version":
		printVersion()
	case "help", "-h", "--help":
		printUsage()
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n\n", command)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println(`Deep Thinking Agent - Schema-driven RAG system

Usage:
  deep-thinking-agent <command> [options]

Commands:
  query       Execute a deep thinking query
  ingest      Ingest documents into the system
  config      Manage configuration
  version     Print version information
  help        Show this help message

Use "deep-thinking-agent <command> -h" for more information about a command.`)
}

func printVersion() {
	fmt.Println("Deep Thinking Agent v0.1.0")
	fmt.Println("Copyright 2025 Gerry Miller <gerry@gerrymiller.com>")
	fmt.Println("Licensed under the MIT License")
}
