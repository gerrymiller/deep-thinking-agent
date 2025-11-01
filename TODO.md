# TODO - Future Improvements

This document tracks planned improvements, enhancements, and technical debt for the Deep Thinking Agent project.

**Last Updated**: 2025-11-01  
**Project Status**: ✅ Complete working reference implementation

---

## Priority: High (Production Readiness)

These items are essential for production deployment and represent the final gaps in the reference implementation.

### Test Coverage Improvements

**Target**: Achieve 90%+ coverage for all packages (currently 7/10 packages meet this goal)

#### 1. pkg/embedding (45.1% → 90%+)
- **Status**: 🔴 Not Started
- **Main Gap**: `Embed()` function success path testing (currently 12.5% coverage)
- **Requirements**: 
  - OpenAI API client mocking infrastructure
  - httptest server or API mock library setup
  - Test fixtures for batch embedding operations
- **Complexity**: High
- **Estimate**: 2-3 days
- **Impact**: Core functionality for document ingestion

#### 2. pkg/nodes (52.9% → 90%+)
- **Status**: 🔴 Not Started
- **Main Gap**: `Execute()` success paths for all 8 node types (currently 0-80% coverage)
- **Requirements**:
  - Complex workflow state setup with plans, steps, and mock agents
  - Comprehensive state fixtures for all node types
  - Integration test patterns for node chains
- **Complexity**: Medium-High
- **Estimate**: 3-4 days
- **Impact**: Workflow orchestration reliability

#### 3. pkg/document/parser (84.1% → 90%+)
- **Status**: 🟡 Partial (close to threshold)
- **Main Gap**: PDF parser at 32.4% coverage
- **Requirements**:
  - PDF test fixtures (actual PDF files or generation)
  - Complex PDF scenarios (multi-page, images, tables)
  - Error handling for malformed PDFs
- **Complexity**: Medium
- **Estimate**: 1-2 days
- **Impact**: Document format support completeness

#### 4. pkg/vectorstore/qdrant (0.0% → 90%+)
- **Status**: 🔴 Not Started
- **Main Gap**: No tests for Qdrant client operations
- **Requirements**:
  - Qdrant test instance via Docker containers (testcontainers-go)
  - Integration tests for CRUD operations
  - Error handling and connection failure scenarios
- **Complexity**: Medium
- **Estimate**: 2-3 days
- **Impact**: Critical for vector storage reliability

**Total Estimated Effort**: 8-12 days to complete all test coverage gaps

---

## Priority: Medium (Planned Features)

These are features documented in the architecture but not yet implemented.

### Additional LLM Providers

#### 1. Anthropic/Claude Support
- **Status**: 🔴 Not Started
- **Requirements**:
  - Implement `llm.Provider` interface in `pkg/llm/anthropic/`
  - Handle Claude-specific API patterns (thinking blocks, tool use)
  - Add to configuration and factory initialization
  - Test coverage for provider implementation
- **Benefits**: 
  - Alternative to OpenAI for users
  - Potentially better reasoning for planning/policy agents
  - Cost optimization options
- **Estimate**: 2-3 days

#### 2. Ollama Support (Local Models)
- **Status**: 🔴 Not Started
- **Requirements**:
  - Implement `llm.Provider` interface in `pkg/llm/ollama/`
  - Handle streaming responses
  - Support for local model deployment
  - Configuration for model selection and endpoints
- **Benefits**:
  - Zero API costs for development/testing
  - Data privacy for sensitive documents
  - Offline operation capability
- **Estimate**: 2-3 days

### Additional Vector Stores

#### 1. Weaviate Support
- **Status**: 🔴 Not Started
- **Requirements**:
  - Implement `vectorstore.Store` interface in `pkg/vectorstore/weaviate/`
  - Schema mapping for metadata
  - Hybrid search support
  - Configuration and connection management
- **Benefits**: Cloud-native option, GraphQL API, built-in ML models
- **Estimate**: 3-4 days

#### 2. Milvus Support
- **Status**: 🔴 Not Started
- **Requirements**:
  - Implement `vectorstore.Store` interface in `pkg/vectorstore/milvus/`
  - Collection and partition management
  - Index optimization for performance
  - Configuration and connection management
- **Benefits**: High performance at scale, active community
- **Estimate**: 3-4 days

### Web Search Integration

#### External Knowledge Integration
- **Status**: 🔴 Not Started
- **Requirements**:
  - Design `pkg/websearch/` interface
  - Implement providers (Tavily, Brave, SerpAPI)
  - Integrate into retrieval supervisor routing logic
  - Add to workflow graph as optional node
  - Configuration for API keys and rate limits
- **Benefits**: 
  - Answer queries requiring current information
  - Supplement document knowledge with web data
  - Fact verification against external sources
- **Complexity**: Medium
- **Estimate**: 4-5 days

### API Interfaces

#### 1. HTTP/REST API
- **Status**: 🔴 Not Started
- **Location**: `cmd/api/`
- **Requirements**:
  - REST endpoints for query, ingest, config
  - OpenAPI/Swagger documentation
  - Authentication/authorization
  - Rate limiting
  - CORS configuration
