# Frontend — Стек и Конвенции

← [Назад к главному README](../../README.md) | [Обзор Frontend →](./README.md) | [Формы →](./forms/) | [Архитектура FSD →](./fsd.md)

---

## ⚡ Фреймворк и Сборка

**React 19 + Vite + TypeScript (Strict Mode)**

- Маршрутизация: **React Router v7** (`react-router-dom`).
- Архитектура: **Feature-Sliced Design (FSD)** (папка `src/`).

---

## 🎨 Стилизация (Tailwind CSS v4)

Используется **Tailwind CSS v4** с подключением через `@import "tailwindcss"` в `src/index.css`.
Файл `index.css` сгенерирован через shadcn CLI (пресет `b2p9XgkXg`, тема `radix-lyra`). Строго запрещено менять существующие токены или добавлять новые без согласования.

### Правило цветов — КРИТИЧЕСКИ ВАЖНО

**Используйте только семантические CSS-переменные из `index.css`. Использование сырых утилит Tailwind запрещено.**

| ❌ Запрещено | ✅ Правильно (Семантика) | Описание |
|-------------|-----------|----------|
| `bg-gray-800` | `bg-primary` | Основные интерактивные элементы |
| `text-white` | `text-primary-foreground` | Текст на primary |
| `bg-green-500` | `bg-secondary` | Вторичные акценты |
| `text-red-500` | `text-destructive` | Ошибки, удаление, блокировки |
| `bg-blue-500` | `bg-info` | Информационные сообщения |
| `bg-amber-500` | `bg-warning` | Предупреждения |
| `bg-purple-500` | `bg-notification` | Уведомления |
| `bg-gray-100` | `bg-muted` | Приглушенный фон |
| `border-gray-200` | `border-border` | Границы и разделители |
| `text-gray-500` | `text-muted-foreground` | Второстепенный текст |

### Утилита `cn()`

Для объединения и условного применения классов всегда используйте `cn()`:
```ts
import { cn } from "@/shared/lib/utils"

cn("base-class", condition && "conditional-class", className)
```

---

## 🧩 UI-компоненты (shadcn/ui)

Компоненты располагаются в `src/shared/ui/`.

Установка новых компонентов:
```bash
npx shadcn@latest add <component-name>
```

---

## 🖼️ Иконки

Используется библиотека **lucide-react** (и `@phosphor-icons/react` при необходимости).

```tsx
import { Bell, Search, Package } from "lucide-react"

<Bell className="h-5 w-5 text-muted-foreground" />
```

---

## 🌐 HTTP-запросы и API

Для взаимодействия с бэкендом настроен глобальный API-клиент на базе **Axios** в **`src/shared/api`**:
- **`apiClient`** — для защищенных запросов (автоматически подставляет заголовок `Authorization: Bearer <accessToken>` и выполняет автоматический фоновый silent-refresh сессии при получении ошибки `401` через HttpOnly-куку с рефреш-токеном).
- **`authClient`** — чистый экземпляр без интерцепторов (для логина, регистрации и обновления токенов).

Глобальное состояние сессии и профиля пользователя управляется через стор **Zustand** в **`src/entities/user`** (с персистентностью безопасных данных в `localStorage`).

```ts
import { apiClient } from "@/shared/api";

// Пример выполнения запроса с авто-авторизацией
const res = await apiClient.get("/api/v1/products");
const data = res.data;
```

---

## 📋 Валидация форм

Используется связка **Zod v4** + `react-hook-form` + компоненты формы shadcn/ui.

```tsx
import { z } from "zod"
import { useForm } from "react-hook-form"
import { zodResolver } from "@hookform/resolvers/zod"

const schema = z.object({ email: z.string().email() })
```

---

## 🗺️ Управление состоянием URL

Для работы с параметрами URL используется хук `useSearchParams` из **React Router**:

```ts
import { useSearchParams } from "react-router-dom"

const [searchParams, setSearchParams] = useSearchParams()
const query = searchParams.get("q") || ""
```

---

## 📅 Работа с датами

Используется **date-fns v4** для форматирования и расчетов дат.

```ts
import { format, differenceInDays } from "date-fns"
import { ru } from "date-fns/locale"

format(new Date(date), "dd.MM.yyyy", { locale: ru })
```
