# Backend — Техническое описание для ИИ

> Этот документ описывает архитектуру, библиотеки и ключевые паттерны бэкенда.  
> Цель — дать контекст ИИ-ассистенту для работы с кодовой базой.  
> API-эндпоинты: [docs/api-endpoints.md](./api-endpoints.md) и папка [docs/api/](./api/)

---

## Стек технологий

| Компонент          | Библиотека / Инструмент                    | Версия  |
|--------------------|---------------------------------------------|---------|
| HTTP-роутер        | `github.com/go-chi/chi/v5`                  | v5.1.0  |
| ORM / БД-доступ    | `gorm.io/gorm` + `gorm.io/driver/postgres`  | v1.31.1 |
| JWT               | `github.com/golang-jwt/jwt/v5`              | v5.3.1  |
| Хэш паролей        | `golang.org/x/crypto/bcrypt`                | —       |
| UUID               | `github.com/google/uuid`                    | v1.6.0  |
| Конфигурация       | `github.com/ilyakaznacheev/cleanenv`         | v1.5.0  |
| Google OAuth       | `google.golang.org/api/idtoken`             | —       |
| Логирование        | stdlib `log/slog`                           | —       |
| Миграции БД        | SQL-файлы (`migrate/`), `golang-migrate`    | —       |
| СУБД               | PostgreSQL (pgx/v5 под капотом GORM)        | —       |

---

## Луковая архитектура (Onion Architecture)

Бэкенд разбит на 4 концентрических слоя. Зависимости идут только **снаружи внутрь**: обработчики зависят от сервисов, сервисы — от репозиториев (через интерфейсы домена), ядро домена не зависит ни от чего.

```
┌──────────────────────────────────────────────────────────────┐
│  Handler  (internal/handler)                                  │
│   HTTP-обработчики, DTOfor request/response, middleware       │
│                                                               │
│  ┌────────────────────────────────────────────────────────┐   │
│  │  Service  (internal/service)                           │   │
│  │   Бизнес-логика, оркестрация репозиториев              │   │
│  │                                                         │   │
│  │  ┌──────────────────────────────────────────────────┐  │   │
│  │  │  Repository  (internal/repository)               │  │   │
│  │  │   GORM-реализации, DAO-объекты, SQL-запросы      │  │   │
│  │  │                                                   │  │   │
│  │  │  ┌────────────────────────────────────────────┐  │  │   │
│  │  │  │  Domain  (internal/domain)                 │  │  │   │
│  │  │  │   Чистые модели, интерфейсы, ошибки        │  │  │   │
│  │  │  └────────────────────────────────────────────┘  │  │   │
│  │  └──────────────────────────────────────────────────┘  │   │
│  └────────────────────────────────────────────────────────┘   │
└──────────────────────────────────────────────────────────────┘
```

### Принцип инверсии зависимостей (DI)

Все слои общаются через **интерфейсы**, определённые в `internal/domain`:

- `domain.UserRepository` — интерфейс для работы с пользователями
- `domain.SessionRepository` — интерфейс для работы с сессиями (refresh-токены)
- `domain.AuthService` — интерфейс бизнес-логики аутентификации
- `domain.TokenService` — интерфейс для создания/валидации JWT
- `domain.EmployeeProfileService` / `EmployeeProfileRepository` — интерфейсы профилей

Конкретные реализации находятся в `repository/` и `service/`. Это позволяет мокировать их в тестах без изменения бизнес-кода.

---

## Точка входа (`cmd/main.go`)

Запускает приложение в строгом порядке:

1. **Логгер** (`internal/pkg/logger`) — `slog.Logger` с контекстом запроса
2. **Конфигурация** (`config.Load()`) — `cleanenv` читает `.env` → env → завершает с ошибкой если обязательные переменные отсутствуют
3. **Graceful shutdown** — `signal.NotifyContext` перехватывает `SIGINT`/`SIGTERM`
4. **База данных** — `repository.NewPostgresDB()` → GORM + pgx
5. **Слои** — `Repository → Service → Handler` (инъекция зависимостей вручную)
6. **HTTP-сервер** — `domain.Server` (обёртка над `net/http.Server`)
7. **Bootstrap** — `bootstrap.SeedAdmin()` создаёт первого admin-пользователя если его нет

---

## Конфигурация (`internal/config/config.go`)

Читается через `cleanenv`. Все env-переменные:

