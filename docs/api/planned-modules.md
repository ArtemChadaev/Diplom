# Планируемые модули (не реализованы в бэкенде)

> Эти эндпоинты описаны в [api-endpoints.md](../api-endpoints.md) и спроектированы, но бэкенд-код для них **ещё не написан**. Ниже — требования к БД и функциям для будущей реализации.

---

## Справочники (`/ref/`)

### GET `/ref/countries`

**БД:** таблица `countries` (code VARCHAR PK, name_ru VARCHAR)  
**Функции:** `refRepo.ListCountries(ctx)` — простой `SELECT * FROM countries ORDER BY name_ru`

### GET `/ref/atc`

**БД:** таблица `atc_codes` (code VARCHAR PK, name VARCHAR, parent_code nullable)  
**Query param:** `?q=` → `WHERE name ILIKE '%?%' OR code ILIKE '%?%'`  
**Функции:** `refRepo.SearchATC(ctx, query string)`

---

## Products (`/products/`)

### Таблица `products` (= `medicaments` в миграциях)

| Поле                  | Тип           |
|-----------------------|---------------|
| `id`                  | UUID PK        |
| `trade_name`          | VARCHAR NOT NULL |
| `mnn`                 | VARCHAR        |
| `sku`                 | VARCHAR UNIQUE |
| `barcode`             | VARCHAR nullable |
| `datamatrix_gtin`     | VARCHAR nullable |
| `ru_number`           | VARCHAR        |
| `atc_codes`           | TEXT[]         |
| `dosage_form`         | VARCHAR        |
| `dosage`              | VARCHAR        |
| `package_multiplicity`| INT            |
| `is_jnvlp`            | BOOLEAN        |
| `is_mdlp`             | BOOLEAN        |
| `is_ns_pv`            | BOOLEAN        |
| `cold_chain`          | BOOLEAN        |
| `temp_min/max`        | NUMERIC nullable |
| `humidity_max`        | NUMERIC nullable |

**GET /products:** `SELECT + ILIKE ?q + WHERE is_jnvlp = ? + LIMIT/OFFSET`  
**POST /products:** INSERT с валидацией уникальности `sku`  
**Функции:** `productRepo.List`, `productRepo.FindByID`, `productRepo.Create`, `productRepo.Update`

---

## Suppliers (`/suppliers/`)

### Таблица `suppliers`

| Поле             | Тип     |
|------------------|---------|
| `id`             | UUID PK |
| `name`           | VARCHAR |
| `inn`            | VARCHAR UNIQUE |
| `license_number` | VARCHAR nullable |

---

## Inbound (`/inbound/`)

### Таблица `inbound_orders`

| Поле                  | Тип           |
|-----------------------|---------------|
| `id`                  | UUID PK        |
| `supplier_id`         | UUID FK        |
| `purchase_type`       | VARCHAR        |
| `invoice_number`      | VARCHAR        |
| `country_of_origin`   | VARCHAR        |
| `manufacturer`        | VARCHAR        |
| `vat_rate`            | VARCHAR        |
| `is_jnvlp_controlled` | BOOLEAN        |
| `jnvlp_markup`        | NUMERIC        |
| `status`              | VARCHAR (quarantine/available/rejected) |
| `created_at`          | TIMESTAMPTZ    |

### Таблица `batches` (серии)

| Поле               | Тип        |
|--------------------|------------|
| `id`               | UUID PK     |
| `inbound_order_id` | UUID FK     |
| `product_id`       | UUID FK     |
| `serial_number`    | VARCHAR     |
| `manufacture_date` | DATE        |
| `expiry_date`      | DATE        |
| `quantity`         | INT         |
| `status`           | VARCHAR     |
| `zone_id`          | UUID FK nullable |

**POST /inbound/quarantine-release:**
```sql
UPDATE batches SET status = 'available' WHERE inbound_order_id = ? AND status = 'quarantine'
-- При result = 'rejected':
UPDATE batches SET status = 'rejected' WHERE inbound_order_id = ?
```

---

## Zones (`/zones/`)

### Таблица `warehouses` (= zones в семантике)

