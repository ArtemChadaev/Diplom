# POST `/auth/verify-code`

> Проверить OTP-код и получить пару токенов (access + refresh cookie).

**Auth:** Нет  
**Route file:** `internal/handler/auth_handler.go`  

---

## Требования к БД

### Таблица `users`

- Найти пользователя по `email` (ReadOnly)
- Статус должен быть `active`, иначе `403`

### Таблица `otp_codes` *(планируемая)*

- Найти активную запись: `SELECT * FROM otp_codes WHERE user_id = ? AND expires_at > NOW() AND revoked_at IS NULL`
- Инкрементировать `attempts++` при каждой проверке
- Если `attempts >= 5` → ответ `429 max_attempts_reached`
- Если код совпадает → пометить запись использованной (`revoked_at = NOW()`)

### Таблица `refresh_tokens`

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
2. Вызвать `service.Auth.VerifyOTPCode(ctx, email, code)` *(не реализовано)*
3. Установить `httpOnly` cookie с refresh token (аналогично существующему `login` handler'у)
4. Вернуть `200 { "access_token": "...", "user": UserDTO }`

### Service

```go
// interface: domain.AuthService
// method: VerifyOTPCode(ctx, email, code string) (*TokenPair, error)
```

**Шаги:**
1. `userRepo.FindByEmail(ctx, email)` → ErrUserNotFound
2. `otpRepo.FindActiveByUserID(ctx, user.ID)` → ErrSessionNotFound если нет активного
3. Сравнить код (bcrypt или константное время `subtle.ConstantTimeCompare`)
4. Если несовпадение: `otpRepo.IncrAttempts(ctx, otp.ID)` → если ≥ 5 → ErrMaxAttempts
5. Если совпадение: `otpRepo.Revoke(ctx, otp.ID)` → `issueTokens(ctx, user, userAgent, ip)`

### Используемые существующие функции

- `authService.issueTokens()` — `internal/service/auth_service.go:209`
- `tokenSvc.GenerateAccessToken()` — `internal/service/token_service.go`
- `tokenSvc.GenerateRefreshToken()` — `internal/service/token_service.go`
- `sessionRepo.Create()` — `internal/repository/session.go:55`
