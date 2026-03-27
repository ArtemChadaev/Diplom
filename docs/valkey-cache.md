# Valkey (Redis) — Схема данных и кэш

> Проект использует **Valkey 8** (Redis-совместимый) как хранилище эфемерных данных:  
> OTP-коды, сессионные данные, кэш, rate limiting, STOP-сигналы.  
> Контейнер: `valkey/valkey:8-alpine`, порт задаётся через `$VALKEY_PORT` → `6379`.

---

## Соглашения по ключам

Формат: `<namespace>:<entity>:<id>`  
Пример: `otp:user:42`, `cache:products:uuid-xxx`

Все TTL устанавливаются через `SET key value EX <seconds>` или `EXPIRE`.

---

## 1. OTP-коды (`otp:`)

Хранятся **только в Valkey**, не в PostgreSQL.

| Ключ | `otp:user:<user_id>` |
|------|----------------------|
| Тип | Hash |
| TTL | **600 секунд (10 минут)** |

**Поля Hash:**
```
code        "482910"   6-значный цифровой код
attempts    "0"        число неверных попыток (макс. 3 до блокировки)
created_at  "1711234567"  Unix timestamp создания
```

**Логика:**
- При запросе нового кода → `DEL otp:user:<id>`, затем `HSET` + `EXPIRE 600`
- При проверке → `HINCRBY otp:user:<id> attempts 1`; если `attempts >= 3` → `DEL` (код сгорает)
- После успешного входа → `DEL otp:user:<id>`
- Повторная отправка кода — не раньше чем через 60 сек (контролируется rate limit ниже)

---

## 2. Rate Limiting (`rl:`)

Защита от брутфорса и спама.

### 2.1 Отправка OTP (по email)

| Ключ | `rl:send_otp:<email>` |
|------|----------------------|
| Тип | Counter (INCR) |
| TTL | **3600 секунд (1 час)** |
| Лимит | 5 запросов за TTL |

```
INCR rl:send_otp:<email>
EXPIRE rl:send_otp:<email> 3600  (только при первом INCR)
→ если значение > 5 → 429 Too Many Requests
```

### 2.2 Попытки ввода OTP (по IP)

| Ключ | `rl:verify_otp:<ip>` |
|------|----------------------|
| Тип | Counter |
| TTL | **300 секунд (5 минут)** |
| Лимит | 10 попыток за 5 минут |

### 2.3 Общий rate limit API (по IP)

| Ключ | `rl:api:<ip>` |
|------|--------------|
| Тип | Counter (скользящее окно) |
| TTL | **60 секунд** |
| Лимит | 120 запросов / мин |

---

## 3. Refresh-токены (Blocklist) (`rt:`)

При logout или компрометации токен добавляется в blocklist до истечения его `expires_at`.

| Ключ | `rt:blocked:<token_hash>` |
|------|--------------------------|
| Тип | String (`"1"`) |
| TTL | Рассчитывается: `expires_at - now()` в секундах |

```
SET rt:blocked:<sha256_of_token> "1" EX <remaining_seconds>
```

Прием токена → проверка `EXISTS rt:blocked:<hash>`. Если `1` → `401 Unauthorized`.

> Полная запись токена хранится в PostgreSQL (`refresh_tokens`). Valkey — только список отозванных.

---

## 4. Кэш данных (`cache:`)

Кэширование часто запрашиваемых, редко меняющихся данных.

### 4.1 Справочник стран

| Ключ | `cache:ref:countries` |
|------|----------------------|
| Тип | String (JSON) |
| TTL | **86400 секунд (24 часа)** |
| Инвалидация | При добавлении новой страны через admin |

### 4.2 ATC-классификация

| Ключ | `cache:ref:atc` |
|------|----------------|
| Тип | String (JSON) |
| TTL | **86400 секунд (24 часа)** |

### 4.3 Настройки системы (МОС и др.)

| Ключ | `cache:settings` |
|------|----------------|
| Тип | Hash |
| TTL | **300 секунд (5 минут)** |
| Инвалидация | `DEL cache:settings` при `PUT /api/v1/settings/*` |

