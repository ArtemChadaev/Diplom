# POST `/auth/login`

> Вход по логину и паролю. Возвращает access token в теле и refresh token в `httpOnly` cookie.

**Auth:** Нет  
**Route file:** `internal/handler/auth_handler.go` → `login()`  

---

## Требования к БД

### Таблица `users`

- `SELECT * FROM users WHERE login = ?` — найти пользователя
- Проверить `status = 'active'`, иначе `403 Forbidden` (`ErrUserUnverified` / `ErrUserBlocked`)
- Bcrypt-сравнение `password_hash`

### Таблица `refresh_tokens`

- `INSERT INTO refresh_tokens (id, user_id, token_hash, expires_at, user_agent, ip_address, metadata)`
- `id` = `gen_random_uuid()` (PostgreSQL)
- `token_hash` = SHA-256 от raw refresh token (не сам токен!)
- `expires_at` = NOW() + 15 дней
- `revoked_at` = NULL

---

## Функции в программе

### Handler

```go
// Файл: internal/handler/auth_handler.go:49
func (h *Handler) login(w http.ResponseWriter, r *http.Request)
```

**Шаги:**
1. Декодировать `dto.LoginRequest{Login, Password}`
2. Извлечь `User-Agent` и IP из запроса (для session metadata)
3. `h.service.Auth.LoginWithPassword(ctx, login, password, userAgent, ip)`
4. Установить `httpOnly` cookie: `refresh_token=<raw>; Path=/auth; HttpOnly; SameSite=Strict; MaxAge=15d`
5. Вернуть `200 dto.TokenResponse{AccessToken, ExpiresIn: 900}`

### Service

```go
// Файл: internal/service/auth_service.go:79
func (s *authService) LoginWithPassword(ctx, login, password, userAgent, ip string) (*TokenPair, error)
```

**Шаги:**
1. `userRepo.FindByLogin(ctx, login)` → любая ошибка → `ErrInvalidCreds` (защита от enumeration)
2. `bcrypt.CompareHashAndPassword(hash, password)` → ошибка → `ErrInvalidCreds`
3. `authService.issueTokens(ctx, user, userAgent, ip)` — приватная функция

### Приватная функция `issueTokens`

```go
// Файл: internal/service/auth_service.go:209
func (s *authService) issueTokens(ctx, u *User, userAgent, ip string) (*TokenPair, error)
```

1. Проверка `u.Status == StatusActive && !u.IsBlocked`
2. `tokenSvc.GenerateRefreshToken()` → (rawToken, tokenHash, err)
3. `sessionRepo.Create(ctx, &RefreshToken{UserID, TokenHash, ExpiresAt, UserAgent, IP})`
4. `tokenSvc.GenerateAccessToken(u, session.ID)` → JWT HS256

### Repository

```go
userRepo.FindByLogin(ctx, login)         // WHERE login = ?
sessionRepo.Create(ctx, refreshToken)   // INSERT INTO refresh_tokens
```
