# Form 5.2: Inventory

← [Back to Forms Index](./index.md) | [← Product Card](./08-product-card.md)

> **UI Spec only.** For Zod schema, API contracts, blind inventory mechanism, and netting act logic → [`forms-logic.md §10`](../forms-logic.md#10-inventory)

## Requirements Checklist

- [ ] Table `inventory_sessions`: `id`, `status: enum(draft, in_progress, completed, cancelled)`, `started_by`, `started_at`, `completed_at`, `zone_id` (null = entire warehouse)
- [ ] Table `inventory_items`: `id`, `session_id`, `product_id`, `batch_id`, `expected_qty` (from system), `actual_qty` (entered by employee), `discrepancy` (calc: actual - expected)
- [ ] Table `inventory_samples`: `id`, `session_id`, `product_id`, `batch_id`, `qty` — lab control samples tracking
- [ ] Table `price_groups`: `id`, `name`, `price_range_from`, `price_range_to` — price groups for recount act
- [ ] Backend: `POST /api/v1/inventory` — (admin/qp/warehouse_manager) start inventory session
- [ ] Backend: `GET /api/v1/inventory/:id` — session details
- [ ] Backend: `PUT /api/v1/inventory/:id/items/:item_id` — enter actual quantity
- [ ] Backend: `POST /api/v1/inventory/:id/complete` — (admin/qp) complete session, record discrepancies
- [ ] Backend: `POST /api/v1/inventory/:id/writeoff-act` — (admin) surplus/deficit netting act
- [ ] **Blind inventory:** employees do not see `expected_qty` until the session is completed
- [ ] shadcn/ui: `npx shadcn@latest add progress alert-dialog`

---

## UI

### Page Layout

- `/inventory` — list of inventory sessions
- `/inventory/:id` — active inventory page (product search + quantity entry)
- `/inventory/:id/results` — results (only after completion)

### shadcn/ui Components

| Component | Purpose |
|-----------|---------|
| `<Table>` | Session list, item entry list |
| `<Badge>` | Session status |
| `<Progress>` | Item entry progress |
| `<Input>` | Actual quantity entry |
| `<Combobox>` | Product search by INN/SKU/DataMatrix |
| `<Alert>` | Large discrepancy warning |
| `<AlertDialog>` | Inventory completion confirmation |
| `<Dialog>` | Recount act form (admin) |
| `<Card>` | Discrepancy summary card |
| `<Tabs>` | Surpluses / Shortfalls / OK — in results |

### UI Structure — Entry (Blind Inventory)

```
┌──────────────────────────────────────────────────────┐
│  Inventory #5 | Full warehouse                        │
│  Progress: ████████░░ 120 of 145 items               │
│  ────────────────────────────────────────────────    │
│  Find product: [Combobox: INN / SKU / DataMatrix]    │
│                                                       │
│  ┌─────────────────────────────────────────────┐     │
│  │ Amoxicillin 500mg | Batch A2025B             │     │
│  │ Actual qty: [____] pcs    [Save]             │     │
│  └─────────────────────────────────────────────┘     │
│  ─────────────────────────────────────────────────   │
│  ⚠ 12 items not yet checked                          │
│                                                       │
│  [Complete inventory]                                 │
└──────────────────────────────────────────────────────┘
```

> **Important:** `expected_qty` is **hidden** until the session is completed.  
> Only after `POST /api/v1/inventory/:id/complete` does it become visible.

### UI Structure — Results (after completion)

```
┌──────────────────────────────────────────────────────┐
│  Inventory #5 Results                                 │
│  [Surpluses (+7)] [Shortfalls (-12)] [OK (126)]      │
├──────────────────────────────────────────────────────┤
│  "Shortfalls" tab:                                    │
│  │ INN │ Batch │ Expected │ Actual │ Diff │          │
│  │ Amo │A2025B │    20    │   17   │  -3  │          │
│  ─────────────────────────────────────────────────   │
│  🔐 [Surplus/Deficit Netting Act] ← admin only       │
└──────────────────────────────────────────────────────┘
```

---

## Admin Block

| Action | Who |
|--------|-----|
| Start inventory session | `admin`, `qp`, `warehouse_manager` |
| Enter actual quantities | all users with warehouse access |
| Complete inventory | `admin`, `qp` |
| Surplus/deficit netting act | **`admin` only** |
| Cancel inventory | **`admin` only** |

In results (`/inventory/:id/results`), the "Netting Act" block is hidden from everyone except `admin`.

---

## Spec Reference

→ [Forms Index — Section 5.2 Inventory](./index.md#52-inventory--09-inventorymd)
→ Logic, Zod schema, blind inventory, API contracts: [`forms-logic.md §10`](../forms-logic.md#10-inventory)
