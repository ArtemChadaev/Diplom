# Git Flow, CI/CD & DevOps

← [Back to Main README](../../README.md)

This document is the **official policy** for branching strategy, merge rules, and CI/CD pipeline behavior. All contributors and AI assistants must follow these rules without exception.

---

## 1. Branching Policy

### Branch Map

```
main                          ← production (stable releases only)
  │
  ├── develop-backend         ← all backend feature development
  │     └── feature/backend/<name>
  │
  └── develop-frontend        ← all frontend feature development
        └── feature/frontend/<name>

hotfix/<scope>/<name>         ← emergency fix, branched from main
```

### Rules Per Branch Type

#### `main` — Production
- **No direct commits.** Ever.
- Only accepts merges from `release/*` or `hotfix/*` branches.
- Every merge into `main` must be tagged with a version (`vX.Y.Z`).
- High-level configuration changes (e.g., `docker-compose.yml`, `.env.example`, root `README.md`) may be committed directly only if they do not affect application logic.

#### `develop-backend` / `develop-frontend` — Active Development
- All feature work happens here.
- Must be kept in a passing state (CI green).
- Direct commits are allowed; prefer short-lived `feature/` branches for larger changes.

#### `feature/<scope>/<name>` — Feature Branches
- Branched from the appropriate develop branch.
- Scope values: `backend`, `frontend`, `all`
- Merged back via Pull Request with at least one review.
- Examples: `feature/backend/fefo-algorithm`, `feature/frontend/inventory-dashboard`

#### `hotfix/<scope>/<name>` — Emergency Fixes
- **Always branched from `main`**, never from develop.
- After fix is verified:
  1. Merge into `main` (and tag)
  2. Merge into **both** `develop-backend` and `develop-frontend` to keep them in sync
- Examples: `hotfix/backend/jwt-expiry-bug`, `hotfix/all/docker-compose-fix`

#### `release/<scope>/<version>` — Release Preparation
- Branched from the appropriate develop branch.
- Only bug fixes allowed — no new features.
- Merged into `main` and back into develop after release.

### Naming Reference

| Type | Pattern | Base Branch |
|------|---------|-------------|
| Feature | `feature/<scope>/<name>` | `develop-<scope>` |
| Release | `release/<scope>/<version>` | `develop-<scope>` |
| Hotfix | `hotfix/<scope>/<name>` | `main` |

Where `<scope>` is one of: `backend`, `frontend`, `all`.

---

## 2. Merge Rules & Pull Requests

- PRs targeting `main` require **at least 1 approval**.
- PRs must pass all CI checks before merge.
- Squash merge is preferred for feature branches to keep `main` history clean.
- Merge commit (no squash) is preferred for release and hotfix branches to preserve audit trail.

---

## 3. Workflow Integrity Rule

> **An AI assistant must not modify existing application logic or documentation unless the change directly relates to the task at hand.**
>
> When documentation refactoring is the task: only add links, restructure, and categorize. Do not rewrite technical content.
> When feature development is the task: do not touch unrelated files.

---

## 4. CI/CD Pipeline

### Pre-commit (local)
- **Backend:** `golangci-lint run`
- **Frontend:** `eslint` via `npm run lint`

### Pull Request Checks
- Backend: `go test ./...` — unit + integration tests
- Frontend: `vitest` — component tests
- Code coverage must not fall below **70%**

### Staging (before release)
- Smoke E2E tests using **Playwright** against the staging environment
- Must pass before any merge into `main`

### Production Deploy
```bash
# From monorepo root
docker compose --profile all up -d --build
```

---

## 5. DevOps & Infrastructure

### Services (docker-compose.yml)

| Service | Image | Profile | Port |
|---------|-------|---------|------|
| `backend` | Custom (Go) | `backend-all`, `all` | `$PORT` |
| `frontend` | Custom (Next.js) | `all` | `$FRONTEND_PORT` |
| `postgres` | `postgres:17-alpine` | `backend-all`, `all` | `$DB_PORT` |
| `valkey` | `valkey/valkey:8-alpine` | `backend-all`, `all` | `$VALKEY_PORT` |
| `migrate` | `migrate/migrate` | `backend-all`, `all` | *(runs once)* |
| `migrate-down` | `migrate/migrate` | `tools` | *(manual rollback)* |

### Running Migrations Manually

```bash
# Apply all pending migrations
docker compose run --rm migrate

# Roll back last migration
docker compose --profile tools run --rm migrate-down
```

### Environment Variables

Copy `.env.example` → `.env`. Required variables:

| Variable | Required | Description |
|----------|----------|-----------|
| `DB_HOST` | **Yes** | PostgreSQL host |
| `DB_USER` | **Yes** | DB username |
| `DB_PASSWORD` | **Yes** | DB password |
| `DB_NAME` | **Yes** | Database name |
| `JWT_SECRET` | **Yes** | JWT signing secret |
| `ADMIN_EMAIL` | **Yes** | Seed admin user email |
| `PORT` | No (8080) | Backend HTTP port |
| `FRONTEND_PORT` | No | Frontend exposed port |
| `GOOGLE_CLIENT_ID` | No | Google OAuth client ID |
| `VALKEY_PORT` | No | Valkey port |
