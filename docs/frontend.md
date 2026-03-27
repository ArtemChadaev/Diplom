# Frontend — Стек и соглашения

> Этот файл — источник правды о технологиях и паттернах фронтенда.  
> Прочитав его, нейронка должна понимать как писать код в этом проекте.  
> Полная спецификация форм: [next-form-ai.md](./next-form-ai.md)

---

## Фреймворк

**Next.js 16 — App Router**

- Директория приложения: `frontend/src/app/`
- Используется App Router (не Pages Router)
- Роут-группы для изоляции лейаутов:
  - `(app)/` — основной интерфейс (Header + Footer)
  - `(auth)/` — страницы входа (упрощённый лейаут)
  - `(admin)/` — административный раздел
- Серверные компоненты по умолчанию; `"use client"` только там, где нужен стейт/интерактивность
- Язык: **TypeScript** (строгий режим)
- React версия: **19**

---

## Стилизация

### Tailwind CSS v4

Проект использует **Tailwind CSS v4** с `@import "tailwindcss"` (не старый `@tailwind base` синтаксис).

Конфиг темы переопределяется через `@theme inline` в `globals.css`.

### Правило цветов — КРИТИЧНО

**Использовать ТОЛЬКО семантические переменные из `globals.css`. Никаких произвольных цветов Tailwind.**

| ❌ Запрещено | ✅ Правильно |
|-------------|-------------|
| `bg-gray-800` | `bg-primary` |
| `text-white` | `text-primary-foreground` |
| `bg-green-500` | `bg-secondary` |
| `text-red-500` | `text-destructive` |
| `bg-gray-100` | `bg-muted` |
| `border-gray-200` | `border-border` |
| `text-gray-500` | `text-muted-foreground` |

### Переменные темы из `globals.css`

```
--background          Основной фон страницы
--foreground          Основной текст
--card                Фон карточек
--card-foreground     Текст на карточках
--primary             Кнопки, активные элементы (#181619)
--primary-foreground  Текст на primary (#ffffff)
--secondary           Акцентный цвет, brand (#006b58)
--secondary-foreground Текст на secondary (#ffffff)
--muted               Приглушённый фон (#f3f3f3)
--muted-foreground    Приглушённый текст (#7a767a)
--accent              Hover-фон для элементов (#e2e2e2)
--accent-foreground   Текст на accent (#49454a)
--destructive         Ошибки, удаление (#ba1a1a)
--destructive-foreground Текст на destructive (#ffffff)
--border              Цвет границ (#e2e2e2)
--input               Цвет обводки инпутов (#e2e2e2)
--ring                Цвет focus-ring (#006b58)
--blue                Информационный (#0077ff)
--blue-foreground     Текст на blue (#ffffff)
```

Углы скруглений через `--radius` (0.625rem); используй `rounded-sm`, `rounded-md`, `rounded-lg`, `rounded-xl` — они привязаны к `--radius` через `@theme`.

### `cn()` — обязательно для условных классов

```ts
import { cn } from "@/lib/utils"
// Внутри: clsx + tailwind-merge
cn("base-class", condition && "conditional-class", className)
```

---

## UI-компоненты — shadcn/ui

Компоненты лежат в `src/components/ui/`. Используются через `@/components/ui/...`.

### Уже установленные компоненты

| Файл | Компоненты |
|------|-----------|
| `avatar.tsx` | `Avatar`, `AvatarImage`, `AvatarFallback` |
| `badge.tsx` | `Badge` |
| `button.tsx` | `Button` |
| `button-group.tsx` | `ButtonGroup` (кастомный) |
| `calendar.tsx` | `Calendar` |
| `card.tsx` | `Card`, `CardHeader`, `CardTitle`, `CardContent`, `CardFooter` |
| `checkbox.tsx` | `Checkbox` |
| `command.tsx` | `Command`, `CommandInput`, `CommandList`, `CommandItem`, `CommandGroup`, `CommandEmpty` |
| `date-picker.tsx` | `DatePicker` (обёртка над Calendar + Popover) |
| `dialog.tsx` | `Dialog`, `DialogContent`, `DialogHeader`, `DialogTitle`, `DialogFooter` |
| `dropdown-menu.tsx` | `DropdownMenu`, `DropdownMenuTrigger`, `DropdownMenuContent`, `DropdownMenuItem` |
| `input.tsx` | `Input` |
| `input-group.tsx` | `InputGroup` (кастомный) |
| `label.tsx` | `Label` |
| `multi-select.tsx` | `MultiSelect` (кастомный, через Command + Popover + Badge) |
| `pagination.tsx` | `Pagination` |
| `popover.tsx` | `Popover`, `PopoverTrigger`, `PopoverContent` |
| `select.tsx` | `Select`, `SelectContent`, `SelectItem`, `SelectTrigger`, `SelectValue` |
| `separator.tsx` | `Separator` |
| `table.tsx` | `Table`, `TableHeader`, `TableBody`, `TableRow`, `TableHead`, `TableCell` |
| `tabs.tsx` | `Tabs`, `TabsList`, `TabsTrigger`, `TabsContent` |
| `textarea.tsx` | `Textarea` |
| `toggle.tsx` | `Toggle` |
| `toggle-group.tsx` | `ToggleGroup`, `ToggleGroupItem` |

