# Form 4.3: Claims & Defects

← [Back to Forms Index](./index.md) | [← Assembly & Shipment](./06-assembly-shipment-fefo.md) | [Product Card →](./08-product-card.md)

> **UI Spec only.** For Zod schema, API contracts, STOP signal polling, and photo upload logic → [`forms-logic.md §8`](../forms-logic.md#8-claims--defects)

## Requirements Checklist

- [ ] Table `claims`: `id`, `type: enum(recall, return_from_pharmacy, return_to_supplier, defect)`, `batch_id`, `product_id`, `status: enum(open, blocked, closed)`, `source`, `notes`, `created_by`, `created_at`
- [ ] Table `claim_photos`: `id`, `claim_id`, `url`, `uploaded_at`
- [ ] Roszdravnadzor integration: table `recalled_batches` (synced by background worker)
- [ ] Backend: `GET /api/v1/claims` — list claims
- [ ] Backend: `POST /api/v1/claims` — create claim
- [ ] Backend: `GET /api/v1/claims/:id` — details
- [ ] Backend: `POST /api/v1/claims/:id/photos` — upload photos (multipart)
- [ ] Backend: `POST /api/v1/claims/:id/close` — (admin/qp) close claim
- [ ] Backend: `GET /api/v1/recalled-batches/sync` — (admin) manual Roszdravnadzor sync trigger
- [ ] STOP signal: when a recalled batch is imported `batch.status → "blocked"` — automatically
- [ ] File storage (S3/MinIO) for defect photos
- [ ] shadcn/ui: `npx shadcn@latest add alert-dialog`

---

## UI

### Page Layout

- `/claims` — claims list
- `/claims/new` — create claim form
- `/claims/:id` — detailed view

### shadcn/ui Components

| Component | Purpose |
|-----------|---------|
| `<Table>` | Claims list |
| `<Badge>` | Claim type and status |
| `<Form>`, `<FormField>` | Creation form |
| `<Select>` | Claim type |
| `<Combobox>` | Product/batch search |
| `<Input>` | Source, description |
| `<Textarea>` | Notes |
| `<Alert variant="destructive">` | STOP signal, block alert |
| `<AlertDialog>` | Close claim confirmation |
| `<Card>` | Uploaded photos block |
| `<Button>` | Upload photo, Save, Close |

### UI Structure — List

```
┌──────────────────────────────────────────────────────┐
│  Claims & Defects                      [+ New]       │
│  [All] [Open] [Blocked] [Closed]                     │
├────┬─────────────┬──────────┬──────────┬────────────┤
│ #  │  Product    │  Batch   │  Type    │ Status     │
├────┼─────────────┼──────────┼──────────┼────────────┤
│ 12 │ Aspirin...  │ RZ2024A  │ Recall   │🔴 BLOCKED  │
│ 11 │ Omeprazole  │ PK2025C  │ Return   │🟡 Open     │
└────┴─────────────┴──────────┴──────────┴────────────┘
```

### STOP Signal (global alert)

When a batch matches the Roszdravnadzor registry — a fixed `<Alert>` is shown **at the top of any page**:
```
🛑 STOP SIGNAL: Batch RZ2024A (Aspirin 500mg) is listed in the 
   Roszdravnadzor recalled batch registry. Stock is blocked.
   [View claim]
```

### UI Structure — Create Claim Form

```
Claim type:       [Select: Roszdravnadzor Recall / Return from pharmacy /
                           Return to supplier / Defect]
Product:          [Combobox search by INN/SKU]
Batch:            [Combobox — only batches of selected product]
Source:           [Input — pharmacy/supplier or RZN order #]
Defect description: [Textarea]
Photo evidence:   [Drag-and-drop zone or Input type=file]
                  [Photo preview grid]
```

---

## Admin Block

| Action | Who |
|--------|-----|
| Close claim | `admin`, `qp` |
| Manual Roszdravnadzor sync | `admin` — "Sync" button in `/admin/settings` |
| Unblock batch | `admin` — mandatory reason required; only for non-RZN blocks |

---

## Spec Reference

→ [Forms Index — Section 4.2 Claims & Defects](./index.md#42-claims--defects--07-claims-defectsmd)
→ Logic, Zod schema, photo upload, STOP signal: [`forms-logic.md §8`](../forms-logic.md#8-claims--defects)
