# CLAUDE.md

This file provides guidance to AI assistants when working with code in this repository.

## Project Overview

**Deep Thinking Agent** is a schema-driven RAG system using Go 1.25.3. It implements an iterative deep thinking loop for complex, multi-hop queries across generic document types using LLM-derived schemas and pluggable components.

**Author**: Gerry Miller <gerry@gerrymiller.com>
**License**: MIT

## Working with AI Assistants

### Proactive Optimization
When working in this repository, AI assistants should:

1. **Be Proactive, Not Reactive** - Don't just complete the task asked. Continuously assess the project state and optimize related configurations, documentation, and code without waiting to be prompted.

2. **Think Holistically** - When touching one file or system, consider all related files that might need updates:
   - When updating CLAUDE.md → Check if .gitignore needs optimization
   - When adding dependencies → Check if CLAUDE.md needs documentation updates
   - When changing architecture → Check if both code AND documentation reflect changes

3. **Verify Before Assuming** - When configuring tooling (Claude Code, Git, CI/CD), always check official documentation before making assumptions about file patterns, naming conventions, or best practices.

4. **Complete the Full Picture** - If setting up infrastructure (documentation, CI/CD, tooling), ensure ALL related components are configured optimally, not just the minimum required.

### AI Assistant Configuration
- `.claude/settings.json` - Team-wide settings (commit to git)
- `.claude/settings.local.json` - Personal preferences (in .gitignore)
- `.claude/commands/` - Custom slash commands (commit to git)

## Code Standards

All code in this repository must adhere to these standards:

### File Headers
Every `.go` source file **authored as part of this project** MUST include this header:

```go
// Copyright 2025 Gerry Miller <gerry@gerrymiller.com>
//
// Licensed under the MIT License.
// See LICENSE file in the project root for full license information.
```

**Important**: Only add headers to files in these directories:
- `pkg/` - Public packages authored for this project
- `internal/` - Internal packages
- `cmd/` - Command-line tools

**Never** add headers to:
- `vendor/` - Third-party vendored code (if used)
- Go module cache files
- Generated files (protobuf, mock files, etc.) unless explicitly maintained
- Third-party code or libraries

For new files created in future years, update the copyright year accordingly.

### Test Coverage

**CRITICAL: Tests are MANDATORY, not optional. Never commit code without tests.**

#### Non-Negotiable Testing Requirements:
1. **EVERY new package MUST have tests before committing**
2. **EVERY new `.go` file MUST have a corresponding `_test.go` file**
3. **Target: 100% test coverage** - measure with `go test -cover ./...`
4. **Tests must be written BEFORE the commit**, not after
5. **Run `go test ./...` and verify ALL tests pass before every commit**

#### Test Quality Standards:
- Write comprehensive unit tests for all exported functions and types
- Use table-driven tests following Go best practices
- Include both positive and negative test cases
- Test error handling paths explicitly
- Test edge cases and boundary conditions
- Add integration tests where appropriate (use build tags like `//go:build integration`)

#### Testing Workflow (MUST FOLLOW):
```
1. Write code
2. Write tests for that code
3. Run: go test ./...
4. Verify: All tests pass
5. Run: go test -cover ./...
6. Verify: Acceptable coverage (aim for 100%)
7. ONLY THEN commit
```

**If you find yourself committing code without tests, STOP and write the tests first.**

#### Current Test Coverage Status

As of latest update, test coverage by package:

**High Coverage (>80%)**:
- `internal/config` - 98.1% ✅
- `pkg/schema` - 96.5% ✅
- `pkg/document/chunker` - 91.7% ✅
- `pkg/agent` - 86.1% ✅

**Moderate Coverage (50-80%)**:
- `pkg/retrieval` - 76.2%
- `pkg/workflow` - 68.2%

**Low Coverage (<50%)**:
- `pkg/llm/openai` - 41.5%
- `pkg/document/parser` - 35.9%
- `pkg/embedding` - 35.3%
- `pkg/nodes` - 16.1%

