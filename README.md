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
- **Qdrant** vector database (for Phase 1)
- **OpenAI API key** (for LLM and embeddings)

### Installation

```bash
# Clone the repository
git clone https://github.com/yourusername/deep-thinking-agent.git
cd deep-thinking-agent

# Install dependencies
go mod download

# Run tests
go test ./...
```

### Configuration

Create a configuration file `config.json`:

```json
{
  "llm": {
    "reasoning_llm": {
      "provider": "openai",
      "api_key": "your-openai-key",
      "model": "gpt-4",
      "default_temperature": 0.7,
      "default_max_tokens": 2048
    },
    "fast_llm": {
      "provider": "openai",
      "api_key": "your-openai-key",
      "model": "gpt-3.5-turbo",
      "default_temperature": 0.5,
      "default_max_tokens": 1024
    }
  },
  "embedding": {
    "provider": "openai",
    "api_key": "your-openai-key",
    "model": "text-embedding-3-small",
    "batch_size": 100
  },
  "vector_store": {
    "type": "qdrant",
    "address": "localhost:6334",
    "default_collection": "documents"
  },
  "workflow": {
    "max_iterations": 10,
    "top_k_retrieval": 10,
    "top_n_reranking": 3,
    "default_strategy": "hybrid"
  }
}
```

Or use environment variables:

```bash
export REASONING_LLM_PROVIDER=openai
export REASONING_LLM_API_KEY=your-key
export REASONING_LLM_MODEL=gpt-4
export VECTOR_STORE_TYPE=qdrant
export VECTOR_STORE_ADDRESS=localhost:6334
```

### Basic Usage (Phase 1)

```go
package main

import (
    "context"
    "log"

    "deep-thinking-agent/internal/config"
    "deep-thinking-agent/pkg/llm/openai"
    "deep-thinking-agent/pkg/embedding"
    "deep-thinking-agent/pkg/vectorstore/qdrant"
)

func main() {
    // Load configuration
    cfg, err := config.LoadFromFile("config.json")
    if err != nil {
        log.Fatal(err)
    }

    // Initialize LLM provider
    llmProvider, err := openai.NewProvider(
        cfg.LLM.ReasoningLLM.APIKey,
        cfg.LLM.ReasoningLLM.Model,
        cfg.ToLLMConfig(),
    )
    if err != nil {
        log.Fatal(err)
    }

    // Initialize embedder
    embedder, err := embedding.NewOpenAIEmbedder(
        cfg.Embedding.APIKey,
        cfg.Embedding.Model,
        cfg.ToEmbeddingConfig(),
    )
    if err != nil {
        log.Fatal(err)
    }

    // Initialize vector store
    vectorStore, err := qdrant.NewStore(
        cfg.VectorStore.Address,
        cfg.ToVectorStoreConfig(),
    )
    if err != nil {
        log.Fatal(err)
    }
    defer vectorStore.Close()

    // Use components...
    ctx := context.Background()

    // Example: Generate completion
    resp, err := llmProvider.Complete(ctx, &llm.CompletionRequest{
        Messages: []llm.Message{
            {Role: "user", Content: "What is deep thinking RAG?"},
        },
    })
    if err != nil {
        log.Fatal(err)
    }

    log.Println("Response:", resp.Content)
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

## Implementation Roadmap

### Phase 1: Foundation ✅ (Completed)
- Core interfaces (LLM, Embedding, VectorStore, Parser)
- OpenAI LLM and embeddings implementation
- Qdrant vector store implementation
- Text and Markdown parsers
- State machine definitions
- Configuration system

### Phase 2: Schema System (In Progress)
- LLM-based schema analyzer
- Multi-level schema storage
- Schema-aware chunking strategies
- Schema registry with predefined patterns
- PDF and HTML parsers

### Phase 3: Agents & Retrieval (Planned)
- All 8 specialized agents
- Vector, keyword, and hybrid retrieval strategies
- Cross-encoder reranking
- Schema-filtered retrieval

### Phase 4: Workflow Execution (Planned)
- Graph-based workflow construction
- Deep thinking loop orchestration
- State machine executor
- Policy decision logic

### Phase 5: Interfaces (Planned)
- CLI tool for interactive queries
- HTTP/gRPC API server
- Usage examples and tutorials

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
