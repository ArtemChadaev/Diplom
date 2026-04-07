# POST `/auth/verify-code`

> Проверить OTP-код и получить пару токенов (access + refresh cookie).

**Auth:** Нет  
**Route file:** `internal/handler/auth_handler.go`

---

## Хранилище

OTP-коды хранятся **исключительно в Valkey (Redis)** — таблицы `otp_codes` в БД нет и не планируется.

### Таблица `users` (PostgreSQL, только чтение)

- Найти пользователя по `email`
- Пользователь должен быть не заблокирован, иначе `403`

### Valkey-ключи (чтение / мутация)

| Ключ | Операция |
|------|----------|
| `otp:hash:{user_id}` | GET — получить сохранённый хеш |
| `otp:attempts:{user_id}` | INCR — увеличить счётчик при неверном коде |
| `otp:hash:{user_id}` + `otp:attempts:{user_id}` | DEL — удалить при успехе |

Если ключ `otp:hash:{user_id}` не существует — значит код истёк или никогда не запрашивался → `404 otp_not_found`.

### Таблица `refresh_tokens` (PostgreSQL, запись)

- После успешной верификации → создать новую запись сессии (см. `issueTokens` в `auth_service.go`)

---

## Функции в программе

### Handler

```go
// Имя (планируется): verifyCode(w http.ResponseWriter, r *http.Request)
// Файл: internal/handler/auth_handler.go
// Route: POST /auth/verify-code
```

**Шаги:**
1. Десериализовать `{ "email": "...", "code": "482910" }`
2. Вызвать `service.Auth.VerifyOTPCode(ctx, email, code, meta)` *(не реализовано)*
3. Установить `httpOnly` cookie с refresh token (аналогично существующему `login` handler'у)
4. Вернуть `200 { "access_token": "...", "expires_in": N }`

### Service

```go
// interface: domain.AuthService
// method: VerifyOTPCode(ctx, email, code string, meta SessionMeta) (*TokenPair, error)
```

**Шаги:**
1. `userRepo.FindByEmail(ctx, email)` → `ErrUserNotFound` → 404
2. Проверить `user.IsBlocked` → `ErrUserBlocked` → 403
3. `otpStore.Get(ctx, user.ID)` → если нет/истёк → `ErrOTPNotFound` → 404
4. Если `attempts >= 5` → `ErrOTPMaxAttempts` → 429
5. Сравнить: `subtle.ConstantTimeCompare(hmacSHA256(code, secret), storedHash)`
6. Если не совпало: `otpStore.IncrAttempts(ctx, user.ID)` → `ErrOTPInvalid` → 401
7. Если совпало: `otpStore.Delete(ctx, user.ID)` → `issueTokens(ctx, user, meta)`

### Используемые существующие функции

- `authService.issueTokens()` — `internal/service/auth_service.go:209`
- `tokenSvc.GenerateAccessToken()` — `internal/service/token_service.go`
- `tokenSvc.GenerateRefreshToken()` — `internal/service/token_service.go`
- `sessionRepo.Create()` — `internal/repository/session.go:55`
