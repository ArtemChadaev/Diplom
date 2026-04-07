# Backend — Логирование

> Часть документации `backend.md`. Описывает систему структурированного логирования.

---

## Система логирования

Используется stdlib `log/slog` со структурированным выводом.

### Компоненты

- `internal/pkg/logger` — обёртка:
  - Настройка формата: `APP_ENV=production` → JSON, иначе text
  - Контекстное обогащение: `logger.WithUserID(ctx, id)`
- `logger.FromContext(ctx)` — извлекает логгер из контекста (создаётся в `AuthRequired` middleware)
- Кастомный `loggingMiddleware` в `Handler` логирует каждый запрос:
  - метод, путь, статус, latency, request_id

### Уровни логов

| Событие                              | Уровень |
|--------------------------------------|---------|
| Отсутствие `.env` файла              | `WARN`  |
| Успешный старт сервера               | `INFO`  |
| Входящий HTTP-запрос                 | `INFO`  |
| Ошибки бизнес-логики (AppError)      | `ERROR` |
| Критические ошибки (panic recovery)  | `ERROR` |

### Пример лога запроса

```json
{
  "time": "2026-03-28T10:00:00Z",
  "level": "INFO",
  "msg": "request completed",
  "method": "PATCH",
  "path": "/api/v1/admin/users/5/role",
  "status": 200,
  "latency_ms": 12,
  "request_id": "abc123",
  "user_id": 1
}
```