| Переменная       | Обязательна | Описание                                |
|------------------|:-----------:|-----------------------------------------|
| `PORT`           | нет (8080)  | Порт HTTP-сервера                       |
| `DB_HOST`        | **да**      | Хост PostgreSQL                         |
| `DB_PORT`        | нет (5432)  | Порт PostgreSQL                         |
| `DB_USER`        | **да**      | Пользователь БД                         |
| `DB_NAME`        | **да**      | Имя БД                                  |
| `DB_PASSWORD`    | **да**      | Пароль БД                               |
| `JWT_SECRET`     | **да**      | Секрет для подписи JWT                  |
| `ADMIN_USER`     | нет (admin) | Логин первого администратора            |
| `ADMIN_PASSWORD` | **да**      | Пароль первого администратора           |
| `GOOGLE_CLIENT_ID` | нет      | Client ID для Google OAuth              |

---

## HTTP-роутер (`chi`)

Пакет `go-chi/chi/v5` — лёгкий идиоматичный роутер для `net/http`.

### Глобальные middleware (в `handler.go`)

```go
r.Use(middleware.RequestID)           // X-Request-Id header
r.Use(middleware.Recoverer)           // Panic recovery → 500
r.Use(middleware.Timeout(60s))        // Таймаут запроса 60 сек
r.Use(middleware.Heartbeat("/ping"))  // Healthcheck endpoint GET /ping
r.Use(h.loggingMiddleware)            // Кастомный slog-логгер запросов
```

### Структура маршрутов

```
POST /auth/register              — регистрация (логин+пароль)
POST /auth/login                 — вход по паролю → access_token + cookie
POST /auth/google                — Google OAuth (ID-token → access_token)
POST /auth/refresh               — обновление токенов по cookie
POST /auth/logout                — выход (отзыв сессии + очистка cookie)

/api/v1/ (middleware: AuthRequired — проверка Bearer JWT)
  DELETE /sessions/{sessionID}   — отзыв своей сессии
  /admin/ (middleware: RequireRole("admin"))
    PATCH  /users/{id}/verify    — верификация пользователя
    PATCH  /users/{id}/role      — назначение роли
    DELETE /sessions/{sessionID} — отзыв любой сессии
    GET    /employees/           — список профилей сотрудников
    GET    /employees/{userID}   — профиль конкретного сотрудника
    PATCH  /employees/{userID}   — обновление профиля сотрудника
```

---

## Middleware (`internal/handler/middleware/`)

### `AuthRequired` (`auth_middleware.go`)

Извлекает JWT из заголовка `Authorization: Bearer <token>`, валидирует через `domain.TokenService`, кладёт в контекст:

```go
ctx.Value(middleware.CtxUserID) → int       // ID пользователя
ctx.Value(middleware.CtxRole)   → UserRole  // Роль пользователя
ctx.Value(middleware.CtxEmail)  → string    // Email
```

### `RequireRole` (`rbac_middleware.go`)

Проверяет роль из контекста. Принимает одну или несколько допустимых `UserRole`. Возвращает `403 Forbidden` если роль не совпадает.

---

## Аутентификация и Сессии

### Схема токенов

| Токен          | Хранение             | TTL    | Назначение                              |
|----------------|----------------------|--------|-----------------------------------------|
| Access Token   | В памяти клиента     | 15 мин | Авторизация запросов (Bearer)           |
| Refresh Token  | `httpOnly` Cookie    | 15 дней| Обновление access token                 |

**Access Token** — JWT, подписанный `HS256` с `JWT_SECRET`. Payload содержит:
- `user_id` — ID пользователя
- `role` — роль (`UserRole`)
- `email` — email
- `session_id` — UUID сессии (refresh token) для отзыва

**Refresh Token** — случайная строка (crypto/rand). В БД хранится **только хэш** (`SHA-256`), сам токен отдаётся клиенту в `httpOnly; Secure; SameSite=Strict` cookie по пути `/auth`.

### Защита от кражи Refresh Token

При повторном использовании **уже отозванного** refresh token происходит **token theft detection**:
- `sessionRepo.RevokeAllForUser()` — отзываются **все** активные сессии пользователя
- Возвращается `401 Unauthorized`

### Ротация

При каждом `/auth/refresh`:
1. Старая сессия помечается `revoked_at = NOW()`
2. Создаётся новая сессия с новым refresh token
3. Выдаётся новый access token

---

## Слой Domain (`internal/domain/`)

Чистые Go-структуры **без тегов ORM или JSON** — только бизнес-поля.

### Ключевые модели

