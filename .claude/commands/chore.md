---
description: Create a new chore branch following gitflow conventions
---

Create a new chore branch from the `develop` branch.

1. Ensure we're starting from an up-to-date `develop` branch
2. Ask the user for a chore name if not provided in the command
3. Create a new branch named `chore/{chore-name}`
4. Confirm the branch was created successfully

Chore branches are used for maintenance tasks such as:
- Refactoring code
- Updating dependencies
- Improving documentation
- Code cleanup
- Performance optimizations (non-feature)

Follow the gitflow workflow documented in CLAUDE.md.