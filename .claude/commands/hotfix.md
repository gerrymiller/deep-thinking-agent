---
description: Create a new hotfix branch following gitflow conventions
---

Create a new hotfix branch from the `main` branch for critical production fixes.

1. Ensure we're starting from an up-to-date `main` branch
2. Ask the user for a hotfix name if not provided in the command
3. Create a new branch named `hotfix/{hotfix-name}`
4. Confirm the branch was created successfully
5. Remind the user that hotfix branches should be merged to both `main` AND `develop` when complete

Hotfix branches are used for critical production issues that need immediate attention.

Follow the gitflow workflow documented in CLAUDE.md.