```go
// User — пользователь системы
type User struct {
    ID           int
    Login        string
    Email        *string   // nullable
    GoogleID     *string   // nullable (OAuth)
    TelegramID   *int64    // nullable (OAuth, не реализован)
    PasswordHash *string   // nullable (social login без пароля)
    Role         UserRole
    Status       UserStatus
    IsBlocked    bool
    CreatedAt    time.Time
}

// Роли пользователей (определены в api-endpoints.md)
RoleAdmin      UserRole = "admin"
RoleEmployee   UserRole = "employee"
RoleUnverified UserRole = "unverified"

// Статусы
StatusUnverified = "unverified"   // после регистрации
StatusActive     = "active"       // после верификации admin'ом
StatusBlocked    = "blocked"      // заблокирован

// RefreshToken — запись об активной сессии
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
```

### Ошибки (`errors.go`)

Sentinel-ошибки (используются с `errors.Is`):

```go
ErrUserNotFound       = "user not found"
ErrLoginTaken         = "login already taken"
ErrEmailTaken         = "email already taken"
ErrInvalidCreds       = "invalid login or password"
ErrUserUnverified     = "account pending administrator verification"
ErrUserBlocked        = "account is blocked"
ErrTokenExpired       = "token expired"
ErrInvalidToken       = "invalid token"
ErrSessionNotFound    = "session not found or terminated"
ErrInsufficientPerms  = "insufficient permissions for this operation"
```

`AppError` — структурированная ошибка с кодом, сообщением и `slog.Attr` для логирования. Поддерживает `errors.Is` через `Unwrap()`.

---

## Слой Repository (`internal/repository/`)

### Подключение к PostgreSQL (`postgres.go`)

```go
// GORM + pgx driver
db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
```

DSN собирается из конфига. После подключения выполняется `sqlDB.Ping()` для проверки.

### DAO-паттерн

Между доменными моделями и GORM стоят **DAO (Data Access Object)** — структуры в `repository/dao/` с GORM-тегами:

- `UserDAO` → таблица `users`
- `SessionDAO` → таблица `refresh_tokens`
- `EmployeeProfileDAO` → таблица `employee_profiles`

Репозитории конвертируют `DAO ↔ domain` через методы `toDomain()` / `fromDomain()`. Это изолирует доменный слой от деталей ORM.

### Реализованные репозитории

#### `UserRepository` (`user.go`)

| Метод                    | SQL / GORM операция                          |
|--------------------------|----------------------------------------------|
| `FindByID`               | `SELECT * FROM users WHERE id = ?`            |
| `FindByLogin`            | `WHERE login = ?`                             |
| `FindByEmail`            | `WHERE email = ?`                             |
| `FindByGoogleID`         | `WHERE google_id = ?`                         |
| `FindByTelegramID`       | `WHERE telegram_id = ?`                       |
| `IsLoginTaken`           | `COUNT(*) WHERE login = ?`                    |
| `Create`                 | `INSERT INTO users`                           |
| `UpdateRole`             | `UPDATE users SET role = ? WHERE id = ?`      |
| `UpdateStatus`           | `UPDATE users SET status = ? WHERE id = ?`    |
| `LinkGoogle`             | `UPDATE users SET google_id = ?`              |
| `LinkTelegram`           | `UPDATE users SET telegram_id = ?`            |
| `SetPasswordHash`        | `UPDATE users SET password_hash = ?`          |
| `FindProfileByUserID`    | JOIN users + employee_profiles                |

#### `SessionRepository` (`session.go`)

| Метод                  | Описание                                        |
|------------------------|-------------------------------------------------|
| `Create`               | Сохранить новую сессию (refresh token hash)     |
| `FindByID`             | Найти сессию по UUID                            |
| `FindByTokenHash`      | Найти сессию по SHA-256 хэшу токена             |
| `FindActiveByUserID`   | Все не истёкшие и не отозванные сессии          |
| `Revoke`               | SET revoked_at = NOW() WHERE id = ?             |
| `RevokeAllForUser`     | SET revoked_at = NOW() WHERE user_id = ?        |
| `DeleteExpired`        | DELETE WHERE expires_at < NOW() (cleanup job)   |

---

## Слой Service (`internal/service/`)

### `Service` (`service.go`)

Агрегатор сервисов. Создаётся один раз в `main.go`:

```go
service.NewService(repos, jwtSecret, googleClientID)
// Настройки токенов:
//   Access Token TTL  = 15 минут
//   Refresh Token TTL = 15 дней
```

### `AuthService` (`auth_service.go`)

Бизнес-логика аутентификации. Методы:

