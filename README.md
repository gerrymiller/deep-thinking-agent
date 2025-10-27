# Deep Thinking Agent

A generic, schema-driven Retrieval-Augmented Generation (RAG) system that uses iterative deep thinking to answer complex, multi-hop queries across any document type.

## Overview

Traditional RAG systems struggle with complex queries that require:
- Multi-step reasoning across multiple sources
- Understanding document structure beyond simple chunks
- Iterative refinement based on intermediate findings
- Integration of both internal documents and external knowledge

Deep Thinking Agent solves these challenges through:

- **Schema-Driven Intelligence**: LLM analyzes documents to derive structural schemas (sections, hierarchy, semantic regions) rather than using hardcoded patterns
- **Deep Thinking Loop**: Graph-based workflow that iteratively plans, retrieves, reflects, and decides whether findings are sufficient
- **Multi-Strategy Retrieval**: Intelligently selects between vector search, keyword search, or hybrid approaches based on query characteristics
- **Pluggable Architecture**: Swap LLM providers, vector stores, and document parsers without changing core logic

## Inspiration

This project is inspired by [deep-thinking-rag](https://github.com/FareedKhan-dev/deep-thinking-rag) ([article](https://levelup.gitconnected.com/building-an-agentic-deep-thinking-rag-pipeline-to-solve-complex-queries-af69c5e044db)) but redesigned from the ground up for:
- Generic document types (not just SEC 10-K filings)
- Pluggable components for maximum flexibility
- Production-ready Go implementation with comprehensive testing

## Key Features

### Schema-Aware Document Processing
- **Dynamic Schema Derivation**: LLM analyzes each document to identify sections, hierarchy, and semantic regions
- **Multi-Level Metadata**: Stores schemas at chunk-level (vector DB), document-level (index), and pattern-level (registry)
- **Predefined Schema Support**: Optional predefined schemas for common document types to skip LLM analysis

### Deep Thinking Workflow
The system implements an iterative workflow with 8 specialized agents:

1. **Planner** - Decomposes queries into sequential substeps
2. **Query Rewriter** - Enhances queries with context and keywords
3. **Retrieval Supervisor** - Selects optimal retrieval strategy (vector/keyword/hybrid)
4. **Retriever** - Executes schema-filtered retrieval across document regions
5. **Reranker** - Applies cross-encoder for precision ranking
6. **Distiller** - Synthesizes retrieved chunks into coherent context
7. **Reflector** - Summarizes findings for accumulating history
8. **Policy Agent** - Decides whether to continue or finish based on sufficiency

### Pluggable Components
- **LLM Providers**: OpenAI (implemented), Anthropic, Ollama (planned)
- **Vector Stores**: Qdrant (implemented), Weaviate, Milvus (planned)
- **Document Parsers**: Text, Markdown (implemented), PDF, HTML (planned)
- **Web Search**: Optional external knowledge integration

## Architecture

```
┌─────────────────────────────────────────────────────────┐
│                    User Interface Layer                 │
│         CLI Tool • HTTP API • Go Library                │
└─────────────────────────────────────────────────────────┘
                         │
┌─────────────────────────────────────────────────────────┐
│              Workflow Orchestration Layer                │
│    Deep Thinking Loop: Plan→Route→Retrieve→             │
│         Rerank→Compress→Reflect→Policy                   │
└─────────────────────────────────────────────────────────┘
                         │
┌─────────────────────────────────────────────────────────┐
│                    Agent Layer                           │
│  8 Specialized Agents (Planner, Rewriter, Supervisor,   │
│  Retriever, Reranker, Distiller, Reflector, Policy)     │
└─────────────────────────────────────────────────────────┘
                         │
┌─────────────────────────────────────────────────────────┐
│              Retrieval & Storage Layer                   │
│    Vector Store • Schema Registry • Web Search          │
└─────────────────────────────────────────────────────────┘
                         │
┌─────────────────────────────────────────────────────────┐
│              Document Processing Layer                   │
│  Parser → Schema Analyzer → Chunker → Embedder          │
└─────────────────────────────────────────────────────────┘
```

For detailed architecture documentation, see [CLAUDE.md](./CLAUDE.md).

## Quick Start

### Prerequisites

- **Go 1.25.3+**
- **Qdrant** vector database
- **OpenAI API key** (for LLM and embeddings)
- **Docker** (optional, for running Qdrant)

### Installation

```bash
# Clone the repository
git clone https://github.com/yourusername/deep-thinking-agent.git
cd deep-thinking-agent

# Install dependencies
go mod download

# Build the CLI
go build -o bin/deep-thinking-agent ./cmd/cli

# Run tests
go test ./...
```

### Start Qdrant Vector Database

```bash
# Using Docker
docker run -d -p 6333:6333 -p 6334:6334 qdrant/qdrant

# Verify it's running
curl http://localhost:6333/collections
```

### Configuration

Initialize a default configuration:

```bash
./bin/deep-thinking-agent config init
```

This creates `config.json`. Edit it to add your API key, or set environment variable:

```bash
export OPENAI_API_KEY="your-openai-key-here"
```

Validate your configuration:

```bash
./bin/deep-thinking-agent config validate config.json
```

Example configuration:

```json
{
  "llm": {
    "reasoning_llm": {
      "provider": "openai",
      "model": "gpt-4",
      "default_temperature": 0.7
    },
    "fast_llm": {
      "provider": "openai",
      "model": "gpt-3.5-turbo",
      "default_temperature": 0.5
    }
  },
  "embedding": {
    "provider": "openai",
    "model": "text-embedding-3-small"
  },
  "vector_store": {
    "type": "qdrant",
    "address": "localhost:6334",
    "default_collection": "documents"
  },
  "workflow": {
    "max_iterations": 10,
    "top_k_retrieval": 10,
    "top_n_reranking": 3
  }
}
```

Environment variables override config file values:

```bash
export OPENAI_API_KEY=your-key
export REASONING_LLM_MODEL=gpt-4
export VECTOR_STORE_ADDRESS=localhost:6334
```

### CLI Usage

#### Ingest Documents

```bash
# Ingest a single document
./bin/deep-thinking-agent ingest document.txt

# Ingest all documents in a directory
./bin/deep-thinking-agent ingest -recursive ./documents

# Ingest with verbose output
./bin/deep-thinking-agent ingest -verbose document.md

# Ingest to a custom collection
./bin/deep-thinking-agent ingest -collection research ./papers
```

#### Query Documents

```bash
# Single query
./bin/deep-thinking-agent query "What are the main risk factors?"

# Interactive mode
./bin/deep-thinking-agent query -interactive

# Verbose mode (shows reasoning steps)
./bin/deep-thinking-agent query -verbose "Summarize the key findings"

# Control max reasoning iterations
./bin/deep-thinking-agent query -max-iterations 15 "Complex question"
```

#### Configuration Management

```bash
# Initialize default config
./bin/deep-thinking-agent config init

# Show current config
./bin/deep-thinking-agent config show

# Validate config
./bin/deep-thinking-agent config validate config.json
```

### Automated Examples

Run the complete workflow with provided examples:

```bash
# 1. Setup and build
./examples/01_setup.sh

# 2. Ingest sample documents
./examples/02_ingest.sh

# 3. Run simple queries
./examples/03_query.sh

# 4. Run advanced multi-hop queries
./examples/04_advanced.sh

# 5. Explore ingestion patterns
./examples/05_ingestion_patterns.sh
```

See [examples/README.md](examples/README.md) for detailed documentation.

### Go Library Usage

You can also use Deep Thinking Agent as a Go library:

```go
package main

import (
    "context"
    "log"

    "deep-thinking-agent/cmd/common"
)

func main() {
    // Load configuration
    config, err := common.LoadConfig("config.json")
    if err != nil {
        log.Fatal(err)
    }

    // Initialize system
    system, err := common.InitializeSystem(config)
    if err != nil {
        log.Fatal(err)
    }
    defer system.Close()

    // Execute a query
    ctx := context.Background()
    state := workflow.NewState("What are the main applications of AI in healthcare?")
    result, err := system.Executor.Execute(ctx, state)
    if err != nil {
        log.Fatal(err)
    }

    log.Println("Answer:", result.FinalAnswer)
}
```

## Development

### Project Structure

```
deep-thinking-agent/
├── pkg/                    # Public packages
│   ├── llm/               # LLM provider abstraction
│   ├── embedding/         # Embedding generation
│   ├── vectorstore/       # Vector database abstraction
│   ├── document/          # Document parsing and chunking
│   ├── schema/            # Schema analysis and management
│   ├── workflow/          # State machine and orchestration
│   ├── agent/             # Specialized agents
│   ├── retrieval/         # Retrieval strategies
│   └── websearch/         # Web search integration
├── internal/              # Private packages
│   ├── config/            # Configuration management
│   └── utils/             # Internal utilities
├── cmd/                   # Command-line tools
│   ├── cli/               # CLI application
│   ├── api/               # API server
│   └── common/            # Shared code
├── examples/              # Usage examples
├── test/                  # Test fixtures and data
└── CLAUDE.md             # AI assistant development guide
```

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests for a specific package
go test ./pkg/llm/openai/

# Run with verbose output
go test -v ./...
```

### Code Standards

All code must:
- Include license headers (see existing files for format)
- Have comprehensive unit tests (target 100% coverage)
- Follow Go conventions and pass `go vet` and `go fmt`
- Include inline documentation for all exported types and functions

### Adding New Components

#### New LLM Provider

1. Create package in `pkg/llm/yourprovider/`
2. Implement the `llm.Provider` interface
3. Add unit tests
4. Update configuration structs
5. Document in README

#### New Vector Store

1. Create package in `pkg/vectorstore/yourstore/`
2. Implement the `vectorstore.Store` interface
3. Add unit tests
4. Update configuration structs
5. Document in README

#### New Document Parser

1. Create file in `pkg/document/parser/`
2. Implement the `parser.Parser` interface
3. Register in `ParserRegistry`
4. Add unit tests
5. Document supported formats

## Implementation Status

### Phase 1: Foundation ✅ (Completed)
- Core interfaces (LLM, Embedding, VectorStore, Parser)
- OpenAI LLM and embeddings implementation
- Qdrant vector store implementation
- Text and Markdown parsers
- State machine definitions
- Configuration system
- Comprehensive test coverage

### Phase 2: Schema System ✅ (Completed)
- LLM-based schema analyzer
- Schema resolver with multiple strategies
- Multi-level schema storage
- Schema-aware chunking
- Schema registry with pattern matching
- Enhanced metadata handling

### Phase 3: Agents & Retrieval ✅ (Completed)
- All 8 specialized agents implemented
- Planner, Rewriter, Supervisor agents
- Retriever with schema-aware filtering
- Reranker with cross-encoder support
- Distiller for context synthesis
- Reflector for step summarization
- Policy agent for decision logic
- Vector, keyword, and hybrid retrieval strategies

### Phase 4: Workflow Execution ✅ (Completed)
- Graph-based workflow construction
- Deep thinking loop orchestration
- State machine executor
- Node-based agent wrappers
- Complete integration testing
- Timeout and error handling

### Phase 5: CLI Interface ✅ (Completed)
- Full-featured CLI tool
- Interactive query mode
- Document ingestion with recursive support
- Configuration management commands
- Comprehensive examples and documentation
- Automated example scripts

## Contributing

Contributions are welcome! Please:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Follow code standards (see Development section)
4. Write comprehensive tests
5. Update documentation
6. Submit a pull request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Author

**Gerry Miller**
Email: gerry@gerrymiller.com

## Acknowledgments

- Inspired by [deep-thinking-rag](https://github.com/FareedKhan-dev/deep-thinking-rag) by Fareed Khan
- Built with Go and leveraging OpenAI, Qdrant, and other open source technologies
