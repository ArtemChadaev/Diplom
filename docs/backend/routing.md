# Backend — HTTP-роутер и Middleware

> Часть документации `backend.md`. Описывает маршруты, middleware и структуру обработчиков.

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
POST /auth/google                — Google OAuth (ID-token → access_token)
POST /auth/refresh               — обновление токенов по cookie
POST /auth/logout                — выход (отзыв сессии + очистка cookie)

/api/v1/ (middleware: AuthRequired — проверка Bearer JWT)
  DELETE /sessions/{sessionID}   — отзыв своей сессии

  /admin/ (middleware: RequireRole("admin"))
    PATCH  /users/{id}/role      — назначение роли
    PATCH  /users/{id}/blocked   — блокировка/разблокировка пользователя
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

## Слой Handler (`internal/handler/`)

### `Handler` (`handler.go`)

Держит ссылки на `service.Service` и `domain.TokenService`. Метод `Router()` собирает chi-роутер.

### DTO (`internal/handler/dto/`)

Структуры запросов/ответов с JSON-тегами (изолированы от domain):

- `auth_dto.go` — `GoogleAuthRequest`, `TokenResponse`
- `user_dto.go` — `AssignRoleRequest`, `SetBlockedRequest`, `UserResponse`
- `employee_profile_dto.go` — `EmployeeProfileDTO`
- `admin_employee_dto.go` — `PatchEmployeeProfileRequest`, `EmployeeProfileResponse`

### Обработчики

| Файл                          | Обработчики                                                        |
|-------------------------------|---------------------------------------------------------------------|
| `auth_handler.go`             | `refresh`, `logout`, `clearRefreshCookie`                           |
| `oauth_handler.go`            | `googleLogin`                                                       |
| `admin_handler.go`            | `adminAssignRole`, `adminSetBlocked`, `revokeSession`, `adminRevokeSession` |
| `admin_employee_handler.go`   | `adminListEmployeeProfiles`, `adminGetEmployeeProfile`, `adminPatchEmployeeProfile` |
