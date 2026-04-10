# Backend — Standard Operating Procedure (SOP)

← [Back to Main README](../../README.md) | [Architecture →](./architecture.md) | [Testing Notes →](./test.md)

> **Authority:** This document is the single source of truth for all backend development decisions.
> When this SOP conflicts with any other document, **this SOP wins**.
> All rules below are non-negotiable unless explicitly revised here.

---

## Table of Contents

1. [Onion Architecture Canon](#1-onion-architecture-canon)
2. [Error Handling](#2-error-handling)
3. [Logging Standard](#3-logging-standard)
4. [TDD / BDD Development Cycle](#4-tdd--bdd-development-cycle)
5. [Swagger / API Synchronization](#5-swagger--api-synchronization)

> **Git & Release Workflow** (merge into `main`, `/merge-all`, `/start-again`) → [`docs/other/git-workflow.md`](../other/git-workflow.md)

---

## 1. Onion Architecture Canon

### 1.1 Layer Call Rules

Dependencies flow **inward only**. No exceptions.

```
Handler  →  Service  →  Repository  →  (Database)
   ↑              ↑              ↑
 Domain         Domain         Domain
(types,       (types,        (types,
 errors,       errors,        errors,
 interfaces)   interfaces)    interfaces)
```

| ✅ Allowed | ❌ Forbidden |
|-----------|------------|
| Handler imports `domain` (errors, types, interfaces) | Handler imports `repository` directly |
| Handler calls `service.*` | Handler calls `repo.*` directly |
| Service imports `domain` | Service imports `handler` or `dto` |
| Repository imports `domain` | Repository imports `service` |

**Rule:** `internal/handler` may only communicate with the outside world via `internal/service`. The only allowed cross-layer import at all layers is `internal/domain`.

### 1.2 Domain Layer Rules

`internal/domain/` is the **core**. It contains:

- Pure Go structs — **no `json:` tags, no `gorm:` tags**
- Business interfaces (`UserRepository`, `AuthService`, etc.)
- Sentinel errors (used via `errors.Is`)
- Custom `AppError` type

If a struct needs JSON tags → it belongs in `internal/handler/dto/`.  
If a struct needs GORM tags → it belongs in `internal/repository/dao/`.

### 1.3 Implementing a New Entity — Mandatory Checklist

Use **`EmployeeProfile`** as the canonical reference implementation.  
Follow this exact sequence. Do not skip steps.

```
[ ] 1. internal/domain/         — Add domain struct (no tags)
[ ] 2. internal/domain/         — Add repository interface
[ ] 3. internal/domain/         — Add service interface
[ ] 4. internal/domain/errors.go — Add sentinel errors (ErrXxxNotFound, etc.)
[ ] 5. internal/repository/dao/ — Add DAO struct (gorm tags + TableName())
[ ] 6. internal/repository/     — Implement the repository interface
[ ]    └── implement toDomain() / fromDomain() converters
[ ] 7. internal/service/        — Implement the service interface
[ ] 8. internal/handler/dto/    — Add request/response DTOs (json tags)
[ ] 9. internal/handler/        — Add handler functions with Swagger annotations
[ ] 10. internal/handler/handler.go — Register routes
[ ] 11. docs/api/               — Run swag init (see §5)
```

> **Shortcut rule:** If a similar entity already exists (e.g., `EmployeeProfile`), duplicate its file structure, rename, and adapt. Do not invent a new layout.

### 1.4 DAO Pattern Invariant

Every DAO must implement two private methods:

```go
// Converts DAO → domain model (used when reading from DB)
func (d *UserDAO) toDomain() *domain.User { ... }

// Converts domain model → DAO (used when writing to DB)
func fromDomain(u *domain.User) *UserDAO { ... }
```

No domain model must ever touch the DB directly. No DAO must ever be returned outside the repository layer.

---

## 2. Error Handling

### 2.1 Forbidden Practices

```go
// ❌ FORBIDDEN — wrapping with fmt.Errorf destroys sentinel error chains
return fmt.Errorf("failed to find user: %w", err)

// ❌ FORBIDDEN — swallowing errors silently
_ = err

// ❌ FORBIDDEN — logging in service or repository layer
log.Error("db error", "err", err)  // ← not here
return err
```

### 2.2 Correct Error Flow

```go
// ✅ Repository: return as-is or map to domain sentinel
func (r *userRepo) FindByID(ctx context.Context, id int) (*domain.User, error) {
    var dao dao.UserDAO
    if err := r.db.WithContext(ctx).First(&dao, id).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, domain.ErrUserNotFound  // map to sentinel
        }
        return nil, err  // unknown DB error — return as-is
    }
    return dao.toDomain(), nil
}

// ✅ Service: propagate sentinel errors unchanged
func (s *authService) GetUser(ctx context.Context, id int) (*domain.User, error) {
    user, err := s.userRepo.FindByID(ctx, id)
    if err != nil {
        return nil, err  // do not wrap, do not log
    }
    return user, nil
}

// ✅ Handler: check sentinels, log, respond
func (h *Handler) getUser(w http.ResponseWriter, r *http.Request) {
    log := logger.FromContext(r.Context())

    user, err := h.service.Auth.GetUser(r.Context(), id)
    if err != nil {
        if errors.Is(err, domain.ErrUserNotFound) {
            writeError(w, http.StatusNotFound, "user not found")
            return
        }
        log.Error("unexpected error getting user", "user_id", id, "error", err)
        writeError(w, http.StatusInternalServerError, "internal error")
        return
    }
    writeJSON(w, http.StatusOK, toUserDTO(user))
}
```

### 2.3 Where Errors Are Logged

| Layer | Log errors? | Return errors? |
|-------|------------|---------------|
| Repository | ❌ No | ✅ Yes — always |
| Service | ❌ No | ✅ Yes — always |
| Handler | ✅ Yes — always before return | ✅ Yes (to caller via HTTP) |

---

## 3. Logging Standard

### 3.1 The Only Correct Logger

```go
import "github.com/ima/diplom-backend/internal/pkg/logger"

// In every function that uses a logger:
log := logger.FromContext(ctx)
```

**Variable name is always `log`**, never `logger`, `l`, or `slog`.  
**Declared at the top of each function**, never inline mid-function.

### 3.2 Forbidden Alternatives

```go
// ❌ All of these are forbidden:
fmt.Println("something happened")
fmt.Printf("user %d logged in", id)
log.Println("...")          // stdlib log
slog.Info("...")            // global slog without context
logger.FromContext(ctx).Info("...")  // inline — declare log first
```

### 3.3 Structured Key-Value Logging

Always use key-value pairs. Never use bare strings for variable data.

```go
// ✅ Correct
log.Info("employee profile updated", "admin_id", callerID, "target_user_id", targetID)
log.Error("failed to revoke session", "session_id", sessionID, "error", err)
log.Warn("session not found during logout", "session_id", id)

// ❌ Incorrect
log.Info("employee profile " + strconv.Itoa(id) + " updated")
log.Error(err.Error())
```

### 3.4 Log Levels

| Level | When to use |
|-------|-------------|
| `Debug` | Dev-only: detailed state, iteration counts |
| `Info` | Successful business operations (profile updated, token issued) |
| `Warn` | Expected but noteworthy failures (session not found during logout) |
| `Error` | Unexpected failures, DB errors, panics |

### 3.5 What the Logger Carries Automatically

`logger.FromContext(ctx)` automatically enriches every log entry with:
- `request_id` — injected by `chi.middleware.RequestID` + `loggingMiddleware`
- `user_id` — injected by `AuthRequired` middleware via `logger.WithUserID(ctx, id)`

Do not manually add these fields — they are injected automatically.

### 3.6 Context Enrichment Functions

```go
// These are called by middleware — NOT by application code directly:
logger.WithRequestID(ctx, requestID)  // called in loggingMiddleware
logger.WithUserID(ctx, userID)        // called in AuthRequired
```

---

## 4. TDD / BDD Development Cycle

### 4.1 The Mandatory Process

For **any non-trivial logic**, the implementation cycle is:

```
Step 1: SCENARIO  — Describe what the function must do in Given/When/Then
Step 2: TEST      — Write the test using the scenario
Step 3: PAUSE     — Present the test to the user for approval
Step 4: CODE      — Only after approval: implement the function
Step 5: VERIFY    — Run tests; all must pass before proceeding
```

**Do not write business code before writing the test.  
Do not present the test for approval without the scenario.  
Do not merge unverified code.**

### 4.2 When Tests Are NOT Required (CRUD Exemption)

A function is exempt from mandatory testing if **all** of the following are true:

- It performs only: Create / Read by ID / Update fields / Delete by ID
- It contains **no conditional logic** (no `if` based on business state)
- It contains **no data transformation** (no calculations, aggregations, sorting)
- It contains **no security-sensitive operations** (no permission checks, no token ops)

If even one condition fails → tests are required.

| Operation | Tests required? |
|-----------|----------------|
| `repo.FindByID` | ❌ No |
| `repo.Create` | ❌ No |
| `service.GetProfile` (permission check) | ✅ Yes |
| `service.RefreshTokens` (theft detection) | ✅ Yes |
| `service.Send OTPCode` (rate limiting) | ✅ Yes |
| `service.VerifyOTPCode` (expiry + attempts) | ✅ Yes |
| Any FEFO algorithm | ✅ Yes |
| Any filter/search with conditions | ✅ Yes |

### 4.3 Scenario Format (BDD)

Write the scenario as a Go comment block **directly above** the test function:

```go
// Scenario: Refresh token theft detection
//   Given:  A refresh token that has already been revoked
//   When:   RefreshTokens is called with that token
//   Then:   All sessions for that user are revoked
//   And:    ErrInvalidToken is returned
func TestAuthService_RefreshTokens_DetectsTheft(t *testing.T) {
    ...
}
```

### 4.4 Test Structure

Use **table-driven tests** with `testify/assert`. One table, multiple scenarios.

```go
func TestAuthService_VerifyOTPCode(t *testing.T) {
    // Scenario: OTP code is verified
    //   Given:  A valid OTP code stored in Valkey
    //   When:   VerifyOTPCode is called with matching email and code
    //   Then:   A token pair is returned
    //   And:    The OTP entry is deleted from Valkey

    tests := []struct {
        name      string
        code      string
        wantErr   error
        wantPair  bool
    }{
        {
            name:     "valid code",
            code:     "123456",
            wantErr:  nil,
            wantPair: true,
        },
        {
            name:    "wrong code",
            code:    "000000",
            wantErr: domain.ErrOTPInvalid,
        },
        {
            name:    "expired code",
            code:    "123456",
            wantErr: domain.ErrOTPNotFound,
        },
        {
            name:    "max attempts exceeded",
            code:    "123456",
            wantErr: domain.ErrOTPMaxAttempts,
        },
    }

    for _, tc := range tests {
        t.Run(tc.name, func(t *testing.T) {
            // ... setup mocks
            pair, err := svc.VerifyOTPCode(ctx, email, tc.code, meta)
            assert.ErrorIs(t, err, tc.wantErr)
            if tc.wantPair {
                assert.NotNil(t, pair)
            }
        })
    }
}
```

### 4.5 Negative Cases Are Mandatory

Every test suite for non-trivial logic **must** include at least one negative case:

- Input that violates a business rule
- Expired or missing resource
- Insufficient permissions
- Duplicate / conflict

A test suite with only happy-path cases will not be approved.

### 4.6 Test Tooling

| Tool | Purpose | Import |
|------|---------|--------|
| `testify/assert` | Assertions for all unit tests | `github.com/stretchr/testify/assert` |
| `testify/mock` | Mock interfaces | `github.com/stretchr/testify/mock` |
| `net/http/httptest` | API-level handler tests | stdlib |

**`httptest` isolation rule:** Handler (API) tests using `httptest.NewRecorder()` are written **only on explicit request**. They are never mixed in the same file as service/unit tests. Naming convention: `handler_test.go` vs `service_test.go`.

---

## 5. Swagger / API Synchronization

### 5.1 The Priority Rule

> **Swagger-first.** Before changing a handler's logic that affects request/response shape — update the annotations first.

Order of operations for any handler change:

```
1. Update @Param / @Success / @Failure annotations in the .go file
2. Run swag init (see §5.2)
3. Verify docs/api/swagger.json reflects the change
4. Implement the logic change
```

### 5.2 Generation Command

Always run from the `backend/` directory:

```bash
swag init -g cmd/main.go -o ../docs/api/ --parseDependency
```

| Flag | Purpose |
|------|---------|
| `-g cmd/main.go` | Entry point for annotation scanning |
| `-o ../docs/api/` | Output directory (relative to `backend/`) |
| `--parseDependency` | Resolves types defined in imported packages |

Output:
- `docs/api/docs.go` — Go code (auto-generated, do not edit)
- `docs/api/swagger.json` — Machine-readable spec
- `docs/api/swagger.yaml` — Human-readable spec

### 5.3 Required Annotations — Minimum Set

Every handler function **must** have all of the following annotations:

```go
// @Summary     Short human-readable action description
// @Tags        auth
// @Accept      json
// @Produce     json
// @Param       body  body  dto.VerifyCodeRequest  true  "OTP code and email"
// @Success     200   {object}  dto.TokenResponse
// @Failure     400   {object}  dto.ErrorResponse  "Invalid request body"
// @Failure     401   {object}  dto.ErrorResponse  "Invalid or expired code"
// @Failure     500   {object}  dto.ErrorResponse  "Internal server error"
// @Router      /auth/verify-code [post]
func (h *Handler) verifyCode(w http.ResponseWriter, r *http.Request) {
```

Handlers without a full annotation set **must not be merged** into `develop-backend` or `main`.

### 5.4 Merge Gate — Swagger Sync Check

Before any merge into `develop-backend` or `main`:

1. Run `swag init -g cmd/main.go -o ../docs/api/ --parseDependency`
2. Check that `docs/api/swagger.json` is not empty (`"paths": {}` = **blocked**)
3. Confirm that every route in `routing.md` has a corresponding entry in swagger

If `docs/api/swagger.json` contains `"paths": {}` after generation → **merge is blocked** until annotations are added.

---

## Appendix: Quick Reference Card

### Forbidden at a Glance

| # | Forbidden | Correct alternative |
|---|-----------|-------------------|
| 1 | `fmt.Println`, `fmt.Printf` | `log := logger.FromContext(ctx); log.Info(...)` |
| 2 | `fmt.Errorf("msg: %w", err)` | `return err` or `return domain.ErrXxx` |
| 3 | Logging in Service/Repository | Return error; log only in Handler |
| 4 | Handler calls `repo.*` directly | Handler calls `service.*` only |
| 5 | Domain struct with `json:` or `gorm:` tags | DTO in `dto/`, DAO in `dao/` |
| 6 | Merging without Swagger update | Run `swag init`, confirm non-empty paths |
| 7 | Direct commits to `main` or wrong branch | See [`git-workflow.md`](../other/git-workflow.md) |
| 8 | Code without scenario/test (for complex logic) | Write scenario → test → get approval → code |
| 9 | Happy-path-only tests | Include at least one negative case per suite |
| 10 | `logger` inline mid-function | `log := logger.FromContext(ctx)` at function top |
