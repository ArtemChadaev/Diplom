# POST `/auth/refresh`

> Обновление access token по refresh token из `httpOnly` cookie. Выполняет ротацию сессии: старый refresh token отзывается, выдаётся новый.

**Auth:** Нет (refresh token читается из cookie)  
**Route file:** `internal/handler/auth_handler.go` → `refresh()`  

---

## Требования к БД

### Таблица `refresh_tokens`

**Чтение:**
```sql
SELECT * FROM refresh_tokens WHERE token_hash = ?
```
- Если не найден → `401 invalid session`
- Если `revoked_at IS NOT NULL` → **theft detection** → `RevokeAllForUser(userID)` → `401`
- Если `expires_at < NOW()` → `401 token expired`

**Отзыв старого токена:**
```sql
UPDATE refresh_tokens SET revoked_at = NOW() WHERE id = ?
```

**Создание нового токена:**
```sql
INSERT INTO refresh_tokens (user_id, token_hash, expires_at, user_agent, ip_address)
VALUES (?, ?, NOW() + interval '15 days', ?, ?)
```

### Таблица `users`

```sql
SELECT * FROM users WHERE id = ?  -- для получения актуальных данных пользователя
```

---

## Функции в программе

### Handler

```go
// Файл: internal/handler/auth_handler.go:92
func (h *Handler) refresh(w http.ResponseWriter, r *http.Request)
```

**Шаги:**
1. `r.Cookie("refresh_token")` → если нет → `401`
2. Собрать `domain.SessionMeta{UserAgent, IPAddress}`
3. `h.service.Auth.RefreshTokens(ctx, cookie.Value, meta)`
4. При ошибке: очистить cookie `MaxAge: -1` → `401`
5. Успех: установить новую cookie + `200 {"access_token": "...", "expires_in": 900}`

### Service

```go
// Файл: internal/service/auth_service.go:153
func (s *authService) RefreshTokens(ctx, oldRefreshToken string, meta SessionMeta) (*TokenPair, error)
```

**Шаги:**
1. `tokenSvc.HashToken(oldRefreshToken)` → получить SHA-256 хэш
2. `sessionRepo.FindByTokenHash(ctx, hash)` → ErrSessionNotFound → 401
3. Проверка `session.RevokedAt != nil` → **theft detected** → `sessionRepo.RevokeAllForUser(userID)` → ErrSessionNotFound
4. `time.Now().After(session.ExpiresAt)` → ErrTokenExpired
5. `sessionRepo.Revoke(ctx, session.ID)` — атомарный отзыв
6. `userRepo.FindByID(ctx, session.UserID)` — актуальные данные пользователя
7. `authService.issueTokens(ctx, user, meta.UserAgent, meta.IPAddress)` — новая пара

### Repository

```go
tokenSvc.HashToken(rawToken) string                      // SHA-256 hex
sessionRepo.FindByTokenHash(ctx, hash) (*RefreshToken, error)
sessionRepo.RevokeAllForUser(ctx, userID) error          // theft protection
sessionRepo.Revoke(ctx, id) error                        // старый токен
userRepo.FindByID(ctx, userID) (*User, error)
sessionRepo.Create(ctx, rt) (*RefreshToken, error)       // новая сессия
```
