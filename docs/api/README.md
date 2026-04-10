# API Contract

← [Back to Main README](../../README.md)

> **⚖️ Authority Rule:** This directory is the **final authority** on all request and response shapes.
> If the Backend implementation and the Frontend expectations ever disagree, the specification in this directory takes precedence. **Update the API doc first, then update the code.**

All endpoints are prefixed with `/api/v1/`.

---

## Implemented Endpoints

### Authentication

| Method | Path | Spec File |
|--------|------|-----------|
| `POST` | `/auth/send-code` | [auth-send-code.md](./auth-send-code.md) |
| `POST` | `/auth/verify-code` | [auth-verify-code.md](./auth-verify-code.md) |
| `POST` | `/auth/login` *(Google OAuth)* | [auth-login.md](./auth-login.md) |
| `POST` | `/auth/register` | [auth-register.md](./auth-register.md) |
| `POST` | `/auth/refresh` | [auth-refresh.md](./auth-refresh.md) |
| `POST` | `/auth/logout` | [auth-logout.md](./auth-logout.md) |
| `POST` | `/auth/google` | [auth-google.md](./auth-google.md) |

### Users & Sessions (Admin)

| Method | Path | Spec File |
|--------|------|-----------|
| `GET/POST/PATCH/DELETE` | `/admin/users` | [admin-users.md](./admin-users.md) |
| `GET/POST/PATCH/DELETE` | `/admin/employees` | [admin-employees.md](./admin-employees.md) |
| `GET/DELETE` | `/sessions` | [sessions.md](./sessions.md) |
| `GET/PATCH` | `/users/me` | [users.md](./users.md) |

---

## Planned (Not Yet Implemented in Backend)

The following modules are **fully specified** but backend handlers have not been written yet.

| Module | Spec File |
|--------|-----------|
| Products, Suppliers, References | [planned-modules.md](./planned-modules.md) |
| Inbound Receiving & Quarantine | [planned-modules.md](./planned-modules.md) |
| Warehouse Zones | [planned-modules.md](./planned-modules.md) |
| Environment Log | [planned-modules.md](./planned-modules.md) |
| Orders & FEFO Assembly | [planned-modules.md](./planned-modules.md) |
| Inventory Sessions | [planned-modules.md](./planned-modules.md) |
| Claims & Recalls | [planned-modules.md](./planned-modules.md) |
| Settings (MOS %) | [planned-modules.md](./planned-modules.md) |

> For full user-flow diagrams of these modules see [`docs/form/`](../form/).

---

## Contract Format

Each endpoint spec file follows this structure:

1. **Overview** — method, path, auth requirement, roles
2. **Request** — headers, path params, query params, body schema
3. **Response** — success shape, error codes
4. **Examples** — curl / JSON samples

---

## Swagger / OpenAPI

Swagger generation is configured in [`docs.go`](./docs.go). The generated `swagger.json` / `swagger.yaml` are placeholders until full annotation is complete.

Legacy monolith reference: [`docs/api-endpoints.md`](../api-endpoints.md) *(kept for historical reference; prefer individual spec files above)*.
