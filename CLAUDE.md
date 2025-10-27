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

## Code Structure

The project is currently in early development. As the codebase grows, this section will be updated with architectural details.
