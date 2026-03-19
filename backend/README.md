# Backend — Архитектура и технологии

## Технологии

| Стек | Описание |
|------|----------|
| **Go 1.23** | Язык разработки |
| **chi v5** | HTTP-роутер |
| **GORM v2** | ORM для работы с PostgreSQL |
| **PostgreSQL 17** | Основная реляционная БД |
| **Valkey (Redis)** | Кэш / сессии (Redis-совместимый) |
| **golang-migrate** | Управление миграциями схемы БД |
| **Docker / Docker Compose** | Контейнеризация и оркестрация |
| **Viper** | Загрузка конфигурации из `.env` |

---

## Архитектура

Проект построен по трёхслойной чистой архитектуре:

```
cmd/main.go
    │
    ├── handler (HTTP)
    │     └── dto/          ← DTO: json + validate теги (слой API)
    │           ↓ маппинг
    ├── service              ← Бизнес-логика, работа с чистыми доменными моделями
    │         ↓
    └── repository           ← Доступ к БД через GORM
          └── dao/           ← DAO: gorm теги, TableName() (слой хранилища)
```

### Слои

| Слой | Пакет | Теги | Зона ответственности |
|------|-------|------|----------------------|
| **Handler (HTTP)** | `internal/handler` | — | Разбор запроса, валидация, вызов сервиса |
| **DTO** | `internal/handler/dto` | `json`, `validate` | Формат входящих/исходящих данных API |
| **Domain** | `internal/domain` | _(нет)_ | Чистые Go-структуры, бизнес-интерфейсы |
| **Service** | `internal/service` | — | Бизнес-логика, оркестрация |
| **Repository** | `internal/repository` | — | Реализация интерфейсов доступа к данным |
| **DAO** | `internal/repository/dao` | `gorm` | Слепки таблиц БД для GORM |

> **Правило:** domain-модели не знают ни о JSON, ни о GORM. Любые теги живут только в DTO или DAO.

---

## Структура директорий

```
backend/
├── cmd/
│   └── main.go
├── internal/
│   ├── config/
│   ├── domain/           ← интерфейсы + чистые модели
│   ├── handler/
│   │   ├── dto/          ← DTO структуры
│   │   ├── handler.go
│   │   └── ...
│   ├── repository/
│   │   ├── dao/          ← DAO структуры
│   │   ├── postgres.go
│   │   └── ...
│   └── service/
├── migrate/              ← SQL-миграции (golang-migrate)
├── Dockerfile
└── .env.example
```

---

## Запуск миграций

Миграции выполняются через официальный Docker-образ [`migrate/migrate`](https://github.com/golang-migrate/migrate).

### Применить все миграции (up)

```bash
docker run --rm \
  --network diplom_app_network \
  -v "$(pwd)/backend/migrate:/migrations" \
  migrate/migrate \
  -path=/migrations \
  -database "postgres://${DB_USER}:${DB_PASSWORD}@postgres:5432/${DB_NAME}?sslmode=disable" \
  up
```

### Откатить последнюю миграцию (down 1)

```bash
docker run --rm \
  --network diplom_app_network \
  -v "$(pwd)/backend/migrate:/migrations" \
  migrate/migrate \
  -path=/migrations \
  -database "postgres://${DB_USER}:${DB_PASSWORD}@postgres:5432/${DB_NAME}?sslmode=disable" \
  down 1
```

> **Перед запуском** убедись, что контейнер `postgres` запущен:
> ```bash
> docker compose up -d postgres
> ```
> Переменные `DB_USER`, `DB_PASSWORD`, `DB_NAME` берутся из `.env` (см. `.env.example`).

---

## Запуск локально

```bash
# из корня монорепо
docker compose up -d          # поднять всё
docker compose logs -f backend # смотреть логи бэкенда
```
