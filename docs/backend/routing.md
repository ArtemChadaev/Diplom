# Backend — HTTP Router & Middleware

← [Back to Main README](../../README.md) | [Domain / Service / Repository →](./domain-repository-service.md)

---

## HTTP Router (`chi`)

Package `go-chi/chi/v5` — lightweight idiomatic router for `net/http`.

### Global Middleware (registered in `handler.go`)

```go
r.Use(middleware.RequestID)           // Adds X-Request-Id header
r.Use(middleware.Recoverer)           // Panic recovery → 500
r.Use(middleware.Timeout(60s))        // 60-second request timeout
r.Use(middleware.Heartbeat("/ping"))  // Healthcheck: GET /ping
r.Use(h.loggingMiddleware)            // Custom slog request logger
```

### Route Structure

```
POST /auth/google                — Google OAuth (ID-token → access_token)
POST /auth/refresh               — Token refresh via cookie
POST /auth/logout                — Logout (revoke session + clear cookie)

/api/v1/ (middleware: AuthRequired — Bearer JWT check)
  DELETE /sessions/{sessionID}   — Revoke own session

  /admin/ (middleware: RequireRole("admin"))
    PATCH  /users/{id}/role      — Assign role
    PATCH  /users/{id}/blocked   — Block / unblock user
    DELETE /sessions/{sessionID} — Revoke any session
    GET    /employees/           — List employee profiles
    GET    /employees/{userID}   — Get employee profile
    PATCH  /employees/{userID}   — Update employee profile
```

---

## Middleware (`internal/handler/middleware/`)

### `AuthRequired` (`auth_middleware.go`)

Extracts JWT from `Authorization: Bearer <token>` header, validates via `domain.TokenService`, and stores in context:

```go
ctx.Value(middleware.CtxUserID) → int       // User ID
ctx.Value(middleware.CtxRole)   → UserRole  // User role
ctx.Value(middleware.CtxEmail)  → string    // Email
```

### `RequireRole` (`rbac_middleware.go`)

Checks role from context. Accepts one or more allowed `UserRole` values. Returns `403 Forbidden` on mismatch.

---

## Handler Layer (`internal/handler/`)

### `Handler` (`handler.go`)

Holds references to `service.Service` and `domain.TokenService`. The `Router()` method assembles the chi router.

### DTOs (`internal/handler/dto/`)

Request/response structs with JSON tags, isolated from domain:

- `auth_dto.go` — `GoogleAuthRequest`, `TokenResponse`
- `user_dto.go` — `AssignRoleRequest`, `SetBlockedRequest`, `UserResponse`
- `employee_profile_dto.go` — `EmployeeProfileDTO`
- `admin_employee_dto.go` — `PatchEmployeeProfileRequest`, `EmployeeProfileResponse`

### Handler Files

| File | Handlers |
|------|----------|
| `auth_handler.go` | `refresh`, `logout`, `clearRefreshCookie` |
| `oauth_handler.go` | `googleLogin` |
| `admin_handler.go` | `adminAssignRole`, `adminSetBlocked`, `revokeSession`, `adminRevokeSession` |
| `admin_employee_handler.go` | `adminListEmployeeProfiles`, `adminGetEmployeeProfile`, `adminPatchEmployeeProfile` |