| Поле           | Тип     |
|----------------|---------|
| `id`           | UUID PK |
| `name`         | VARCHAR |
| `type`         | ENUM (general, cold_chain, flammable, safe_strong) |
| `temp_min/max` | NUMERIC nullable |
| `humidity_max` | NUMERIC nullable |

**GET /zones:** фильтрация `safe_strong` → если у пользователя нет `ns_pv_access`, эта зона скрывается

---

## Environment Log (`/environment-log/`, `/zones/:id/environment-log`)

### Таблица `environment_logs`

| Поле          | Тип          |
|---------------|--------------|
| `id`          | UUID PK       |
| `zone_id`     | UUID FK       |
| `shift`       | ENUM (morning, evening) |
| `temperature` | NUMERIC       |
| `humidity`    | NUMERIC       |
| `notes`       | TEXT nullable |
| `recorded_by` | INT FK → users |
| `recorded_at` | TIMESTAMPTZ   |

**Уникальность:** `UNIQUE(zone_id, shift, DATE(recorded_at))` — одна запись на зону за смену в день. При дубле → `400 already_recorded_for_shift`.

---

## Orders (`/orders/`)

### Таблица `orders`

| Поле             | Тип     |
|------------------|---------|
| `id`             | UUID PK |
| `type`           | VARCHAR (regular, cito) |
| `destination_id` | UUID FK |
| `status`         | VARCHAR (new, assembled, shipped) |
| `created_by`     | INT FK  |
| `created_at`     | TIMESTAMPTZ |

### FEFO-сборка

При `GET /orders/:id` → для каждой позиции вычислить рекомендации по FEFO:
```sql
SELECT b.* FROM batches b
WHERE b.product_id = ? AND b.status = 'available'
ORDER BY b.expiry_date ASC  -- First Expired First Out
LIMIT 1
```

### MOS-блокировка при отгрузке

```sql
-- Получить настройку МОС
SELECT mos_percent FROM settings WHERE key = 'mos_percent'
-- Проверить остатки по каждой серии в заказе
-- Серии с MOS < threshold → заблокировать отгрузку
```

---

## Inventory (`/inventory/`)

### Таблица `inventory_sessions`

| Поле         | Тип     |
|--------------|---------|
| `id`         | UUID PK |
| `zone_id`    | UUID FK nullable |
| `status`     | VARCHAR (in_progress, completed) |
| `created_by` | INT FK  |
| `created_at` | TIMESTAMPTZ |

### Таблица `inventory_items`

| Поле           | Тип     |
|----------------|---------|
| `id`           | UUID PK |
| `session_id`   | UUID FK |
| `product_id`   | UUID FK |
| `batch_id`     | UUID FK |
| `expected_qty` | INT     |
| `actual_qty`   | INT nullable |

**Слепая инвентаризация:** `expected_qty` скрыт пока `status = 'in_progress'`. При `POST /inventory/:id/complete` → раскрыть `expected_qty`, вычислить расхождения.

---

## Claims (`/claims/`)

### Таблица `claims`

| Поле         | Тип     |
|--------------|---------|
| `id`         | UUID PK |
| `type`       | VARCHAR (defect, recall) |
| `product_id` | UUID FK |
| `batch_id`   | UUID FK |
| `source`     | VARCHAR |
| `notes`      | TEXT    |
| `status`     | VARCHAR (open, closed) |
| `resolution` | TEXT nullable |
| `created_by` | INT FK  |
| `created_at` | TIMESTAMPTZ |

### STOP-сигналы и заблокированные серии

**GET /stop-signals:** серии, перешедшие в статус `blocked` из-за рекламации типа `recall`
```sql
SELECT b.serial_number, m.trade_name, c.id AS claim_id
FROM batches b
JOIN claims c ON c.batch_id = b.id
JOIN medicaments m ON m.id = b.product_id
WHERE b.status = 'blocked' AND c.type = 'recall'
```

---

## Settings (`/settings/`)

### Таблица `settings`

| Поле    | Тип     | Описание                       |
|---------|---------|--------------------------------|
| `key`   | VARCHAR PK | Ключ настройки (`mos_percent`) |
| `value` | JSONB   | Значение                       |

**GET/PUT /settings/mos:**
```sql
SELECT value FROM settings WHERE key = 'mos_percent'
UPDATE settings SET value = ? WHERE key = 'mos_percent'
```
