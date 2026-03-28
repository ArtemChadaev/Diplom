# Database Schema

> Основан на текущих миграциях `000001`–`000011` + расширения, требуемые спецификацией ERP.  
> Текущие миграции — минимальный скелет. Ниже — **целевая схема** с необходимыми изменениями.

---

## Enums

```sql
-- Роли пользователей (заменить текущий минимальный ENUM)
CREATE TYPE user_role AS ENUM (
  'admin',
  'qp',               -- Уполномоченное лицо / QP
  'warehouse_manager',
  'storekeeper',
  'pharmacist'
);

-- Статус серии (batch)
CREATE TYPE batch_status AS ENUM (
  'quarantine',   -- только принят, ожидает протокола приёмки
  'available',    -- выпущен в обращение
  'rejected',     -- забракован при приёмке
  'blocked'       -- заблокирован (STOP-сигнал / Росздравнадзор)
);

-- Тип зоны склада
CREATE TYPE zone_type AS ENUM (
  'general',
  'cold_chain',
  'flammable',
  'safe_strong'   -- сейфовая, НС/ПВ
);

-- Смена для журнала среды
CREATE TYPE shift_type AS ENUM ('morning', 'evening');

-- Тип рекламации
CREATE TYPE claim_type AS ENUM (
  'recall',               -- изъятие Росздравнадзором
  'return_from_pharmacy', -- возврат из аптеки
  'return_to_supplier',   -- возврат поставщику
  'defect'                -- производственный брак
);

-- Тип закупки
CREATE TYPE purchase_type AS ENUM ('direct', 'tender', 'state');

-- Тип заказа
CREATE TYPE order_type AS ENUM ('regular', 'cito');

-- Статус заказа
CREATE TYPE order_status AS ENUM ('new', 'assembling', 'ready', 'shipped', 'cancelled');

-- Статус инвентаризации
CREATE TYPE inventory_status AS ENUM ('draft', 'in_progress', 'completed', 'cancelled');

-- Результат обучения GDP
CREATE TYPE training_result AS ENUM ('pass', 'fail');
```

---

## Таблицы

### `users`
> Изменения vs. текущей миграции: убрать `password_hash`, переработать `role`, убрать `login` → использовать `email`.

```sql
CREATE TABLE users (
  id           SERIAL PRIMARY KEY,
  email        VARCHAR(255) UNIQUE NOT NULL,
  google_id    VARCHAR(255) UNIQUE,
  telegram_id  BIGINT UNIQUE,
  role         user_role NOT NULL DEFAULT 'pharmacist',
  ns_pv_access BOOLEAN NOT NULL DEFAULT false,  -- допуск к НС/ПВ
  ukep_bound   BOOLEAN NOT NULL DEFAULT false,  -- привязана УКЭП (электронная подпись)
  is_blocked   BOOLEAN NOT NULL DEFAULT false,
  created_at   TIMESTAMPTZ DEFAULT NOW(),
  updated_at   TIMESTAMPTZ DEFAULT NOW()
);
```

