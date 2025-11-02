# Pre-Commit Hook Setup

## What Happened

On 2025-11-01, compilation errors were introduced and committed without running build/test verification first. This violated the testing workflow documented in AGENTS.md and temporarily broke the build.

**Errors**:
1. Qdrant store: Type mismatch (`uint32` vs `*uint32`)
2. Agent tests: Missing `List()` method in mock

**Root Cause**: Failed to run `go build ./...` and `go test ./...` before committing.

## Solution: Automated Pre-Commit Hook

A pre-commit hook has been installed at `.git/hooks/pre-commit` that automatically runs before every commit:

### What It Checks

1. ✅ **Format** - Runs `gofmt -l .` to catch formatting violations
2. ✅ **Build** - Runs `go build ./...` to catch compilation errors  
3. ✅ **Tests** - Runs `go test ./...` to catch test failures
4. ✅ **Vet** - Runs `go vet ./...` to catch common mistakes

### How It Works

```bash
# When you try to commit:
git commit -m "Your message"

# Hook runs automatically:
🔍 Running pre-commit checks...
  → Checking gofmt...
  ✅ Format check passed
  → Checking build...
  ✅ Build check passed
  → Running tests...
  ✅ Tests passed
  → Running go vet...
  ✅ go vet passed
✨ All pre-commit checks passed! Proceeding with commit...

# If any check fails, commit is blocked with helpful error message
```

### Benefits

- **Prevents** broken commits from entering history
- **Catches** compilation errors immediately
- **Enforces** code quality standards
- **Matches** the workflow in AGENTS.md
- **Saves time** by catching issues before push/PR

## Installation for Contributors

The hook is already installed in this repository's `.git/hooks/` directory and is executable.

**For contributors cloning the repo**, the hook needs to be installed manually (git hooks don't transfer via clone):

```bash
# After cloning, copy the hook:
cp .git/hooks/pre-commit.sample .git/hooks/pre-commit

# Or create it fresh:
cat > .git/hooks/pre-commit << 'EOF'
#!/bin/sh
# See PRE_COMMIT_HOOK_SETUP.md for full content
echo "🔍 Running pre-commit checks..."
# ... full script ...
EOF

chmod +x .git/hooks/pre-commit
```

**Note**: We cannot commit hooks to git (they don't transfer), but we document them here.

## Bypassing the Hook (Emergency Only)

If you absolutely must commit without running checks (not recommended):

```bash
git commit --no-verify -m "Emergency commit"
```

**Use sparingly** - only for:
- Documentation-only changes
- Urgent hotfixes where you'll fix tests in next commit
- Work-in-progress branches

## Manual Verification (Alternative)

If you prefer not to use hooks, run manually before committing:

```bash
# Quick check
go build ./... && go test ./... && gofmt -l .

# Full check (matches pre-commit hook)
gofmt -l . && \
go build ./... && \
go test ./... && \
go vet ./...
```

Add this to your workflow or create a shell alias:

```bash
# Add to ~/.bashrc or ~/.zshrc
alias go-verify='gofmt -l . && go build ./... && go test ./... && go vet ./...'

# Then run before commits:
go-verify && git commit -m "Message"
```

## CI/CD Catches Too

Even if local checks are bypassed, our GitHub Actions CI workflow will catch issues:
- `.github/workflows/ci.yml` runs on every push/PR
- Includes same checks: format, build, test, lint
- Blocks merging if checks fail

**Defense in depth**: Local hooks + CI/CD = double protection

## Lesson Learned

**From AGENTS.md Testing Requirements**:
```
### Completing Work
# Ensure code is tested and formatted
go test ./...
go fmt ./...
go vet ./...

# Commit changes
git add .
git commit -m "Description of changes"
```

**We now enforce this automatically with the pre-commit hook.**

## Statistics

**Issue**: 2 compilation errors in 2 files  
**Detection**: After commit (too late)  
**Fix Time**: 5 minutes  
**Prevention**: Pre-commit hook (automated)  
**Future Risk**: Eliminated by automation  

## Related Files

- `.git/hooks/pre-commit` - The actual hook script (local only)
- `AGENTS.md` - Testing workflow documentation
- `.github/workflows/ci.yml` - CI/CD verification
- `CHANGELOG.md` - Records the incident and fix

---

**Acknowledgment**: Thanks to Claude Code for catching this and providing clear feedback about the build failure.
