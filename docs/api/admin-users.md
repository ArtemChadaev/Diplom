# Admin — Верификация и Роли

## PATCH `/api/v1/admin/users/{id}/verify`

> Верифицировать нового пользователя: перевести статус `unverified → active`.

**Auth:** `admin`  
**Route file:** `internal/handler/admin_handler.go` → `adminVerifyUser()`  

---

## Требования к БД

### Таблица `users`

```sql
UPDATE users SET status = 'active' WHERE id = ?
```

| Поле     | До верификации | После верификации |
|----------|---------------|-------------------|
| `status` | `unverified`  | `active`          |

**Ожидаемый результат:** `RowsAffected > 0`, иначе `ErrUserNotFound`.

---

## Функции в программе

### Handler

```go
// Файл: internal/handler/admin_handler.go:15
func (h *Handler) adminVerifyUser(w http.ResponseWriter, r *http.Request)
```

1. `chi.URLParam(r, "id")` → `strconv.Atoi` → userID
2. Извлечь `callerID`, `callerRole` из контекста middleware
3. `h.service.Auth.VerifyUser(ctx, callerID, callerRole, userID)`
4. `200 {"status": "active"}`

### Service

```go
// Файл: internal/service/auth_service.go:194
func (s *authService) VerifyUser(ctx, adminID int, adminRole UserRole, targetUserID int) error
```

1. `adminRole != RoleAdmin` → `ErrInsufficientPerms` (дополнительная проверка, middleware уже проверила)
2. `userRepo.UpdateStatus(ctx, targetUserID, StatusActive)`

### Repository

```go
// Файл: internal/repository/user.go:148
func (r *userRepository) UpdateStatus(ctx, userID int, status UserStatus) error
// UPDATE users SET status = ? WHERE id = ?
// RowsAffected == 0 → ErrUserNotFound
```

---

## PATCH `/api/v1/admin/users/{id}/role`

> Назначить роль пользователю.

**Auth:** `admin`  
**Route file:** `internal/handler/admin_handler.go` → `adminAssignRole()`  

---

## Требования к БД

### Таблица `users`

```sql
UPDATE users SET role = ? WHERE id = ?
-- role ∈ {'admin', 'employee', 'unverified'}
-- В api-endpoints.md расширенный список: qp, warehouse_manager, storekeeper, pharmacist
```

> **Важно:** В текущем коде в `domain/user.go` определены только 3 роли. В `api-endpoints.md` описана более широкая ролевая модель — требуется синхронизация.

---

## Функции в программе

### Handler

```go
// Файл: internal/handler/admin_handler.go:36
func (h *Handler) adminAssignRole(w http.ResponseWriter, r *http.Request)
```

1. Парсить `userID` из URL
2. Декодировать `dto.AssignRoleRequest{Role string}` из тела
3. `h.service.Auth.AssignRole(ctx, callerID, callerRole, userID, domain.UserRole(req.Role))`
4. `200 {"role": "<новая роль>"}`

### Service

```go
// Файл: internal/service/auth_service.go:201
func (s *authService) AssignRole(ctx, adminID int, adminRole UserRole, targetUserID int, role UserRole) error
```

1. `adminRole != RoleAdmin` → `ErrInsufficientPerms`
2. `userRepo.UpdateRole(ctx, targetUserID, role)`

### Repository

```go
// Файл: internal/repository/user.go:137
func (r *userRepository) UpdateRole(ctx, userID int, role UserRole) error
// UPDATE users SET role = ? WHERE id = ?
```
