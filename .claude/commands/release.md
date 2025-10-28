---
description: Create a new release branch following gitflow conventions
---

Create a new release branch from the `develop` branch.

1. Ensure we're starting from an up-to-date `develop` branch
2. Ask the user for a version number if not provided in the command
3. Create a new branch named `release/v{version}` (e.g., release/v1.0.0)
4. Confirm the branch was created successfully
5. Remind the user that release branches should:
   - Only receive bug fixes, documentation, and release-oriented tasks
   - Be merged to both `main` AND `develop` when complete
   - Be tagged with the version number after merging to main

Release branches are used for preparing a new production release.

Follow semantic versioning (MAJOR.MINOR.PATCH) and the gitflow workflow documented in CLAUDE.md.