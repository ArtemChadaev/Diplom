# Git Workflow & Release Process

← [Back to Other Docs](../../README.md) | See also: [`/merge-all` workflow](../../.agents/workflows/merge-all.md) | [`/start-again` workflow](../../.agents/workflows/start-again.md)

> **This is the authoritative guide for merging into `main` and managing releases.**
> Procedures below are executed by the AI agent via slash-commands — not manually by developers.

---

## Branch Structure

```
main                 ← stable releases only; no direct commits
├── develop-backend  ← all backend feature work goes here
└── develop-frontend ← all frontend feature work goes here
```

- **Direct commits to `main` are forbidden.**
- All changes reach `main` exclusively via the `/merge-all` workflow.
- Code for the backend **only** goes into `develop-backend`.
- Code for the frontend **only** goes into `develop-frontend`.

---

## Merge-All Workflow (`/merge-all`)

Triggered when a stable release is ready. Full procedure:

1. **Sync remotes**
   ```bash
   git checkout develop-backend && git pull origin develop-backend
   git checkout develop-frontend && git pull origin develop-frontend
   git checkout main && git pull origin main
   ```

2. **Sequential merge**
   ```bash
   git merge develop-backend
   git merge develop-frontend
   ```
   > If any merge conflict occurs → **STOP**. Resolve manually before continuing.

3. **Post-merge validation**
   ```bash
   # From backend/
   swag init -g cmd/main.go -o ../docs/api/ --parseDependency
   go test ./...
   ```

4. **On success** → proceed to reporting / start-again
   **On failure** → create hotfix branch:
   ```bash
   git checkout -b hotfix/release-repair
   # fix all failures
   go test ./...
   git checkout main
   git merge hotfix/release-repair
   git branch -d hotfix/release-repair
   ```

---

## Start-Again Workflow (`/start-again`)

Triggered **only after** `main` is confirmed stable (all tests green, swagger synced).

**Version Tag Format:** `v[Status].[Release].[Hotfix]`

| Digit | Increment when |
|-------|---------------|
| Status | Major structural or architectural change |
| Release | Standard merge from develop branches |
| Hotfix | Direct fix applied to `main` |

```bash
# Step 1: Tag
git describe --tags --abbrev=0          # check last tag
git tag -a vX.X.X -m "Finalized Release vX.X.X"

# Step 2: Reset develop branches to main
git checkout develop-backend && git reset --hard main
git checkout develop-frontend && git reset --hard main

# Step 3: Push everything
git push origin main --tags
git push origin develop-backend --force
git push origin develop-frontend --force
```

> ⚠️ `--force` permanently rewrites the history of both develop branches. Only run after confirming `main` is stable.
