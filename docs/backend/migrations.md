# Backend — Миграции БД

> Часть документации `backend.md`. Описывает структуру и порядок применения миграций.

---

## Применение миграций

SQL-файлы пронумерованы и применяются через `golang-migrate`:

```bash
migrate -path ./backend/migrate -database "postgres://..." up
```

Каждая миграция имеет файлы `*.up.sql` (применить) и `*.down.sql` (откатить).

---

## Список миграций

| Миграция | Таблица / изменение                                        | Фаза |
|----------|------------------------------------------------------------|------|
| 000001   | ENUMs: `user_role`, `batch_status`, `zone_type`, `shift_type`, `claim_type`, `purchase_type`, `order_type`, `order_status`, `inventory_status`, `training_result` | 1 |
| 000002   | `users` — email-based, без логина/пароля, с ролями ERP     | 1    |
| 000003   | `categories`                                               | 1    |
| 000004   | `products` (бывший `medicaments`) — полная переработка     | 1    |
| 000005   | `warehouse_zones` (бывший `warehouses`) — зонирование      | 1    |
| 000006   | `suppliers` — картотека поставщиков                        | 1    |
| 000007   | `batches` (бывший `stock_items`) — учёт серий              | 1    |
| 000008   | `audit_logs` — с полями hash-цепочки (Phase 3, mock)       | 3    |
| 000009   | `employee_profiles` — расширен GDP, медкнижка, допуск      | 1    |
| 000010   | `refresh_tokens` — без изменений                           | 1    |
| 000011   | `inbound_receipts` + `inbound_positions` — приёмка         | 1    |
| 000012   | `product_photos` — фото препаратов                         | 1    |
| 000013   | `settings` — настройки системы (seeded: mos_percent=60)    | 1    |
| 000014   | `orders` + `order_items` — заказы (FEFO)                   | 2    |
| 000015   | `inventory_sessions` + `inventory_items` + `inventory_samples` | 2 |
| 000016   | `claims` + `claim_photos` + `recalled_batches`             | 3    |
| 000017   | `environment_logs` — журнал среды склада                   | 3    |

---

## Удалённые (устаревшие) миграции

Следующие таблицы из старых миграций были заменены или удалены:

| Старая таблица     | Замена                                      |
|--------------------|---------------------------------------------|
| `medicaments`      | `products`                                  |
| `warehouses`       | `warehouse_zones`                           |
| `stock_items`      | `batches`                                   |
| `stock_operations` | `batches` + `order_items`                   |
| `users.login`      | Убран, основной идентификатор — `email`     |
| `users.password_hash` | Убран, только OAuth (Google / Telegram)  |
| `users.status`     | Убран, статус заменён на `is_blocked`       |

Подробная схема БД: [docs/database-schema.md](../database-schema.md)
