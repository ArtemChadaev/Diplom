# Frontend — Stack & Conventions (Full Reference)

← [Back to Main README](../../README.md) | [Frontend Overview →](./README.md) | [Forms →](./forms/)

---

## Framework

**Next.js 16 — App Router**

- Application directory: `frontend/src/app/`
- Uses App Router (not Pages Router)
- Route groups for layout isolation:
  - `(app)/` — main interface (Header + Footer)
  - `(auth)/` — login pages (simplified layout)
  - `(admin)/` — admin section
- Server Components by default; `"use client"` only where state/interactivity is needed
- Language: **TypeScript** (strict mode)
- React version: **19**

---

## Styling

### Tailwind CSS v4

Uses **Tailwind CSS v4** with `@import "tailwindcss"` (not the legacy `@tailwind base` syntax).

Theme configuration is overridden via `@theme inline` in `globals.css`.

### Color Rule — CRITICAL

**Only use semantic variables from `globals.css`. No raw Tailwind color utilities.**

| ❌ Forbidden | ✅ Correct |
|-------------|-----------|
| `bg-gray-800` | `bg-primary` |
| `text-white` | `text-primary-foreground` |
| `bg-green-500` | `bg-secondary` |
| `text-red-500` | `text-destructive` |
| `bg-gray-100` | `bg-muted` |
| `border-gray-200` | `border-border` |
| `text-gray-500` | `text-muted-foreground` |

### Theme Variables (`globals.css`)

```
--background          Main page background
--foreground          Main text color
--card                Card background
--card-foreground     Text on cards
--primary             Buttons, active elements (#181619)
--primary-foreground  Text on primary (#ffffff)
--secondary           Accent / brand color (#006b58)
--secondary-foreground Text on secondary (#ffffff)
--muted               Muted background (#f3f3f3)
--muted-foreground    Muted text (#7a767a)
--accent              Hover background (#e2e2e2)
--accent-foreground   Text on accent (#49454a)
--destructive         Errors, delete actions (#ba1a1a)
--destructive-foreground Text on destructive (#ffffff)
--border              Border color (#e2e2e2)
--input               Input border color (#e2e2e2)
--ring                Focus ring color (#006b58)
--blue                Informational (#0077ff)
--blue-foreground     Text on blue (#ffffff)
```

Corner radius via `--radius` (0.625rem); use `rounded-sm`, `rounded-md`, `rounded-lg`, `rounded-xl` — they are bound to `--radius` via `@theme`.

### `cn()` — required for conditional classes

```ts
import { cn } from "@/lib/utils"
// Internally: clsx + tailwind-merge
cn("base-class", condition && "conditional-class", className)
```

---

## UI Components — shadcn/ui

Components live in `src/components/ui/`. Import via `@/components/ui/...`.

### Installed Components

| File | Components |
|------|-----------|
| `avatar.tsx` | `Avatar`, `AvatarImage`, `AvatarFallback` |
| `badge.tsx` | `Badge` |
| `button.tsx` | `Button` |
| `button-group.tsx` | `ButtonGroup` (custom) |
| `calendar.tsx` | `Calendar` |
| `card.tsx` | `Card`, `CardHeader`, `CardTitle`, `CardContent`, `CardFooter` |
| `checkbox.tsx` | `Checkbox` |
| `command.tsx` | `Command`, `CommandInput`, `CommandList`, `CommandItem`, `CommandGroup`, `CommandEmpty` |
| `date-picker.tsx` | `DatePicker` (wrapper over Calendar + Popover) |
| `dialog.tsx` | `Dialog`, `DialogContent`, `DialogHeader`, `DialogTitle`, `DialogFooter` |
| `dropdown-menu.tsx` | `DropdownMenu`, `DropdownMenuTrigger`, `DropdownMenuContent`, `DropdownMenuItem` |
| `input.tsx` | `Input` |
| `input-group.tsx` | `InputGroup` (custom) |
| `label.tsx` | `Label` |
| `multi-select.tsx` | `MultiSelect` (custom, via Command + Popover + Badge) |
| `pagination.tsx` | `Pagination` |
| `popover.tsx` | `Popover`, `PopoverTrigger`, `PopoverContent` |
| `select.tsx` | `Select`, `SelectContent`, `SelectItem`, `SelectTrigger`, `SelectValue` |
| `separator.tsx` | `Separator` |
| `table.tsx` | `Table`, `TableHeader`, `TableBody`, `TableRow`, `TableHead`, `TableCell` |
| `tabs.tsx` | `Tabs`, `TabsList`, `TabsTrigger`, `TabsContent` |
| `textarea.tsx` | `Textarea` |
| `toggle.tsx` | `Toggle` |
| `toggle-group.tsx` | `ToggleGroup`, `ToggleGroupItem` |

