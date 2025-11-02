# Branch Protection Setup Guide

Quick reference for setting up branch protection rules via GitHub web UI.

## ⚠️ Current Status

**Issue**: "Your main branch isn't protected" warning on GitHub

**Solution**: Configure branch protection rules via repository settings (web UI only - cannot be done via CLI or code)

---

## Quick Setup (5 Minutes)

### Step 1: Navigate to Settings

Go to: https://github.com/gerrymiller/deep-thinking-agent/settings/branches

Or manually:
1. Open repository on GitHub
2. Click **Settings** tab (requires admin access)
3. Click **Branches** in left sidebar

### Step 2: Add Protection Rule for `main`

Click **"Add branch protection rule"** button

**Branch name pattern**: `main`

#### Basic Protection (Minimum - Recommended)

Check these boxes:

- ✅ **Require a pull request before merging**
  - Set "Required approvals" to: **1**
  
- ✅ **Require conversation resolution before merging**

- ✅ **Require linear history**

- ✅ **Do not allow bypassing the above settings**
  - For solo dev: UNCHECK "Allow specified actors to bypass" or add yourself if you need flexibility

- ✅ **Allow force pushes** → DISABLE (leave unchecked)

- ✅ **Allow deletions** → DISABLE (leave unchecked)

Click **"Create"** at the bottom

#### Advanced Protection (Optional - Add Later)

After your first workflow run, you can add:

