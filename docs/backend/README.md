# Backend — Local Setup & Developer Guide

← [Back to Main README](../../README.md) | [Full Sub-Docs Index →](../../docs/backend.md)

> ⛔ **Branch Rule:** All backend code **must** be written in `develop-backend` only.
> Direct commits to `main` or `develop-frontend` are strictly forbidden.
> Merging into `main` is done exclusively via the [`/merge-all`](../other/git-workflow.md) workflow.

> **AI Context & Navigation:** This README is the entry point for the backend. Before writing or editing code, you MUST know what rules apply:
> - **Architecture & Commits:** Read [**`sop.md`**](./sop.md). It is the supreme authority.
> - **Logging:** You MUST use the custom logger, not `fmt`. Read [**`logging.md`**](./logging.md).
> - **DB Schema:** Check [`migrations.md`](./migrations.md) or [`../database-schema.md`](../database-schema.md).
> - **Domain/Repo details:** See [`domain-repository-service.md`](./domain-repository-service.md) and [`architecture.md`](./architecture.md).

---

## Tech Stack

| Component | Library / Tool | Version |
|-----------|---------------|---------|
| HTTP Router | `github.com/go-chi/chi/v5` | v5.1.0 |
| ORM | `gorm.io/gorm` + `gorm.io/driver/postgres` | v1.31.1 |
| JWT | `github.com/golang-jwt/jwt/v5` | v5.3.1 |
| UUID | `github.com/google/uuid` | v1.6.0 |
| Config | `github.com/ilyakaznacheev/cleanenv` | v1.5.0 |
| Google OAuth | `google.golang.org/api/idtoken` | — |
| Logger | stdlib `log/slog` (custom context wrapper) | — |
| DB Migrations | SQL files in `migrate/`, `golang-migrate` | — |
| Database | PostgreSQL 17 via pgx/v5 (under GORM) | — |
| Cache | Valkey 8 (Redis-compatible) | — |

> **Logging Rule:** Never use `fmt.Println` / `fmt.Printf`. Always use the context-aware logger:
> ```go
> log := logger.FromContext(ctx)
> log.Info("message", "key", value)
> ```

---

## Directory Structure

```
backend/
├── cmd/
│   └── main.go           ← entry point
├── internal/
│   ├── config/           ← env config (cleanenv)
│   ├── domain/           ← interfaces + pure models (no JSON/GORM tags)
│   ├── handler/
│   │   ├── dto/          ← request/response DTOs (json + validate tags)
│   │   ├── middleware/   ← AuthRequired, RequireRole
│   │   └── handler.go    ← chi router assembly
│   ├── repository/
│   │   ├── dao/          ← GORM DAO structs (gorm tags, TableName())
│   │   └── postgres.go   ← GORM + pgx connection
│   └── service/          ← business logic
├── migrate/              ← SQL migration files (*.up.sql / *.down.sql)
├── Dockerfile
└── .env.example
```

---

## Running Locally

```bash
# From the monorepo root — start all services (PostgreSQL, Valkey, Backend, Frontend)
docker compose --profile all up -d

# Backend + DB + cache only (for backend development)
docker compose --profile backend-all up -d

# Watch backend logs
docker compose logs -f backend
```

---

## Database Migrations

Migrations run automatically via the `migrate` Docker Compose service on startup.

### Apply all pending migrations (manual)

```bash
docker compose run --rm migrate
```

### Roll back last migration

```bash
docker compose --profile tools run --rm migrate-down
```

### Run via CLI directly

```bash
docker run --rm \
  --network my_diplom_network \
  -v "$(pwd)/backend/migrate:/migrations" \
  migrate/migrate \
  -path=/migrations \
  -database "postgres://${DB_USER}:${DB_PASSWORD}@postgres:5432/${DB_NAME}?sslmode=disable" \
  up
```

> Before running: ensure the `postgres` container is healthy:
> ```bash
> docker compose up -d postgres
> ```
> Variables `DB_USER`, `DB_PASSWORD`, `DB_NAME` are loaded from `.env` (see `.env.example`).

---

## Environment Variables

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `PORT` | No | 8080 | HTTP server port |
| `DB_HOST` | **Yes** | — | PostgreSQL host |
| `DB_PORT` | No | 5432 | PostgreSQL port |
| `DB_USER` | **Yes** | — | DB username |
| `DB_NAME` | **Yes** | — | Database name |
| `DB_PASSWORD` | **Yes** | — | DB password |
| `JWT_SECRET` | **Yes** | — | JWT signing secret (HS256) |
| `ADMIN_EMAIL` | **Yes** | — | Seed admin user email (bootstrap) |
| `GOOGLE_CLIENT_ID` | No | — | Google OAuth client ID |
| `VALKEY_PORT` | No | 6379 | Valkey port |

---

## Navigation Directory

To prevent context bloat, read only the files you need for your current task:

→ **SOP & Core Rules (Errors, TDD, Flow):** [`sop.md`](./sop.md) (Mandatory read before coding)
→ **Logging Rules (No fmt allowed):** [`logging.md`](./logging.md) (Strictly enforced)
→ **Architecture & Folder Info:** [`architecture.md`](./architecture.md)
→ **Routing & Handlers:** [`routing.md`](./routing.md)
→ **Domain & Service Layer:** [`domain-repository-service.md`](./domain-repository-service.md)
→ **Auth Rules (Tokens, Roles):** [`auth.md`](./auth.md)
→ **Migrations List:** [`migrations.md`](./migrations.md)
→ **Test Rules:** [`test.md`](./test.md)