### Installing New Components

```bash
npx shadcn@latest add <component-name>
```

> Components are placed in `src/components/ui/` and automatically use variables from `globals.css`.

---

## Icons

**lucide-react** — the only icon library used.

```tsx
import { Bell, Search, Package } from "lucide-react"
<Bell className="h-5 w-5 text-muted-foreground" />
```

Sizes: `h-4 w-4` (small), `h-5 w-5` (standard), `h-6 w-6` (large).

---

## Typography

**Inter** (Google Fonts) — connected via `next/font/google` in `layout.tsx`, applied via the `--font-sans` CSS variable.

```tsx
// layout.tsx
const inter = Inter({ subsets: ['latin'], variable: '--font-sans' })
<html className={cn("font-sans", inter.variable)}>
```

No other fonts allowed.

---

## HTTP Requests

Use the native **`fetch`** API (Next.js extends it with caching).

**Server Component (RSC):**
```ts
export default async function Page() {
  const data = await fetch("/api/v1/products", {
    headers: { Authorization: `Bearer ${token}` },
    next: { revalidate: 60 }, // ISR optional
  }).then(r => r.json())
}
```

**Client Component (`"use client"`):**
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

- No axios, no react-query — native `fetch` only
- Tokens: `httpOnly` cookies (set by server) or `Authorization: Bearer` header
- Errors: always check `res.ok`, parse `{ error: string }` from JSON

---

## Form Validation

**Zod v4** + shadcn/ui Form components.

```ts
import { z } from "zod"

const schema = z.object({
  email: z.string().email("Invalid email"),
})
type FormData = z.infer<typeof schema>
```

Form built with shadcn/ui Form components:
```tsx
<Form {...form}>
  <FormField
    control={form.control}
    name="email"
    render={({ field }) => (
      <FormItem>
        <FormLabel>Email</FormLabel>
        <FormControl><Input {...field} /></FormControl>
        <FormMessage />  {/* auto-shows validation error */}
      </FormItem>
    )}
  />
</Form>
```

> `react-hook-form` and `@hookform/resolvers` are required for the `Form` component: `npx shadcn@latest add form`

---

## URL State Management

**nuqs** — URL search param state management (`useQueryState`, `useQueryStates`).

```ts
import { parseAsString, parseAsArrayOf, parseAsInteger } from "nuqs"
import { useQueryStates } from "nuqs"

const [params, setParams] = useQueryStates({
  q: parseAsString.withDefault(""),
  page: parseAsInteger.withDefault(1),
})
```

`NuqsAdapter` is wrapped around `<body>` in the root `layout.tsx`.
Use for: search, filters, pagination — anything that must persist in the URL.

---

## Working with Dates

**date-fns v4** for formatting and calculations.

```ts
import { format, differenceInDays } from "date-fns"
import { ru } from "date-fns/locale"

format(new Date(expiryDate), "dd.MM.yyyy", { locale: ru })
differenceInDays(new Date(expiryDate), new Date())
```

**react-day-picker v9** is used internally by the `DatePicker` component (already in the project).

---

## Key Conventions Summary

1. **`"use client"`** — only when `useState`, `useEffect`, or event handlers are needed
2. **Colors** — only from `globals.css` (semantic variables, see table above)
3. **Imports** — use `@/` alias (`@/` = `src/`)
4. **Components** — named export (`export function Foo`), not default
5. **Types** — declare explicitly, no `any`
6. **Icons** — `lucide-react` only
7. **Dates** — `date-fns` only for formatting and calculations
