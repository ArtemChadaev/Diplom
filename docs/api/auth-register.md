# POST `/auth/register`

> Регистрация нового пользователя по логину и паролю. Аккаунт создаётся со статусом `unverified` и ожидает верификации администратором.

**Auth:** Нет  
**Route file:** `internal/handler/auth_handler.go` → `register()`  

---

## Требования к БД

### Таблица `users`

| Поле            | Действие при регистрации                          |
|-----------------|---------------------------------------------------|
| `login`         | Уникальный, проверяется `IsLoginTaken` перед INSERT |
| `email`         | Опциональный, уникальный, проверяется `FindByEmail` |
| `password_hash` | `bcrypt.GenerateFromPassword(password, cost=12)`   |
| `role`          | Устанавливается `unverified`                      |
| `status`        | Устанавливается `unverified`                      |
| `is_blocked`    | `false` (default)                                 |
| `created_at`    | `NOW()` (autoCreateTime GORM)                     |

**Constraints:** `UNIQUE(login)`, `UNIQUE(email)` — нарушение вернёт ошибку `409 Conflict`.

---

## Функции в программе

### Handler

```go
// Файл: internal/handler/auth_handler.go:16
func (h *Handler) register(w http.ResponseWriter, r *http.Request)
```

**Шаги:**
1. Декодировать `dto.RegisterRequest{Login, Email, Password}`
2. Базовая валидация: `login != ""` и `password != ""`
3. Вызов `h.service.Auth.Register(ctx, domain.RegisterInput{Login, Email, Password})`
4. `ErrLoginTaken` / `ErrEmailTaken` → `409 Conflict`
5. Успех → `201 Created` + `{"message": "account pending admin approval"}`

### Service

```go
// Файл: internal/service/auth_service.go:38
func (s *authService) Register(ctx, req RegisterInput) (*User, error)
```

**Шаги:**
1. `userRepo.IsLoginTaken(ctx, login)` → `ErrLoginTaken`
2. `userRepo.FindByEmail(ctx, email)` (если email задан) → `ErrEmailTaken`
3. `bcrypt.GenerateFromPassword([]byte(password), 12)`
4. `userRepo.Create(ctx, &User{..., Role: RoleUnverified, Status: StatusUnverified})`

### Repository

```go
// Файл: internal/repository/user.go
userRepo.IsLoginTaken(ctx, login) (bool, error)   // COUNT WHERE login = ?
userRepo.FindByEmail(ctx, email) (*User, error)    // WHERE email = ?
userRepo.Create(ctx, u) (*User, error)             // INSERT INTO users
```
