# Contributing to Deep Thinking Agent

Thank you for your interest in contributing to Deep Thinking Agent! This document provides guidelines and instructions for contributing to the project.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [How Can I Contribute?](#how-can-i-contribute)
- [Development Setup](#development-setup)
- [Development Workflow](#development-workflow)
- [Code Standards](#code-standards)
- [Testing Requirements](#testing-requirements)
- [Pull Request Process](#pull-request-process)
- [Issue Guidelines](#issue-guidelines)
- [Community and Communication](#community-and-communication)

## Code of Conduct

This project adheres to a Code of Conduct that all contributors are expected to follow. Please read [CODE_OF_CONDUCT.md](CODE_OF_CONDUCT.md) before contributing.

## How Can I Contribute?

### Reporting Bugs

Before creating bug reports, please check existing issues to avoid duplicates. When creating a bug report, include:

- **Clear descriptive title**
- **Steps to reproduce** the behavior
- **Expected vs actual behavior**
- **Environment details** (Go version, OS, dependencies)
- **Relevant logs or screenshots**
- **Possible solutions** if you have ideas

Use the bug report template when creating issues.

### Suggesting Enhancements

Enhancement suggestions are tracked as GitHub issues. When creating an enhancement suggestion:

- **Use a clear descriptive title**
- **Provide detailed description** of the proposed functionality
- **Explain why this enhancement would be useful** to most users
- **List any similar features** in other projects
- **Consider implementation approach** if you have ideas

Check [TODO.md](TODO.md) for planned features before suggesting new ones.

### Your First Code Contribution

Unsure where to start? Look for issues labeled:

- `good first issue` - Simple issues perfect for beginners
- `help wanted` - Issues where maintainers need assistance
- `documentation` - Documentation improvements
- `tests` - Test coverage improvements

### Pull Requests

We actively welcome pull requests for:

- Bug fixes
- Documentation improvements
- Test coverage improvements
- Performance optimizations
- New features (after discussion in issues)

## Development Setup

### Prerequisites

Ensure you have:

- **Go 1.25.3+** - [Installation guide](https://go.dev/doc/install)
- **Docker** - For running Qdrant locally
- **Git** - For version control
- **OpenAI API key** - For testing (or use mocks)

### Setup Steps

1. **Fork the repository** on GitHub

2. **Clone your fork**:
   ```bash
   git clone https://github.com/YOUR_USERNAME/deep-thinking-agent.git
   cd deep-thinking-agent
   ```

3. **Add upstream remote**:
   ```bash
   git remote add upstream https://github.com/gerrymiller/deep-thinking-agent.git
   ```

4. **Install dependencies**:
   ```bash
   go mod download
   go mod verify
   ```

5. **Start Qdrant** (for integration tests):
   ```bash
   docker run -p 6333:6333 -p 6334:6334 qdrant/qdrant
   ```

6. **Run tests** to verify setup:
   ```bash
   go test ./...
   ```

For detailed setup instructions, see [SETUP.md](SETUP.md).

## Development Workflow

This project follows **gitflow** branching strategy:

### Branch Types

- `main` - Production-ready releases (protected)
- `develop` - Integration branch for features (protected)
- `feature/*` - New features and enhancements
- `bugfix/*` - Bug fixes during development
- `hotfix/*` - Critical production fixes
- `chore/*` - Maintenance, refactoring, dependencies

### Creating a Branch

1. **Sync with upstream**:
   ```bash
   git checkout develop
   git pull upstream develop
   ```

2. **Create your branch**:
   ```bash
   # For features
   git checkout -b feature/your-feature-name
   
   # For bug fixes
   git checkout -b bugfix/issue-123-description
   
   # For chores
   git checkout -b chore/update-dependencies
   ```

3. **Make your changes** following code standards

4. **Commit frequently** with clear messages:
   ```bash
   git add .
   git commit -m "Add schema validation for custom models"
   ```

### Commit Message Format

Follow conventional commits style:

```
<type>: <description>

[optional body]

[optional footer]
```

**Types**:
- `feat:` New feature
- `fix:` Bug fix
- `docs:` Documentation changes
- `test:` Test additions or updates
- `refactor:` Code refactoring
- `perf:` Performance improvements
- `chore:` Maintenance tasks

**Examples**:
```
feat: Add Anthropic LLM provider support

Implements the llm.Provider interface for Claude models with
support for thinking blocks and tool use.

Closes #42
```

```
fix: Handle nil pointer in schema resolver

Adds nil check before accessing schema fields to prevent panic
when processing documents without metadata.

Fixes #128
```

### Keeping Your Branch Updated

```bash
# Fetch latest changes
git fetch upstream

# Rebase your branch on develop
git checkout your-branch-name
git rebase upstream/develop

# Resolve conflicts if any, then continue
git add .
git rebase --continue
```

## Code Standards

All contributions must adhere to project standards. See [AGENTS.md](AGENTS.md) for complete details.

### File Headers

Every new `.go` file must include:

```go
// Copyright 2025 Gerry Miller <gerry@gerrymiller.com>
//
// Licensed under the MIT License.
// See LICENSE file in the project root for full license information.
```

### Code Quality

- **Run `go fmt`** before committing:
  ```bash
  go fmt ./...
  ```

- **Run `go vet`** to catch common issues:
  ```bash
  go vet ./...
  ```

- **Run linter** if available:
  ```bash
  golangci-lint run
  ```

- **Follow Go conventions**:
  - [Effective Go](https://go.dev/doc/effective_go)
  - [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)

### Documentation

- **Document all exported types and functions** with godoc comments
- **Update README.md** for user-facing changes
- **Update AGENTS.md** for architectural changes
- **Add examples** for complex functionality

Example:

```go
// Analyzer uses an LLM to derive document schemas.
// It analyzes document structure to identify sections, hierarchy,
// semantic regions, and recommend chunking strategies.
type Analyzer struct {
    llmProvider llm.Provider
    temperature float32
    maxTokens   int
    timeout     time.Duration
}

// AnalyzeDocument performs LLM-based analysis of a document.
// Returns a DocumentSchema with identified structure and metadata.
//
// Example:
//   schema, err := analyzer.AnalyzeDocument(ctx, "doc-1", content, "markdown")
//   if err != nil {
//       return err
//   }
func (a *Analyzer) AnalyzeDocument(ctx context.Context, docID, content, format string) (*DocumentSchema, error) {
    // Implementation
}
```

## Testing Requirements

**CRITICAL**: All code contributions must include tests.

### Requirements

- **90% minimum test coverage** for new packages
- **Unit tests** for all exported functions
- **Table-driven tests** for multiple scenarios
- **Error case testing** for all error paths
- **Edge case testing** for boundary conditions

### Running Tests

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run specific package
go test ./pkg/schema

# Run with verbose output
go test -v ./...

# Run specific test
go test -run TestAnalyzer ./pkg/schema
```

### Test Coverage Check

Before submitting PR:

```bash
go test -cover ./...
```

All modified packages should show ≥90% coverage.

### Writing Tests

Use table-driven tests:

```go
func TestAnalyzer_BuildPrompt(t *testing.T) {
    tests := []struct {
        name    string
        content string
        format  string
        want    string
        wantErr bool
    }{
        {
            name:    "valid markdown",
            content: "# Title\n\nContent",
            format:  "markdown",
            want:    "Analyze the following markdown document...",
            wantErr: false,
        },
        {
            name:    "empty content",
            content: "",
            format:  "text",
            want:    "",
            wantErr: true,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := buildPrompt(tt.content, tt.format)
            if (err != nil) != tt.wantErr {
                t.Errorf("buildPrompt() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if got != tt.want {
                t.Errorf("buildPrompt() = %v, want %v", got, tt.want)
            }
        })
    }
}
```

## Pull Request Process

### Before Submitting

1. **Ensure all tests pass**:
   ```bash
   go test ./...
   ```

2. **Check test coverage**:
   ```bash
   go test -cover ./...
   ```

3. **Format code**:
   ```bash
   go fmt ./...
   ```

4. **Run linter**:
   ```bash
   go vet ./...
   ```

5. **Update documentation** if needed

6. **Add entry to TODO.md** if completing a tracked item

### Submitting Pull Request

1. **Push to your fork**:
   ```bash
   git push origin your-branch-name
   ```

2. **Create Pull Request** on GitHub:
   - Target the `develop` branch (not `main`)
   - Use the PR template
   - Reference related issues (e.g., "Closes #123")
   - Provide clear description of changes

3. **Complete the PR checklist**:
   - [ ] Tests added/updated and passing
   - [ ] Documentation updated
   - [ ] Code formatted and linted
   - [ ] Branch rebased on latest develop
   - [ ] Reviewed own changes

### PR Review Process

- **Maintainers will review** within 3-5 business days
- **Address feedback** by pushing new commits
- **Keep discussions professional** and constructive
- **Be patient** - quality takes time

### After Approval

- Maintainers will merge using **squash and merge**
- Your contribution will be credited in release notes
- Branch will be deleted automatically

## Issue Guidelines

### Creating Issues

- **Search existing issues** first to avoid duplicates
- **Use appropriate template** (bug report, feature request, question)
- **Provide complete information** - incomplete issues may be closed
- **Use clear, descriptive titles**
- **Add relevant labels** if you have permissions

### Issue Labels

- `bug` - Something isn't working
- `enhancement` - New feature or request
- `documentation` - Improvements or additions to docs
- `good first issue` - Good for newcomers
- `help wanted` - Extra attention is needed
- `question` - Further information is requested
- `wontfix` - This will not be worked on
- `duplicate` - This issue already exists
- `priority:high` - High priority items
- `priority:medium` - Medium priority items
- `priority:low` - Low priority items

## Community and Communication

### Getting Help

- **GitHub Issues** - For bugs, features, and questions
- **Discussions** - For general questions and ideas
- **Documentation** - Check [README.md](README.md), [SETUP.md](SETUP.md), and [AGENTS.md](AGENTS.md)

### Code Review Philosophy

Reviews focus on:

- **Correctness** - Does it work as intended?
- **Test coverage** - Are edge cases tested?
- **Readability** - Is code clear and maintainable?
- **Performance** - Are there obvious optimizations?
- **Documentation** - Is functionality documented?

We value:

- **Constructive feedback** - Focus on code, not people
- **Learning opportunities** - Reviews are teaching moments
- **Respectful disagreement** - Different perspectives are valuable

### Recognition

Contributors are recognized through:

- **GitHub contributors page**
- **Release notes** - Credits for significant contributions
- **README acknowledgments** - For major features

## License

By contributing, you agree that your contributions will be licensed under the MIT License.

## Questions?

If you have questions not covered here:

1. Check existing documentation
2. Search closed issues for similar questions
3. Open a new issue with the `question` label
4. Be specific about what you need help with

Thank you for contributing to Deep Thinking Agent! 🚀
