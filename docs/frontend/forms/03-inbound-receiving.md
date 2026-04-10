# Form 3.1: Inbound Receiving

← [Back to Forms Index](./index.md) | [← Employee Profile](./02-employee-profile.md) | [Warehouse Zoning →](./04-warehouse-zoning.md)

> **UI Spec only.** For Zod schema, API contracts, post-selection logic, and quarantine release → [`forms-logic.md §4`](../forms-logic.md#4-inbound-receiving)

## Requirements Checklist

- [ ] Table `suppliers`: `id`, `name`, `inn`, `license_number`
- [ ] Table `products`: `id`, `mnn`, `ru_number`, `sku`, `atc_code`, `is_jnvlp`, `is_mdlp`, `is_ns_pv`, `cold_chain`
- [ ] Table `batches`: `id`, `product_id`, `serial_number`, `manufacture_date`, `expiry_date`, `quantity`, `status: enum(quarantine, available, rejected)`
- [ ] Table `inbound_receipts`: `id`, `supplier_id`, `invoice_number`, `country_of_origin`, `manufacturer`, `vat_rate`, `markup_jnvlp`, `created_by`, `created_at`
- [ ] Countries reference (`GET /api/v1/ref/countries`)
- [ ] Suppliers reference (`GET /api/v1/suppliers`)
- [ ] Products search (`GET /api/v1/products/search?q=...`)
- [ ] Backend: `POST /api/v1/inbound` — create receipt
- [ ] Backend: `GET /api/v1/inbound/:id` — view receipt
- [ ] Backend: `POST /api/v1/inbound/:id/quarantine-release` — (admin/QP) release from quarantine
- [ ] shadcn/ui components: `npx shadcn@latest add form select combobox date-picker textarea toast`
- [ ] `<DatePicker>` is already in the project at `components/ui/date-picker.tsx`
- [ ] `<MultiSelect>` is already in the project at `components/ui/multi-select.tsx`

---

## UI

### Page Layout

- Page `/receiving/new` — a multi-step master form (Stepper)
- Steps: **Step 1 → Step 2 → Step 3 → Confirmation**
- Responsive: vertical layout on mobile

### shadcn/ui Components

| Component | Purpose |
|-----------|---------|
| `<Form>` | Root wrapper (react-hook-form) |
| `<Select>`, `<SelectContent>`, `<SelectItem>` | Supplier, Purchase type, Country, VAT |
| `<Input>` | Invoice number, Manufacturer |
| `<Combobox>` | Product search by INN/SKU (autocomplete) |
| `<DatePicker>` | Manufacture date, Expiry date |
| `<Textarea>` | Notes |
| `<Table>` | Line items in the receipt |
| `<Badge>` | Batch status (Quarantine / Available) |
| `<Button>` | Next step, Save, Add item |
| `<Dialog>` | Quarantine release confirmation |
| `<Toast>` | Notifications |
| `<Alert>` | Quarantine warning |

### UI Structure (by steps)

**Step 1 — Document Details:**
```
Supplier:          [Combobox — search by name/INN]
Purchase type:     [Select: Direct / Tender / State procurement]
Invoice #:         [Input]
Country of origin: [Select — countries reference]
Manufacturer:      [Input]
```

**Step 2 — Finance:**
```
VAT rate:            [Select: 0% / 10% / 20%]
JNVLP markup control: [Switch — active if product is_jnvlp]
Markup (%):          [Input type=number — if JNVLP enabled]
```

**Step 3 — Line Items (repeating block):**
```
Product (INN/SKU): [Combobox with search]
Batch #:           [Input]
Manufacture date:  [DatePicker]
Expiry date:       [DatePicker] ← CRITICAL
Quantity:          [Input type=number]
[+ Add item]

─────────────────────────────────────────
Items table:
│ INN │ Batch │ Mfg date │ Expiry │ Qty │ Delete │
```

**Step 4 — Confirmation:**
```
⚠️  Goods will be placed in QUARANTINE until the acceptance protocol is signed.
[Summary view of all data]
[Save and place in quarantine]
```

---

## Admin Block (admin / QP)

### Acceptance Protocol & Quarantine Release

Visible **only** to users with role `admin` or `qp`:

```
┌──────────────────────────────────────────┐
│  📋 Acceptance Protocol                  │
│                                          │
│  Qualified Person: [Auto-fill]           │
│  Inspection date: [DatePicker]           │
│  Result: [✓ Compliant / ✗ Rejected]     │
│  Notes: [Textarea]                       │
│                                          │
│  [🔐 Sign with e-signature & Release]    │
└──────────────────────────────────────────┘
```

Button “Sign & Release” → `POST /api/v1/inbound/:id/quarantine-release`

If `result = "rejected"` → batches transition to `status: "rejected"` and a "Create return to supplier" button appears.

---

## Spec Reference

→ [Forms Index — Section 3.1 Inbound Receiving](./index.md#31-inbound-receiving--03-inbound-receivingmd)
→ Logic, Zod schema, API contracts: [`forms-logic.md §4`](../forms-logic.md#4-inbound-receiving)