> **OTP-коды хранятся в Valkey, не в PostgreSQL.**  
> Структура ключей: `otp:user:<user_id>` (Hash), TTL 600 сек.  
> Подробности: [valkey-cache.md](./valkey-cache.md#1-otp-коды-otp)

---

### `refresh_tokens`
> Текущая миграция 000011 — оставить без изменений.

```sql
-- Без изменений (миграция 000011)
CREATE TABLE refresh_tokens (
  id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id     INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  token_hash  TEXT NOT NULL UNIQUE,
  expires_at  TIMESTAMPTZ NOT NULL,
  user_agent  TEXT,
  ip_address  INET,
  metadata    JSONB DEFAULT '{}',
  created_at  TIMESTAMPTZ DEFAULT NOW(),
  revoked_at  TIMESTAMPTZ
);
```

---

### `employee_profiles`
> Текущая миграция 000009. Добавить GDP-обучение, медкнижку, спец-допуск.

```sql
CREATE TABLE employee_profiles (
  id                    SERIAL PRIMARY KEY,
  user_id               INT UNIQUE NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  employee_code         VARCHAR(100) UNIQUE NOT NULL,
  full_name             VARCHAR(255) NOT NULL,
  corporate_email       VARCHAR(255) UNIQUE NOT NULL,
  phone                 VARCHAR(20) UNIQUE NOT NULL,
  position              VARCHAR(255) NOT NULL,
  department            VARCHAR(255) NOT NULL,
  birth_date            DATE NOT NULL,
  avatar_url            TEXT,
  hire_date             DATE NOT NULL,
  dismissal_date        DATE,
  -- Новые поля (vs. 000009):
  medical_book_scan_url TEXT, -- Скан медкнижки для подтверждения работы
  special_zone_access   BOOLEAN NOT NULL DEFAULT false, -- Допуск к спец зонам
  gdp_training_history  JSONB NOT NULL DEFAULT '[]' -- История прохождения Good Distribution Practice которые надо делать раз в год
  -- Структура JSONB-элемента:
  -- { date, course_name, result: "pass"|"fail", certificate_url }
);
```

---

### `categories`
> Текущая миграция 000003 — без изменений (используется для medicaments).

```sql
CREATE TABLE categories (
  id   SERIAL PRIMARY KEY,
  name VARCHAR(255) UNIQUE NOT NULL
);
```

---

### `products` (бывший `medicaments`)
> Текущая миграция 000004 — полностью переработать: переименовать, добавить все поля спецификации.

```sql
CREATE TABLE products (
  id                   UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  trade_name           VARCHAR(255) NOT NULL, -- Торговое наименование препарата
  mnn                  VARCHAR(255) NOT NULL, -- Международное непатентованное наименование
  sku                  VARCHAR(100) UNIQUE, -- Артикул
  barcode              VARCHAR(50), -- Штрихкод
  datamatrix_gtin      VARCHAR(50), -- Часть кода маркировки Честный Знак
  ru_number            VARCHAR(100) NOT NULL,       -- № Регистрационного удостоверения
  atc_codes            TEXT[] NOT NULL DEFAULT '{}', -- массив ATC-кодов
  dosage_form          VARCHAR(100), -- Лекарственная форма
  dosage               VARCHAR(100), -- Дозировка
  package_multiplicity INT NOT NULL DEFAULT 1, -- Количество в упаковке
  -- Флаги
  is_jnvlp             BOOLEAN NOT NULL DEFAULT false, -- Поизнак жизненно необходимых и важнейших лекарственных препаратов
  is_mdlp              BOOLEAN NOT NULL DEFAULT false,  -- Подлежит маркировке
  is_ns_pv             BOOLEAN NOT NULL DEFAULT false,  -- НС/ПВ (ПКУ)
  cold_chain           BOOLEAN NOT NULL DEFAULT false, -- Требует холодовой цепи
  -- Условия хранения
  temp_min             NUMERIC(5,2), -- Минимальная температура хранения
  temp_max             NUMERIC(5,2), -- Максимальная температура хранения
  humidity_max         NUMERIC(5,2),
  -- ВГХ (Весогаборитные характеристики)
  weight_g             NUMERIC(10,3),
  width_cm             NUMERIC(8,2),
  height_cm            NUMERIC(8,2),
  depth_cm             NUMERIC(8,2),
  -- Текст описания
  description          TEXT,
  -- Мета
  created_at           TIMESTAMPTZ DEFAULT NOW(),
  updated_at           TIMESTAMPTZ DEFAULT NOW(),
  deleted_at           TIMESTAMPTZ        -- мягкое удаление
);

CREATE INDEX idx_products_mnn ON products(mnn);
CREATE INDEX idx_products_sku ON products(sku);
CREATE INDEX idx_products_ru_number ON products(ru_number);
```

---

### `product_photos`
> **Новая таблица.**

```sql
CREATE TABLE product_photos (
  id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
  url        TEXT NOT NULL,
  is_primary BOOLEAN NOT NULL DEFAULT false,
  created_at TIMESTAMPTZ DEFAULT NOW()
);
```

---

### `suppliers`
> **Новая таблица.**
> Картотека поставщиков
```sql
CREATE TABLE suppliers (
  id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name           VARCHAR(255) NOT NULL,
  inn            VARCHAR(12) UNIQUE NOT NULL,
  license_number VARCHAR(100),
  created_at     TIMESTAMPTZ DEFAULT NOW()
);
```

---

### `warehouse_zones` (бывший `warehouses`)
> Текущая миграция 000005 — полностью переработать под систему зонирования.

```sql
CREATE TABLE warehouse_zones (
  id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name         VARCHAR(255) NOT NULL,
  type         zone_type NOT NULL DEFAULT 'general',
  description  TEXT,
  temp_min     NUMERIC(5,2),
  temp_max     NUMERIC(5,2),
  humidity_max NUMERIC(5,2),
  created_at   TIMESTAMPTZ DEFAULT NOW()
);
```

---

### `batches` (серии)
> Бывший `stock_items` — переработать полностью.

```sql
CREATE TABLE batches (
  id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  product_id       UUID NOT NULL REFERENCES products(id) ON DELETE RESTRICT,
  zone_id          UUID REFERENCES warehouse_zones(id) ON DELETE SET NULL,
  serial_number    VARCHAR(100) NOT NULL,
  manufacture_date DATE NOT NULL,
  expiry_date      DATE NOT NULL,
  quantity         INT NOT NULL CHECK (quantity >= 0),
  status           batch_status NOT NULL DEFAULT 'quarantine',
  updated_at       TIMESTAMPTZ DEFAULT NOW(),
  UNIQUE (product_id, serial_number)
);

CREATE INDEX idx_batches_product_id ON batches(product_id);
CREATE INDEX idx_batches_expiry_date ON batches(expiry_date);
CREATE INDEX idx_batches_status ON batches(status);
```

---

### `inbound_receipts` (накладные приёмки)
> **Новая таблица.**

```sql
CREATE TABLE inbound_receipts (
  id                   UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  supplier_id          UUID NOT NULL REFERENCES suppliers(id),
  purchase_type        purchase_type NOT NULL,
  invoice_number       VARCHAR(100) NOT NULL,
  country_of_origin    VARCHAR(3) NOT NULL,   -- ISO-код страны
  manufacturer         VARCHAR(255) NOT NULL,
  vat_rate             SMALLINT NOT NULL,     -- 0, 10, 20
  is_jnvlp_controlled  BOOLEAN NOT NULL DEFAULT false,
  jnvlp_markup         NUMERIC(5,2),
  -- Протокол приёмки
  qp_user_id           INT REFERENCES users(id),
  inspection_date      DATE,
  inspection_result    VARCHAR(20),           -- 'approved' | 'rejected'
  inspection_notes     TEXT,
  -- Мета
  created_by           INT NOT NULL REFERENCES users(id),
  created_at           TIMESTAMPTZ DEFAULT NOW()

  photo_urls           TEXT[] DEFAULT '{}', -- Фото поврежденных коробок/паллет при разгрузке
  digital_signature_id UUID,                -- Ссылка на файл открепленной подписи (УКЭП) провизора
);
```

---

### `inbound_positions` (позиции накладной)
> **Новая таблица.**

```sql
CREATE TABLE inbound_positions (
  id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  inbound_id  UUID NOT NULL REFERENCES inbound_receipts(id) ON DELETE CASCADE,
  batch_id    UUID NOT NULL REFERENCES batches(id)
);
```
---


### `digital_signatures` (цифровые подписи)
> **Новая таблица.**
> Единое хранилище подписей, на которое будут ссылаться все документы (приемка, логи, списания). 
> MOCK: Реальная криптографическая проверка подписи заменена на логическую фиксацию. signature_url будет просто строкой 'mock_signature_file_path', а signed_hash — простым md5.
> Вообще не буду делать в коде скорее всего

```sql
CREATE TABLE digital_signatures (
  id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id       INT NOT NULL REFERENCES users(id),
  
  -- Ссылка на файл подписи (обычно .p7s или .sig) в S3-хранилище
  signature_url TEXT NOT NULL,
  
  -- Хэш данных, которые были подписаны (чтобы нельзя было подменить данные в логах)
  signed_hash   TEXT NOT NULL,
  
  -- Данные сертификата (для истории, даже если срок действия выйдет)
  cert_serial   VARCHAR(100),
  cert_issuer   TEXT,
  signed_at     TIMESTAMPTZ DEFAULT NOW(),
  
  -- Результат проверки (валидна ли подпись на момент записи)
  is_valid      BOOLEAN DEFAULT true
);
```

---

### `environment_logs` (журнал среды)
> **Новая таблица.**

```sql
CREATE TABLE environment_logs (
  id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  zone_id     UUID NOT NULL REFERENCES warehouse_zones(id) ON DELETE CASCADE,
  digital_signature_id UUID REFERENCES digital_signatures(id), -- ЭЦП сотрудника, зафиксировавшего параметры (не реализовывать)
  recorded_by INT NOT NULL REFERENCES users(id),
  shift       shift_type NOT NULL,
  temperature NUMERIC(5,2) NOT NULL,
  humidity    NUMERIC(5,2) NOT NULL,
  notes       TEXT,
  recorded_at TIMESTAMPTZ DEFAULT NOW(),
  UNIQUE (zone_id, recorded_at::DATE, shift)   -- одна запись за смену в день
);

CREATE INDEX idx_env_logs_zone_date ON environment_logs(zone_id, recorded_at);
```

---

### `orders` (заказы)
> **Новая таблица.**

```sql
CREATE TABLE orders (
  id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  type           order_type NOT NULL DEFAULT 'regular',
  status         order_status NOT NULL DEFAULT 'new',
  destination_id UUID,   -- FK → destinations (аптеки, будущая таблица)
  -- КТО заказал (подпись аптеки/заказчика) (не реализовывать)
  customer_signature_id UUID REFERENCES digital_signatures(id),
  -- КТО разрешил отгрузку (подпись Уполномоченного лица склада) (не реализовывать)
  outbound_signature_id UUID REFERENCES digital_signatures(id),
  assembled_by   INT REFERENCES users(id),
  shipped_by     INT REFERENCES users(id),
  shipped_at     TIMESTAMPTZ,
  ttn_url        TEXT,
  created_by     INT NOT NULL REFERENCES users(id),
  created_at     TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_orders_status ON orders(status);
CREATE INDEX idx_orders_type ON orders(type);
```

---

### `order_items` (позиции заказа)
> **Новая таблица.**

```sql
CREATE TABLE order_items (
  id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  order_id      UUID NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
  product_id    UUID NOT NULL REFERENCES products(id),
  requested_qty INT NOT NULL CHECK (requested_qty > 0),
  batch_id      UUID REFERENCES batches(id),  -- заполняется при сборке (FEFO)
  assembled_qty INT,
  status        VARCHAR(20) NOT NULL DEFAULT 'pending'
  -- 'pending', 'mos_blocked', 'assembled'
);
```

---

### `claims` (рекламации)
> **Новая таблица.**

```sql
CREATE TABLE claims (
  id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  digital_signature_id UUID; -- ЭЦП лица, утвердившего блокировку/возврат серии (не реализовывать)
  type       claim_type NOT NULL,
  batch_id   UUID REFERENCES batches(id),
  product_id UUID NOT NULL REFERENCES products(id),
  status     VARCHAR(20) NOT NULL DEFAULT 'open',
  -- 'open', 'blocked', 'closed'
  source     TEXT,
  notes      TEXT,
  resolution TEXT,
  created_by INT NOT NULL REFERENCES users(id),
  created_at TIMESTAMPTZ DEFAULT NOW(),
  closed_at  TIMESTAMPTZ
);
```

---

### `claim_photos`
> **Новая таблица.**

```sql
CREATE TABLE claim_photos (
  id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  claim_id    UUID NOT NULL REFERENCES claims(id) ON DELETE CASCADE,
  url         TEXT NOT NULL,
  uploaded_at TIMESTAMPTZ DEFAULT NOW()
);
```

---

### `recalled_batches` (изъятые серии Росздравнадзора)
> **Новая таблица.**
> MOCK: Прям подключать не думаю что возможно просто сам добавлю пару штук и мб сделаю функцию тоже MOCK.

```sql
CREATE TABLE recalled_batches (
  id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  serial_number VARCHAR(100) NOT NULL,
  product_name  VARCHAR(255),
  ru_number     VARCHAR(100),
  recall_reason TEXT,
  issued_at     DATE,
  synced_at     TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_recalled_batches_serial ON recalled_batches(serial_number);
```

---

### `inventory_sessions` (сессии инвентаризации)
> **Новая таблица.**

```sql
CREATE TABLE inventory_sessions (
  id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  status       inventory_status NOT NULL DEFAULT 'draft',
  zone_id      UUID REFERENCES warehouse_zones(id),  -- NULL=весь склад
  started_by   INT NOT NULL REFERENCES users(id),
  completed_by INT REFERENCES users(id),
  started_at   TIMESTAMPTZ DEFAULT NOW(),
  completed_at TIMESTAMPTZ
);
```

---

### `inventory_items` (позиции инвентаризации)
> **Новая таблица.**

```sql
CREATE TABLE inventory_items (
  id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  session_id   UUID NOT NULL REFERENCES inventory_sessions(id) ON DELETE CASCADE,
  product_id   UUID NOT NULL REFERENCES products(id),
  batch_id     UUID REFERENCES batches(id),
  expected_qty INT NOT NULL,      -- из системы, скрыто до завершения
  actual_qty   INT,               -- вводит сотрудник
  discrepancy  INT GENERATED ALWAYS AS (actual_qty - expected_qty) STORED
);
```

---

### `inventory_samples` (контрольные образцы)
> **Новая таблица.**

```sql
CREATE TABLE inventory_samples (
  id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  session_id UUID NOT NULL REFERENCES inventory_sessions(id) ON DELETE CASCADE,
  product_id UUID NOT NULL REFERENCES products(id),
  batch_id   UUID REFERENCES batches(id),
  qty        INT NOT NULL CHECK (qty > 0),
  created_at TIMESTAMPTZ DEFAULT NOW()
);
```

---

### `settings`
> **Новая таблица.** Конфигурационные параметры системы.

```sql
CREATE TABLE settings (
  key        VARCHAR(100) PRIMARY KEY,
  value      TEXT NOT NULL,
  updated_by INT REFERENCES users(id),
  updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Начальные данные:
INSERT INTO settings (key, value) VALUES ('mos_percent', '60');
```

---

### `audit_logs`
> Текущая миграция 000008. Доработать — добавить хэш для Immutable Logs.

```sql
CREATE TABLE audit_logs (
  id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id      INT REFERENCES users(id),
  action       VARCHAR(255) NOT NULL,
  entity       VARCHAR(100) NOT NULL,
  entity_id    VARCHAR(100),
  old_values   JSONB,
  new_values   JSONB,
  ip_address   INET,
  -- Вполне возможно не буду делать ибо слишком сложно, как мок оставлю, но фиг его знает
  -- Иммутабельность (Фаза 3):
  prev_hash    TEXT,    -- хэш предыдущей записи в цепочке
  log_hash     TEXT,    -- SHA-256(prev_hash + user_id + action + entity_id + new_values + created_at)
  created_at   TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_audit_logs_entity ON audit_logs(entity, entity_id);
CREATE INDEX idx_audit_logs_user_id ON audit_logs(user_id);
```

---

## Статус таблиц по фазам Roadmap

| Таблица | Фаза | Статус vs. текущих миграций |
|---------|------|-----------------------------|
| `users` | 1 | ⚠️ Переработать (убрать `password_hash`, `login`, расширить `role`) |
| `otp:user:*` (Valkey) | 1 | 🔴 В Valkey, не в PostgreSQL — см. [valkey-cache.md](./valkey-cache.md) |
| `employee_profiles` | 1 | ⚠️ Добавить `medical_book_scan_url`, `special_zone_access`, `gdp_training_history` |
| `products` | 1 | ⚠️ Переработать `medicaments` полностью |
| `product_photos` | 1 | 🆕 Новая |
| `suppliers` | 1 | 🆕 Новая |
| `warehouse_zones` | 1 | ⚠️ Переработать `warehouses` |
| `batches` | 1 | ⚠️ Переработать `stock_items` |
| `inbound_receipts` | 1 | 🆕 Новая |
| `inbound_positions` | 1 | 🆕 Новая |
| `settings` | 1 | 🆕 Новая |
| `refresh_tokens` | 1 | ✅ Без изменений (миграция 000011) |
| `environment_logs` | 3 | 🆕 Новая |
| `orders` | 2 | 🆕 Новая (FEFO) |
| `order_items` | 2 | 🆕 Новая |
| `claims` | 3 | 🆕 Новая |
| `claim_photos` | 3 | 🆕 Новая |
| `recalled_batches` | 3 | 🆕 Новая |
| `inventory_sessions` | 2 | 🆕 Новая |
| `inventory_items` | 2 | 🆕 Новая |
| `inventory_samples` | 2 | 🆕 Новая |
| `audit_logs` | 3 | ⚠️ Добавить `prev_hash`, `log_hash` |
| `categories` | — | ✅ Без изменений (вспомогательная) |
| `stock_operations` | — | ♻️ Заменяется `batches` + `order_items` |
