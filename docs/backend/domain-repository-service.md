# Backend — Domain, Repository & Service

← [Back to Main README](../../README.md) | [Routing & Middleware →](./routing.md)

---

## Domain Layer (`internal/domain/`)

Pure Go structs **without ORM or JSON tags** — business fields only.

### Key Models

```go
// User — system user
type User struct {
    ID          int
    Email       string    // primary identifier
    GoogleID    *string   // nullable (OAuth)
    TelegramID  *int64    // nullable (OAuth)
    Role        UserRole
    NsPvAccess  bool      // access to NS/PV (narcotic/psychotropic substances)
    UkepBound   bool      // qualified e-signature bound
    IsBlocked   bool
    CreatedAt   time.Time
    UpdatedAt   time.Time
}

// User roles (mirror of PostgreSQL ENUM user_role)
RoleAdmin            UserRole = "admin"
RoleQP               UserRole = "qp"
RoleWarehouseManager UserRole = "warehouse_manager"
RoleStorekeeper      UserRole = "storekeeper"
RolePharmacist       UserRole = "pharmacist"

// RefreshToken — active session
type RefreshToken struct {
    ID        uuid.UUID
    UserID    int
    TokenHash string       // SHA-256 of raw token
    ExpiresAt time.Time
    UserAgent string
    IPAddress string
    Metadata  map[string]any
    CreatedAt time.Time
    RevokedAt *time.Time   // nil = active
}

// EmployeeProfile — staff profile
type EmployeeProfile struct {
    ID                 uint
    UserID             uint
    EmployeeCode       string
    FullName           string
    CorporateEmail     string
    Phone              string
    Position           string
    Department         string
    BirthDate          time.Time
    AvatarURL          string
    HireDate           time.Time
    DismissalDate      *time.Time
    MedicalBookScanURL string
    SpecialZoneAccess  bool
    GDPTrainingHistory []GDPTrainingRecord  // JSONB array
}
```

### Errors (`errors.go`)

Sentinel errors (used with `errors.Is`):

```go
ErrUserNotFound            = "user not found"
ErrEmailTaken              = "email already taken"
ErrInvalidCreds            = "invalid credentials"
ErrUserBlocked             = "account is blocked"
ErrTokenExpired            = "token expired"
ErrInvalidToken            = "invalid token"
ErrSessionNotFound         = "session not found or terminated"
ErrInsufficientPerms       = "insufficient permissions for this operation"
ErrEmployeeProfileNotFound = "employee profile not found"
```

`AppError` — structured error with a code, message, and `slog.Attr` for logging.

---

## Repository Layer (`internal/repository/`)

### PostgreSQL Connection (`postgres.go`)

```go
// GORM + pgx driver
db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
```

### DAO Pattern

Domain models and GORM are separated by **DAO structs** in `repository/dao/` which carry GORM tags:

- `UserDAO` → table `users`
- `SessionDAO` → table `refresh_tokens`
- `EmployeeProfileDAO` → table `employee_profiles`

Repositories convert between `DAO ↔ domain` via `toDomain()` / `fromDomain()` methods.

### `UserRepository` (`user.go`)

| Method | SQL / GORM |
|--------|-----------|
| `FindByID` | `SELECT * FROM users WHERE id = ?` |
| `FindByEmail` | `WHERE email = ?` |
| `FindByGoogleID` | `WHERE google_id = ?` |
| `FindByTelegramID` | `WHERE telegram_id = ?` |
| `IsEmailTaken` | `COUNT(*) WHERE email = ?` |
| `Create` | `INSERT INTO users` |
| `UpdateRole` | `UPDATE users SET role = ? WHERE id = ?` |
| `LinkGoogle` | `UPDATE users SET google_id = ?` |
| `LinkTelegram` | `UPDATE users SET telegram_id = ?` |
| `SetNsPvAccess` | `UPDATE users SET ns_pv_access = ?` |
| `SetBlocked` | `UPDATE users SET is_blocked = ?` |
| `FindProfileByUserID` | JOIN users + employee_profiles |

### `SessionRepository` (`session.go`)

| Method | Description |
|--------|-------------|
| `Create` | Store new session (refresh token hash) |
| `FindByID` | Find session by UUID |
| `FindByTokenHash` | Find session by SHA-256 token hash |
| `FindActiveByUserID` | All non-expired, non-revoked sessions |
| `Revoke` | `SET revoked_at = NOW() WHERE id = ?` |
| `RevokeAllForUser` | `SET revoked_at = NOW() WHERE user_id = ?` |
| `DeleteExpired` | `DELETE WHERE expires_at < NOW()` (cleanup job) |

---

## Service Layer (`internal/service/`)

### `Service` (`service.go`)

Service aggregator. Instantiated once in `main.go`:

```go
service.NewService(repos, jwtSecret, googleClientID)
// Access Token TTL  = 15 minutes
// Refresh Token TTL = 15 days
```

### `AuthService` (`auth_service.go`)

| Method | Description |
|--------|-------------|
| `LoginWithGoogle` | `idtoken.Validate` → find/create user → `issueTokens` |
| `LoginWithTelegram` | **Stub** (`panic("implement...")`), needs bot token hash verification |
| `RefreshTokens` | `HashToken` → `FindByTokenHash` → theft check → revoke old → `issueTokens` |
| `RevokeSession` | Verify ownership or admin → `sessionRepo.Revoke` |
| `AssignRole` | Admin only → `userRepo.UpdateRole` |
| `SetBlocked` | Admin only → `userRepo.SetBlocked` |

**Private `issueTokens`:**
1. Verify user is not blocked
2. Generate raw refresh token (`crypto/rand`) + its SHA-256 hash
3. Create `RefreshToken` record in DB
4. Generate JWT access token with `session_id = RefreshToken.ID`

### `TokenService` (`token_service.go`)

- `GenerateAccessToken(user, sessionID)` → JWT HS256
- `ParseAccessToken(token)` → `TokenClaims{UserID, Role, Email, SessionID}`
- `GenerateRefreshToken()` → `(rawToken, hashToken, error)`
- `HashToken(raw)` → SHA-256 hex

### `EmployeeProfileService` (`employee_profile_service.go`)

CRUD for employee profiles. Checks caller permissions (admin or the user themselves).
