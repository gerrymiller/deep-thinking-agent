# AI Safety Checklist for Deep Thinking Agent

## Purpose

This checklist ensures AI assistants (Claude Code, Droid, etc.) follow safe git workflows and never commit directly to protected branches.

---

## Before EVERY Commit

### Step 1: Verify Current Branch

Run this command and check the output:

```bash
git rev-parse --abbrev-ref HEAD
```

**Decision tree:**
- ✅ Output is `develop` → Safe to commit
- ✅ Output starts with `feature/` → Safe to commit
- ✅ Output starts with `bugfix/` → Safe to commit
- ✅ Output starts with `chore/` → Safe to commit
- ✅ Output starts with `hotfix/` → Safe to commit (but must merge to both main and develop)
- ❌ Output is `main` → **STOP! DO NOT COMMIT!**
- ❌ Output is `origin/main` → **STOP! Detached HEAD state!**

### Step 2: If on Main Branch (ERROR State)

**You should NEVER be on main when committing. Fix this immediately:**

```bash
# DO NOT commit! Switch to proper workflow:
git checkout develop
git pull origin develop
git checkout -b feature/your-feature-name

# Now you can commit safely
```

### Step 3: Run Quality Checks

The pre-commit hook will automatically run these checks, but you can run them manually:

```bash
# Format check
gofmt -l .

# Build check
go build ./...

# Test check
go test ./...

# Vet check
go vet ./...
```

### Step 4: Commit Safely

```bash
git add .
git commit -m "type: description"
```

The pre-commit hook will:
- ✅ Block commits to `main` automatically
- ✅ Check code formatting
- ✅ Verify build succeeds
- ✅ Run all tests
- ✅ Run go vet

---

## Required Workflow Patterns

### Starting New Work

```bash
# 1. Ensure develop is up to date
git checkout develop
git pull origin develop

# 2. Create appropriate branch type
git checkout -b feature/add-new-feature    # For new features
git checkout -b bugfix/fix-parser-error    # For bug fixes
git checkout -b chore/update-dependencies  # For maintenance
git checkout -b hotfix/security-patch      # For production fixes (from main)

# 3. Verify you're on the new branch
git rev-parse --abbrev-ref HEAD

# 4. Make your changes and commit
```

### Making Commits

```bash
# ALWAYS verify branch first!
echo "Current branch: $(git rev-parse --abbrev-ref HEAD)"

# If not on main, proceed:
git add .
git commit -m "feat: add schema validation"

# Pre-commit hook will verify quality
```

### Completing Work

```bash
# Push your feature branch to remote
git push origin feature/your-feature-name

# Create pull request on GitHub:
# feature → develop → main
# OR
# hotfix → main AND develop

# NEVER merge directly to main without PR
```

---

## Never Do This

### ❌ Direct Commit to Main

```bash
# WRONG - Will be blocked by pre-commit hook
git checkout main
git commit -m "Quick fix"  # ❌ FORBIDDEN
```

### ❌ Commit Without Branch Check

```bash
# WRONG - Don't assume you're on the right branch
git commit -m "Changes"  # ❌ Check branch first!
```

### ❌ Bypass Pre-Commit Hook

```bash
# WRONG - Only for absolute emergencies
git commit --no-verify  # ❌ Bypasses all safety checks
```

### ❌ Force Push to Protected Branches

```bash
# WRONG - Destructive and blocked by GitHub
git push --force origin main  # ❌ FORBIDDEN
```

---

## Emergency Procedures

### If You Accidentally Committed to Main

**STOP! Don't push!**

If you committed to main locally but haven't pushed:

```bash
# Option 1: Move commit to a new branch
git branch feature/accidental-commit
git reset --hard origin/main
git checkout feature/accidental-commit
# Now push the feature branch and create PR

# Option 2: Reset to remote main
git reset --hard origin/main
# Your changes are lost - recommit on proper branch
```

### If You Pushed to Main (Branch Protection Should Block This)

If GitHub branch protection is properly configured, the push will be rejected. If it somehow got through:

```bash
# Create revert commit
git checkout main
git revert HEAD
git push origin main

# Create proper feature branch with original work
git checkout develop
git checkout -b feature/proper-branch
git cherry-pick <original-commit-sha>
git push origin feature/proper-branch
# Create PR: feature → develop
```

---

## Defense Layers

This repository has 5 layers of defense against direct main commits:

