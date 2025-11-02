<div align="center">

```
╔═══════════════════════════════════════════════════════════════════════╗
║                                                                       ║
║   ____                     ________    _       __   _                 ║
║  / __ \___  ___  ____     /_  __/ /_  (_)___  / /__(_)___  ____ _     ║
║ / / / / _ \/ _ \/ __ \     / / / __ \/ / __ \/ //_/ / __ \/ __ `/     ║
║/ /_/ /  __/  __/ /_/ /    / / / / / / / / / / ,< / / / / / /_/ /      ║
║\____/\___/\___/ .___/    /_/ /_/ /_/_/_/ /_/_/|_/_/_/ /_/\__, /       ║
║              /_/            Agent                        /____/       ║
║                                                                       ║
║        Schema-Driven RAG with Iterative Multi-Hop Reasoning           ║
║                                                                       ║
╚═══════════════════════════════════════════════════════════════════════╝
```

[![Go Version][go-badge]][go-url]
[![License][license-badge]][license-url]
[![Release][release-badge]][release-url]
[![CI Status][ci-badge]][ci-url]
[![Coverage][coverage-badge]][coverage-url]
[![Go Report][report-badge]][report-url]
[![Security][security-badge]][security-url]

[![GitHub Stars][stars-badge]][stars-url]
[![Issues][issues-badge]][issues-url]
[![PRs][prs-badge]][prs-url]
[![Contributors][contributors-badge]][contributors-url]
[![Last Commit][commit-badge]][commit-url]

---

**[Features](#key-features) • [Quickstart](#5-minute-quickstart) • [Architecture](#architecture) • [Contributing](#contributing) • [Examples](examples/)**

---

</div>

## 🎯 Overview

Traditional RAG systems struggle with complex queries that require:
- Multi-step reasoning across multiple sources
- Understanding document structure beyond simple chunks
- Iterative refinement based on intermediate findings
- Integration of both internal documents and external knowledge

**Deep Thinking Agent** solves these challenges with a schema-driven approach and specialized AI agents that work together through an iterative deep thinking loop.

```
┌──────────────────────────────────────────────────────────────────┐
│                     Deep Thinking Loop                           │
│                                                                  │
│  Query ──▶ Plan ──▶ Route ──▶ Retrieve ──▶ Rerank ──▶ Answer     │
│              ↑                                      ↓            │
│              └──── Reflect ◀──── Compress ◀────────┘             │
│                                                                  │
│  8 Specialized Agents • Schema-Aware Retrieval • Multi-Strategy  │
└──────────────────────────────────────────────────────────────────┘
```

### ✨ What Makes This Different?

- **🔍 Schema-Driven**: LLM analyzes documents to derive structure (sections, hierarchy, semantic regions) for targeted retrieval
- **🤔 Deep Thinking Loop**: Iterative workflow that plans, retrieves, reflects, and decides whether findings are sufficient
- **🎯 Multi-Strategy Retrieval**: Intelligently selects between vector search, BM25 keyword search, or hybrid approaches
- **🔌 Pluggable Architecture**: Swap LLM providers, vector stores, and document parsers without changing core logic
- **✅ Production-Ready**: 88% test coverage, comprehensive error handling, extensive documentation

## 🌟 Inspiration

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
│              Workflow Orchestration Layer               │
│    Deep Thinking Loop: Plan→Route→Retrieve→             │
│         Rerank→Compress→Reflect→Policy                  │
└─────────────────────────────────────────────────────────┘
                         │
┌─────────────────────────────────────────────────────────┐
│                    Agent Layer                          │
│  8 Specialized Agents (Planner, Rewriter, Supervisor,   │
│  Retriever, Reranker, Distiller, Reflector, Policy)     │
└─────────────────────────────────────────────────────────┘
                         │
┌─────────────────────────────────────────────────────────┐
│              Retrieval & Storage Layer                  │
│    Vector Store • Schema Registry • Web Search          │
└─────────────────────────────────────────────────────────┘
                         │
┌─────────────────────────────────────────────────────────┐
│              Document Processing Layer                  │
│  Parser → Schema Analyzer → Chunker → Embedder          │
└─────────────────────────────────────────────────────────┘
```

