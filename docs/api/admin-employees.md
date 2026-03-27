# Admin — Профили сотрудников

## GET `/api/v1/admin/employees`

> Список профилей сотрудников с базовой информацией.

**Auth:** `admin`  
**Route file:** `internal/handler/admin_employee_handler.go` → `adminListEmployeeProfiles()`  

---

## Требования к БД

### Таблица `users` + `employee_profiles`

```sql
SELECT u.id, u.login, u.email, u.role, u.status, u.is_blocked,
       ep.employee_code, ep.full_name, ep.position, ep.department
FROM users u
LEFT JOIN employee_profiles ep ON ep.user_id = u.id
ORDER BY u.id
```

---

## Функции в программе

### Handler

```go
// Файл: internal/handler/admin_employee_handler.go
func (h *Handler) adminListEmployeeProfiles(w http.ResponseWriter, r *http.Request)
```

### Service

```go
// Файл: internal/service/employee_profile_service.go
// domain interface: domain.EmployeeProfileService
func (s *employeeProfileService) ListAll(ctx) ([]*domain.UserProfile, error)
```

### Repository

```go
// Файл: internal/repository/employee_profile.go
// domain interface: domain.EmployeeProfileRepository
employeeProfileRepo.List(ctx) ([]*EmployeeProfile, error)
```

---

## GET `/api/v1/admin/employees/{userID}`

> Детальный профиль конкретного сотрудника.

**Auth:** `admin`  
**Route file:** `internal/handler/admin_employee_handler.go` → `adminGetEmployeeProfile()`  

### Требования к БД

```sql
SELECT ep.*, u.login, u.email, u.role, u.status
FROM employee_profiles ep
JOIN users u ON u.id = ep.user_id
WHERE ep.user_id = ?
```

### Функции в программе

```go
// Handler: chi.URLParam(r, "userID") → int
// Service/Repo: employeeProfileRepo.FindByUserID(ctx, userID)
// Существует: userRepo.FindProfileByUserID(ctx, userID) — internal/repository/user.go:192
```

---

## PATCH `/api/v1/admin/employees/{userID}`

> Обновить поля профиля сотрудника (частичное обновление).

**Auth:** `admin`  
**Route file:** `internal/handler/admin_employee_handler.go` → `adminPatchEmployeeProfile()`  

### Требования к БД

### Таблица `employee_profiles`

| Поле                    | Тип           | Описание                                     |
|-------------------------|---------------|----------------------------------------------|
| `user_id`               | INT (FK)       | Первичный ключ по внешнему ключу             |
| `employee_code`         | VARCHAR        | Табельный номер                              |
| `full_name`             | VARCHAR        | ФИО сотрудника                               |
| `position`              | VARCHAR        | Должность                                    |
| `department`            | VARCHAR        | Отдел                                        |
| `avatar_url`            | TEXT nullable  | URL аватарки                                 |
| `medical_book_scan_url` | TEXT nullable  | URL скана медкнижки                          |
| `special_zone_access`   | BOOLEAN        | Доступ в зону НС/ПВ                          |
| `gdp_training_history`  | JSONB          | История обучений GDP (массив дат/документов) |

```sql
UPDATE employee_profiles
SET full_name = COALESCE(?, full_name),
    position  = COALESCE(?, position),
    department = COALESCE(?, department),
    -- и т.д.
WHERE user_id = ?
```

### Функции в программе

```go
// Handler: парсить dto.PatchEmployeeProfileRequest из тела
// Service/Repo: employeeProfileRepo.Update(ctx, userID, changes)
// Файл: internal/repository/employee_profile.go
```
