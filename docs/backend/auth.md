# Backend — Authentication & Sessions

← [Back to Main README](../../README.md) | [Routing & Middleware →](./routing.md)

---

## Token Scheme

| Token | Storage | TTL | Purpose |
|-------|---------|-----|---------|
| Access Token | Client memory | 15 min | API authorization (Bearer header) |
| Refresh Token | `httpOnly` Cookie | 15 days | Access token renewal |

**Access Token** — JWT signed with `HS256` using `JWT_SECRET`. Payload contains:
- `user_id` — user ID
- `role` — user role (`UserRole`)
- `email` — email address
- `session_id` — refresh token session UUID (used for revocation)

**Refresh Token** — random string (`crypto/rand`). Only its **SHA-256 hash** is stored in the DB. The raw token is sent to the client via `httpOnly; Secure; SameSite=Strict` cookie scoped to `/auth`.

---

## User Roles

Roles are defined as a PostgreSQL ENUM `user_role`:

| Role | Description |
|------|-------------|
| `admin` | System administrator |
| `qp` | Qualified Person (QP) — authorizes batch release |
| `warehouse_manager` | Warehouse manager |
| `storekeeper` | Storekeeper |
| `pharmacist` | Pharmacist (default role on user creation) |

---

## Google OAuth Flow

1. Client obtains a Google ID Token (via Google Sign-In SDK)
2. `POST /auth/google` → backend validates the ID Token via `idtoken.Validate`
3. If email exists in DB → issue tokens
4. If email not found → create new user with role `pharmacist` → issue tokens

---

## Refresh Token Theft Detection

When a **revoked** refresh token is reused:
- `sessionRepo.RevokeAllForUser()` — **all** active sessions for the user are revoked
- Returns `401 Unauthorized`

---

## Token Rotation

On every `/auth/refresh`:
1. Old session is marked `revoked_at = NOW()`
2. New session is created with a new refresh token
3. New access token is issued

---

## OTP Codes (Valkey)

OTP codes are stored **only in Valkey**, never in PostgreSQL.
- Key: `otp:user:<user_id>` (Hash), TTL 600 sec
- Details: [../../docs/valkey-cache.md](../valkey-cache.md)
