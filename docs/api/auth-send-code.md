# POST `/auth/send-code`

> Отправить 6-значный OTP-код на email сотрудника.

**Auth:** Нет  
**Route file:** `internal/handler/auth_handler.go`  

---

## Требования к БД

### Таблица `users`

- Найти пользователя по `email` → `SELECT * FROM users WHERE email = ?`
- **Нет** в таблице → вернуть `404 user_not_found`

### Таблица `otp_codes` *(планируемая, ещё не создана)*

| Поле         | Тип           | Описание                               |
|--------------|---------------|----------------------------------------|
| `id`         | UUID / SERIAL  | PK                                    |
| `user_id`    | INT            | FK → users.id                         |
| `code`       | VARCHAR(6)     | 6-значный цифровой код                |
| `expires_at` | TIMESTAMPTZ    | Текущее время + 10 мин (TTL = 600 с)  |
| `attempts`   | SMALLINT       | Счётчик неверных попыток              |
| `created_at` | TIMESTAMPTZ    | Время создания                        |

**Rate-limit:** перед созданием кода проверить, нет ли существующего кода для этого `user_id` созданного < 1 часа назад и уже превышающего попытки (`attempts >= 5`). Если да — вернуть `429`.

---

## Функции в программе

### Handler

```go
// Имя (планируется): sendCode(w http.ResponseWriter, r *http.Request)
// Файл: internal/handler/auth_handler.go
// Route: POST /auth/send-code
```

**Шаги:**
1. Десериализовать `{ "email": "..." }` из тела запроса
2. Вызвать `service.Auth.SendOTPCode(ctx, email)` *(не реализовано)*
3. Вернуть `200 { "message": "code_sent", "expires_in": 600 }`

### Service

```go
// interface: domain.AuthService
// method: SendOTPCode(ctx, email) error
```

**Шаги:**
1. `userRepo.FindByEmail(ctx, email)` → `ErrUserNotFound` → 404
2. Rate-limit check по `user_id` из хранилища OTP (Redis или БД)
3. Сгенерировать 6 цифр: `crypto/rand → fmt.Sprintf("%06d", n%1000000)`
4. Сохранить хэш кода в `otp_codes` с TTL 10 мин
5. Отправить email через SMTP / сервис рассылки

### Repository

```go
// Нужно реализовать: OTPRepository
// методы: Create(ctx, otp), FindActiveByUserID(ctx, userID), IncrAttempts(ctx, id)
```
