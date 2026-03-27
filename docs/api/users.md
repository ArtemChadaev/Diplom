# Users API

## GET `/users/me`

> Получить профиль текущего авторизованного пользователя.

**Auth:** Bearer (все роли)  
**Route file:** *(планируется)* `internal/handler/user_handler.go`  

### Требования к БД

#### Таблица `users` + `employee_profiles`

```sql
SELECT u.*, ep.*
FROM users u
LEFT JOIN employee_profiles ep ON ep.user_id = u.id
WHERE u.id = ?   -- из JWT claims.UserID
```

Реализовано через `userRepo.FindProfileByUserID(ctx, userID)` — `internal/repository/user.go:192`.

### Функции в программе

**Handler (планируется):**
```go
func (h *Handler) getMe(w http.ResponseWriter, r *http.Request) {
    userID := r.Context().Value(middleware.CtxUserID).(int)
    profile, err := h.service.User.GetProfile(ctx, userID)
    // → 200 UserDTO
}
```

**Repository (существует):**
```go
// Файл: internal/repository/user.go:192
userRepo.FindProfileByUserID(ctx, userID) (*domain.UserProfile, error)
// SELECT users + LEFT JOIN employee_profiles WHERE user_id = ?
```

---

## POST `/users/me/medical-book`

> Загрузить PDF/фото медицинской книжки.

**Auth:** Bearer (все роли)  
**Content-Type:** `multipart/form-data`

### Требования к БД

#### Таблица `employee_profiles`

```sql
UPDATE employee_profiles SET medical_book_scan_url = ? WHERE user_id = ?
```

| Поле                   | Действие                                          |
|------------------------|---------------------------------------------------|
| `medical_book_scan_url`| URL в S3/локальное хранилище после загрузки файла |

### Функции в программе

**Высокоуровневый порядок:**
1. Принять файл из `multipart/form-data`
2. Сохранить файл в хранилище (S3 или локально), получить URL
3. `employeeProfileRepo.Update(ctx, userID, &EmployeeProfile{MedicalBookScanURL: url})`

---

## GET `/users` *(admin)*

> Список всех пользователей с фильтрацией и пагинацией.

**Auth:** `admin`

### Требования к БД

```sql
SELECT u.*, ep.*
FROM users u
LEFT JOIN employee_profiles ep ON ep.user_id = u.id
WHERE (u.login ILIKE ? OR u.email ILIKE ? OR ep.full_name ILIKE ?)  -- ?q
  AND (u.role = ?)  -- ?role
ORDER BY u.created_at DESC
LIMIT ? OFFSET ?   -- ?page, ?limit
```

### Функции в программе

**Repository (нужно добавить):**
```go
userRepo.List(ctx, filter UserListFilter) ([]*UserProfile, int, error)
// filter: { Query, Role, Page, Limit }
```

---

## GET `/users/:id` *(admin)*

> Профиль конкретного пользователя.

### Требования к БД

```sql
-- Аналогично /users/me, но по произвольному id
SELECT u.*, ep.* FROM users u LEFT JOIN employee_profiles ep ON ep.user_id = u.id WHERE u.id = ?
```

**Repository:** `userRepo.FindProfileByUserID(ctx, id)` — `internal/repository/user.go:192`

---

## PATCH `/users/:id` *(admin)*

> Обновить роль / допуски сотрудника.

### Требования к БД

#### Таблица `users`

```sql
UPDATE users SET role = ? WHERE id = ?
-- если role передан
```

#### Таблица `employee_profiles`

```sql
UPDATE employee_profiles SET special_zone_access = ?, ns_pv_access = ? WHERE user_id = ?
-- если ns_pv_access / special_zone_access переданы
```

**Реализованные функции:**
```go
userRepo.UpdateRole(ctx, userID, role)     // internal/repository/user.go:137
// + нужно добавить обновление employee_profiles
```

---

## POST `/users/:id/send-login-link` *(admin)*

> Выслать новый OTP-код для входа сотруднику (аналог `/auth/send-code`, триггер — администратор).

### Требования к БД

Аналогично [auth-send-code.md](./auth-send-code.md) — создание записи в `otp_codes` для указанного `user_id`.

### Функции в программе

```go
// Handler вызывает: service.Auth.SendOTPCode(ctx, user.Email)
// Предварительно: userRepo.FindByID(ctx, id) — получить email пользователя
```
