# Form 3.2a: Warehouse Zoning

← [Back to Forms Index](./index.md) | [← Inbound Receiving](./03-inbound-receiving.md) | [Environment Log →](./05-environment-log.md)

> **UI Spec only.** For Zod schema, API contracts, and access control logic → [`forms-logic.md §5`](../forms-logic.md#5-warehouse-zoning)

## Requirements Checklist

- [ ] Table `warehouse_zones`: `id`, `name`, `type: enum(general, cold_chain, flammable, safe_strong)`, `description`, `temp_min`, `temp_max`, `humidity_max`
- [ ] Link `storage_locations`: `id`, `zone_id`, `rack`, `shelf`, `cell` (address storage — Phase 4)
- [ ] Backend: `GET /api/v1/zones` — list zones
- [ ] Backend: `POST /api/v1/zones` — (admin) create zone
- [ ] Backend: `PUT /api/v1/zones/:id` — (admin) edit zone
- [ ] Backend: `GET /api/v1/zones/:id/stock` — stock in zone
- [ ] shadcn/ui components: `npx shadcn@latest add card badge table tooltip`
- [ ] Zone type `cold_chain` — only products with `cold_chain: true` can be placed here
- [ ] Zone type `safe_strong` — only users with `ns_pv_access: true` can see stock in this zone

---

## UI

### Page Layout

- Page `/warehouse/zones` — accessible to `admin`, `qp`, `warehouse_manager`
- Grid layout of zone cards (`2x2` or `3x...`)
- Each card shows a brief zone summary with a link to detailed view

### shadcn/ui Components

| Component | Purpose |
|-----------|---------|
| `<Card>`, `<CardHeader>`, `<CardContent>`, `<CardFooter>` | Zone card |
| `<Badge>` | Zone type (color-coded) |
| `<Table>` | Zone stock items |
| `<Tooltip>` | Zone requirements hints |
| `<Dialog>` | Create/edit zone form (admin) |
| `<Form>`, `<Input>`, `<Select>` | Fields inside edit dialog |

### Zone Type Color Coding

| Zone Type | Badge color | Icon |
|-----------|------------|------|
| `general` | `default` | 📦 |
| `cold_chain` | `blue` | ❄️ |
| `flammable` | `orange` | 🔥 |
| `safe_strong` | `red` | 🔒 |

### UI Structure

```
┌──────────────────────────────────────────┐
│   Warehouse: Zones               [+ Add]  │
├──────────────────────────────────────────┤
│ ┌──────────────┐  ┌──────────────┐       │
│ │ 📦 General  │  │ ❄️ Cold chain │       │
│ │ zone        │  │               │       │
│ │             │  │ 2-8°C | ≤60% │       │
│ │ 142 items   │  │ 38 items      │       │
│ │ [View]      │  │ [View]        │       │
│ └──────────────┘  └──────────────┘       │
│ ┌──────────────┐  ┌──────────────┐       │
│ │ 🔥 Flammable│  │ 🔒 Safe zone │       │
│ │             │  │ (NS/PV)      │       │
│ │ 8 items     │  │ [Access restr│       │
│ └──────────────┘  └──────────────┘       │
└──────────────────────────────────────────┘
```

**Zone detail page `/warehouse/zones/:id`:**
```
Zone: Cold Chain  [Edit]
Requirements: 2–8°C | Humidity: ≤60%

Stock table:
│ INN │ Batch │ Expiry │ Qty │ Status │
```

---

## Admin Block

- `[+ Add zone]` button — visible to `admin` only
- Zone card edit dialog — `admin` only
- `safe_strong` zone card stock data visible only to `admin` and users with `ns_pv_access: true`

---

## Spec Reference

→ [Forms Index — Section 3.2 Warehouse Zoning](./index.md#32-warehouse-zoning--04-warehouse-zoningmd)
→ Logic, Zod schema, API contracts: [`forms-logic.md §5`](../forms-logic.md#5-warehouse-zoning)