**No Tests (0%)**:
- `cmd/cli` - Requires integration testing framework
- `cmd/common` - Requires integration testing framework
- `pkg/vectorstore/qdrant` - Requires Qdrant test instance
- `pkg/llm` - Interface-only package
- `pkg/vectorstore` - Interface-only package

**Priority for Improvement**:
1. Add integration test framework for `cmd/` packages
2. Increase coverage for `pkg/document/parser` (add PDF and HTML parser tests)
3. Add `pkg/vectorstore/qdrant` tests with test containers
4. Improve `pkg/nodes` coverage with better mocking
5. Increase `pkg/embedding` and `pkg/llm/openai` coverage

### Code Quality
- All code must pass `go fmt`
- All code must pass `go vet`
- Run `golangci-lint run` if available
- Follow [Effective Go](https://go.dev/doc/effective_go) conventions
- Document all exported types, functions, and constants
- Keep functions focused and testable

### Documentation
- Update CLAUDE.md when adding new architectural components
- Update README.md for user-facing changes
- Include inline comments for complex logic
- Write clear commit messages following conventional commits style

## Git Workflow

This project follows a gitflow branching strategy:

- **`main`** - Production-ready code. Never commit directly to main.
- **`develop`** - Root development branch. Most branches stem from here.
- **Branch Types**:
  - **`feature/`** - New features and enhancements (from `develop`)
    - Example: `feature/add-thinking-loop`
  - **`bugfix/`** - Bug fixes during development (from `develop`)
    - Example: `bugfix/fix-schema-parsing`
  - **`hotfix/`** - Critical production fixes (from `main`)
    - Example: `hotfix/security-patch`
    - Must merge to both `main` AND `develop`
  - **`chore/`** - Maintenance, refactoring, dependencies (from `develop`)
    - Example: `chore/update-dependencies`
  - **`release/`** - Release preparation (from `develop`)
    - Example: `release/v1.0.0`
    - Must merge to both `main` AND `develop`

### Creating Branches

Use the custom slash commands (see [Claude Code Custom Commands](#claude-code-custom-commands)) or manually:

```bash
# Ensure you're on develop and up to date
git checkout develop
git pull origin develop

# Create and checkout a new branch
git checkout -b <branch-type>/<description>
```

### Completing Work
```bash
# Ensure code is tested and formatted
go test ./...
go fmt ./...
go vet ./...

# Commit and push
git add .
git commit -m "Description of changes"
git push origin <branch-name>

# Merge back to develop (after review if needed)
git checkout develop
git merge <branch-name>
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

This project includes custom slash commands to streamline common gitflow workflows:

- **`/new-feature`** - Creates a new feature branch from `develop` for new features and enhancements
- **`/bugfix`** - Creates a new bugfix branch from `develop` for fixing bugs during development
- **`/hotfix`** - Creates a new hotfix branch from `main` for critical production fixes
- **`/chore`** - Creates a new chore branch from `develop` for maintenance, refactoring, and dependencies
- **`/release`** - Creates a new release branch from `develop` for preparing production releases

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
│   ├── anthropic/          # Claude (planned)
│   └── ollama/             # Local models (planned)
│
├── embedding/              # Embedding generation
│   ├── embedder.go         # Embedder interface
│   └── openai_embedder.go  # OpenAI implementation
│
├── vectorstore/            # Vector database abstraction
│   ├── interface.go        # Store interface
│   ├── qdrant/             # Qdrant implementation
│   ├── weaviate/           # Weaviate (planned)
│   └── milvus/             # Milvus (planned)
│
├── document/
│   ├── parser/             # Format-specific parsers
│   │   ├── interface.go    # Parser interface + registry
│   │   ├── text.go         # Plain text
│   │   ├── markdown.go     # Markdown
│   │   ├── pdf.go          # PDF parsing
│   │   └── html.go         # HTML parsing
│   └── chunker/            # Schema-aware chunking
│       ├── chunker.go      # Chunker interface
│       ├── section.go      # Section-based chunking
│       ├── hierarchical.go # Hierarchical chunking
│       ├── semantic.go     # Semantic region chunking
│       └── sliding_window.go # Sliding window chunking
│
├── schema/                 # Schema analysis & management
│   ├── types.go            # Schema data structures
│   ├── analyzer.go         # LLM-based schema analysis
│   ├── resolver.go         # Schema resolution strategies
│   ├── registry.go         # Pattern storage and matching
│   └── metadata.go         # Metadata generation helpers
│
├── workflow/               # State machine & orchestration
│   ├── state.go            # State definitions and helpers
│   ├── graph.go            # Graph construction
│   └── executor.go         # Execution engine
│
├── nodes/                  # Workflow node wrappers
│   └── nodes.go            # Node implementations for all agents
│
├── agent/                  # Specialized agents
│   ├── planner.go          # Query decomposition
│   ├── rewriter.go         # Query enhancement
│   ├── supervisor.go       # Strategy selection
│   ├── retriever.go        # Schema-aware retrieval
│   ├── reranker.go         # Precision ranking
│   ├── distiller.go        # Context synthesis
│   ├── reflector.go        # Step summarization
│   └── policy.go           # Continue/finish decisions
│
├── retrieval/              # Retrieval strategies
│   ├── vector.go           # Semantic search
│   ├── keyword.go          # BM25 search
│   ├── hybrid.go           # RRF combination
│   └── schema.go           # Schema-filtered retrieval
│
└── websearch/              # Optional web search (planned)
    ├── interface.go
    └── providers/

internal/
├── config/                 # Configuration management
│   └── config.go           # JSON/env config loading
└── utils/                  # Internal utilities (reserved)

cmd/
├── cli/                    # CLI tool
│   ├── main.go             # CLI entry point
│   ├── query.go            # Query commands
│   ├── ingest.go           # Ingestion commands
│   └── config.go           # Config commands
├── api/                    # HTTP API server (planned)
└── common/                 # Shared CLI/API infrastructure
    └── system.go           # System initialization
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

#### Phase 1: Foundation ✅ (Completed)

**Goal**: Establish core abstractions and basic implementations

**Completed Components**:
- ✅ `pkg/llm/` - LLM provider interface with OpenAI implementation
- ✅ `pkg/embedding/` - Embedding interface with OpenAI embedder (batch processing, configurable dimensions)
- ✅ `pkg/vectorstore/` - Vector store interface with Qdrant implementation (full CRUD, metadata filtering)
- ✅ `pkg/document/parser/` - Parser interface with text and markdown implementations
- ✅ `pkg/schema/types.go` - Complete type definitions for schema system
- ✅ `pkg/workflow/state.go` - State machine types and helper methods
- ✅ `internal/config/` - JSON and environment-based configuration
- ✅ Comprehensive unit tests for all components

**Key Files**: 18 Go source files, 3 test files, all passing

---

#### Phase 2: Schema System ✅ (Completed)

**Goal**: Implement LLM-based document schema analysis and multi-level storage

**Completed Components**:
- ✅ `pkg/schema/analyzer.go` - LLM-powered schema derivation with section, hierarchy, and semantic region identification
- ✅ `pkg/schema/resolver.go` - Resolution strategy implementation (explicit → pattern → LLM → hybrid)
- ✅ `pkg/schema/registry.go` - Schema pattern storage and retrieval with pattern matching
- ✅ `pkg/schema/metadata.go` - Chunk-level and document-level metadata generation
- ✅ `pkg/document/chunker/` - Complete chunking implementations:
  - `section.go` - Section-based chunking
  - `hierarchical.go` - Hierarchical chunking
  - `semantic.go` - Semantic region chunking
  - `sliding_window.go` - Sliding window with schema boundaries
  - `chunker.go` - Main chunker interface
- ✅ `pkg/document/parser/pdf.go` - PDF parsing using pdfcpu
- ✅ `pkg/document/parser/html.go` - HTML parsing using golang.org/x/net/html
- ✅ Comprehensive unit tests for analyzer, resolver, registry, and chunking logic (96.5% coverage)

**Key Files**: 11 Go source files in schema and chunker packages, comprehensive test coverage

**Success Criteria**: ✅ Can ingest any document, derive schema, chunk appropriately, and store with metadata

---

#### Phase 3: Agents & Retrieval ✅ (Completed)

**Goal**: Implement all specialized agents and retrieval strategies

**Completed Components**:
1. **Agent Implementations** (`pkg/agent/`) - All 8 specialized agents:
   - ✅ `planner.go` - Query decomposition using reasoning LLM
   - ✅ `rewriter.go` - Query enhancement using fast LLM
   - ✅ `supervisor.go` - Strategy selection (vector/keyword/hybrid)
   - ✅ `retriever.go` - Schema-aware retrieval with filtering
   - ✅ `reranker.go` - Cross-encoder reranking implementation
   - ✅ `distiller.go` - Context synthesis and compression
   - ✅ `reflector.go` - Step summarization
   - ✅ `policy.go` - Continue/finish decision logic

2. **Retrieval Strategies** (`pkg/retrieval/`) - All strategies implemented:
   - ✅ `vector.go` - Pure semantic search
   - ✅ `keyword.go` - BM25 implementation
   - ✅ `hybrid.go` - RRF (Reciprocal Rank Fusion) combination
   - ✅ `schema.go` - Schema-filtered retrieval with metadata

3. **Integration**:
   - ✅ Agents connected to workflow state
   - ✅ Agent factory pattern implemented
   - ✅ Configuration for agent behavior

**Tests**: Unit tests for all agents (86.1% coverage), retrieval strategies (76.2% coverage)

**Key Files**: 12 Go source files (8 agents + 4 retrieval strategies), comprehensive test coverage

**Success Criteria**: ✅ All 8 agents working independently with comprehensive tests

---

#### Phase 4: Workflow Execution ✅ (Completed)

**Goal**: Implement graph-based workflow orchestration

**Completed Components**:
1. **Graph Construction** (`pkg/workflow/graph.go`)
   - ✅ Workflow graph structure defined
   - ✅ Node definitions and edges
   - ✅ Conditional routing logic

2. **Executor** (`pkg/workflow/executor.go`)
   - ✅ State machine execution engine
   - ✅ Node invocation and result handling
   - ✅ Error propagation and recovery

3. **Workflow Nodes** (`pkg/nodes/nodes.go`)
   - ✅ Node wrappers for all 8 agents
   - ✅ State transformation logic
   - ✅ Route decision functions

4. **Deep Thinking Loop**
   - ✅ Plan → Route → Retrieve → Rerank → Compress → Reflect → Policy
   - ✅ Iteration management
   - ✅ History accumulation

5. **End-to-End Integration**
   - ✅ All phases connected (parsing → schema → agents → workflow)
   - ✅ Complete query execution pipeline

**Tests**: Workflow package tests (68.2% coverage), integration tests implemented

**Key Files**: 3 Go source files in workflow, 1 in nodes package

**Success Criteria**: ✅ Can execute complex multi-hop queries end-to-end

---

#### Phase 5: CLI Interface ✅ (Completed)

**Goal**: Provide user-facing command-line interface

**Completed Tasks**:
1. **CLI Tool** (`cmd/cli/`)
   - Interactive query mode
   - Single-shot query mode
   - Document ingestion with recursive directory support
   - Configuration management (init, show, validate)
   - Verbose output mode
   - Max iterations control

2. **Common Infrastructure** (`cmd/common/`)
   - Configuration loading from JSON and environment variables
   - Complete system initialization
   - All 8 agents instantiation
   - Workflow graph construction

3. **Examples** (`examples/`)
   - Comprehensive README with usage instructions
   - Example configuration file
   - Sample documents (research, business, technical)
   - 5 automated example scripts:
     - 01_setup.sh - Build and initialization
     - 02_ingest.sh - Document ingestion
     - 03_query.sh - Simple queries
     - 04_advanced.sh - Multi-hop reasoning queries
     - 05_ingestion_patterns.sh - Advanced ingestion patterns

4. **Documentation**
   - Updated main README with CLI usage
   - Examples README with detailed instructions
   - Troubleshooting guide
   - Configuration reference
   - Performance tips

**Success Criteria**: ✅ Production-ready CLI with comprehensive documentation and examples

**Note**: HTTP API and gRPC API can be added in future phases if needed

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
