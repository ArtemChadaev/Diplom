# Diplom — Git Flow Guide

Проект использует адаптированный **Git Flow** для монорепозитория с раздельными ветками backend и frontend.

---

## Структура веток

```
main                         ← продакшн (стабильный релиз)
develop                      ← интеграционная ветка (backend + frontend)
  ├── develop/backend        ← интеграция backend-изменений
  └── develop/frontend       ← интеграция frontend-изменений

feature/all/<name>           ← фича, затрагивающая обе части
feature/backend/<name>       ← фича только для backend
feature/frontend/<name>      ← фича только для frontend

release/all/<version>        ← общий релиз
release/backend/<version>    ← релиз только backend
release/frontend/<version>   ← релиз только frontend

hotfix/all/<name>            ← срочный фикс (backend + frontend)
hotfix/backend/<name>        ← срочный фикс backend
hotfix/frontend/<name>       ← срочный фикс frontend
```

---

## Правила именования

| Тип        | Шаблон                          | Базовая ветка         |
|------------|---------------------------------|-----------------------|
| Feature    | `feature/<scope>/<name>`       | `develop/<scope>`     |
| Release    | `release/<scope>/<version>`    | `develop`             |
| Hotfix     | `hotfix/<scope>/<name>`        | `main`                |

Где `<scope>` — одно из: `all`, `backend`, `frontend`.

---

## Проект

| Папка      | Описание              |
|------------|-----------------------|
| `backend/` | Go-сервер             |
| `frontend/`| Next.js приложение    |

