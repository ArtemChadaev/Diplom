# Form 3.2b: Environment Log

← [Back to Forms Index](./index.md) | [← Warehouse Zoning](./04-warehouse-zoning.md) | [Assembly & Shipment →](./06-assembly-shipment-fefo.md)

> **UI Spec only.** For Zod schema (dynamic zone-aware), API contracts, and pending badge logic → [`forms-logic.md §6`](../forms-logic.md#6-environment-log)

## Requirements Checklist

- [ ] Table `environment_logs`: `id`, `zone_id`, `recorded_by`, `recorded_at`, `temperature`, `humidity`, `shift: enum(morning, evening)`, `notes`
- [ ] Backend: `POST /api/v1/zones/:id/environment-log` — create entry
- [ ] Backend: `GET /api/v1/zones/:id/environment-log` — entry history (paginated, date-filtered)
- [ ] Backend: `GET /api/v1/environment-log/today` — all entries for today (aggregated by zone)
- [ ] Rule: **2 entries per day** — morning (`morning`) and evening (`evening`) for each zone
- [ ] Alert: if entry for the current shift is missing → notification in the Header notification center
- [ ] Zones with `temp_min`/`temp_max` — valid range enforced on both server and client
- [ ] shadcn/ui: `npx shadcn@latest add form input select textarea toast alert`

---

## UI

### Page Layout

- Page `/warehouse/environment-log` — accessible to `warehouse_manager`, `storekeeper`, `qp`, `admin`
- Two sections:
  1. **Quick input panel** — zone cards with "Enter data" button
  2. **History** — log table with date/zone filter

### shadcn/ui Components

| Component | Purpose |
|-----------|---------|
| `<Card>` | Zone card in quick input panel |
| `<Badge>` | Entry status (✓ recorded / ⚠ pending) |
| `<Dialog>` | Data input form |
| `<Form>`, `<FormField>` | Input fields |
| `<Input type="number">` | Temperature, humidity |
| `<Select>` | Shift selection (morning/evening) |
| `<Textarea>` | Notes |
| `<Table>` | Entry history |
| `<Alert>` | Out-of-range warning |
| `<Toast>` | Save confirmation |

### UI Structure

**Quick input panel:**
```
┌──────────────────────────────────────────────────────┐
│  Temperature & Humidity Log — Today                   │
│  27.03.2026 | Shift: [Morning ▼]                     │
├──────────┬──────────┬──────────┬────────────────────┤
│ 📦 Gen.  │ ❄️ Cold  │ 🔥 Flamm.│  🔒 NS/PV         │
│  ✓ Done  │  ⚠ Pend. │  ✓ Done  │  ⚠ Pending        │
│ [View]   │ [Enter]  │ [View]   │  [Enter]           │
└──────────┴──────────┴──────────┴────────────────────┘
```

**Data entry dialog:**
```
┌──────────────────────────────────────┐
│  ❄️ Cold Chain — Morning Shift        │
│  Range: 2–8°C | Humidity ≤60%       │
├──────────────────────────────────────┤
│  Temperature (°C): [____]            │
│  Humidity (%):     [____]            │
│  Notes:            [_____________]   │
│                                      │
│  ⚠ 9.2°C — ABOVE LIMIT! (2–8°C)    │
│                                      │
│  [Cancel]            [Save]          │
└──────────────────────────────────────┘
```

**History table:**
```
│ Zone │ Date │ Shift │ Temp. │ Hum. │ Employee │ Notes │
```

---

## Admin Block / QP

- **Can edit** already saved log (`PUT /api/v1/zones/:id/environment-log/:log_id`) — a "Correct" button
- **Can export** data (`GET /api/v1/environment-log/export?from=...&to=...`) — Excel by zone and date range
- Admin sees **all zones** (including `safe_strong`) without access restriction

---

## Spec Reference

→ [Forms Index — Section 3.2b Environment Log](./index.md#32b-environment-log--05-environment-logmd)
→ Logic, dynamic Zod schema, API contracts: [`forms-logic.md §6`](../forms-logic.md#6-environment-log)
