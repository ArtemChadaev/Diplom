# Backend — Domain, Repository и Service

> Часть документации `backend.md`. Описывает ключевые модели, репозитории и сервисы.

---

## Слой Domain (`internal/domain/`)

Чистые Go-структуры **без тегов ORM или JSON** — только бизнес-поля.

### Ключевые модели

```go
// User — пользователь системы
type User struct {
    ID          int
    Email       string    // основной идентификатор
    GoogleID    *string   // nullable (OAuth)
    TelegramID  *int64    // nullable (OAuth)
    Role        UserRole
    NsPvAccess  bool      // допуск к НС/ПВ
    UkepBound   bool      // привязана УКЭП
    IsBlocked   bool
    CreatedAt   time.Time
    UpdatedAt   time.Time
}

// Роли пользователей (ENUM user_role в PostgreSQL)
RoleAdmin            UserRole = "admin"
RoleQP               UserRole = "qp"
RoleWarehouseManager UserRole = "warehouse_manager"
RoleStorekeeper      UserRole = "storekeeper"
RolePharmacist       UserRole = "pharmacist"

// RefreshToken — активная сессия
type RefreshToken struct {
    ID        uuid.UUID
    UserID    int
    TokenHash string       // SHA-256 от raw token
    ExpiresAt time.Time
    UserAgent string
    IPAddress string
    Metadata  map[string]any
    CreatedAt time.Time
    RevokedAt *time.Time   // nil = активна
}

// EmployeeProfile — профиль сотрудника
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
    GDPTrainingHistory []GDPTrainingRecord  // JSONB массив
}
```

### Ошибки (`errors.go`)

Sentinel-ошибки (используются с `errors.Is`):

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

`AppError` — структурированная ошибка с кодом, сообщением и `slog.Attr` для логирования.

---

## Слой Repository (`internal/repository/`)

### Подключение к PostgreSQL (`postgres.go`)

```go
// GORM + pgx driver
db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
```

### DAO-паттерн

Между доменными моделями и GORM стоят **DAO** в `repository/dao/` с GORM-тегами:

- `UserDAO` → таблица `users`
- `SessionDAO` → таблица `refresh_tokens`
- `EmployeeProfileDAO` → таблица `employee_profiles`

Репозитории конвертируют `DAO ↔ domain` через методы `toDomain()` / `fromDomain()`.

### `UserRepository` (`user.go`)

| Метод                | SQL / GORM операция                          |
|----------------------|----------------------------------------------|
| `FindByID`           | `SELECT * FROM users WHERE id = ?`            |
| `FindByEmail`        | `WHERE email = ?`                             |
| `FindByGoogleID`     | `WHERE google_id = ?`                         |
| `FindByTelegramID`   | `WHERE telegram_id = ?`                       |
| `IsEmailTaken`       | `COUNT(*) WHERE email = ?`                    |
| `Create`             | `INSERT INTO users`                           |
| `UpdateRole`         | `UPDATE users SET role = ? WHERE id = ?`      |
| `LinkGoogle`         | `UPDATE users SET google_id = ?`              |
| `LinkTelegram`       | `UPDATE users SET telegram_id = ?`            |
| `SetNsPvAccess`      | `UPDATE users SET ns_pv_access = ?`           |
| `SetBlocked`         | `UPDATE users SET is_blocked = ?`             |
| `FindProfileByUserID`| JOIN users + employee_profiles                |

### `SessionRepository` (`session.go`)

| Метод                | Описание                                        |
|----------------------|-------------------------------------------------|
| `Create`             | Сохранить новую сессию (refresh token hash)     |
| `FindByID`           | Найти сессию по UUID                            |
| `FindByTokenHash`    | Найти сессию по SHA-256 хэшу токена             |
| `FindActiveByUserID` | Все не истёкшие и не отозванные сессии          |
| `Revoke`             | SET revoked_at = NOW() WHERE id = ?             |
| `RevokeAllForUser`   | SET revoked_at = NOW() WHERE user_id = ?        |
| `DeleteExpired`      | DELETE WHERE expires_at < NOW() (cleanup job)   |

---

## Слой Service (`internal/service/`)

### `Service` (`service.go`)

Агрегатор сервисов. Создаётся один раз в `main.go`:

```go
service.NewService(repos, jwtSecret, googleClientID)
// Access Token TTL  = 15 минут
// Refresh Token TTL = 15 дней
```

### `AuthService` (`auth_service.go`)

| Метод               | Описание                                                                  |
|---------------------|---------------------------------------------------------------------------|
| `LoginWithGoogle`   | idtoken.Validate → найти/создать user → issueTokens                       |
| `LoginWithTelegram` | Заглушка (`panic("implement...")`), нужна верификация hash по bot token    |
| `RefreshTokens`     | HashToken → FindByTokenHash → проверка theft → Revoke old → issueTokens  |
| `RevokeSession`     | Проверка владельца или admin → sessionRepo.Revoke                         |
| `AssignRole`        | Только admin → userRepo.UpdateRole                                        |
| `SetBlocked`        | Только admin → userRepo.SetBlocked                                        |

**Приватный метод `issueTokens`:**
1. Проверяет что пользователь не заблокирован
2. Генерирует raw refresh token (crypto/rand) + его SHA-256 хэш
3. Создаёт запись `RefreshToken` в БД
4. Генерирует JWT access token с `session_id = RefreshToken.ID`

### `TokenService` (`token_service.go`)

- `GenerateAccessToken(user, sessionID)` → JWT HS256
- `ParseAccessToken(token)` → `TokenClaims{UserID, Role, Email, SessionID}`
- `GenerateRefreshToken()` → (rawToken, hashToken, error)
- `HashToken(raw)` → SHA-256 hex

### `EmployeeProfileService` (`employee_profile_service.go`)

CRUD операций профиля сотрудника. Проверяет права вызывающего (admin или сам пользователь).