### Layer 1: Pre-Commit Hook (Local)
- File: `.git/hooks/pre-commit`
- Blocks commits to `main` before they happen
- Runs format, build, test, vet checks
- Can be bypassed with `--no-verify` (don't do this)

### Layer 2: GitHub Branch Protection (Remote)
- Configured at: https://github.com/gerrymiller/deep-thinking-agent/settings/branches
- Blocks direct pushes to `main`
- Requires pull requests with approval
- Cannot be bypassed (even by owner if configured properly)

### Layer 3: Documentation (Education)
- This file: `AI_SAFETY_CHECKLIST.md`
- `AGENTS.md` - Critical warning section at top
- `BRANCH_PROTECTION_GUIDE.md` - Setup instructions
- `PRE_COMMIT_HOOK_SETUP.md` - Hook documentation

### Layer 4: Git Configuration (Visual)
- Branch name highlighted in git output
- Status shows current branch prominently
- Helps maintain branch awareness

### Layer 5: Code Review Process
- All changes require pull requests
- Changes reviewed before merging
- CI/CD checks must pass
- Conversation resolution required

---

## Verification Commands

### Quick Safety Check

Run before any commit:

```bash
# One-liner safety check
BRANCH=$(git rev-parse --abbrev-ref HEAD) && \
echo "Current branch: $BRANCH" && \
if [ "$BRANCH" = "main" ]; then echo "❌ DO NOT COMMIT TO MAIN"; exit 1; else echo "✅ Safe to commit"; fi
```

### Full Status Check

```bash
# Comprehensive status
echo "=== Git Status ==="
git status

echo -e "\n=== Current Branch ==="
git rev-parse --abbrev-ref HEAD

echo -e "\n=== Branch Protection Status ==="
if [ "$(git rev-parse --abbrev-ref HEAD)" = "main" ]; then
    echo "❌ WARNING: You are on main branch!"
else
    echo "✅ Safe branch for commits"
fi

echo -e "\n=== Recent Commits ==="
git log --oneline -5

echo -e "\n=== Quality Checks ==="
gofmt -l . && echo "✅ Format OK" || echo "❌ Format issues"
go build ./... &>/dev/null && echo "✅ Build OK" || echo "❌ Build failed"
go test ./... &>/dev/null && echo "✅ Tests OK" || echo "❌ Tests failed"
```

---

## Incident History

### 2025-11-08: Direct Commits to Main

**What Happened:**
- AI assistant made 3 commits directly to `main` branch
- Commits: b5f03c8, de46e2b, f5bfc35
- Violated gitflow workflow

**Impact:**
- Lost a day of development work
- Created orphaned commits (886cb12, 3c47bf8)
- Required manual recovery and revert
- Broke synchronization between main and develop

**Root Cause:**
- Failed to run `git rev-parse --abbrev-ref HEAD` before committing
- No pre-commit hook to block main commits (at the time)
- AI assistant didn't follow documented gitflow workflow

**Resolution:**
1. Reverted incorrect commit (f5bfc35)
2. Synchronized main and develop branches
3. Implemented this safety checklist
4. Enhanced pre-commit hook with branch check
5. Updated AGENTS.md with critical warning section

**Lessons Learned:**
- **ALWAYS** verify current branch before committing
- **NEVER** assume you're on the right branch
- **NEVER** commit directly to main, even for "quick fixes"
- Automation (pre-commit hooks) prevents human/AI error
- Documentation must be prominent and impossible to miss

**Prevention:**
- This checklist created
- Pre-commit hook enhanced with branch check
- AGENTS.md updated with critical warning at top
- Git configuration enhanced
- 5 layers of defense implemented

---

## Quick Reference

### Safe Commit Workflow

```bash
# 1. Check branch
git rev-parse --abbrev-ref HEAD

# 2. If on main, switch
git checkout develop && git checkout -b feature/name

# 3. Make changes
# ... edit files ...

# 4. Commit (hook will verify)
git add .
git commit -m "type: description"

# 5. Push feature branch
git push origin feature/name

# 6. Create PR on GitHub
```

### Branch Type Prefixes

- `feature/` - New features (from develop)
- `bugfix/` - Bug fixes (from develop)
- `chore/` - Maintenance, refactoring, dependencies (from develop)
- `hotfix/` - Critical production fixes (from main, merge to both main and develop)
- `release/` - Release preparation (from develop, merge to both main and develop)

### Commit Message Format

```
type: description

Optional longer explanation

Co-authored-by: Gerry Miller <gerry@gerrymiller.com>
```

**NEVER include:**
- `Co-authored-by: Claude` or any AI assistant attribution
- `Co-authored-by: factory-droid[bot]` or similar

All work is authored by Gerry Miller. AI assistants are tools, not co-authors.

---

## Testing This Checklist

### Test 1: Verify Pre-Commit Hook Blocks Main

```bash
# Should be blocked:
git checkout main
echo "test" >> /tmp/test.txt
git add /tmp/test.txt
git commit -m "test: Should be blocked"

# Expected: ❌ BLOCKED: Direct commits to 'main' branch are FORBIDDEN!
```

### Test 2: Verify Feature Branch Works

```bash
# Should succeed:
git checkout develop
git checkout -b test/safety-check
echo "# Test" >> /tmp/test.txt
git add /tmp/test.txt
git commit -m "test: Feature branch commit"

# Expected: ✅ All pre-commit checks passed!

# Cleanup:
git checkout develop
git branch -D test/safety-check
rm /tmp/test.txt
```

---

## Summary

**Before EVERY commit:**
1. Run: `git rev-parse --abbrev-ref HEAD`
2. Verify output is NOT `main`
3. If on `main`, switch to feature branch
4. Commit with confidence - hook will verify quality

**The pre-commit hook is your safety net, but YOU must verify the branch first.**

**Remember: One check prevents hours of recovery work.**