For detailed architecture documentation, see [AGENTS.md](./AGENTS.md).

## ⚠️ Cost Warning

This system uses OpenAI APIs for LLM and embedding operations. **Running the example scripts will incur API costs** (approximately $2-5 total for all examples). Individual queries cost $0.06-0.31 depending on complexity. See [examples/README.md](examples/README.md#cost-estimates) for detailed cost breakdown.

**Recommendations:**
- Set [spending limits](https://platform.openai.com/account/limits) in your OpenAI account
- Start with simple queries to understand costs
- Use the `--no-schema` flag for ingestion to reduce costs during testing

## 5-Minute Quickstart

Get up and running quickly:

**1. Prerequisites:** Go 1.25.3+, Docker, OpenAI API key
   *For detailed setup instructions, see [SETUP.md](SETUP.md)*

**2. Clone & Build:**
```bash
git clone https://github.com/yourusername/deep-thinking-agent.git
cd deep-thinking-agent
go build -o bin/deep-thinking-agent ./cmd/cli
```

**3. Start Qdrant:**
```bash
docker run -d --name qdrant -p 6333:6333 -p 6334:6334 qdrant/qdrant
```

**4. Configure:**
```bash
export OPENAI_API_KEY="sk-your-key-here"
./bin/deep-thinking-agent config init
```

**5. Run Examples:**
```bash
cd examples
./01_setup.sh && ./02_ingest.sh && ./03_query.sh
```

**6. When Done:** Clean up resources
```bash
./examples/06_cleanup.sh
```

For detailed instructions, troubleshooting, and cleanup, see:
- [SETUP.md](SETUP.md) - Comprehensive installation guide
- [CLEANUP.md](CLEANUP.md) - Teardown and cleanup instructions

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
      "model": "gpt-4o",
      "default_temperature": 0.7
    },
    "fast_llm": {
      "provider": "openai",
      "model": "gpt-4o-mini",
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
└── AGENTS.md              # AI coding agent development guide
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

We welcome contributions! Deep Thinking Agent is open source and built with community collaboration in mind.

### Quick Start for Contributors

1. **Read** [CONTRIBUTING.md](CONTRIBUTING.md) for detailed guidelines
2. **Check** [TODO.md](TODO.md) for planned features and improvement areas
3. **Look for** issues labeled `good first issue` or `help wanted`
4. **Follow** our [Code of Conduct](CODE_OF_CONDUCT.md)

### Ways to Contribute

- **Report bugs** - Help us identify and fix issues
- **Suggest features** - Share ideas for improvements
- **Improve documentation** - Make the project easier to understand
- **Add tests** - Help us reach 100% test coverage
- **Fix bugs** - Submit pull requests for existing issues
- **Add features** - Implement planned features from TODO.md

### Development Setup

See [CONTRIBUTING.md](CONTRIBUTING.md#development-setup) for complete setup instructions.

Quick start:
```bash
git clone https://github.com/gerrymiller/deep-thinking-agent.git
cd deep-thinking-agent
go mod download
go test ./...
```

### Code Standards

- **90% minimum test coverage** for all new code
- **Go fmt and go vet** must pass
- **Comprehensive documentation** for exported functions
- **Conventional commit messages**
- **Gitflow branching** strategy

See [AGENTS.md](AGENTS.md) for detailed code standards and [CONTRIBUTING.md](CONTRIBUTING.md) for contribution workflow.

## 🛠️ Built With

[![Go](https://img.shields.io/badge/Go-1.25.3-00ADD8?style=for-the-badge&logo=go&logoColor=white)](https://go.dev/)
[![OpenAI](https://img.shields.io/badge/OpenAI-API-412991?style=for-the-badge&logo=openai&logoColor=white)](https://openai.com/)
[![Qdrant](https://img.shields.io/badge/Qdrant-Vector_DB-DC244C?style=for-the-badge)](https://qdrant.tech/)
[![Docker](https://img.shields.io/badge/Docker-Container-2496ED?style=for-the-badge&logo=docker&logoColor=white)](https://www.docker.com/)

**Core Technologies:**
- **LLM Integration**: OpenAI GPT-4 (extensible to Anthropic Claude, Ollama)
- **Vector Database**: Qdrant (with support for Weaviate, Milvus planned)
- **Embeddings**: OpenAI text-embedding-3-small
- **Document Parsing**: Native Go parsers for Text, Markdown, PDF, HTML

## 🤖 Development Tools

This project was developed using modern AI-assisted development tools to accelerate implementation while maintaining high code quality standards:

- **[Claude Code](https://claude.ai/code)** (Anthropic) - AI pair programming assistant
- **[Droid](https://factory.ai/)** (Factory) - AI software engineering agent

All code is authored by Gerry Miller with AI assistance as a productivity multiplier, similar to how developers use advanced IDEs, linters, and code completion tools. This transparency reflects our commitment to honest attribution while recognizing that AI tools are transforming software development workflows.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Author

**Gerry Miller**
Email: gerry@gerrymiller.com
GitHub: [@gerrymiller](https://github.com/gerrymiller)

## Acknowledgments

- Inspired by [deep-thinking-rag](https://github.com/FareedKhan-dev/deep-thinking-rag) by Fareed Khan
- Built with Go and leveraging OpenAI, Qdrant, and other open source technologies
- Community contributors who help improve and extend this project

---

<div align="center">

Made with ❤️ and 🧠 by [Gerry Miller](https://github.com/gerrymiller)

**[⭐ Star this repo](https://github.com/gerrymiller/deep-thinking-agent)** if you find it helpful!

</div>

<!-- Badge Definitions -->
[go-badge]: https://img.shields.io/github/go-mod/go-version/gerrymiller/deep-thinking-agent?style=for-the-badge&logo=go&logoColor=white
[go-url]: https://go.dev/
[license-badge]: https://img.shields.io/badge/License-MIT-blue.svg?style=for-the-badge
[license-url]: LICENSE
[release-badge]: https://img.shields.io/github/v/release/gerrymiller/deep-thinking-agent?style=for-the-badge&logo=github
[release-url]: https://github.com/gerrymiller/deep-thinking-agent/releases
[ci-badge]: https://img.shields.io/github/actions/workflow/status/gerrymiller/deep-thinking-agent/ci.yml?branch=main&style=for-the-badge&logo=github-actions&logoColor=white&label=CI
[ci-url]: https://github.com/gerrymiller/deep-thinking-agent/actions/workflows/ci.yml
[coverage-badge]: https://img.shields.io/codecov/c/github/gerrymiller/deep-thinking-agent?style=for-the-badge&logo=codecov&logoColor=white
[coverage-url]: https://codecov.io/gh/gerrymiller/deep-thinking-agent
[report-badge]: https://goreportcard.com/badge/github.com/gerrymiller/deep-thinking-agent?style=for-the-badge
[report-url]: https://goreportcard.com/report/github.com/gerrymiller/deep-thinking-agent
[security-badge]: https://img.shields.io/snyk/vulnerabilities/github/gerrymiller/deep-thinking-agent?style=for-the-badge&logo=snyk
[security-url]: https://snyk.io/test/github/gerrymiller/deep-thinking-agent
[stars-badge]: https://img.shields.io/github/stars/gerrymiller/deep-thinking-agent?style=social
[stars-url]: https://github.com/gerrymiller/deep-thinking-agent/stargazers
[issues-badge]: https://img.shields.io/github/issues/gerrymiller/deep-thinking-agent?style=flat-square
[issues-url]: https://github.com/gerrymiller/deep-thinking-agent/issues
[prs-badge]: https://img.shields.io/github/issues-pr/gerrymiller/deep-thinking-agent?style=flat-square&logo=github
[prs-url]: https://github.com/gerrymiller/deep-thinking-agent/pulls
[contributors-badge]: https://img.shields.io/github/contributors/gerrymiller/deep-thinking-agent?style=flat-square
[contributors-url]: https://github.com/gerrymiller/deep-thinking-agent/graphs/contributors
[commit-badge]: https://img.shields.io/github/last-commit/gerrymiller/deep-thinking-agent?style=flat-square
[commit-url]: https://github.com/gerrymiller/deep-thinking-agent/commits/main