- ✅ **Require status checks to pass before merging**
  - Search and add these checks (they'll appear after workflows run):
    - `test`
    - `lint`
    - `format`
    - `coverage`

To add status checks:
1. Edit the branch protection rule
2. Check "Require status checks to pass before merging"
3. Search for status check names in the search box
4. Select each one
5. Check "Require branches to be up to date before merging"
6. Save changes

### Step 3: Add Protection Rule for `develop` (Recommended)

Click **"Add branch protection rule"** again

**Branch name pattern**: `develop`

#### Recommended Settings for Develop

- ✅ **Require a pull request before merging**
  - Set "Required approvals" to: **0** (for solo dev) or **1** (for team)
  
- ✅ **Require conversation resolution before merging**

- ✅ **Require linear history**

- ✅ **Allow force pushes** → DISABLE

- ✅ **Allow deletions** → DISABLE

Click **"Create"**

---

## Verification

### Test Main Branch Protection

Try to push directly to main (should be blocked):

```bash
git checkout main
echo "test" >> test.txt
git add test.txt
git commit -m "test: Direct push to main"
git push origin main
```

**Expected Result**: Push rejected with message about branch protection

### Test Proper Workflow

Create a feature branch instead:

```bash
git checkout develop
git pull origin develop
git checkout -b test/branch-protection
echo "# Test" >> TEST.md
git add TEST.md
git commit -m "test: Verify branch protection workflow"
git push origin test/branch-protection
```

**Expected Result**: Push succeeds

Then create a PR on GitHub to merge into `develop` or `main`.

---

## Status Check Configuration (After First Workflow Run)

Once your CI workflows have run at least once, status checks will be available:

### 1. Trigger First Workflow Run

Push a commit to trigger workflows:

```bash
git checkout develop
echo "# Trigger CI" >> .github/README.md
git add .github/README.md
git commit -m "ci: Trigger first workflow run"
git push origin develop
```

### 2. Wait for Workflows to Complete

Go to: https://github.com/gerrymiller/deep-thinking-agent/actions

Wait for all workflows to finish (usually 2-5 minutes)

### 3. Add Status Checks to Branch Protection

1. Go to Settings → Branches
2. Click **Edit** on the `main` branch protection rule
3. Check **"Require status checks to pass before merging"**
4. In the search box, type and select:
   - `test` (from ci.yml)
   - `lint` (from ci.yml)
   - `format` (from ci.yml)
   - `coverage` (from ci.yml)
5. Check **"Require branches to be up to date before merging"**
6. Scroll down and click **"Save changes"**

Repeat for `develop` branch if desired.

---

## Branch Protection Rules Summary

### Main Branch (Production)

| Setting | Value | Reason |
|---------|-------|--------|
| Require PR | Yes (1 approval) | Code review before production |
| Required status checks | test, lint, format, coverage | Ensure quality |
| Conversation resolution | Yes | Resolve discussions |
| Linear history | Yes | Clean git history |
| Force pushes | No | Prevent history rewriting |
| Deletions | No | Protect production code |

### Develop Branch (Development)

| Setting | Value | Reason |
|---------|-------|--------|
| Require PR | Optional (0 or 1 approval) | Flexible for solo dev |
| Required status checks | test, lint, format | Quality gate |
| Conversation resolution | Yes | Good practice |
| Linear history | Yes | Clean history |
| Force pushes | No | Prevent accidents |
| Deletions | No | Protect development work |

---

## Solo Development Considerations

### Option A: Strict Protection (Recommended)

- Require PRs even for yourself
- Builds good habits
- Creates audit trail
- Can review changes before merging

**Workflow**:
```bash
# Feature work
git checkout -b feature/my-work
# ... make changes ...
git push origin feature/my-work
# Create PR on GitHub → Review → Merge
```

### Option B: Flexible Protection

- Uncheck "Do not allow bypassing the above settings"
- Add yourself to "Allow specified actors to bypass"
- Can push directly when needed, but protection still warns

**When to use**: Hotfixes, urgent documentation updates

---

## Troubleshooting

### "Cannot enable status checks - No status checks found"

**Cause**: Workflows haven't run yet on this repository

**Solution**: 
1. Push a commit to trigger workflows
2. Wait for workflows to complete
3. Return to branch protection settings
4. Status checks will now appear in search

### "I can't access Settings tab"

**Cause**: You don't have admin permissions

**Solution**: 
- For your own repo: Make sure you're logged in as the owner
- For organization repo: Ask an admin to add you

### "I want to bypass protection temporarily"

**Solution**: 
1. Edit branch protection rule
2. Add yourself to "Allow specified actors to bypass required pull requests"
3. Make your emergency change
4. Remove yourself from bypass list after

### "Required checks are failing"

**Cause**: Code quality issues or test failures

**Solution**:
1. Fix the issues in your branch
2. Push fixes
3. Wait for checks to pass
4. Then merge

---

## Best Practices

### ✅ DO:

- Enable branch protection on `main` and `develop`
- Require status checks after first workflow run
- Use PRs for code review (even solo projects)
- Require conversation resolution
- Maintain linear history
- Test branch protection setup with a test branch

### ❌ DON'T:

- Allow force pushes to protected branches
- Skip PR requirements for "quick fixes"
- Bypass protection rules casually
- Delete protected branches
- Disable protection "temporarily" and forget to re-enable

---

## Quick Reference URLs

- **Branch Settings**: https://github.com/gerrymiller/deep-thinking-agent/settings/branches
- **Actions/Workflows**: https://github.com/gerrymiller/deep-thinking-agent/actions
- **PRs**: https://github.com/gerrymiller/deep-thinking-agent/pulls
- **GitHub Docs**: https://docs.github.com/en/repositories/configuring-branches-and-merges-in-your-repository/managing-protected-branches

---

## After Setup Checklist

- [ ] Main branch protection rule created
- [ ] Develop branch protection rule created
- [ ] Test attempted direct push to main (should fail)
- [ ] Test feature branch creation and push (should succeed)
- [ ] First workflow run triggered
- [ ] Status checks added to protection rules
- [ ] Test PR creation and merge process
- [ ] Warning message gone from GitHub UI

---

**Estimated Time**: 5 minutes for basic setup, +5 minutes after first workflow run to add status checks

**Note**: Branch protection can only be configured via GitHub web UI, not through git CLI or API (unless using GitHub CLI with specific commands).
