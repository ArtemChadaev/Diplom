# Form 4.1–4.2: Order Assembly & Shipment Control (FEFO)

← [Back to Forms Index](./index.md) | [← Environment Log](./05-environment-log.md) | [Claims & Defects →](./07-claims-defects.md)

> **UI Spec only.** For FEFO algorithm, API contracts, MOS calculation, and role-based shipment logic → [`forms-logic.md §7`](../forms-logic.md#7-assembly--shipment-fefo)

## Requirements Checklist

- [ ] Table `orders`: `id`, `type: enum(regular, cito)`, `status`, `destination` (pharmacy), `created_at`, `assembled_by`
- [ ] Table `order_items`: `id`, `order_id`, `product_id`, `requested_qty`, `batch_id` (auto-selected by FEFO), `assembled_qty`, `status`
- [ ] Table `batches` contains `expiry_date` — backend implements FEFO batch selection algorithm
- [ ] Backend: `POST /api/v1/orders` — create order
- [ ] Backend: `GET /api/v1/orders` — list orders (with filters: status, type, date)
- [ ] Backend: `GET /api/v1/orders/:id` — order details with FEFO recommendations
- [ ] Backend: `POST /api/v1/orders/:id/confirm-assembly` — confirm assembly
- [ ] Backend: `POST /api/v1/orders/:id/ship` — (admin/qp/warehouse_manager) confirm shipment + generate TTN
- [ ] Backend: `GET /api/v1/orders/:id/quality-registry` — quality certificate registry for TTN
- [ ] MOS setting (minimum remaining shelf life): `GET/PUT /api/v1/settings/mos` — percentage, `admin` only
- [ ] shadcn/ui: `npx shadcn@latest add progress alert-dialog badge`

---

## UI

### Page Layout

- `/orders` — list of all orders with filters
- `/orders/:id/assemble` — assembly page for a specific order

### shadcn/ui Components

| Component | Purpose |
|-----------|---------|
| `<Table>` | Order list, order line items |
| `<Badge>` | Order status, type (Cito!/Regular) |
| `<Progress>` | Assembly progress (X of Y items assembled) |
| `<Alert>` | MOS block, warnings |
| `<AlertDialog>` | Shipment / cancellation confirmation |
| `<Tabs>` | Filter by status (New / Assembling / Ready / Shipped) |
| `<Button>` | Accept to work, Confirm assembly, Ship |
| `<Card>` | FEFO recommendation block per item |
| `<Tooltip>` | Remaining shelf life details |

### UI Structure — Order List

```
┌──────────────────────────────────────────────────────┐
│  Orders                                               │
│  [New] [Assembling] [Ready to ship] [Shipped]        │
├─────┬─────────────┬────────────┬──────────┬─────────┤
│  #  │  Pharmacy   │  Date      │  Type    │ Status  │
├─────┼─────────────┼────────────┼──────────┼─────────┤
│ 101 │ Pharmacy #3 │ 27.03.2026 │ 🚨 Cito! │ New    │
│ 100 │ Pharmacy #7 │ 27.03.2026 │ Regular  │ Assembly│
└─────┴─────────────┴────────────┴──────────┴─────────┘
```

### UI Structure — Assembly Page

```
┌──────────────────────────────────────────────────────┐
│  Order #101 | 🚨 CITO! | Pharmacy #3                 │
│  Assembly progress: ████████░░ 8 of 10 items         │
├──────────────────────────────────────────────────────┤
│  Item 1: Amoxicillin 500mg (SKU: AM500)              │
│  Requested: 20 pcs                                    │
│                                                       │
│  ⚡ FEFO recommendation:                              │
│  ┌─────────────────────────────────────────────┐     │
│  │ Batch: A2025B | Exp: 15.06.2026             │     │
│  │ Remaining: 80 days (27%)  ⚠ MOS BLOCK!      │     │
│  │ [🚫 SHIPMENT BLOCKED]                       │     │
│  └─────────────────────────────────────────────┘     │
│  ┌─────────────────────────────────────────────┐     │
│  │ Batch: C2025X | Exp: 15.01.2027             │     │
│  │ Remaining: 294 days (81%)  ✓                │     │
│  │ Assemble: [20] pcs    [✓ Confirm]           │     │
│  └─────────────────────────────────────────────┘     │
└──────────────────────────────────────────────────────┘
```

---

## Admin Block

### MOS Configuration

Page `/admin/settings`:
```
MOS (Minimum Remaining Shelf Life):
[60] % ← Input + [Save]
Applies to all recipient pharmacies.
```
`PUT /api/v1/settings/mos { mos_percent: number }` — visible to `admin` only.

### Shipment Authorization

"Ship" button — only for `admin`, `qp`, `warehouse_manager`.
`storekeeper` and `pharmacist` can only **assemble** (confirm line items), but cannot ship.

---

## Spec Reference

→ [Forms Index — Section 4.1 Assembly & Shipment (FEFO)](./index.md#41-assembly--shipment-fefo--06-assembly-shipment-fefomd)
→ Logic, MOS calculation, API contracts: [`forms-logic.md §7`](../forms-logic.md#7-assembly--shipment-fefo)
