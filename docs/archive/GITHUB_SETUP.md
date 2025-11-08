# GitHub Repository Setup Guide

This document provides step-by-step instructions for configuring the Deep Thinking Agent repository on GitHub with all recommended settings, branch protection rules, and integrations.

## Table of Contents

- [Prerequisites](#prerequisites)
- [Initial Repository Setup](#initial-repository-setup)
- [Branch Protection Rules](#branch-protection-rules)
- [Repository Settings](#repository-settings)
- [Security Configuration](#security-configuration)
- [Third-Party Integrations](#third-party-integrations)
- [Labels and Automations](#labels-and-automations)
- [Verification](#verification)

## Prerequisites

Before starting, ensure you have:

- [x] GitHub account with admin access to the repository
- [x] Repository created (or ready to create)
- [x] All code committed to `main` and `develop` branches
- [x] OpenAI API key (for testing)
- [ ] Snyk account (sign up at https://snyk.io/)
- [ ] Codecov account (optional, sign up at https://codecov.io/)

## Initial Repository Setup

### 1. Create/Configure Repository

If creating a new repository:

```bash
# On GitHub.com
1. Go to https://github.com/new
2. Repository name: deep-thinking-agent
3. Description: Schema-Driven RAG with Iterative Multi-Hop Reasoning
4. Visibility: Public
5. Initialize: Do NOT initialize (we have existing code)
6. Create repository
```

If repository already exists locally without remote:

```bash
# Add remote
git remote add origin https://github.com/gerrymiller/deep-thinking-agent.git

# Push main branch
git push -u origin main

# Push develop branch
git checkout develop
git push -u origin develop

# Push all tags
git push --tags
```

### 2. Repository Configuration

Navigate to **Settings** in your GitHub repository.

#### General Settings

**Features:**
- ✅ Issues
- ✅ Projects
- ✅ Discussions (for Q&A and community)
- ❌ Wikis (we use docs/ folder instead)
- ✅ Sponsorships (optional)

**Pull Requests:**
- ✅ Allow squash merging (with custom message)
- ❌ Allow merge commits
- ✅ Allow rebase merging
- ✅ Always suggest updating pull request branches
- ✅ Automatically delete head branches

**Archives:**
- ❌ Do not archive this repository

**Default Branch:**
- Set to: `develop`

#### Topics (Repository Settings)

Add relevant topics for discoverability:
```
go, golang, rag, llm, ai, vector-database, qdrant, openai, 
deep-thinking, schema-driven, multi-hop-reasoning, agentic-ai,
retrieval-augmented-generation, machine-learning
```

#### Description

Set repository description:
```
Schema-Driven RAG with Iterative Multi-Hop Reasoning • 8 Specialized Agents • Production-Ready Go Implementation
```

#### Website

Set to:
```
https://github.com/gerrymiller/deep-thinking-agent
```

## Branch Protection Rules

### Main Branch Protection

Navigate to **Settings → Branches → Add branch protection rule**

**Branch name pattern:** `main`

**Protection Rules:**

- [x] **Require a pull request before merging**
  - Required approvals: `1`
  - [x] Dismiss stale pull request approvals when new commits are pushed
  - [x] Require review from Code Owners
  - [x] Require approval of the most recent reviewable push

- [x] **Require status checks to pass before merging**
  - [x] Require branches to be up to date before merging
  - **Required status checks:** (add these after first workflow runs)
    - `test`
    - `lint`
    - `format`
    - `coverage`
    - `CodeQL`
    - `Snyk Security Scan`

- [x] **Require conversation resolution before merging**

- [x] **Require signed commits**

- [x] **Require linear history**

- [x] **Include administrators** (for solo development)
  - ❌ **Uncheck this for team development**

- [x] **Restrict who can push to matching branches**
  - Add: Repository admins/maintainers only

- [x] **Allow force pushes** → ❌ Disable

- [x] **Allow deletions** → ❌ Disable

### Develop Branch Protection

**Branch name pattern:** `develop`

**Protection Rules:**

- [x] **Require a pull request before merging**
  - Required approvals: `0` (for solo) or `1` (for team)
  - [x] Dismiss stale pull request approvals when new commits are pushed

- [x] **Require status checks to pass before merging**
  - [x] Require branches to be up to date before merging
  - **Required status checks:**
    - `test`
    - `lint`
    - `format`
    - `coverage`

- [x] **Require conversation resolution before merging**

- [x] **Require signed commits**

- [x] **Allow force pushes** → ❌ Disable

- [x] **Allow deletions** → ❌ Disable

## Repository Settings

### Code Security and Analysis

Navigate to **Settings → Code security and analysis**

**Security:**

- [x] **Dependency graph** → Enable
- [x] **Dependabot alerts** → Enable
- [x] **Dependabot security updates** → Enable
- [x] **Code scanning** → Set up CodeQL (already configured in workflows)
- [x] **Secret scanning** → Enable
- [x] **Secret scanning push protection** → Enable

**Private vulnerability reporting:**
- [x] Enable (allows security researchers to privately report vulnerabilities)

## Security Configuration

### GitHub Secrets

Navigate to **Settings → Secrets and variables → Actions**

Add the following secrets:

1. **SNYK_TOKEN**
   - Get from: https://app.snyk.io/account
   - Click: Settings → General → API Token
   - Copy token and add to GitHub secrets

2. **CODECOV_TOKEN** (optional)
   - Get from: https://codecov.io/gh/gerrymiller/deep-thinking-agent/settings
   - Copy token and add to GitHub secrets

### CODEOWNERS File

Create `.github/CODEOWNERS`:

```
# Default owners for everything in the repo
* @gerrymiller

# Package-specific owners (add as team grows)
/pkg/llm/ @gerrymiller
/pkg/agent/ @gerrymiller
/pkg/schema/ @gerrymiller

# Documentation
*.md @gerrymiller
/docs/ @gerrymiller

# CI/CD
/.github/ @gerrymiller
```

Commit and push:
```bash
git add .github/CODEOWNERS
git commit -m "chore: Add CODEOWNERS file"
git push
```

## Third-Party Integrations

### 1. Snyk Integration

**Setup:**
1. Go to https://snyk.io/
2. Sign up with GitHub account
3. Click "Add project"
4. Select `gerrymiller/deep-thinking-agent`
5. Import repository
6. Generate API token (Settings → General → API Token)
7. Add token to GitHub secrets as `SNYK_TOKEN`

**Configuration:**
- Enable: Automatic dependency updates
- Set: High severity threshold
- Enable: PR checks

### 2. Codecov Integration (Optional)

**Setup:**
1. Go to https://codecov.io/
2. Sign up with GitHub account
3. Click "Add new repository"
4. Select `gerrymiller/deep-thinking-agent`
5. Copy upload token
6. Add token to GitHub secrets as `CODECOV_TOKEN`

**Badge:**
Already added to README.md - will activate after first CI run with coverage upload.

### 3. Go Report Card

**Setup:**
1. Visit: https://goreportcard.com/
2. Enter: `github.com/gerrymiller/deep-thinking-agent`
3. Click "Generate Report"
4. Wait for analysis (happens automatically)

No configuration needed - updates automatically on each push.

## Labels and Automations

### Create Labels

Navigate to **Issues → Labels** and create:

**Type Labels:**
```
bug (color: d73a4a) - Something isn't working
enhancement (color: a2eeef) - New feature or request
documentation (color: 0075ca) - Documentation improvements
question (color: d876e3) - Further information requested
```

**Priority Labels:**
```
priority: critical (color: d73a4a) - Blocking production use
priority: high (color: ff9800) - Important but not blocking
priority: medium (color: ffeb3b) - Nice to have
priority: low (color: 8bc34a) - Future consideration
```

**Status Labels:**
```
good first issue (color: 7057ff) - Good for newcomers
help wanted (color: 008672) - Community help requested
wontfix (color: 6c757d) - Won't be fixed/implemented
duplicate (color: 6c757d) - Duplicate issue
```

**Package Labels:**
```
pkg: llm (color: e0e0e0)
pkg: embedding (color: e0e0e0)
pkg: vectorstore (color: e0e0e0)
pkg: document (color: e0e0e0)
pkg: schema (color: e0e0e0)
pkg: agent (color: e0e0e0)
pkg: workflow (color: e0e0e0)
pkg: retrieval (color: e0e0e0)
pkg: nodes (color: e0e0e0)
```

**Area Labels:**
```
cli (color: 0366d6)
api (color: 0366d6)
tests (color: 0366d6)
ci/cd (color: 0366d6)
dependencies (color: 0366d6)
examples (color: 0366d6)
configuration (color: 0366d6)
```

**Automated Label:**
```
automated (color: 6c757d) - Created by automation
```

### Enable GitHub Discussions

Navigate to **Settings → General → Features**

- [x] Enable Discussions

Create categories:
1. **Announcements** - Project updates and releases
2. **General** - General discussion
3. **Ideas** - Feature requests and suggestions
4. **Q&A** - Questions and answers
5. **Show and Tell** - Community projects using Deep Thinking Agent

## Verification

### 1. Test Branch Protection

Try to push directly to main:
```bash
git checkout main
touch test.txt
git add test.txt
git commit -m "test: Branch protection test"
git push origin main
```

Should fail with: "protected branch"

### 2. Test CI Workflows

Create a test PR:
```bash
git checkout develop
git checkout -b test/ci-verification
echo "# Test" >> TEST.md
git add TEST.md
git commit -m "test: Verify CI workflows"
git push origin test/ci-verification
```

Open PR on GitHub and verify:
- [x] CI workflow runs
- [x] Security workflow runs
- [x] Auto-labeler applies labels
- [x] Status checks appear

### 3. Verify Security Scanning

Navigate to **Security → Code scanning**

Should show:
- [x] CodeQL analysis results
- [x] Gosec scan results
- [x] Snyk vulnerabilities (if any)

### 4. Verify Dependabot

Navigate to **Insights → Dependency graph → Dependabot**

Should show:
- [x] Dependabot enabled
- [x] Monitoring Go modules
- [x] Monitoring GitHub Actions

### 5. Test Release Workflow

Create a test release:
```bash
git checkout main
git tag -a v0.1.0-beta -m "Beta release for testing"
git push origin v0.1.0-beta
```

Verify:
- [x] Release workflow triggers
- [x] Binaries are built for all platforms
- [x] Release is created with notes
- [x] Checksums are generated

### 6. Check Badges

Visit your README on GitHub and verify all badges display correctly:
- [x] Go version badge
- [x] License badge
- [x] Release badge
- [x] CI status badge
- [x] Coverage badge (after first upload)
- [x] Go Report Card badge
- [x] Security badge
- [x] Social badges

## Post-Setup Tasks

### 1. First Release

Create your first official release:

```bash
git checkout main
git tag -a v0.1.0 -m "Initial public release"
git push origin v0.1.0
```

### 2. Announcement

Create an announcement in Discussions:
```markdown
# 🎉 Deep Thinking Agent is now open source!

We're excited to announce the public release of Deep Thinking Agent...

[Link to README and features]
```

### 3. Social Media

Share your project:
- Twitter/X
- LinkedIn
- Reddit (r/golang, r/MachineLearning)
- Hacker News (Show HN)
- Product Hunt (optional)

### 4. Documentation Site (Future)

Consider setting up GitHub Pages:
- Go to **Settings → Pages**
- Source: Deploy from branch `main`
- Folder: `/docs`

### 5. Monitor Activity

Set up notifications:
- Watch your repository
- Enable email notifications for:
  - Issues
  - Pull requests
  - Security alerts
  - Discussions

## Troubleshooting

### Workflows Not Running

**Issue:** GitHub Actions workflows not triggering

**Solution:**
1. Check workflow files are in `.github/workflows/`
2. Verify YAML syntax with: https://www.yamllint.com/
3. Check Actions permissions: Settings → Actions → General
4. Enable: "Allow all actions and reusable workflows"

### Badges Not Displaying

**Issue:** Badges show "unknown" or don't load

**Solution:**
1. **Go Report Card**: Visit goreportcard.com and generate report manually
2. **Codecov**: Wait for first CI run with coverage upload
3. **CI Badge**: Wait for first workflow run
4. Clear browser cache and refresh

### Protected Branch Issues

**Issue:** Can't push to protected branch

**Solution:**
This is expected! Protected branches require PRs. Create feature branches and open PRs.

### Dependabot PRs Not Appearing

**Issue:** No Dependabot PRs after enabling

**Solution:**
1. Check `dependabot.yml` syntax
2. Wait 24-48 hours for first scan
3. Manually trigger: Security → Dependabot → Check for updates

## Maintenance

### Weekly Tasks

- [ ] Review and merge Dependabot PRs
- [ ] Check security alerts
- [ ] Review open issues and PRs
- [ ] Update project board (if using)

### Monthly Tasks

- [ ] Review and update documentation
- [ ] Check test coverage trends
- [ ] Review contributor activity
- [ ] Plan next release

### Quarterly Tasks

- [ ] Review and update branch protection rules
- [ ] Audit security settings
- [ ] Review label usage and cleanup
- [ ] Update TODO.md with community feedback

## Resources

- [GitHub Docs: Branch Protection](https://docs.github.com/en/repositories/configuring-branches-and-merges-in-your-repository/managing-protected-branches/about-protected-branches)
- [GitHub Actions Documentation](https://docs.github.com/en/actions)
- [Snyk Documentation](https://docs.snyk.io/)
- [Codecov Documentation](https://docs.codecov.com/)
- [Dependabot Documentation](https://docs.github.com/en/code-security/dependabot)

---

**Questions or Issues?**

Open an issue or start a discussion in the repository!