### Установка новых компонентов

```bash
npx shadcn@latest add <component-name>
```

> При установке компоненты попадают в `src/components/ui/` и автоматически используют переменные из `globals.css`.

---

## Иконки

**lucide-react** — единственная используемая библиотека иконок.

```tsx
import { Bell, Search, Package } from "lucide-react"
<Bell className="h-5 w-5 text-muted-foreground" />
```

Размеры: `h-4 w-4` (мелкие), `h-5 w-5` (стандарт), `h-6 w-6` (крупные).

---

## Типографика

**Inter** (Google Fonts) — подключён через `next/font/google` в `layout.tsx`, применяется через CSS-переменную `--font-sans`.

```tsx
// layout.tsx
const inter = Inter({ subsets: ['latin'], variable: '--font-sans' })
<html className={cn("font-sans", inter.variable)}>
```

Не использовать другие шрифты.

---

## HTTP-запросы

Используется нативный **`fetch`** (Next.js расширяет его кэшированием).

**Серверный компонент (RSC):**
```ts
// Данные фетчатся прямо в async-компоненте
export default async function Page() {
  const data = await fetch("/api/v1/products", {
    headers: { Authorization: `Bearer ${token}` },
    next: { revalidate: 60 }, // ISR опционально
  }).then(r => r.json())
}
```

**Клиентский компонент (`"use client"`):**
```ts
const res = await fetch("/api/v1/auth/send-code", {
  method: "POST",
  headers: { "Content-Type": "application/json" },
  body: JSON.stringify({ email }),
})
if (!res.ok) {
  const err = await res.json()
  throw new Error(err.error)
}
const data = await res.json()
```

- Нет axios, нет react-query — только `fetch`
- Токены: `httpOnly`-cookie (добавляются сервером) или `Authorization: Bearer` header из cookie/localStorage
- Ошибки: всегда проверять `res.ok`, парсить `{ error: string }` из JSON

---

## Валидация форм

**Zod v4** + паттерн react-hook-form (с shadcn/ui Form-компонентами).

```ts
import { z } from "zod"

const schema = z.object({
  email: z.string().email("Некорректный email"),
})
type FormData = z.infer<typeof schema>
```

Форма строится через shadcn/ui Form-компоненты:
```tsx
<Form {...form}>
  <FormField
    control={form.control}
    name="email"
    render={({ field }) => (
      <FormItem>
        <FormLabel>Email</FormLabel>
        <FormControl><Input {...field} /></FormControl>
        <FormMessage />  {/* авто-показ ошибки */}
      </FormItem>
    )}
  />
</Form>
```

> `react-hook-form` нет в `package.json` — нужно добавить вместе с `@hookform/resolvers` при установке `Form` компонента shadcn: `npx shadcn@latest add form`

---

## URL-параметры и поиск

**nuqs** — управление поиском через URL (`useQueryState`, `useQueryStates`).

```ts
import { parseAsString, parseAsArrayOf, parseAsInteger } from "nuqs"
import { useQueryStates } from "nuqs"

const [params, setParams] = useQueryStates({
  q: parseAsString.withDefault(""),
  page: parseAsInteger.withDefault(1),
})
```

`NuqsAdapter` обёрнут вокруг `<body>` в корневом `layout.tsx`.  
Используется для: поиска, фильтров, пагинации — всего что должно сохраняться в URL.

---

## Работа с датами

**date-fns v4** для форматирования и вычислений.

```ts
import { format, differenceInDays } from "date-fns"
import { ru } from "date-fns/locale"

format(new Date(expiryDate), "dd.MM.yyyy", { locale: ru })
differenceInDays(new Date(expiryDate), new Date())
```

**react-day-picker v9** используется внутри `DatePicker` компонента (уже есть в проекте).

---

## Структура папок

```
frontend/src/
├── app/
│   ├── (app)/          # Основные страницы (с Header)
│   │   ├── layout.tsx
│   │   ├── page.tsx    # Dashboard
│   │   └── search/
│   ├── (auth)/         # Страницы входа (без Header)
│   │   ├── layout.tsx
│   │   └── auth/       # /auth/login, /auth/verify
│   ├── (admin)/        # Админ-страницы
│   │   └── admin/
│   ├── globals.css     # Единственный CSS-файл, цвета и тема
│   └── layout.tsx      # Корневой layout (шрифт, NuqsAdapter)
├── components/
│   ├── ui/             # shadcn/ui компоненты
│   ├── header.tsx
│   └── footer.tsx
└── lib/
    └── utils.ts        # cn() хелпер
```

---

## Ключевые соглашения

1. **`"use client"`** — только если нужен `useState`, `useEffect`, обработчик событий
2. **Цвета** — только из `globals.css` (семантические переменные, см. таблицу выше)
3. **Импорты** — алиас `@/` = `src/`
4. **Компоненты** — именованный экспорт (`export function Foo`), не дефолтный
5. **Типы** — объявлять явно, `any` не использовать
6. **Иконки** — только `lucide-react`
7. **Даты** — только `date-fns` для форматирования и расчётов
