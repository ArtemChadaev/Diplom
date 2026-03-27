# Sessions API

## DELETE `/api/v1/sessions/{sessionID}`

> Отозвать конкретную сессию текущего пользователя.

**Auth:** Bearer (все роли)  
**Route file:** `internal/handler/admin_handler.go` → `revokeSession()`  

---

## Требования к БД

### Таблица `refresh_tokens`

```sql
SELECT * FROM refresh_tokens WHERE id = ?  -- проверка владельца
UPDATE refresh_tokens SET revoked_at = NOW() WHERE id = ?  -- отзыв
```

**Проверка:** `session.UserID` должен совпадать с `callerID` из JWT. Иначе `ErrInsufficientPerms`.

---

## Функции в программе

### Handler

```go
// Файл: internal/handler/admin_handler.go:83
func (h *Handler) revokeSession(w http.ResponseWriter, r *http.Request)
```

1. `chi.URLParam(r, "sessionID")` → `uuid.Parse(...)`
2. Извлечь `callerID` и `callerRole` из контекста
3. `h.service.Auth.RevokeSession(ctx, sessionID, callerID, callerRole)`
4. `204 No Content`

### Service

```go
// Файл: internal/service/auth_service.go:183
func (s *authService) RevokeSession(ctx, sessionID uuid.UUID, callerID int, callerRole UserRole) error
```

1. `sessionRepo.FindByID(ctx, sessionID)`
2. Если `sess.UserID != callerID && callerRole != RoleAdmin` → `ErrInsufficientPerms`
3. `sessionRepo.Revoke(ctx, sessionID)` → `SET revoked_at = NOW()`

---

## DELETE `/api/v1/admin/sessions/{sessionID}` *(admin)*

> Администратор может отозвать любую сессию любого пользователя.

**Auth:** `admin`  
**Route file:** `internal/handler/admin_handler.go` → `adminRevokeSession()`  

### Требования к БД

Идентично обычному `revokeSession`, но проверка роли уже выполнена middleware `RequireRole("admin")`, поэтому `callerRole == admin` всегда выполняется.

### Функции в программе

```go
// Файл: internal/handler/admin_handler.go:63
func (h *Handler) adminRevokeSession(w http.ResponseWriter, r *http.Request)
// Вызывает ту же h.service.Auth.RevokeSession(ctx, sessionID, callerID, callerRole)
// callerRole = admin → проверка владельца bypassed
```