```
HSET cache:settings mos_percent "60"
```

### 4.4 Карточка товара

| Ключ | `cache:product:<product_id>` |
|------|------------------------------|
| Тип | String (JSON) |
| TTL | **120 секунд (2 минуты)** |
| Инвалидация | `DEL cache:product:<id>` при `PUT /api/v1/products/:id` |

### 4.5 Список зон склада

| Ключ | `cache:zones` |
|------|--------------|
| Тип | String (JSON) |
| TTL | **300 секунд** |
| Инвалидация | `DEL cache:zones` при `POST/PUT /api/v1/zones` |

---

## 5. STOP-сигналы (`stop:`)

Активные STOP-сигналы Росздравнадзора — быстрый доступ без запроса к PostgreSQL.

| Ключ | `stop:signals` |
|------|---------------|
| Тип | Set |
| TTL | **Нет** (обновляется при синхронизации) |

```
SADD stop:signals "RZ2024A" "PK2022X"
SISMEMBER stop:signals <serial_number>   → 0 или 1
```

При синхронизации с Росздравнадзором:
1. `DEL stop:signals`
2. `SADD stop:signals <serial1> <serial2> ...`

Фронтенд делает polling `GET /api/v1/stop-signals` → сервер отвечает из Valkey, не обращаясь к PostgreSQL.

---

## 6. Очередь Cito!-заказов (`queue:`)

Приоритетная очередь заказов типа `cito` для реализации срочной сборки.

| Ключ | `queue:cito_orders` |
|------|---------------------|
| Тип | Sorted Set |
| Score | Unix timestamp создания заказа |
| TTL | Нет (элементы удаляются при отгрузке) |

```
ZADD queue:cito_orders <timestamp> <order_id>
ZRANGE queue:cito_orders 0 -1      → все необработанные Cito! заказы по приоритету
ZREM queue:cito_orders <order_id>  → при отгрузке
```

---

## 7. Сессии журнала среды (`env_lock:`)

Блокировка: предотвращает двойную запись в журнал за одну смену.

| Ключ | `env_lock:<zone_id>:<date>:<shift>` |
|------|-------------------------------------|
| Тип | String (`"1"`) |
| TTL | **До конца дня (авто-расчёт)** |
| Пример | `env_lock:uuid-zone-1:2026-03-27:morning` |

```
SET env_lock:<zone>:<date>:<shift> "1" EX <seconds_until_midnight> NX
→ NX: только если ещё не установлен → 0 = уже записано (skip), 1 = OK
```

---

## Сводная таблица TTL

| Namespace | Данные | TTL |
|-----------|--------|-----|
| `otp:user:*` | OTP-код входа | 10 мин |
| `rl:send_otp:*` | Rate limit отправки кода | 1 час |
| `rl:verify_otp:*` | Rate limit проверки кода | 5 мин |
| `rl:api:*` | Rate limit API | 1 мин |
| `rt:blocked:*` | Отозванные refresh-токены | До `expires_at` токена |
| `cache:ref:countries` | Справочник стран | 24 часа |
| `cache:ref:atc` | ATC-классификация | 24 часа |
| `cache:settings` | Системные настройки | 5 мин |
| `cache:product:*` | Карточка товара | 2 мин |
| `cache:zones` | Список зон склада | 5 мин |
| `stop:signals` | STOP-сигналы Росздравнадзора | Без TTL |
| `queue:cito_orders` | Очередь срочных заказов | Без TTL |
| `env_lock:*` | Блокировка записи журнала среды | До конца дня |

---

## Конфигурация подключения (Go)

```go
// Пример конфига клиента
type ValkeyConfig struct {
    Addr     string // VALKEY_HOST:VALKEY_PORT из .env
    Password string // VALKEY_PASSWORD (пустой если нет)
    DB       int    // 0
    PoolSize int    // 10
}
```

Используемый клиент: [`valkey-go`](https://github.com/valkey-io/valkey-go) или совместимый Redis-клиент (`go-redis/v9`).
