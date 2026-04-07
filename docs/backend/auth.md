# Backend — Аутентификация и сессии

> Часть документации `backend.md`. Описывает схему токенов, защиту от кражи сессий и OAuth-потоки.

---

## Схема токенов

| Токен          | Хранение             | TTL     | Назначение                              |
|----------------|----------------------|---------|-----------------------------------------|
| Access Token   | В памяти клиента     | 15 мин  | Авторизация запросов (Bearer)           |
| Refresh Token  | `httpOnly` Cookie    | 15 дней | Обновление access token                 |

**Access Token** — JWT, подписанный `HS256` с `JWT_SECRET`. Payload содержит:
- `user_id` — ID пользователя
- `role` — роль (`UserRole`)
- `email` — email
- `session_id` — UUID сессии (refresh token) для отзыва

**Refresh Token** — случайная строка (crypto/rand). В БД хранится **только хэш** (`SHA-256`), сам токен отдаётся клиенту в `httpOnly; Secure; SameSite=Strict` cookie по пути `/auth`.

---

## Роли пользователей

Роли определены как PostgreSQL ENUM `user_role`:

| Роль               | Описание                                      |
|--------------------|-----------------------------------------------|
| `admin`            | Администратор системы                         |
| `qp`               | Уполномоченное лицо (Qualified Person)        |
| `warehouse_manager`| Начальник склада                              |
| `storekeeper`      | Кладовщик                                     |
| `pharmacist`       | Провизор (роль по умолчанию при создании)     |

---

## Аутентификация через Google OAuth

Поток:
1. Клиент получает Google ID Token (через Google Sign-In SDK)
2. `POST /auth/google` → backend валидирует ID Token через `idtoken.Validate`
3. Если email найден в БД → выдаются токены
4. Если email не найден → создаётся новый пользователь с ролью `pharmacist` и выдаются токены

---

## Защита от кражи Refresh Token

При повторном использовании **уже отозванного** refresh token происходит **token theft detection**:
- `sessionRepo.RevokeAllForUser()` — отзываются **все** активные сессии пользователя
- Возвращается `401 Unauthorized`

---

## Ротация токенов

При каждом `/auth/refresh`:
1. Старая сессия помечается `revoked_at = NOW()`
2. Создаётся новая сессия с новым refresh token
3. Выдаётся новый access token

---

## OTP-коды (Valkey)

OTP-коды **не хранятся в PostgreSQL** — только в Valkey.
- Ключ: `otp:user:<user_id>` (Hash), TTL 600 сек
- Подробности: [valkey-cache.md](../valkey-cache.md)
