# Changelog

All notable changes to Deep Thinking Agent will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- `VectorStore.List()` method for efficient document enumeration without vector similarity
- Detailed test coverage status section in README
- Known limitations documentation for BM25 and integration testing
- Branch protection setup guide (BRANCH_PROTECTION_GUIDE.md)
- **Pre-commit hook** for automated quality checks (gofmt, build, test, vet)
- Pre-commit hook setup documentation (PRE_COMMIT_HOOK_SETUP.md)

### Changed
- **BREAKING**: `VectorStore` interface now requires `List()` method implementation
- BM25 KeywordRetriever now uses `List()` instead of dummy vector workaround
- README coverage claims updated from "88% production-ready" to "Production-Grade Core: 7 packages with 90%+"
- Updated project description from "production-ready" to "production-grade architecture with strong reference implementation"

### Fixed
- Formatting violations in 3 test files (gofmt compliance)
- Security badge now links to SECURITY.md instead of non-configured Snyk service
- BM25 architectural issue (no longer requires dummy embedding vectors)
- **Compilation errors** from VectorStore.List() addition (Qdrant pointer type, agent mock)

### Documentation
- Added honest test coverage breakdown by package
- Documented which packages are production-ready (90%+) vs need hardening
- Added transparency note about integration layer testing status
- Clarified BM25 in-memory limitation for large datasets
- Documented the build failure incident and prevention measures

---

## History

### 2025-11-01 - GitHub Publication
- Initial public release on GitHub
- Repository URL: https://github.com/gerrymiller/deep-thinking-agent
- Complete CI/CD workflows (test, security, release, labeler)
- Comprehensive documentation and examples
- 7 core packages with 90%+ test coverage

### 2025-10-28 - Complete Reference Implementation
- All 5 phases implemented and working
- 8 specialized agents operational
- Real BM25 keyword search
- Schema-aware chunking
- Full workflow orchestration
- 88% coverage for tested packages

### Prior Development
- Phase 1: Foundation (LLM, embedding, vector store interfaces)
- Phase 2: Schema system (analyzer, resolver, registry, chunking)
- Phase 3: Agents & Retrieval (8 agents, 3 retrieval strategies)
- Phase 4: Workflow Execution (graph-based orchestration)
- Phase 5: CLI Interface (full-featured command-line tool)

---

## Migration Notes

### List() Method Addition

If you have custom `VectorStore` implementations, you must add the `List()` method:

```go
func (s *Store) List(ctx context.Context, collectionName string, filter vectorstore.Filter, limit int, offset int) ([]vectorstore.Document, error) {
    // Implementation that retrieves documents without vector similarity
    // Should support:
    // - Optional metadata filtering
    // - Pagination via limit/offset
    // - Return documents with embeddings and metadata
}
```

**For Qdrant users**: The provided implementation uses the `Scroll` API.

**For other vector stores**: Implement using equivalent list/scan operations available in your store.

---

## Versioning Strategy

- **0.x.x** - Pre-1.0 development, API may change
- **1.0.0** - First stable release with frozen API
- **1.x.x** - Backward-compatible additions
- **2.x.x** - Breaking changes

Current status: **0.x.x** (development/reference implementation phase)

---

## Links

- **Repository**: https://github.com/gerrymiller/deep-thinking-agent
- **Issues**: https://github.com/gerrymiller/deep-thinking-agent/issues
- **Documentation**: See README.md, AGENTS.md, CONTRIBUTING.md
- **Examples**: See examples/ directory
