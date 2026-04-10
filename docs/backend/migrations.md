# Backend — Database Migrations

← [Back to Main README](../../README.md) | [Local Setup →](./README.md)

---

## Applying Migrations

SQL files are numbered and applied via `golang-migrate`. Each migration has `*.up.sql` (apply) and `*.down.sql` (rollback) files.

```bash
# Apply all (via Docker Compose — preferred)
docker compose run --rm migrate

# Apply via CLI
migrate -path ./backend/migrate -database "postgres://..." up
```

---

## Migration List

| Migration | Table / Change | Phase |
|-----------|---------------|-------|
| 000001 | ENUMs: `user_role`, `batch_status`, `zone_type`, `shift_type`, `claim_type`, `purchase_type`, `order_type`, `order_status`, `inventory_status`, `training_result` | 1 |
| 000002 | `users` — email-based, no login/password, with ERP roles | 1 |
| 000003 | `categories` | 1 |
| 000004 | `products` (formerly `medicaments`) — full redesign | 1 |
| 000005 | `warehouse_zones` (formerly `warehouses`) — zone management | 1 |
| 000006 | `suppliers` — supplier registry | 1 |
| 000007 | `batches` (formerly `stock_items`) — batch/series tracking | 1 |
| 000008 | `audit_logs` — with hash chain fields (Phase 3, mock) | 3 |
| 000009 | `employee_profiles` — extended GDP, medical book, zone access | 1 |
| 000010 | `refresh_tokens` — unchanged | 1 |
| 000011 | `inbound_receipts` + `inbound_positions` — goods receiving | 1 |
| 000012 | `product_photos` — drug product photos | 1 |
| 000013 | `settings` — system settings (seeded: `mos_percent=60`) | 1 |
| 000014 | `orders` + `order_items` — orders (FEFO) | 2 |
| 000015 | `inventory_sessions` + `inventory_items` + `inventory_samples` | 2 |
| 000016 | `claims` + `claim_photos` + `recalled_batches` | 3 |
| 000017 | `environment_logs` — warehouse environment journal | 3 |

---

## Deprecated / Removed Tables

The following tables from old migrations were replaced or removed:

| Old Table | Replacement |
|-----------|------------|
| `medicaments` | `products` |
| `warehouses` | `warehouse_zones` |
| `stock_items` | `batches` |
| `stock_operations` | `batches` + `order_items` |
| `users.login` | Removed — primary identifier is now `email` |
| `users.password_hash` | Removed — OAuth only (Google / Telegram) |
| `users.status` | Removed — replaced by `is_blocked` boolean |

→ Full SQL schema: [`docs/other/migrate.md`](../Diplom/migrate.md)
