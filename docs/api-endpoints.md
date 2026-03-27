# API Endpoints Reference

**Base URL:** `/api/v1`  
**Auth:** Bearer JWT (access token) в заголовке `Authorization: Bearer <token>`  
**Refresh token:** передаётся и возвращается **исключительно через `httpOnly`-cookie** (`Set-Cookie: refresh_token=...; HttpOnly; Secure; SameSite=Strict; Path=/api/v1/auth/refresh`). В JSON-ответах и телах запросов не передаётся.  
**Роли:** `admin` > `qp` > `warehouse_manager` > `storekeeper` > `pharmacist`

> 📖 **Архитектура бэкенда:** [backend.md](./backend.md)  
> 📁 **Детальные требования по каждому эндпоинту (БД + функции):** [docs/api/](./api/)

---

## Auth

> Подробные требования к БД и функциям: [api/auth-register.md](./api/auth-register.md) · [api/auth-login.md](./api/auth-login.md) · [api/auth-send-code.md](./api/auth-send-code.md) · [api/auth-verify-code.md](./api/auth-verify-code.md) · [api/auth-google.md](./api/auth-google.md) · [api/auth-refresh.md](./api/auth-refresh.md) · [api/auth-logout.md](./api/auth-logout.md)

### POST `/auth/send-code`
Отправить 6-значный OTP-код на email.

**Auth:** Нет  
**Body:**
```json
{ "email": "user@example.com" }
```
**Responses:**
| Code | Body |
|------|------|
| 200 | `{ "message": "code_sent", "expires_in": 600 }` |
| 404 | `{ "error": "user_not_found" }` |
| 429 | `{ "error": "rate_limit_exceeded", "retry_after": 3600 }` |

---

### POST `/auth/verify-code`
Проверить OTP-код, получить токены.

**Auth:** Нет  
**Body:**
```json
{ "email": "user@example.com", "code": "482910" }
```
**Responses:**
| Code | Body / Headers |
|------|---------------|
| 200 | `{ "access_token": "...", "user": UserDTO }` + `Set-Cookie: refresh_token=...; HttpOnly; Secure` |
| 400 | `{ "error": "invalid_code", "attempts_left": 2 }` |
| 410 | `{ "error": "code_expired" }` |
| 429 | `{ "error": "max_attempts_reached" }` |

---

### GET `/auth/google`
Редирект на Google OAuth.

**Auth:** Нет  
**Response:** `302 Location: https://accounts.google.com/...`

---

### GET `/auth/google/callback`
Callback после Google OAuth.

**Auth:** Нет  
**Query:** `code=...&state=...`  
**Response 200:**
```json
{ "access_token": "...", "user": UserDTO }
```
**Headers:** `Set-Cookie: refresh_token=...; HttpOnly; Secure; SameSite=Strict; Path=/api/v1/auth/refresh`

---

### POST `/auth/refresh`
Обновить access token по refresh token из cookie.

**Auth:** Нет  
**Body:** пустой — refresh token читается автоматически из `httpOnly`-cookie  
**Response 200:** `{ "access_token": "..." }` + `Set-Cookie: refresh_token=...; HttpOnly; Secure` (обновлённый)  
**Response 401:** `{ "error": "invalid_or_expired_refresh_token" }`

---

### POST `/auth/logout`
Отозвать refresh token (server-side revocation).

**Auth:** Bearer  
**Body:** пустой — refresh token читается из `httpOnly`-cookie  
**Response 204:** No content + `Set-Cookie: refresh_token=; Max-Age=0` (очистка cookie)

---

## Users

> Подробные требования к БД и функциям: [api/users.md](./api/users.md)

### GET `/users/me`
Текущий пользователь.

**Auth:** Bearer (все роли)  
**Response 200:**
```json
{
  "id": 1,
  "email": "user@example.com",
  "role": "pharmacist",
  "ns_pv_access": false,
  "ukep_bound": false,
  "profile": {
    "full_name": "Иванов Иван Иванович",
    "position": "Фармацевт",
    "avatar_url": "https://...",
    "medical_book_scan_url": null,
    "special_zone_access": false,
    "gdp_training_history": []
  }
}
```

