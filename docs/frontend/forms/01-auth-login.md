# Form 2.1: Authentication — Login

← [Back to Forms Index](./index.md) | [Next: OTP Verify →](./01b-email-otp-verify.md)

> **UI Spec only.** For Zod schema, API contracts, and error handling → [`forms-logic.md §1`](../forms-logic.md#1-auth--login)

## Requirements Checklist

- [ ] Table `users`, fields: `email`, `role`, `ns_pv_access`, `ukep_bound`
- [ ] Valkey: OTP code at `otp:user:<id>` (Hash), TTL 600 sec — see [valkey-cache.md](../../valkey-cache.md)
- [ ] Backend: `POST /api/v1/auth/send-code` — send 6-digit code to email
- [ ] Backend: `POST /api/v1/auth/verify-code` — verify code, issue tokens
- [ ] Backend: `GET /api/v1/auth/google` — OAuth redirect
- [ ] Backend: `GET /api/v1/auth/google/callback` — receive token after Google
- [ ] Codes: 6 digits, TTL 10 minutes, single-use
- [ ] Rate limit: max 5 send-code requests per hour
- [ ] shadcn/ui components needed: `npx shadcn@latest add form alert`

---

## UI

### shadcn/ui Components

| Component | Purpose |
|-----------|---------|
| `<Card>`, `<CardHeader>`, `<CardContent>` | Form container |
| `<Form>` | Wrapper (react-hook-form + zod) |
| `<FormField>`, `<FormItem>`, `<FormLabel>`, `<FormControl>`, `<FormMessage>` | Email field |
| `<Input>` | Email address |
| `<Button>` | "Get code", "Login with Google" |
| `<Alert>`, `<AlertDescription>` | Errors (invalid email, rate limit) |
| `<Separator>` | Divider between OAuth and email sections |

---

## Admin Block

No admin-specific elements on this form.
After login with role `admin` → redirect to `/admin`.

---

## Spec Reference

→ [Forms Index — Section 2.1 Auth](./index.md#21-auth--login--01-auth-loginmd)
→ Logic, Zod schema, API contracts: [`forms-logic.md §1`](../forms-logic.md#1-auth--login)
