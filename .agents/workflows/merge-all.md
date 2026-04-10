---
description: This workflow is responsible for the technical migration of features from the development environment to the production branch.
---

Role: Release Integrator

Step 1: Raw Integration
Sync Remotes: * git checkout develop-backend && git pull origin develop-backend

git checkout develop-frontend && git pull origin develop-frontend

Switch to Main: git checkout main and ensure it is up to date (git pull origin main).

Sequential Merge:

Execute git merge develop-backend.

Execute git merge develop-frontend.

CRITICAL: If any merge conflicts occur during these steps, STOP immediately and notify the user. Do not proceed to testing until conflicts are resolved manually.

Step 2: Post-Merge Validation (Bug Hunting)
API Documentation: From the backend root, run swag init -o ../docs/api/. Verify if the code is compatible with the documentation.

Backend Testing: Execute go test ./... inside the backend directory.

Result Analysis:

If Success: Proceed to Step 4 (Reporting).

If Bugs/Failures Found: Proceed to Step 3 (Hotfix Cycle).

Step 3: Hotfix Cycle (Conditional)
If bugs were detected in Step 2:

Create Hotfix Branch: git checkout -b hotfix/release-repair.

Fixing: Address all failing tests and bugs within this branch.

Verification: Run go test ./... again.

Final Integration: Once all tests are green, switch back to main and merge the fix:

git checkout main

git merge hotfix/release-repair

git branch -d hotfix/release-repair (Delete local hotfix branch).

Step 4: Reporting
Provide a summary of the merge process (commit count, status of tests).

Report whether the integration was direct or required a hotfix cycle.

Confirm that the main branch is now integrated and verified.