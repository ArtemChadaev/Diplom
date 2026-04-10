# Form 5.1: Product Card (Medicament)

← [Back to Forms Index](./index.md) | [← Claims & Defects](./07-claims-defects.md) | [Inventory →](./09-inventory.md)

> **UI Spec only.** For Zod schema, API contracts, flag rendering logic, and access control → [`forms-logic.md §9`](../forms-logic.md#9-product-card)

## Requirements Checklist

- [ ] Table `products`: `id`, `mnn`, `trade_name`, `sku`, `barcode`, `datamatrix_gtin`, `ru_number`, `atc_code`, `dosage_form`, `dosage`, `package_multiplicity`, `weight_g`, `width_cm`, `height_cm`, `depth_cm`
- [ ] Flag fields: `is_jnvlp: bool`, `is_mdlp: bool`, `is_ns_pv: bool`, `cold_chain: bool`
- [ ] Storage fields: `temp_min`, `temp_max`, `humidity_max`
- [ ] Table `product_photos`: `id`, `product_id`, `url`, `is_primary: bool`
- [ ] Backend: `GET /api/v1/products` — product list (search + pagination)
- [ ] Backend: `GET /api/v1/products/:id` — product card
- [ ] Backend: `POST /api/v1/products` — (admin) create product
- [ ] Backend: `PUT /api/v1/products/:id` — (admin) edit product
- [ ] Backend: `DELETE /api/v1/products/:id` — (admin) soft delete (`deleted_at`)
- [ ] Backend: `POST /api/v1/products/:id/photos` — upload photo
- [ ] ATC classification reference (`GET /api/v1/ref/atc`)
- [ ] shadcn/ui: `npx shadcn@latest add tooltip separator`

---

## UI

### Page Layout

- `/products` — searchable product table (accessible to all)
- `/products/:id` — product card (accessible to all)
- `/products/new` — (admin) create new product
- `/products/:id/edit` — (admin) edit product

### shadcn/ui Components

| Component | Purpose |
|-----------|---------|
| `<Card>` | Main card container |
| `<Badge>` | Flags (JNVLP, MDLP, NS/PV, ❄️ Cold Chain) |
| `<Tabs>` | Sections: Basic / Storage / Stock / History |
| `<Table>` | Current batches with stock |
| `<Tooltip>` | Flag and code hints |
| `<Separator>` | Section dividers |
| `<Form>`, `<Input>`, `<Select>`, `<Textarea>` | Create/edit form |
| `<Button>` | Edit, Upload photo |
| `<MultiSelect>` | Multiple ATC code selection (already in project) |

### UI Structure — Product Card

```
┌──────────────────────────────────────────────────────────┐
│  [Photo]   Amoxicillin 500mg                              │
│            [JNVLP] [MDLP] [❄️ Cold Chain]                │
│            INN: Amoxicillin                               │
│            RU#: ЛП-002712  |  SKU: AM500CAP              │
│            ATC: J01CA04                                    │
├──────────────────────────────────────────────────────────┤
│  [Basic] [Storage conditions] [Stock] [History]           │
├──────────────────────────────────────────────────────────┤
│  "Basic" tab:                                             │
│  Dosage form: Capsules 500 mg                             │
│  Package multiplicity: 16 pcs                             │
│  Barcode: 4607053860207                                   │
│  GTIN (DataMatrix): 04607053860207                        │
│  Dimensions: 120×80×30 mm | 45 g                         │
├──────────────────────────────────────────────────────────┤
│  "Storage conditions" tab:                                │
│  Temperature: 2–8°C                                       │
│  Humidity: max 60%                                        │
├──────────────────────────────────────────────────────────┤
│  "Stock" tab (current batches):                           │
│  │ Batch │ Lot │ Expiry │ Qty │ Zone │ Status │           │
└──────────────────────────────────────────────────────────┘
```

### Create/Edit Form (admin)

```
Trade name:         [Input]
INN:                [Input]
SKU:                [Input] (auto-generate or manual)
Barcode:            [Input]
GTIN:               [Input]
Registration #:     [Input] ← Required
ATC code:           [MultiSelect — from reference]
Dosage form:        [Select: Tablets / Capsules / Solution / ...]
Dosage:             [Input]
Package mult.:      [Input type=number]

Flags:
  [Switch] JNVLP
  [Switch] MDLP labeling required
  [Switch] NS/PV (Schedule II/III)
  [Switch] Cold Chain

Storage conditions:
  Temperature from: [Input] °C
  Temperature to:   [Input] °C
  Max humidity:     [Input] %

Dimensions:
  Weight: [Input] g
  W × H × D: [Input] × [Input] × [Input] mm

Photo:
  [Upload primary photo] ← drag-and-drop
```

---

## Admin Block

- "Edit" button — `admin` only
- "Create product" button — `admin` only
- Products with `is_ns_pv = true` are visible to **all users**, but the "Stock" tab is visible only to `admin` and `ns_pv_access: true`
- For `is_ns_pv = true` products, the edit form shows an additional "Schedule II/III requirements" section — `admin` only

---

## Spec Reference

→ [Forms Index — Section 5.1 Product Card](./index.md#51-product-card--08-product-cardmd)
→ Logic, Zod schema, API contracts, flag rendering: [`forms-logic.md §9`](../forms-logic.md#9-product-card)
