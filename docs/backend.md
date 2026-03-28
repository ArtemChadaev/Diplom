# Backend — Техническое описание для ИИ

> Этот документ является **оглавлением** бэкенд-документации.  
> Цель — дать AI-ассистенту быстрый контекст и направить к детальным разделам.  
> API-эндпоинты: [docs/api-endpoints.md](./api-endpoints.md) и папка [docs/api/](./api/)

---

## Разделы документации

| Раздел | Описание |
|--------|----------|
| [Архитектура и стек](./backend/architecture.md) | Луковая архитектура, технологии, точка входа, конфигурация |
| [Аутентификация и сессии](./backend/auth.md) | OAuth потоки, JWT схема, роли, защита от кражи токенов |
| [HTTP-роутер и Middleware](./backend/routing.md) | Маршруты, chi middleware, DTO и обработчики |
| [Domain, Repository, Service](./backend/domain-repository-service.md) | Доменные модели, интерфейсы, DAO-паттерн, бизнес-логика |
| [Миграции БД](./backend/migrations.md) | Список миграций, устарелые таблицы, применение |
| [Логирование](./backend/logging.md) | slog, уровни логов, структура записей |

---

## Быстрый обзор

**Язык:** Go · **Роутер:** chi · **ORM:** GORM + pgx · **БД:** PostgreSQL · **JWT:** HS256

### Аутентификация
- Вход только через **Google OAuth** или **Telegram** (Telegram — заглушка)
- Нет логина/пароля: таблица `users` содержит `email`, `google_id`, `telegram_id`
- Access Token (JWT, 15 мин) + Refresh Token (cookie, 15 дней)

### Роли
`admin` | `qp` | `warehouse_manager` | `storekeeper` | `pharmacist`

### Схема БД
17 миграций покрывают все сущности ERP: продукты, серии, зоны склада, поставщики, накладные, заказы, инвентаризация, рекламации, журнал среды, аудит.  
Подробная схема: [docs/database-schema.md](./database-schema.md)

### Что не реализовано
- `LoginWithTelegram` — метод-заглушка (`panic("implement...")`), нужна верификация hash по bot token
- Valkey/Redis кэш сессий — описан в [docs/valkey-cache.md](./valkey-cache.md), в коде отсутствует
- Immutable audit log hash chain — поля `prev_hash`, `log_hash` добавлены в схему, логика не реализована
- Бизнес-модули (Inbound, Orders, Inventory, Claims и т.д.) — только схема в БД, handlers не написаны