- **Benefits**: Web UI integration, microservices deployment
- **Estimate**: 5-7 days

#### 2. gRPC API
- **Status**: 🔴 Not Started
- **Location**: `cmd/api/`
- **Requirements**:
  - Protocol buffer definitions
  - Streaming support for long queries
  - Service implementation
  - Client SDK generation
- **Benefits**: High performance, strong typing, cross-language support
- **Estimate**: 4-5 days

---

## Priority: Low (Enhancements)

These improvements enhance production operations, performance, and user experience but aren't blocking deployment.

### Observability & Monitoring

#### 1. Structured Logging
- **Status**: 🔴 Not Started
- **Current State**: Basic logging with fmt.Printf
- **Requirements**:
  - Integrate structured logging library (zerolog, zap, or slog)
  - Add log levels (DEBUG, INFO, WARN, ERROR)
  - Context propagation with request IDs
  - Log correlation across agents and nodes
- **Benefits**:
  - Production debugging capabilities
  - Log aggregation and analysis
  - Performance troubleshooting
- **Estimate**: 2-3 days

#### 2. Metrics & Instrumentation
- **Status**: 🔴 Not Started
- **Requirements**:
  - Prometheus metrics integration
  - Key metrics:
    - Query latency by agent/node
    - LLM API call duration and costs
    - Retrieval strategy performance
    - Token usage per query
    - Error rates by component
  - Metrics endpoint for scraping
- **Benefits**:
  - Performance monitoring
  - Cost tracking and optimization
  - SLA compliance verification
- **Estimate**: 3-4 days

#### 3. Distributed Tracing
- **Status**: 🔴 Not Started
- **Requirements**:
  - OpenTelemetry integration
  - Span creation for all agents and nodes
  - Trace context propagation
  - Integration with Jaeger or Tempo
- **Benefits**:
  - End-to-end query visualization
  - Bottleneck identification
  - Multi-hop reasoning flow analysis
- **Estimate**: 3-4 days

### Cost Optimization

#### 1. Query Result Caching
- **Status**: 🔴 Not Started
- **Requirements**:
  - Cache layer design (Redis or in-memory)
  - Cache key generation from query + context
  - TTL configuration
  - Cache invalidation on document updates
  - Metrics for cache hit rates
- **Benefits**:
  - Reduce duplicate LLM API calls
  - Faster response for repeated queries
  - Significant cost savings for common questions
- **Estimate**: 3-4 days

#### 2. LLM Response Caching
- **Status**: 🔴 Not Started
- **Requirements**:
  - Agent-specific caching (especially planner, rewriter)
  - Semantic similarity matching for near-duplicate queries
  - Persistent cache storage
  - Cache warming for common queries
- **Benefits**:
  - 50-80% cost reduction for repeated patterns
  - Faster response times
- **Estimate**: 2-3 days

#### 3. Batch Operations
- **Status**: 🔴 Not Started
- **Requirements**:
  - Batch embedding API calls (already partially implemented)
  - Batch LLM completion for multiple steps
  - Optimize Qdrant batch upload
  - Parallel document processing during ingestion
- **Benefits**:
  - Lower per-operation costs
  - Higher throughput for bulk operations
  - Better API rate limit utilization
- **Estimate**: 2-3 days

### Performance Optimization

#### 1. Parallel Retrieval Strategies
- **Status**: 🔴 Not Started
- **Current State**: Sequential execution of vector/keyword/hybrid search
- **Requirements**:
  - Goroutine-based concurrent retrieval
  - Result aggregation with timeouts
  - Error handling for partial failures
  - Configuration for max concurrency
- **Benefits**:
  - 2-3x faster retrieval for hybrid search
  - Reduced query latency
  - Better resource utilization
- **Estimate**: 1-2 days

#### 2. Streaming Responses
- **Status**: 🔴 Not Started
- **Requirements**:
  - Stream LLM responses to CLI in real-time
  - Progressive result display during workflow execution
  - Status updates for long-running operations
  - Graceful handling of interruptions
- **Benefits**:
  - Better UX for long queries (10-30s responses)
  - Perceived performance improvement
  - Early cancellation capability
- **Estimate**: 2-3 days

#### 3. Schema Registry Persistence
- **Status**: 🔴 Not Started
- **Current State**: In-memory schema pattern storage
- **Requirements**:
  - Persistent storage backend (SQLite, PostgreSQL, or JSON files)
  - Schema versioning and migration
  - CRUD API for schema management
  - Pattern matching query optimization
- **Benefits**:
  - Preserve schema patterns across restarts
  - Share schemas across instances
  - Schema evolution tracking
- **Estimate**: 3-4 days

### User Experience

#### 1. Interactive CLI Improvements
- **Status**: 🔴 Not Started
- **Requirements**:
  - Query history and navigation
  - Autocomplete for commands
  - Syntax highlighting for output
  - Progress bars for long operations
  - Rich formatting with colors and tables
- **Libraries**: bubbletea, lipgloss, cobra enhancements
- **Estimate**: 2-3 days