---

### POST `/users/me/medical-book`
Загрузить скан медкнижки.

**Auth:** Bearer (все роли)  
**Content-Type:** `multipart/form-data`  
**Body:** `file: <binary>`  
**Response 200:** `{ "url": "https://..." }`

---

### GET `/users` *(admin)*
Список всех сотрудников.

**Auth:** `admin`  
**Query:** `?q=&role=&page=1&limit=20`  
**Response 200:** `{ "items": [UserDTO], "total": 42, "page": 1 }`

---

### GET `/users/:id` *(admin)*
Профиль сотрудника.

**Auth:** `admin`  
**Response 200:** `UserDTO` (полный, как `/users/me`)

---

### PATCH `/users/:id` *(admin)*
Обновить роль / допуски.

**Auth:** `admin`  
**Body (частичное обновление):**
```json
{
  "role": "qp",
  "ns_pv_access": true,
  "special_zone_access": true
}
```
**Response 200:** `UserDTO`

---

### POST `/users/:id/send-login-link` *(admin)*
Выслать новый код входа для сотрудника.

**Auth:** `admin`  
**Response 200:** `{ "message": "code_sent" }`

---

## Справочники (References)

> Подробные требования к БД и функциям: [api/planned-modules.md](./api/planned-modules.md)

### GET `/ref/countries`
Список стран.

**Auth:** Bearer  
**Response 200:** `[{ "code": "RU", "name_ru": "Россия" }]`

---

### GET `/ref/atc`
ATC-классификация.

**Auth:** Bearer  
**Query:** `?q=`  
**Response 200:** `[{ "code": "J01CA04", "name": "Amoxicillin" }]`

---

## Products (Товары / Медикаменты)

