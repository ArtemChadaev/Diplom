# POST `/auth/send-code`

> Отправить 6-значный OTP-код на email сотрудника.

**Auth:** Нет  
**Route file:** `internal/handler/auth_handler.go`

---

## Хранилище

OTP-коды хранятся **исключительно в Valkey (Redis)** — БД не используется.  
TTL истечения управляется самим Valkey автоматически.

### Таблица `users` (PostgreSQL, только чтение)

- Найти пользователя по `email` → `SELECT * FROM users WHERE email = ?`
- **Нет** в таблице → вернуть `404 user_not_found`

### Valkey-ключи

| Ключ | Значение | TTL | Назначение |
|------|----------|-----|------------|
| `otp:hash:{user_id}` | HMAC-SHA256 hex кода | 600 с | Хранит хеш для верификации |
| `otp:attempts:{user_id}` | INT (счётчик) | 600 с | Счётчик неверных попыток |
| `otp:cooldown:{user_id}` | `1` | 60 с | Rate-limit: не даёт запросить новый код чаще раза в минуту |

**Rate-limit:** перед созданием кода проверить наличие `otp:cooldown:{user_id}`.  
Если ключ существует → вернуть `429 too_many_requests`.

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
2. Rate-limit: `otpStore.HasCooldown(ctx, user.ID)` → если есть → `ErrOTPCooldown` → 429
3. Сгенерировать 6 цифр: `crypto/rand → fmt.Sprintf("%06d", n%1000000)`
4. Хеш: `hmacSHA256(code, cfg.OTPHMACSecret)` → `hex.EncodeToString`
5. `otpStore.Set(ctx, user.ID, hash, 600s)` — сохранить хеш в Valkey
6. `otpStore.SetCooldown(ctx, user.ID, 60s)` — выставить кулдаун
7. Отправить email через UniSender

### OTP Store (Valkey)

```go
// Нужно реализовать: domain.OTPStore (интерфейс)
// Файл: internal/repository/cache/otp_store.go
// методы:
//   Set(ctx, userID, hash string, ttl time.Duration) error
//   Get(ctx, userID int) (hash string, attempts int, error)
//   IncrAttempts(ctx, userID int) (int, error)
//   Delete(ctx, userID int) error
//   HasCooldown(ctx, userID int) (bool, error)
//   SetCooldown(ctx, userID int, ttl time.Duration) error
```
