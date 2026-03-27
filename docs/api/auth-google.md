# GET `/auth/google` и POST `/auth/google`

> Google OAuth: клиент получает Google ID Token (через Google Sign-In SDK) и передаёт его на бэкенд. Бэкенд валидирует токен и выдаёт пару токенов или сообщает о необходимости верификации.

**Auth:** Нет  
**Route file:** `internal/handler/oauth_handler.go` → `googleLogin()`  
**Route mapping:** `POST /auth/google`

> **Примечание:** В `api-endpoints.md` описаны `GET /auth/google` (редирект) и `GET /auth/google/callback` (OAuth code flow). В реальном коде реализован **ID Token flow** — клиент сам получает ID-токен от Google и посылает его на `POST /auth/google`.

---

## Требования к БД

### Таблица `users`

**Сценарий A: пользователь уже зарегистрирован через Google**
- `SELECT * FROM users WHERE google_id = ?`

**Сценарий B: есть аккаунт с таким email, но ещё не привязан к Google**
- `SELECT * FROM users WHERE email = ?`
- `UPDATE users SET google_id = ? WHERE id = ?` (привязка)

**Сценарий C: новый пользователь**
- `INSERT INTO users (login, email, google_id, role='unverified', status='unverified')`
- Токены НЕ выдаются → возвращается `ErrUserUnverified`

### Таблица `refresh_tokens`

- При успешном входе (статус `active`) → `INSERT` новой сессии (как в `/auth/login`)

---

## Функции в программе

### Handler

```go
// Файл: internal/handler/oauth_handler.go
func (h *Handler) googleLogin(w http.ResponseWriter, r *http.Request)
```

**Шаги:**
1. Декодировать тело `{ "id_token": "eyJ..." }`
2. Взять User-Agent и IP
3. `h.service.Auth.LoginWithGoogle(ctx, idToken, userAgent, ip)`
4. `ErrUserUnverified` → `403 {"error": "account_pending_verification"}`
5. Успех → cookie + `200 {"access_token": "...", "user": UserDTO}`

### Service

```go
// Файл: internal/service/auth_service.go:92
func (s *authService) LoginWithGoogle(ctx, idToken, userAgent, ip string) (*TokenPair, error)
```

**Шаги:**
1. `idtoken.Validate(ctx, idToken, s.googleClientID)` — Google SDK валидирует токен
2. Извлечь `googleID = payload.Subject`, `email = payload.Claims["email"]`
3. `userRepo.FindByGoogleID(ctx, googleID)` → найден → `issueTokens`
4. `userRepo.FindByEmail(ctx, email)` → найден → `userRepo.LinkGoogle(...)` → `issueTokens`
5. Не найден → создать нового user со статусом `unverified` → `ErrUserUnverified`

### Repository

```go
userRepo.FindByGoogleID(ctx, googleID) (*User, error)
userRepo.FindByEmail(ctx, email) (*User, error)
userRepo.LinkGoogle(ctx, userID, googleID) error  // UPDATE SET google_id = ?
userRepo.Create(ctx, u) (*User, error)
sessionRepo.Create(ctx, rt) (*RefreshToken, error)
```
