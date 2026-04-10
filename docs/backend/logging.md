# Backend — Logging

← [Back to Main README](../../README.md) | [Architecture →](./architecture.md)

---

## Logging System

Uses stdlib `log/slog` with structured output via a custom context-aware wrapper.

> **Critical Rule:** Never use `fmt.Println`, `fmt.Printf`, or any `fmt` print functions for logs or debug output.
> Always use the project logger:
> ```go
> log := logger.FromContext(ctx)
> log.Info("...", "key", value)
> ```
> Name the logger variable `log`, not `logger`.

### Components

- `internal/pkg/logger` — wrapper:
  - Format selection: `APP_ENV=production` → JSON output, otherwise text
  - Context enrichment: `logger.WithUserID(ctx, id)`
- `logger.FromContext(ctx)` — extracts logger from context (injected by `AuthRequired` middleware)
- Custom `loggingMiddleware` in `Handler` logs every request:
  - method, path, status, latency, request_id

### Log Levels

| Event | Level |
|-------|-------|
| Missing `.env` file | `WARN` |
| Successful server start | `INFO` |
| Incoming HTTP request | `INFO` |
| Business logic errors (`AppError`) | `ERROR` |
| Critical errors (panic recovery) | `ERROR` |

### Request Log Example

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