> Подробные требования к БД и функциям: [api/planned-modules.md](./api/planned-modules.md#products-products)

### GET `/products`
Список товаров с поиском.

**Auth:** Bearer (все роли)  
**Query:** `?q=&is_jnvlp=&page=1&limit=20`  
**Response 200:** `{ "items": [ProductShortDTO], "total": 100 }`

---

### GET `/products/:id`
Карточка товара.

**Auth:** Bearer  
**Response 200:** `ProductDTO` (полная карточка с сериями и фото)

---

### POST `/products` *(admin)*
Создать товар.

**Auth:** `admin`  
**Body:** `ProductCreateDTO`  
**Response 201:** `ProductDTO`

---

### PUT `/products/:id` *(admin)*
Обновить товар.

**Auth:** `admin`  
**Body:** `ProductCreateDTO`  
**Response 200:** `ProductDTO`

---

### POST `/products/:id/photos` *(admin)*
Загрузить фото товара.

**Auth:** `admin`  
**Content-Type:** `multipart/form-data`  
**Body:** `file: <binary>`, `is_primary: bool`  
**Response 201:** `{ "id": "...", "url": "...", "is_primary": true }`

---

## Suppliers (Поставщики)

> Подробные требования к БД и функциям: [api/planned-modules.md](./api/planned-modules.md#suppliers-suppliers)

### GET `/suppliers`
Список поставщиков.

**Auth:** Bearer  
**Query:** `?q=`  
**Response 200:** `[{ "id": "...", "name": "...", "inn": "..." }]`

---

### POST `/suppliers` *(admin)*
Создать поставщика.

**Auth:** `admin`  
**Body:** `{ "name": "...", "inn": "...", "license_number": "..." }`  
**Response 201:** `SupplierDTO`

---

## Inbound (Приход)

> Подробные требования к БД и функциям: [api/planned-modules.md](./api/planned-modules.md#inbound-inbound)

### POST `/inbound`
Создать приход (все серии → статус `quarantine`).

**Auth:** `warehouse_manager`, `qp`, `admin`  
**Body:**
```json
{
  "supplier_id": "uuid",
  "purchase_type": "direct",
  "invoice_number": "ТТН-2026-001",
  "country_of_origin": "RU",
  "manufacturer": "Фармстандарт",
  "vat_rate": "10",
  "is_jnvlp_controlled": true,
  "jnvlp_markup": 15.5,
  "positions": [
    {
      "product_id": "uuid",
      "serial_number": "A2025B",
      "manufacture_date": "2025-01-15",
      "expiry_date": "2027-01-15",
      "quantity": 100
    }
  ]
}
```
**Response 201:** `InboundDTO`

---

### GET `/inbound`
Список приходов.

**Auth:** `warehouse_manager`, `qp`, `admin`  
**Query:** `?status=quarantine&page=1&limit=20`  
**Response 200:** `{ "items": [InboundShortDTO], "total": 5 }`

---

### GET `/inbound/:id`
Детали прихода.

**Auth:** `warehouse_manager`, `qp`, `admin`  
**Response 200:** `InboundDTO`

---

### POST `/inbound/:id/quarantine-release` *(qp, admin)*
Подписать протокол и выпустить из карантина.

**Auth:** `qp`, `admin`  
**Body:**
```json
{
  "inspection_date": "2026-03-27",
  "result": "approved",
  "notes": "Соответствует"
}
```
**Response 200:** `{ "released_batches": 5, "status": "available" }`  
**Note:** При `result = "rejected"` → серии переходят в `rejected`.

---

## Zones (Зоны склада)

> Подробные требования к БД и функциям: [api/planned-modules.md](./api/planned-modules.md#zones-zones)

### GET `/zones`
Список зон.

**Auth:** Bearer (все роли; `safe_strong` скрыта без `ns_pv_access`)  
**Response 200:** `[ZoneDTO]`

---

### GET `/zones/:id/stock`
Остатки в зоне.

**Auth:** Bearer  
**Response 200:** `{ "zone": ZoneDTO, "items": [StockItemDTO] }`

---

### POST `/zones` *(admin)*
Создать зону.

**Auth:** `admin`  
**Body:** `{ "name": "...", "type": "cold_chain", "temp_min": 2, "temp_max": 8, "humidity_max": 60 }`  
**Response 201:** `ZoneDTO`

---

### PUT `/zones/:id` *(admin)*
Обновить зону.

**Auth:** `admin`  
**Body:** Частичное, те же поля  
**Response 200:** `ZoneDTO`

---

## Environment Log (Журнал среды)

> Подробные требования к БД и функциям: [api/planned-modules.md](./api/planned-modules.md#environment-log-environment-log-zonesidenvironment-log)

### GET `/environment-log/today`
Сводка записей журнала за сегодня.

**Auth:** `warehouse_manager`, `storekeeper`, `qp`, `admin`  
**Response 200:**
```json
[{
  "zone_id": "...",
  "zone_name": "Холодовая цепь",
  "temp_min": 2, "temp_max": 8, "humidity_max": 60,
  "morning_log": null,
  "evening_log": { "temperature": 5.2, "humidity": 45, "recorded_by": {...} }
}]
```

---

### POST `/zones/:id/environment-log`
Внести запись журнала.

**Auth:** `warehouse_manager`, `storekeeper`, `qp`, `admin`  
**Body:**
```json
{
  "shift": "morning",
  "temperature": 5.2,
  "humidity": 45,
  "notes": null
}
```
**Response 201:** `EnvLogDTO`  
**Response 400:** `{ "error": "already_recorded_for_shift" }`

---

### GET `/zones/:id/environment-log`
История журнала зоны.

**Auth:** `warehouse_manager`, `qp`, `admin`  
**Query:** `?from=2026-01-01&to=2026-03-31&page=1&limit=50`  
**Response 200:** `{ "items": [EnvLogDTO], "total": 120 }`

---

### GET `/environment-log/export` *(admin, qp)*
Экспорт в Excel.

**Auth:** `admin`, `qp`  
**Query:** `?from=&to=`  
**Response 200:** `Content-Type: application/vnd.openxmlformats-officedocument.spreadsheetml.sheet`

---

## Orders (Заказы / Сборка FEFO)

> Подробные требования к БД и функциям: [api/planned-modules.md](./api/planned-modules.md#orders-orders)

### GET `/orders`
Список заказов.

**Auth:** Bearer (все роли)  
**Query:** `?status=new&type=cito&page=1&limit=20`  
**Response 200:** `{ "items": [OrderShortDTO], "total": 10 }`

---

### POST `/orders`
Создать заказ.

**Auth:** `warehouse_manager`, `qp`, `admin`  
**Body:**
```json
{
  "type": "regular",
  "destination_id": "uuid",
  "items": [{ "product_id": "uuid", "requested_qty": 20 }]
}
```
**Response 201:** `OrderDTO`

---

### GET `/orders/:id`
Детали заказа с FEFO-рекомендациями.

**Auth:** Bearer  
**Response 200:** `OrderDTO` (включает `fefo_recommendations` по каждой позиции)

---

### POST `/orders/:id/confirm-assembly`
Подтвердить сборку позиций.

**Auth:** Bearer (все)  
**Body:**
```json
{
  "assembled_items": [
    { "item_id": "uuid", "batch_id": "uuid", "assembled_qty": 20 }
  ]
}
```
**Response 200:** `{ "status": "ready" }`

---

### POST `/orders/:id/ship` *(warehouse_manager, qp, admin)*
Подтвердить отгрузку, сгенерировать ТТН.

**Auth:** `warehouse_manager`, `qp`, `admin`  
**Response 200:** `{ "ttn_url": "https://...", "quality_registry_url": "https://..." }`  
**Response 422:** `{ "error": "mos_blocked_items", "items": ["uuid"] }`

---

### GET `/orders/:id/quality-registry`
Реестр сертификатов качества к ТТН.

**Auth:** `warehouse_manager`, `qp`, `admin`  
**Response 200:** PDF

---

## Settings (Настройки)

> Подробные требования к БД и функциям: [api/planned-modules.md](./api/planned-modules.md#settings-settings)

### GET `/settings/mos` *(admin)*
Текущий МОС-порог.

**Auth:** `admin`  
**Response 200:** `{ "mos_percent": 60 }`

---

### PUT `/settings/mos` *(admin)*
Обновить МОС-порог.

**Auth:** `admin`  
**Body:** `{ "mos_percent": 60 }`  
**Response 200:** `{ "mos_percent": 60 }`

---

## Claims (Рекламации и Брак)

> Подробные требования к БД и функциям: [api/planned-modules.md](./api/planned-modules.md#claims-claims)

### GET `/claims`
Список рекламаций.

**Auth:** Bearer  
**Query:** `?status=open&type=recall&page=1&limit=20`  
**Response 200:** `{ "items": [ClaimShortDTO], "total": 15 }`

---

### POST `/claims`
Создать рекламацию.

**Auth:** Bearer  
**Body:**
```json
{
  "type": "defect",
  "product_id": "uuid",
  "batch_id": "uuid",
  "source": "Аптека №5",
  "notes": "Повреждена упаковка"
}
```
**Response 201:** `ClaimDTO`

---

### GET `/claims/:id`
Детали рекламации.

**Auth:** Bearer  
**Response 200:** `ClaimDTO`

---

### POST `/claims/:id/photos`
Загрузить фото брака.

**Auth:** Bearer  
**Content-Type:** `multipart/form-data`  
**Body:** `photos: [<binary>]`  
**Response 201:** `[{ "id": "...", "url": "..." }]`

---

### POST `/claims/:id/close` *(admin, qp)*
Закрыть рекламацию.

**Auth:** `admin`, `qp`  
**Body:** `{ "resolution": "..." }`  
**Response 200:** `ClaimDTO`

---

### GET `/stop-signals`
Активные STOP-сигналы (изъятые серии Росздравнадзора).

**Auth:** Bearer  
**Response 200:** `[{ "batch_serial": "...", "product_name": "...", "claim_id": "..." }]`

---

### POST `/recalled-batches/sync` *(admin)*
Ручная синхронизация с реестром Росздравнадзора.

**Auth:** `admin`  
**Response 200:** `{ "synced_count": 3, "newly_blocked": 1 }`

---

## Inventory (Инвентаризация)

> Подробные требования к БД и функциям: [api/planned-modules.md](./api/planned-modules.md#inventory-inventory)

### GET `/inventory`
Список сессий инвентаризации.

**Auth:** `warehouse_manager`, `qp`, `admin`  
**Response 200:** `{ "items": [InventorySessionShortDTO], "total": 8 }`

---

### POST `/inventory`
Начать новую сессию инвентаризации.

**Auth:** `warehouse_manager`, `qp`, `admin`  
**Body:** `{ "zone_id": "uuid | null" }` (null = весь склад)  
**Response 201:** `InventorySessionDTO`

---

### GET `/inventory/:id`
Детали сессии. В статусе `in_progress` — без `expected_qty`.

**Auth:** Bearer  
**Response 200:** `InventorySessionDTO`

---

### PUT `/inventory/:id/items/:item_id`
Ввести фактическое количество (слепая инвентаризация).

**Auth:** Bearer  
**Body:** `{ "actual_qty": 17 }`  
**Response 200:** `InventoryItemDTO`

---

### POST `/inventory/:id/complete` *(qp, admin)*
Завершить инвентаризацию (раскрыть `expected_qty`, вычислить расхождения).

**Auth:** `qp`, `admin`  
**Response 200:** `InventorySessionDTO` (с расхождениями)

---

### POST `/inventory/:id/writeoff-act` *(admin)*
Акт зачёта пересортицы.

**Auth:** `admin`  
**Body:**
```json
{
  "price_group_id": "uuid",
  "surplus_item_ids": ["uuid"],
  "deficit_item_ids": ["uuid"]
}
```
**Response 201:** `{ "act_url": "https://..." }`

---

### POST `/inventory/:id/samples`
Добавить контрольный образец.

**Auth:** Bearer  
**Body:** `{ "product_id": "uuid", "batch_id": "uuid", "qty": 1 }`  
**Response 201:** `InventorySampleDTO`

---

## Общие DTO

```ts
type UserDTO = {
  id: number
  email: string
  role: "admin" | "qp" | "warehouse_manager" | "storekeeper" | "pharmacist"
  ns_pv_access: boolean
  ukep_bound: boolean
  profile: EmployeeProfileDTO
}

type ProductDTO = {
  id: string
  trade_name: string
  mnn: string
  sku: string
  barcode: string | null
  datamatrix_gtin: string | null
  ru_number: string
  atc_codes: string[]
  dosage_form: string
  dosage: string
  package_multiplicity: number
  is_jnvlp: boolean
  is_mdlp: boolean
  is_ns_pv: boolean
  cold_chain: boolean
  temp_min: number | null
  temp_max: number | null
  humidity_max: number | null
  photos: PhotoDTO[]
  batches: BatchDTO[]
}

type BatchDTO = {
  id: string
  serial_number: string
  manufacture_date: string   // ISO
  expiry_date: string        // ISO
  quantity: number
  status: "quarantine" | "available" | "rejected" | "blocked"
  zone_name: string
}

type ZoneDTO = {
  id: string
  name: string
  type: "general" | "cold_chain" | "flammable" | "safe_strong"
  temp_min: number | null
  temp_max: number | null
  humidity_max: number | null
  stock_count: number
}

type EnvLogDTO = {
  id: string
  zone_id: string
  shift: "morning" | "evening"
  temperature: number
  humidity: number
  notes: string | null
  recorded_by: { id: number; full_name: string }
  recorded_at: string
}
```
