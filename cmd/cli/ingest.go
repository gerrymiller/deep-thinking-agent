// Copyright 2025 Gerry Miller <gerry@gerrymiller.com>
//
// Licensed under the MIT License.
// See LICENSE file in the project root for full license information.

package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"deep-thinking-agent/cmd/common"
)

func runIngest(args []string) error {
	fs := flag.NewFlagSet("ingest", flag.ExitOnError)
	configPath := fs.String("config", "config.json", "Path to configuration file")
	recursive := fs.Bool("recursive", false, "Recursively process directories")
	collection := fs.String("collection", "documents", "Target collection name")
	deriveSchema := fs.Bool("derive-schema", true, "Derive document schema using LLM")
	verbose := fs.Bool("verbose", false, "Show detailed processing information")

	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, `Usage: deep-thinking-agent ingest [options] <file-or-directory>...

Ingest documents into the vector store with schema analysis.

Options:
  -config string
        Path to configuration file (default "config.json")
  -recursive
        Recursively process directories
  -collection string
        Target collection name (default "documents")
  -derive-schema
        Derive document schema using LLM (default true)
  -verbose
        Show detailed processing information

Examples:
  # Ingest a single file
  deep-thinking-agent ingest document.txt

  # Ingest a directory
  deep-thinking-agent ingest -recursive ./documents

  # Ingest with custom collection
  deep-thinking-agent ingest -collection research_papers ./papers
`)
	}

	if err := fs.Parse(args); err != nil {
		return err
	}

	if fs.NArg() < 1 {
		return fmt.Errorf("at least one file or directory path is required")
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

	ctx := context.Background()

	// Process each path
	var totalFiles, totalChunks int
	for _, path := range fs.Args() {
		files, chunks, err := processPath(ctx, system, path, *recursive, *collection, *deriveSchema, *verbose)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to process %s: %v\n", path, err)
			continue
		}
		totalFiles += files
		totalChunks += chunks
	}

	fmt.Printf("\nIngestion complete:\n")
	fmt.Printf("  Files processed: %d\n", totalFiles)
	fmt.Printf("  Chunks created: %d\n", totalChunks)
	fmt.Printf("  Collection: %s\n", *collection)

	return nil
}

func processPath(ctx context.Context, system *common.System, path string, recursive bool, collection string, deriveSchema bool, verbose bool) (int, int, error) {
	info, err := os.Stat(path)
	if err != nil {
		return 0, 0, err
	}

	if info.IsDir() {
		return processDirectory(ctx, system, path, recursive, collection, deriveSchema, verbose)
	}

	return processFile(ctx, system, path, collection, deriveSchema, verbose)
}

func processDirectory(ctx context.Context, system *common.System, dirPath string, recursive bool, collection string, deriveSchema bool, verbose bool) (int, int, error) {
	var totalFiles, totalChunks int

	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return 0, 0, err
	}

	for _, entry := range entries {
		fullPath := filepath.Join(dirPath, entry.Name())

		if entry.IsDir() {
			if recursive {
				files, chunks, err := processDirectory(ctx, system, fullPath, recursive, collection, deriveSchema, verbose)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Warning: failed to process directory %s: %v\n", fullPath, err)
					continue
				}
				totalFiles += files
				totalChunks += chunks
			}
			continue
		}

		files, chunks, err := processFile(ctx, system, fullPath, collection, deriveSchema, verbose)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to process file %s: %v\n", fullPath, err)
			continue
		}
		totalFiles += files
		totalChunks += chunks
	}

	return totalFiles, totalChunks, nil
}

func processFile(ctx context.Context, system *common.System, filePath string, collection string, deriveSchema bool, verbose bool) (int, int, error) {
	// Check if file extension is supported
	ext := strings.ToLower(filepath.Ext(filePath))
	supportedExts := []string{".txt", ".md", ".markdown"}
	supported := false
	for _, supportedExt := range supportedExts {
		if ext == supportedExt {
			supported = true
			break
		}
	}

	if !supported {
		if verbose {
			fmt.Printf("Skipping unsupported file: %s\n", filePath)
		}
		return 0, 0, nil
	}

	if verbose {
		fmt.Printf("Processing: %s\n", filePath)
	}

	// Read file
	content, err := os.ReadFile(filePath)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to read file: %w", err)
	}

	// Ingest document
	chunks, err := system.IngestDocument(ctx, filePath, string(content), deriveSchema)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to ingest: %w", err)
	}

	if verbose {
		fmt.Printf("  Created %d chunks\n", chunks)
	}

	return 1, chunks, nil
}
