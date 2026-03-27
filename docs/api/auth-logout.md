# POST `/auth/logout`

> Выйти из системы: отозвать сессию на сервере и очистить `httpOnly` cookie с refresh token.

**Auth:** Bearer (опционально — если заголовок есть, сессия отзывается; если нет — только очищается cookie)  
**Route file:** `internal/handler/auth_handler.go` → `logout()`  

---

## Требования к БД

### Таблица `refresh_tokens`

**При наличии валидного Bearer токена:**
```sql
UPDATE refresh_tokens SET revoked_at = NOW() WHERE id = ?
-- WHERE id = session_id из JWT claims
```

**Если Bearer токен невалиден/отсутствует:** запрос в БД не выполняется — просто очищается cookie.

**Проверка владельца сессии** (`RevokeSession`):
```sql
SELECT * FROM refresh_tokens WHERE id = ?
-- если session.UserID != callerID AND callerRole != 'admin' → ErrInsufficientPerms
```

---

## Функции в программе

### Handler

```go
// Файл: internal/handler/auth_handler.go:136
func (h *Handler) logout(w http.ResponseWriter, r *http.Request)
```

**Шаги:**
1. Если `Authorization` заголовок отсутствует → только `clearRefreshCookie(w)` → `204`
2. Парсинг `Bearer <token>` из заголовка
3. `tokenSvc.ParseAccessToken(token)` → если ошибка → `clearRefreshCookie` → `204` (logout всё равно выполнен)
4. `service.Auth.RevokeSession(ctx, claims.SessionID, claims.UserID, claims.Role)`
5. При ошибке revoke: логировать `WARN`, но продолжить (не возвращать ошибку клиенту)
6. `clearRefreshCookie(w)` → `204 No Content`

### Вспомогательная функция `clearRefreshCookie`

```go
// Файл: internal/handler/auth_handler.go:173
func (h *Handler) clearRefreshCookie(w http.ResponseWriter)
// Set-Cookie: refresh_token=; Path=/auth; HttpOnly; MaxAge=-1
```

### Service

```go
// Файл: internal/service/auth_service.go:183
func (s *authService) RevokeSession(ctx, sessionID uuid.UUID, callerID int, callerRole UserRole) error
```

**Шаги:**
1. `sessionRepo.FindByID(ctx, sessionID)` — найти сессию
2. Проверить: `sess.UserID == callerID || callerRole == RoleAdmin`
3. Иначе → `ErrInsufficientPerms`
4. `sessionRepo.Revoke(ctx, sessionID)` → `SET revoked_at = NOW()`

### Repository

```go
sessionRepo.FindByID(ctx, id) (*RefreshToken, error)
sessionRepo.Revoke(ctx, id) error   // UPDATE SET revoked_at = NOW()
```