| Метод               | Описание                                                                 |
|---------------------|--------------------------------------------------------------------------|
| `Register`          | Bcrypt(cost=12) пароля → создание user со статусом `unverified`          |
| `LoginWithPassword` | Bcrypt.Compare → issueTokens (если статус `active`)                      |
| `LoginWithGoogle`   | idtoken.Validate → найти/создать user → issueTokens или ErrUserUnverified|
| `RefreshTokens`     | HashToken → FindByTokenHash → проверка theft → Revoke old → issueTokens |
| `RevokeSession`     | Проверка владельца или admin → sessionRepo.Revoke                        |
| `VerifyUser`        | Только admin → userRepo.UpdateStatus(active)                             |
| `AssignRole`        | Только admin → userRepo.UpdateRole                                       |

**Приватный метод `issueTokens`:**
1. Проверяет статус пользователя (active, не заблокирован)
2. Генерирует raw refresh token (crypto/rand) + его SHA-256 хэш
3. Создаёт запись `RefreshToken` в БД
4. Генерирует JWT access token с `session_id = RefreshToken.ID`

### `TokenService` (`token_service.go`)

Реализация JWT:
- `GenerateAccessToken(user, sessionID)` → JWT HS256
- `ParseAccessToken(token)` → `TokenClaims{UserID, Role, Email, SessionID}`
- `GenerateRefreshToken()` → (rawToken, hashToken, error)
- `HashToken(raw)` → SHA-256 hex

### `EmployeeProfileService` (`employee_profile_service.go`)

CRUD операций профиля сотрудника.

---

## Слой Handler (`internal/handler/`)

### `Handler` (`handler.go`)

Держит ссылки на `service.Service` и `domain.TokenService`. Метод `Router()` собирает chi-роутер.

### DTO (`internal/handler/dto/`)

Структуры запросов/ответов с JSON-тегами (изолированы от domain):

- `auth_dto.go` — `RegisterRequest`, `LoginRequest`, `TokenResponse`
- `user_dto.go` — `AssignRoleRequest`
- `employee_profile_dto.go` — ответ профиля
- `admin_employee_dto.go` — запросы/ответы для admin-управления профилями

### Обработчики

| Файл                          | Обработчики                                                       |
|-------------------------------|-------------------------------------------------------------------|
| `auth_handler.go`             | `register`, `login`, `refresh`, `logout`, `clearRefreshCookie`    |
| `oauth_handler.go`            | `googleLogin`                                                     |
| `admin_handler.go`            | `adminVerifyUser`, `adminAssignRole`, `revokeSession`, `adminRevokeSession` |
| `admin_employee_handler.go`   | `adminListEmployeeProfiles`, `adminGetEmployeeProfile`, `adminPatchEmployeeProfile` |

---

## Миграции (`migrate/`)

SQL-файлы пронумерованы. Применяются через `golang-migrate` или вручную.

| Миграция | Таблица / изменение                    |
|----------|----------------------------------------|
| 000001   | Создание ENUM-типов PostgreSQL          |
| 000002   | `users`                                |
| 000003   | `categories`                           |
| 000004   | `medicaments`                          |
| 000005   | `warehouses`                           |
| 000006   | `stock_items`                          |
| 000007   | `stock_operations`                     |
| 000008   | `audit_logs`                           |
| 000009   | `employee_profiles`                    |
| 000010   | Расширение `users` (google_id, telegram_id, is_blocked, status) |
| 000011   | `refresh_tokens`                       |

Подробная схема БД — в [docs/database-schema.md](./database-schema.md).

---

## Логирование

Используется stdlib `log/slog` со структурированным выводом.

- `internal/pkg/logger` — обёртка: настройка формата (`APP_ENV=production` → JSON, иначе text), контекстное обогащение (`logger.WithUserID(ctx, id)`)
- `logger.FromContext(ctx)` — извлекает логгер из контекста (создаётся в `AuthRequired` middleware)
- Кастомный `loggingMiddleware` в `Handler` логирует каждый запрос: метод, путь, статус, latency, request_id

---

## Что ещё запланировано / не реализовано

- `LoginWithTelegram` — метод-заглушка (`panic("implement...")`), нужно реализовать верификацию hash по bot token
- Valkey/Redis кэш сессий — описан в [docs/valkey-cache.md](./valkey-cache.md), в коде отсутствует
- Продвинутые бизнес-модули (Inbound, Orders, Inventory и т.д.) — только спроектированы в [api-endpoints.md](./api-endpoints.md), бэкенд не реализован
