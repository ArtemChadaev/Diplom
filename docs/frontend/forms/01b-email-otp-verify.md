# Form 2.1b: Email OTP Verification

← [Back to Forms Index](./index.md) | [← Auth Login](./01-auth-login.md) | [Employee Profile →](./02-employee-profile.md)

> **UI Spec only.** For Zod schema, API contracts, timer logic, and error handling → [`forms-logic.md §2`](../forms-logic.md#2-auth--email-otp-verify)

## Requirements Checklist

- [ ] Depends on form 01 (Auth Login): user has already entered email and clicked "Get code"
- [ ] Valkey: key `otp:user:<user_id>` (Hash), TTL 600 sec — see [valkey-cache.md](../../valkey-cache.md)
- [ ] Backend: `POST /api/v1/auth/verify-code` — verify 6-digit code
- [ ] Backend: `POST /api/v1/auth/send-code` — resend code (from the same form)
- [ ] Code: 6 digits, TTL 10 minutes, single-use; after use → `DEL otp:user:<id>` in Valkey
- [ ] Max 3 wrong attempts (`attempts` field in Hash) → `DEL`, show "Send new code" button
- [ ] Countdown timer (10:00 → 0:00) until code expiry
- [ ] shadcn/ui additional component: `npx shadcn@latest add input-otp`

---

## UI

### shadcn/ui Components

| Component | Purpose |
|-----------|---------|
| `<InputOTP>`, `<InputOTPGroup>`, `<InputOTPSlot>` | 6 individual digit cells |
| `<Form>`, `<FormField>` | Wrapper with validation |
| `<Button>` | "Confirm", "Resend code" |
| `<Alert>` | Errors (wrong code, expired) |

---

## Admin Block

No admin-specific elements on this form.

---

## Spec Reference

→ [Forms Index — Section 2.1b OTP Verify](./index.md#21b-auth--email-otp-verify--01b-email-otp-verifymd)
→ Logic, Zod schema, API contracts: [`forms-logic.md §2`](../forms-logic.md#2-auth--email-otp-verify)