#### 2. Configuration Validation
- **Status**: 🔴 Not Started
- **Requirements**:
  - Comprehensive validation on CLI startup
  - Helpful error messages for misconfigurations
  - Config schema documentation
  - Examples for common setups
- **Estimate**: 1-2 days

#### 3. Query Templates
- **Status**: 🔴 Not Started
- **Requirements**:
  - Pre-defined query patterns for common use cases
  - Template variables for document-specific queries
  - Template library management
  - CLI command for template listing and execution
- **Benefits**: Easier onboarding, consistent results
- **Estimate**: 2-3 days

---

## Technical Debt

### Code Quality

#### 1. Error Handling Consistency
- **Status**: 🟡 Mostly Good
- **Issues**:
  - Some places use fmt.Errorf, others use errors.New
  - Error wrapping not always consistent
  - Need custom error types for better handling
- **Action**: Audit error handling patterns, define standards
- **Estimate**: 1-2 days

#### 2. Configuration Validation
- **Status**: 🟡 Basic Validation
- **Issues**:
  - Runtime failures for missing config rather than startup validation
  - No schema validation for config.json
  - Environment variable precedence unclear
- **Action**: Add config validation package
- **Estimate**: 1-2 days

### Documentation

#### 1. API Documentation
- **Status**: 🟡 Partial
- **Missing**:
  - Godoc coverage for all exported functions (currently ~70%)
  - Package-level documentation for some packages
  - Examples in godoc
- **Action**: Complete godoc coverage, add pkg.go.dev examples
- **Estimate**: 2-3 days

#### 2. Architecture Decision Records (ADRs)
- **Status**: 🔴 Not Started
- **Need**:
  - Document key architectural decisions (schema-driven approach, agent specialization, RRF fusion)
  - Rationale for technology choices
  - Trade-offs and alternatives considered
- **Action**: Create `docs/adr/` directory with ADRs
- **Estimate**: 2-3 days

---

## Completed ✅

### Phase 1: Foundation (2025-10-27)
- ✅ LLM provider abstraction with OpenAI implementation
- ✅ Embedding interface with OpenAI embedder
- ✅ Vector store interface with Qdrant implementation
- ✅ Document parser interface with text/markdown parsers
- ✅ Configuration management
- ✅ Test coverage: 90%+ for all Phase 1 packages

### Phase 2: Schema System (2025-10-27)
- ✅ LLM-based schema analyzer
- ✅ Schema resolver with multiple strategies
- ✅ Schema registry for pattern storage
- ✅ Metadata generation helpers
- ✅ Schema-aware chunking (section, hierarchical, semantic, sliding window)
- ✅ PDF and HTML parser implementations
- ✅ Test coverage: 96.5% for schema package

### Phase 3: Agents & Retrieval (2025-10-27)
- ✅ All 8 specialized agents (planner, rewriter, supervisor, retriever, reranker, distiller, reflector, policy)
- ✅ Vector search implementation
- ✅ Keyword search (BM25) implementation
- ✅ Hybrid search with RRF fusion
- ✅ Schema-filtered retrieval
- ✅ Test coverage: 91.2% for agents, 95%+ for retrieval

### Phase 4: Workflow Execution (2025-10-27)
- ✅ Workflow graph construction
- ✅ Execution engine with state management
- ✅ Node wrappers for all agents
- ✅ Deep thinking loop implementation
- ✅ Test coverage: 89.9% for workflow

### Phase 5: CLI Interface (2025-10-27)
- ✅ Interactive query mode
- ✅ Single-shot query mode
- ✅ Document ingestion (recursive directories, multiple formats)
- ✅ Configuration management commands
- ✅ Example scripts and documentation
- ✅ Cost transparency documentation

### Documentation (2025-10-27 - 2025-11-01)
- ✅ Comprehensive CLAUDE.md (695 lines)
- ✅ Detailed README with architecture
- ✅ SETUP.md with installation guide
- ✅ CLEANUP.md for maintenance
- ✅ Examples with 6 runnable scripts
- ✅ Multi-tool support (Droid + Claude Code)

---

## Notes

### Prioritization Philosophy

1. **High Priority**: Items blocking production deployment or representing significant quality gaps
2. **Medium Priority**: Planned architectural features that expand capabilities
3. **Low Priority**: Enhancements that improve operations but aren't blocking

### Estimation Approach

Estimates assume:
- Senior-level Go developer
- Familiarity with the existing codebase
- Includes implementation, testing, and documentation
- Does not include extensive exploratory research

### Contributing

When working on items from this TODO:

1. Create a feature/bugfix/chore branch per gitflow conventions
2. Update this file to mark items as 🟡 In Progress
3. Ensure 90%+ test coverage for new code
4. Update relevant documentation (CLAUDE.md, README.md)
5. Move completed items to the "Completed" section with date
6. Follow all standards in CLAUDE.md

### Review Cycle

This TODO should be reviewed and updated:
- After completing each major feature
- Monthly for priority adjustments
- When new requirements emerge
- Before planning new development sprints
