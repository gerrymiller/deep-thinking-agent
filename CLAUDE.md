# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Go project named "deep-thinking-agent" using Go 1.25.3. The project appears to be in early stages of development.

## Working with Claude Code

### Proactive Optimization
When working in this repository, Claude Code instances should:

1. **Be Proactive, Not Reactive** - Don't just complete the task asked. Continuously assess the project state and optimize related configurations, documentation, and code without waiting to be prompted.

2. **Think Holistically** - When touching one file or system, consider all related files that might need updates:
   - When updating CLAUDE.md → Check if .gitignore needs optimization
   - When adding dependencies → Check if CLAUDE.md needs documentation updates
   - When changing architecture → Check if both code AND documentation reflect changes

3. **Verify Before Assuming** - When configuring tooling (Claude Code, Git, CI/CD), always check official documentation before making assumptions about file patterns, naming conventions, or best practices.

4. **Complete the Full Picture** - If setting up infrastructure (documentation, CI/CD, tooling), ensure ALL related components are configured optimally, not just the minimum required.

### Claude Code Specific Configuration
- `.claude/settings.json` - Team-wide settings (commit to git)
- `.claude/settings.local.json` - Personal preferences (in .gitignore)
- `.claude/commands/` - Custom slash commands (commit to git)
- Always verify configuration patterns against [official documentation](https://docs.claude.com/en/docs/claude-code/)

## Git Workflow

This project follows a gitflow-like branching strategy:

- **`main`** - Production-ready code. Never commit directly to main.
- **`develop`** - Root development branch. All feature branches stem from here.
- **Feature branches** - All development work happens in feature branches off `develop`
  - Naming: `feature/description-of-feature`
  - Example: `feature/add-thinking-loop`

### Creating a Feature Branch
```bash
# Ensure you're on develop and up to date
git checkout develop
git pull origin develop

# Create and checkout a new feature branch
git checkout -b feature/your-feature-name
```

### Completing a Feature
```bash
# Ensure code is tested and formatted
go test ./...
go fmt ./...
go vet ./...

# Commit and push
git add .
git commit -m "Description of changes"
git push origin feature/your-feature-name

# Merge back to develop (after review if needed)
git checkout develop
git merge feature/your-feature-name
git push origin develop
```

## Development Commands

### Building
```bash
go build -o bin/deep-thinking-agent
```

### Running
```bash
go run .
```

### Testing
```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run a single test
go test -run TestName ./path/to/package

# Run tests with verbose output
go test -v ./...
```

### Code Quality
```bash
# Format code
go fmt ./...

# Run linter (requires golangci-lint)
golangci-lint run

# Vet code for common issues
go vet ./...

# Tidy dependencies
go mod tidy

# Verify dependencies
go mod verify
```

## Security

This project uses Snyk for security scanning. When adding or modifying code:

1. Run security scans on new first-party code using snyk_code_scan tool
2. Fix any security issues found using Snyk's context
3. Rescan after fixes to verify issues are resolved
4. Repeat until no new issues are found

## Claude Code Custom Commands

This project includes custom slash commands to streamline common workflows:

- **`/new-feature`** - Creates a new feature branch from `develop` following gitflow conventions

Custom commands are defined in `.claude/commands/` and can be extended as needed.

## Architecture

### Overview

This is a **schema-driven deep-thinking RAG system** inspired by [deep-thinking-rag](https://github.com/FareedKhan-dev/deep-thinking-rag) but designed for:
- **Generic document types** (not just 10-Ks) using LLM-derived schemas
- **Pluggable components** for LLMs, vector stores, and search providers
- **Schema-aware retrieval** using multi-level metadata

### Core Concepts

#### Deep Thinking Loop
The system uses a graph-based workflow that iteratively:
1. **Plans** - Decomposes queries into sequential steps
2. **Routes** - Selects retrieval tools (doc_search, web_search, schema_filter)
3. **Retrieves** - Executes broad retrieval with strategy selection (vector/keyword/hybrid)
4. **Reranks** - Applies cross-encoder for precision filtering
5. **Compresses** - Distills retrieved chunks into coherent context
6. **Reflects** - Summarizes findings and adds to history
7. **Policy** - Decides whether to continue or finish

#### Schema-Driven Metadata
Unlike traditional RAG that uses fixed document structures:
- LLM analyzes documents to derive schemas (sections, hierarchy, semantic regions)
- Schemas stored at multiple levels (vector DB chunks, document index, registry)
- Retrieval uses schema metadata for targeted searches
- Supports predefined schemas for efficiency (designed but not yet implemented)

### Package Structure

```
pkg/
├── llm/                    # LLM provider abstraction
│   ├── interface.go        # Provider interface
│   ├── openai/             # OpenAI implementation
│   ├── anthropic/          # Claude (Phase 2)
│   └── ollama/             # Local models (Phase 2)
│
├── embedding/              # Embedding generation
│   ├── embedder.go         # Embedder interface
│   └── openai_embedder.go  # OpenAI implementation
│
├── vectorstore/            # Vector database abstraction
│   ├── interface.go        # Store interface
│   ├── qdrant/             # Qdrant implementation
│   ├── weaviate/           # Weaviate (Phase 2)
│   └── milvus/             # Milvus (Phase 2)
│
├── document/
│   ├── parser/             # Format-specific parsers
│   │   ├── interface.go    # Parser interface + registry
│   │   ├── text.go         # Plain text
│   │   ├── markdown.go     # Markdown
│   │   ├── pdf.go          # PDF (Phase 2)
│   │   └── html.go         # HTML (Phase 2)
│   └── chunker/            # Schema-aware chunking (Phase 2)
│
├── schema/                 # Schema analysis & management
│   ├── types.go            # Schema data structures
│   ├── analyzer.go         # LLM-based analysis (Phase 2)
│   ├── registry.go         # Pattern storage (Phase 2)
│   └── metadata.go         # Metadata helpers (Phase 2)
│
├── workflow/               # State machine & orchestration
│   ├── state.go            # State definitions
│   ├── graph.go            # Graph construction (Phase 3)
│   ├── executor.go         # Execution engine (Phase 3)
│   └── nodes.go            # Node implementations (Phase 3)
│
├── agent/                  # Specialized agents (Phase 3)
│   ├── planner.go          # Query decomposition
│   ├── rewriter.go         # Query enhancement
│   ├── supervisor.go       # Strategy selection
│   ├── retriever.go        # Schema-aware retrieval
│   ├── reranker.go         # Precision ranking
│   ├── distiller.go        # Context synthesis
│   ├── reflector.go        # Step summarization
│   └── policy.go           # Continue/finish decisions
│
├── retrieval/              # Retrieval strategies (Phase 3)
│   ├── vector.go           # Semantic search
│   ├── keyword.go          # BM25 search
│   ├── hybrid.go           # RRF combination
│   └── schema.go           # Schema-filtered retrieval
│
└── websearch/              # Optional web search (Phase 4)
    ├── interface.go
    └── providers/

internal/
├── config/                 # Configuration management
│   └── config.go           # JSON/env config loading
└── utils/                  # Internal utilities

cmd/
├── cli/                    # CLI tool (Phase 5)
├── api/                    # HTTP API server (Phase 5)
└── common/                 # Shared CLI/API code
```

### Key Data Structures

#### Workflow State (pkg/workflow/state.go)
```go
type State struct {
    OriginalQuestion   string
    Plan               *Plan
    CurrentStepIndex   int
    PastSteps          []PastStep
    RetrievedDocs      []Document
    RerankedDocs       []Document
    SynthesizedContext string
    FinalAnswer        string
    RelevantSchemas    map[string]*DocumentSchema  // Schema context
    ActiveFilters      *SchemaFilters              // Metadata filters
    ShouldContinue     bool
}
```

#### Document Schema (pkg/schema/types.go)
```go
type DocumentSchema struct {
    DocID            string
    Format           string
    Sections         []Section           // Logical divisions
    Hierarchy        *HierarchyTree      // Structural hierarchy
    SemanticRegions  []SemanticRegion    // Topic-based regions
    CustomAttributes map[string]interface{}
    ChunkingStrategy string
}
```

#### Chunk Metadata
```go
type ChunkMetadata struct {
    DocID         string
    SectionID     string
    SectionType   string      // e.g., "risk_factors", "methodology"
    HierarchyPath string      // e.g., "1.2.3"
    SemanticTags  []string    // LLM-identified tags
}
```

### Configuration

Configuration via JSON file or environment variables (see `internal/config/config.go`):

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
    "top_n_reranking": 3,
    "default_strategy": "hybrid"
  }
}
```

### Implementation Status

**Phase 1: Foundation** ✅ (Current)
- Core interfaces (LLM, Embedding, VectorStore, Parser)
- Basic implementations (OpenAI LLM/embeddings, Qdrant)
- Document parsers (text, markdown; PDF/HTML in Phase 2)
- State definitions
- Configuration management
- Unit tests

**Phase 2: Schema System** (Next)
- LLM-based schema analyzer
- Multi-level schema storage
- Schema-aware chunking
- Schema registry with predefined patterns
- PDF and HTML parsers

**Phase 3: Agents & Retrieval** (Planned)
- All 8 agent implementations
- Retrieval strategies (vector, keyword, hybrid)
- Cross-encoder reranking
- Schema-filtered retrieval

**Phase 4: Workflow Execution** (Planned)
- Graph construction
- Deep thinking loop orchestration
- State machine execution
- Policy decisions

**Phase 5: Interfaces** (Planned)
- CLI tool
- HTTP/gRPC API
- Usage examples

### Adding New Components

#### New LLM Provider
1. Implement `llm.Provider` interface in `pkg/llm/yourprovider/`
2. Add provider type to config
3. Update factory in initialization code

#### New Vector Store
1. Implement `vectorstore.Store` interface in `pkg/vectorstore/yourstore/`
2. Add store type to config
3. Update factory in initialization code

#### New Document Parser
1. Implement `parser.Parser` interface in `pkg/document/parser/`
2. Register in `ParserRegistry`

### References

- Inspired by: https://github.com/FareedKhan-dev/deep-thinking-rag
- Article: https://levelup.gitconnected.com/building-an-agentic-deep-thinking-rag-pipeline-to-solve-complex-queries-af69c5e044db
