# Form 2.2: Employee Profile

← [Back to Forms Index](./index.md) | [← OTP Verify](./01b-email-otp-verify.md) | [Inbound Receiving →](./03-inbound-receiving.md)

> **UI Spec only.** For API contracts, upload logic, and admin endpoints → [`forms-logic.md §3`](../forms-logic.md#3-employee-profile)

## Requirements Checklist

- [ ] Table `users`: `full_name`, `email`, `role`, `ukep_bound`, `ns_pv_access`
- [ ] Table `employee_profiles`: `medical_book_scan_url`, `gdp_training_history (jsonb)`, `special_zone_access: bool`
- [ ] Backend: `GET /api/v1/users/me` — current user data
- [ ] Backend: `PUT /api/v1/users/me` — update profile
- [ ] Backend: `POST /api/v1/users/me/medical-book` — upload medical book scan (multipart)
- [ ] Backend: `GET /api/v1/users/:id` — (admin) view any employee
- [ ] File storage (S3/MinIO) for medical book scans
- [ ] shadcn/ui additional components: `npx shadcn@latest add avatar tabs toast`

---

## UI

### Page Layout

- Page `/profile` — available to all authenticated users
- Page `/admin/users/:id` — admin only
- Split into tabs via `<Tabs>`

### shadcn/ui Components

| Component | Purpose |
|-----------|---------|
| `<Avatar>`, `<AvatarImage>`, `<AvatarFallback>` | Employee avatar |
| `<Tabs>`, `<TabsList>`, `<TabsTrigger>`, `<TabsContent>` | Section tabs |
| `<Card>` | Section containers |
| `<Form>`, `<FormField>`, etc. | Editable fields |
| `<Input>` | Name, email |
| `<Badge>` | Role, access status |
| `<Button>` | Save, upload file |
| `<Toast>` | Success notification |
| `<Separator>` | Section dividers |

---

## Admin Block

On the `/admin/users/:id` page, admin sees additionally:

| Element | Action |
|---------|--------|
| Toggle **"NS/PV Access"** | `PATCH /api/v1/admin/users/:id` |
| Toggle **"Special Zone Access"** | same endpoint |
| **Role selector** (`<Select>`) | `PATCH /api/v1/admin/users/:id` |
| Button **"Send new login code"** | `POST /api/v1/admin/users/:id/send-login-link` |

Toggles use shadcn/ui `<Switch>` (install: `npx shadcn@latest add switch`).

---

## Spec Reference

→ [Forms Index — Section 2.2 Employee Profile](./index.md#22-employee-profile--02-employee-profilemd)
→ Logic, API contracts, upload snippets: [`forms-logic.md §3`](../forms-logic.md#3-employee-profile)
