# Forms — Business Logic Reference

← [Back to Frontend Overview](./README.md) | [Form UI Specs →](./forms/) | [Conventions →](./conventions.md)

> **Authority:** This is the single source of truth for all form **logic**: Zod schemas, API contracts, error handling, and post-submit behavior.
> UI design, component selection, and layout are documented in the individual [form files](./forms/).
> **API Contracts Note:** If you need exact API endpoints or schemas and they seem outdated here, always refer to the generated [`docs/api/swagger.json`](../../docs/api/swagger.json) as the ultimate source of truth for the API layer.

---

## Table of Contents

1. [Auth — Login](#1-auth--login)
2. [Auth — Email OTP Verify](#2-auth--email-otp-verify)
3. [Employee Profile](#3-employee-profile)
4. [Inbound Receiving](#4-inbound-receiving)
5. [Warehouse Zoning](#5-warehouse-zoning)
6. [Environment Log](#6-environment-log)
7. [Assembly & Shipment (FEFO)](#7-assembly--shipment-fefo)
8. [Claims & Defects](#8-claims--defects)
9. [Product Card](#9-product-card)
10. [Inventory](#10-inventory)

---

## Hook Location Convention

Every form's logic lives in a dedicated hook:

```
frontend/src/hooks/
├── use-auth-login.ts          → form 01
├── use-auth-login.test.ts
├── use-otp-verify.ts          → form 01b
├── use-otp-verify.test.ts
├── use-employee-profile.ts    → form 02
├── use-inbound-receiving.ts   → form 03
├── use-warehouse-zones.ts     → form 04
├── use-environment-log.ts     → form 05
├── use-order-assembly.ts      → form 06
├── use-claims.ts              → form 07
├── use-product-card.ts        → form 08
└── use-inventory.ts           → form 09
```

> Hook test file lives in the **same directory** as the hook (`use-xxx.test.ts`).
> Component test file lives in the **same directory** as the component (`my-form.test.tsx`).

---

## 1. Auth — Login

UI spec: [forms/01-auth-login.md](./forms/01-auth-login.md) | Hook: `hooks/use-auth-login.ts`

### Required Fields

- `email` — required, valid email format

### Zod Schema

```ts
const emailSchema = z.object({
  email: z.string().email("Invalid email address"),
})
```

### API Endpoints

**Request OTP code:**
```
POST /api/v1/auth/send-code
Body:    { email: string }
200:     { message: "Code sent", expires_in: 600 }
404:     { error: "User with that email not found" }
429:     { error: "Too many requests. Try again later" }
```

**Google OAuth:**
```
GET /api/v1/auth/google              → redirect to Google
GET /api/v1/auth/google/callback     → { access_token } + Set-Cookie: refresh_token
```

### User DTO

```ts
type UserDTO = {
  id: string
  email: string
  full_name: string
  role: "admin" | "qp" | "warehouse_manager" | "storekeeper" | "pharmacist"
  ns_pv_access: boolean
  ukep_bound: boolean
}
```

### Redirects After Successful Login

| Role | Redirect |
|------|----------|
| `admin` | `/admin/users` |
| `qp` | `/receiving` |
| `warehouse_manager` | `/warehouse` |
| `storekeeper`, `pharmacist` | `/search` |

### Error Handling

| HTTP | Display |
|------|---------|
| `404` | `<Alert>` "User with that email was not found" |
| `429` | `<Alert>` "Too many requests. Try again later" |

---

## 2. Auth — Email OTP Verify

UI spec: [forms/01b-email-otp-verify.md](./forms/01b-email-otp-verify.md) | Hook: `hooks/use-otp-verify.ts`

### Required Fields

- `code` — exactly 6 digits
- `email` — passed from previous step via URL query param `?email=...`

### Zod Schema

```ts
const otpSchema = z.object({
  code: z.string().length(6, "Code must be 6 digits").regex(/^\d{6}$/, "Digits only"),
})
```

### API Endpoints

**Verify code:**
```
POST /api/v1/auth/verify-code
Body:    { email: string, code: string }
200:     { access_token: string, expires_in: number }
         + Set-Cookie: refresh_token=...; HttpOnly; Secure; Path=/auth
400:     { error: "Invalid code", attempts_left: number }
410:     { error: "Code expired" }
429:     { error: "Max attempts reached" }
```

**Resend code:**
```
POST /api/v1/auth/send-code
Body: { email: string }
```

### Countdown Timer Logic

```ts
// Timer starts on component mount. expires_at comes from the send-code response.
const [secondsLeft, setSecondsLeft] = useState(
  Math.max(0, Math.floor((expiresAt - Date.now()) / 1000))
)

useEffect(() => {
  const interval = setInterval(() => {
    setSecondsLeft(s => Math.max(0, s - 1))
  }, 1000)
  return () => clearInterval(interval)
}, [])
```

"Resend" button is disabled while `secondsLeft > 0`.
When `secondsLeft === 0` → show "Code expired" and activate the button.

### Auto-Submit

When all 6 digits are entered, the form submits automatically:

```ts
<InputOTP maxLength={6} onComplete={() => form.handleSubmit(onSubmit)()} >
```

### Email State Between Steps

```ts
// In Login form after successful send-code:
router.push(`/auth/verify?email=${encodeURIComponent(email)}`)
```

### Error Handling

| HTTP | Display |
|------|---------|
| `400` with `attempts_left` | `<Alert>` "Wrong code. Attempts remaining: N" |
| `410` | `<Alert>` "Code expired. Request a new one" + activate Resend |
| `429` | `<Alert>` "Max attempts reached. Request a new code" + activate Resend |
| `200` | Store access token → redirect by role (see §1) |

---

## 3. Employee Profile

UI spec: [forms/02-employee-profile.md](./forms/02-employee-profile.md) | Hook: `hooks/use-employee-profile.ts`

### Required Fields

- All fields are read-only for the user themselves
- Admin can edit: `role`, `ns_pv_access`, `special_zone_access`

### API Endpoints

```
GET /api/v1/users/me
Response: {
  id, full_name, email, role, ukep_bound, ns_pv_access,
  profile: {
    medical_book_scan_url: string | null,
    special_zone_access: boolean,
    gdp_training_history: Array<{
      date: string,        // ISO 8601
      course_name: string,
      result: "pass" | "fail",
      certificate_url: string | null
    }>
  }
}

PUT /api/v1/users/me          → update own profile (non-role fields)
POST /api/v1/users/me/medical-book → multipart upload
GET /api/v1/users/:id         → (admin) view any employee
```

### Medical Book Upload

```ts
const onUpload = async (file: File) => {
  const formData = new FormData()
  formData.append("file", file)
  await fetch("/api/v1/users/me/medical-book", { method: "POST", body: formData })
}
```

### Role Label Mapping

```ts
const roleLabels: Record<string, string> = {
  admin: "Administrator",
  qp: "Qualified Person (QP)",
  warehouse_manager: "Warehouse Manager",
  storekeeper: "Storekeeper",
  pharmacist: "Pharmacist",
}
```

### Conditional Rendering

- `medical_book_scan_url === null` → show "Upload scan" button
- `medical_book_scan_url !== null` → show "View" (opens in new tab) + "Update"
- `gdp_training_history` → rendered as table

### Admin API

```
PATCH /api/v1/admin/users/:id   { ns_pv_access: boolean }
PATCH /api/v1/admin/users/:id   { role: RoleEnum }
PATCH /api/v1/admin/users/:id   { special_zone_access: boolean }
POST  /api/v1/admin/users/:id/send-login-link
```

---

## 4. Inbound Receiving

UI spec: [forms/03-inbound-receiving.md](./forms/03-inbound-receiving.md) | Hook: `hooks/use-inbound-receiving.ts`

### Required Fields

- `supplier_id`, `purchase_type`, `invoice_number`, `country_of_origin`, `manufacturer`
- `vat_rate`
- At least one position with: `product_id`, `serial_number`, `manufacture_date`, `expiry_date`, `quantity`

### Zod Schema

```ts
const inboundSchema = z.object({
  supplier_id: z.string().min(1),
  purchase_type: z.enum(["direct", "tender", "state"]),
  invoice_number: z.string().min(1),
  country_of_origin: z.string().min(1),
  manufacturer: z.string().min(1),
  vat_rate: z.enum(["0", "10", "20"]),
  is_jnvlp_controlled: z.boolean(),
  jnvlp_markup: z.number().min(0).max(100).optional(),
  positions: z.array(z.object({
    product_id: z.string().min(1),
    serial_number: z.string().min(1),
    manufacture_date: z.date(),
    expiry_date: z.date(),
    quantity: z.number().int().min(1),
  })).min(1, "Add at least one item"),
})
  .refine((d) => d.positions.every(p => p.expiry_date > p.manufacture_date), {
    message: "Expiry date must be after manufacture date",
    path: ["positions"],
  })
```

### API Endpoints

```
GET /api/v1/suppliers?q={query}
Response: Array<{ id, name, inn }>

GET /api/v1/ref/countries
Response: Array<{ code: string, name_en: string }>

GET /api/v1/products/search?q={mnn_or_sku}
Response: Array<{
  id, mnn, sku, ru_number,
  is_jnvlp, is_mdlp, is_ns_pv, cold_chain,
  temp_min, temp_max
}>

POST /api/v1/inbound
Body: InboundSchema
Response 201: { id: string, status: "quarantine" }

POST /api/v1/inbound/:id/quarantine-release     (admin / QP only)
Body: { qp_user_id, inspection_date, result: "approved"|"rejected", notes }
```

### Post-Selection Logic (product chosen from Combobox)

1. `is_jnvlp = true` → auto-enable `is_jnvlp_controlled` + show markup field
2. `is_ns_pv = true` → show warning Alert "NS/PV product. Ensure required permits are in place"
3. `cold_chain = true` → show Cold Chain badge on line item

### Status After Save

After `POST /api/v1/inbound`, all batches get `status: "quarantine"` — blocked from all movements.

### Quarantine Release (admin / QP)

```
POST /api/v1/inbound/:id/quarantine-release
Body: { qp_user_id, inspection_date, result: "approved"|"rejected", notes }
→ "approved": batches.status → "available"
→ "rejected": batches.status → "rejected" + show "Create return to supplier" button
```

---

## 5. Warehouse Zoning

UI spec: [forms/04-warehouse-zoning.md](./forms/04-warehouse-zoning.md) | Hook: `hooks/use-warehouse-zones.ts`

### Zod Schema (zone creation, admin only)

```ts
const zoneSchema = z.object({
  name: z.string().min(2),
  type: z.enum(["general", "cold_chain", "flammable", "safe_strong"]),
  description: z.string().optional(),
  temp_min: z.number().nullable(),
  temp_max: z.number().nullable(),
  humidity_max: z.number().min(0).max(100).nullable(),
})
```

### API Endpoints

```
GET /api/v1/zones
Response: Array<{
  id, name,
  type: "general" | "cold_chain" | "flammable" | "safe_strong",
  description, temp_min, temp_max, humidity_max,
  stock_count: number
}>

GET /api/v1/zones/:id/stock
Response: { zone: ZoneDTO, items: Array<{ product_id, mnn, serial_number, expiry_date, quantity, status }> }

POST /api/v1/zones       (admin)
PUT  /api/v1/zones/:id   (admin)
```

### Access Control Logic

- `safe_strong` zone: if `user.ns_pv_access === false` → show placeholder instead of stock data
- Zone creation and editing — `admin` only

### Conditional Rendering

- `temp_min === null && temp_max === null` → display "No temperature restriction"
- `humidity_max === null` → display "No humidity restriction"

---

## 6. Environment Log

UI spec: [forms/05-environment-log.md](./forms/05-environment-log.md) | Hook: `hooks/use-environment-log.ts`

### Required Fields

- `shift` — morning or evening
- `temperature`, `humidity` — must be within zone limits
- `zone_id` — from context (selected zone)

### Zod Schema (dynamic — receives zone config)

```ts
const envLogSchema = (zone: ZoneDTO) => z.object({
  shift: z.enum(["morning", "evening"]),
  temperature: z.number()
    .min(-50).max(100)
    .refine(
      (v) => zone.temp_min == null || v >= zone.temp_min,
      `Below minimum (${zone.temp_min}°C)`
    )
    .refine(
      (v) => zone.temp_max == null || v <= zone.temp_max,
      `Above maximum (${zone.temp_max}°C)`
    ),
  humidity: z.number().min(0).max(100)
    .refine(
      (v) => zone.humidity_max == null || v <= zone.humidity_max,
      `Humidity above limit (${zone.humidity_max}%)`
    ),
  notes: z.string().optional(),
})
```

### API Endpoints

```
GET /api/v1/environment-log/today
Response: Array<{
  zone_id, zone_name, zone_type,
  temp_min, temp_max, humidity_max,
  morning_log: EnvLogDTO | null,
  evening_log: EnvLogDTO | null
}>

POST /api/v1/zones/:id/environment-log
Body: { shift, temperature, humidity, notes? }

GET /api/v1/zones/:id/environment-log?from=YYYY-MM-DD&to=YYYY-MM-DD&page=1&limit=50

GET /api/v1/environment-log/export?from=...&to=...  (admin / QP only → Excel)
```

### Pending Entry Logic

- If `morning_log === null` and `Date.now() < 14:00 local` → Badge "Pending"
- If a value exceeds `temp_min/temp_max` or `humidity_max` → show `<Alert variant="destructive">` inline in form
- After save → `<Toast>` "Data saved. Recorded by: [Name]"

### Admin / QP

- Can **edit** an already-saved entry ("Correct" button) → `PUT /api/v1/zones/:id/environment-log/:log_id`
- Admin sees **all zones** including `safe_strong`

---

## 7. Assembly & Shipment (FEFO)

UI spec: [forms/06-assembly-shipment-fefo.md](./forms/06-assembly-shipment-fefo.md) | Hook: `hooks/use-order-assembly.ts`

### API Endpoints

```
POST /api/v1/orders
GET  /api/v1/orders?status=new&type=cito&page=1&limit=20
GET  /api/v1/orders/:id
Response: {
  id, type: "regular"|"cito",
  status: "new"|"assembling"|"ready"|"shipped",
  destination: { id, name, address },
  items: Array<{
    id, product_id, mnn, sku, requested_qty,
    fefo_recommendations: Array<{
      batch_id, serial_number,
      manufacture_date, expiry_date,
      available_qty,
      remaining_life_percent: number,
      is_mos_blocked: boolean
    }>
  }>
}

POST /api/v1/orders/:id/confirm-assembly
POST /api/v1/orders/:id/ship            (admin / qp / warehouse_manager)
GET  /api/v1/orders/:id/quality-registry
GET  /api/v1/settings/mos               → { mos_percent: number }
PUT  /api/v1/settings/mos               (admin) { mos_percent: number }
```

### MOS Calculation (client-side)

```ts
const calculateMos = (expiryDate: Date, manufactureDate: Date, mosPercent: number): boolean => {
  const totalLife = expiryDate.getTime() - manufactureDate.getTime()
  const remaining = expiryDate.getTime() - Date.now()
  const remainingPercent = (remaining / totalLife) * 100
  return remainingPercent >= mosPercent  // true = NOT blocked
}
```

### Business Rules

- Data sorted by `expiry_date ASC` server-side (FEFO). Client only displays.
- `is_mos_blocked = true` → confirm button disabled for that batch
- `type: "cito"` orders appear at the top of the list (sorted server-side)
- Shipment allowed only for: `admin`, `qp`, `warehouse_manager`
- `storekeeper` and `pharmacist` can only confirm assembly (line items)

### Post-Ship

```
POST /api/v1/orders/:id/ship
→ Server generates TTN PDF
→ GET /api/v1/orders/:id/quality-registry → quality certificate PDF
→ UI shows download links for both PDFs
```

---

## 8. Claims & Defects

UI spec: [forms/07-claims-defects.md](./forms/07-claims-defects.md) | Hook: `hooks/use-claims.ts`

### Required Fields

- `type`, `product_id`, `batch_id`

### Zod Schema

```ts
const claimSchema = z.object({
  type: z.enum(["recall", "return_from_pharmacy", "return_to_supplier", "defect"]),
  product_id: z.string().min(1, "Select a product"),
  batch_id: z.string().min(1, "Select a batch"),
  source: z.string().optional(),
  notes: z.string().optional(),
  photos: z.array(z.instanceof(File)).optional(),
})
```

### API Endpoints

```
GET /api/v1/claims?status=open&page=1&limit=20
Response: Array<{
  id, type, status, created_at,
  product: { id, mnn, sku },
  batch: { id, serial_number, expiry_date },
  source: string | null,
  photos_count: number,
  created_by: { full_name }
}>

POST /api/v1/claims
POST /api/v1/claims/:id/photos     → multipart (photos: File[])
POST /api/v1/claims/:id/close      (admin / QP)
GET  /api/v1/stop-signals          → { active: boolean, signals: Array<StopSignalDTO> }
GET  /api/v1/recalled-batches/sync (admin) → manual Roszdravnadzor sync
PATCH /api/v1/batches/:id/unblock  (admin) → { reason: string }
```

### Status Config

```ts
const statusConfig = {
  open:    { label: "Open",    variant: "warning" },
  blocked: { label: "BLOCKED", variant: "destructive" },
  closed:  { label: "Closed",  variant: "secondary" },
}
```

### STOP Signal (Automatic Block)

Background worker syncs `recalled_batches` → matching batch: `status → "blocked"`.
Client polls `GET /api/v1/stop-signals`. When `active: true` → global Alert rendered outside main layout.

### Photo Upload

```ts
const uploadPhotos = async (claimId: string, files: File[]) => {
  const formData = new FormData()
  files.forEach(f => formData.append("photos", f))
  await fetch(`/api/v1/claims/${claimId}/photos`, { method: "POST", body: formData })
}
```

### Unblock Rule

"Unblock" button visible to `admin` only, and only for batches where source is **not Roszdravnadzor**.
Unblock requires mandatory `reason` field.

---

## 9. Product Card

UI spec: [forms/08-product-card.md](./forms/08-product-card.md) | Hook: `hooks/use-product-card.ts`

### Required Fields

- `trade_name`, `mnn`, `ru_number` (registration #), `atc_codes` (at least one), `dosage_form`, `dosage`, `package_multiplicity`
- If `cold_chain = true` → `temp_min` and `temp_max` are required

### Zod Schema

```ts
const productSchema = z.object({
  trade_name: z.string().min(2),
  mnn: z.string().min(2),
  sku: z.string().optional(),
  barcode: z.string().optional(),
  datamatrix_gtin: z.string().optional(),
  ru_number: z.string().min(1, "Registration # is required"),
  atc_codes: z.array(z.string()).min(1, "At least one ATC code required"),
  dosage_form: z.string().min(1),
  dosage: z.string().min(1),
  package_multiplicity: z.number().int().min(1),
  is_jnvlp: z.boolean().default(false),
  is_mdlp: z.boolean().default(false),
  is_ns_pv: z.boolean().default(false),
  cold_chain: z.boolean().default(false),
  temp_min: z.number().nullable(),
  temp_max: z.number().nullable(),
  humidity_max: z.number().min(0).max(100).nullable(),
  weight_g: z.number().min(0).optional(),
  width_cm: z.number().min(0).optional(),
  height_cm: z.number().min(0).optional(),
  depth_cm: z.number().min(0).optional(),
})
.refine(
  (d) => !d.cold_chain || (d.temp_min !== null && d.temp_max !== null),
  { message: "Cold chain products must have a temperature range specified" }
)
```

### API Endpoints

```
GET /api/v1/products?q={query}&is_jnvlp=true&page=1&limit=20
GET /api/v1/products/:id
Response: {
  id, trade_name, mnn, sku, barcode, datamatrix_gtin,
  ru_number, atc_codes: string[],
  dosage_form, dosage, package_multiplicity,
  is_jnvlp, is_mdlp, is_ns_pv, cold_chain,
  temp_min, temp_max, humidity_max,
  weight_g, width_cm, height_cm, depth_cm,
  photos: Array<{ id, url, is_primary }>,
  batches: Array<{ id, serial_number, manufacture_date, expiry_date, quantity, zone_name, status }>
}

GET  /api/v1/ref/atc        → ATC code reference
POST /api/v1/products        (admin)
PUT  /api/v1/products/:id    (admin)
DELETE /api/v1/products/:id  (admin → soft delete via deleted_at)
POST /api/v1/products/:id/photos  → multipart
```

### Flag Rendering Logic

```ts
const flags = [
  { key: "is_jnvlp", label: "JNVLP",      variant: "outline" },
  { key: "is_mdlp",  label: "MDLP",       variant: "secondary" },
  { key: "is_ns_pv", label: "NS/PV",      variant: "destructive" },
  { key: "cold_chain", label: "❄️ Cold Chain", variant: "blue" },
]
// Only render flags where product[flag.key] === true
```

### Search Supports

INN, trade name, SKU, barcode, registration #, DataMatrix (GTIN).

### Access Control

- `is_ns_pv = true` products: "Stock" tab visible only to `admin` and users with `ns_pv_access: true`
- Create, Edit, Delete — `admin` only

---

## 10. Inventory

UI spec: [forms/09-inventory.md](./forms/09-inventory.md) | Hook: `hooks/use-inventory.ts`

### Zod Schema (item entry)

```ts
const inventoryItemSchema = z.object({
  actual_qty: z.number().int().min(0, "Cannot be negative"),
})
```

### API Endpoints

```
POST /api/v1/inventory                               (admin / qp / warehouse_manager)
GET  /api/v1/inventory/:id

GET /api/v1/inventory/:id/items   (status = in_progress)
Response: Array<{
  id, product_id, mnn, sku,
  batch: { serial_number, expiry_date },
  actual_qty: number | null
  // expected_qty NOT returned while in_progress
}>

GET /api/v1/inventory/:id/items   (status = completed)
Response: Array<{
  ...same + expected_qty: number, discrepancy: number
}>

PUT  /api/v1/inventory/:id/items/:item_id   { actual_qty: number }
POST /api/v1/inventory/:id/complete         (admin / qp)
POST /api/v1/inventory/:id/writeoff-act     (admin)
     Body: { surplus_item_ids: string[], deficit_item_ids: string[], price_group_id: string }
POST /api/v1/inventory/:id/samples
     Body: { product_id, batch_id, qty }
```

### Blind Inventory Rule

`expected_qty` is **not returned** by the server while session `status !== "completed"`.
Client must **not** attempt to show or cache it during the active session.

### Item Save

```ts
const onSaveItem = async (sessionId: string, itemId: string, actualQty: number) => {
  await fetch(`/api/v1/inventory/${sessionId}/items/${itemId}`, {
    method: "PUT",
    body: JSON.stringify({ actual_qty: actualQty }),
  })
}
```

### Access Control

| Action | Who |
|--------|-----|
| Start session | `admin`, `qp`, `warehouse_manager` |
| Enter actual quantities | all users with warehouse access |
| Complete session | `admin`, `qp` |
| Netting act | `admin` only |
| Cancel session | `admin` only |

### Surplus/Deficit Netting Act

Offsets surpluses and shortfalls **within a single price group** only.

```
POST /api/v1/inventory/:id/writeoff-act
Body: { surplus_item_ids: string[], deficit_item_ids: string[], price_group_id: string }
```
