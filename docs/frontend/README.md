# Frontend — Overview & Development Standards

← [Back to Main README](../../README.md) | [Full Conventions →](./conventions.md) | [Forms UI Specs →](./forms/) | [Forms Logic →](./forms-logic.md)

> ⛔ **Branch Rule:** All frontend code **must** be written in `develop-frontend` only.
> Direct commits to `main` or `develop-backend` are strictly forbidden.
> Merging into `main` is done exclusively via the [`/merge-all`](../other/git-workflow.md) workflow.

Next.js 16 App Router application for the Pharmaceutical ERP system. TypeScript-strict, server-components-first.

> **AI Context & Navigation:** This README is your starting point. It contains the most critical mandates. For specific technical details, jump to the relevant sub-document:
> - Writing components, styles, fetching data or dates? Read [**`conventions.md`**](./conventions.md).
> - Need form UI / layouts? Check [**`forms/`**](./forms/).
> - Writing form logic, API calls, or schemas? Use [**`forms-logic.md`**](./forms-logic.md) and verify with [`docs/api/swagger.json`](../../docs/api/swagger.json).
> - Writing hooks / business logic? Read [**`hooks.md`**](./hooks.md) for Colocation and TDD rules.

---

## Tech Stack

| Technology | Version | Role |
|-----------|---------|------|
| **Next.js** | 16 (App Router) | Framework |
| **React** | 19 | UI runtime |
| **TypeScript** | Strict mode | Language |
| **Tailwind CSS** | v4 | Styling |
| **shadcn/ui** | Latest | Component library |
| **lucide-react** | Latest | Icons (only) |
| **Zod** | v4 | Schema validation |
| **nuqs** | Latest | URL state management |
| **date-fns** | v4 | Date formatting & calculations |

---

## Local Development

```bash
cd frontend
pnpm install
pnpm dev    # http://localhost:3000
```

Or via Docker Compose from the monorepo root:

```bash
docker compose --profile all up -d
```

---

## global.css — Immutability Mandate

> ⛔ **`src/app/globals.css` is frozen. Do not modify it under any circumstances.**

Rules:
1. **Zero new variables** — do not add colors, shadows, or border radii
2. **Zero raw Tailwind utilities** — never use `bg-gray-800`, `text-red-500`, etc.
3. **Semantic variables only** — use only what is already defined: `bg-primary`, `text-foreground`, `rounded-xl`, etc.
4. **Modifying shadcn base components** — requires explicit approval before any change

```
✅  bg-primary          ❌  bg-gray-800
✅  text-destructive    ❌  text-red-500
✅  border-border       ❌  border-gray-200
✅  rounded-xl          ❌  rounded-[14px]
```

→ Full variable reference: [`conventions.md — Styling`](./conventions.md#styling)

---

## Project Structure

```
frontend/src/
├── app/
│   ├── (app)/          # Main pages (Header + Footer layout)
│   │   ├── layout.tsx
│   │   ├── page.tsx    # Dashboard
│   │   └── search/
│   ├── (auth)/         # Login pages (minimal layout)
│   │   ├── layout.tsx
│   │   └── auth/       # /auth/login, /auth/verify
│   ├── (admin)/        # Admin section
│   │   └── admin/
│   ├── globals.css     # ⛔ FROZEN — all theme tokens (do not touch)
│   └── layout.tsx      # Root layout (font, NuqsAdapter)
├── components/
│   ├── ui/             # shadcn/ui components (base — do not modify without approval)
│   ├── header.tsx      # "Dumb" — no fetch, no Zod inside
│   └── footer.tsx
├── hooks/              # ← All business logic lives here (see §Folder Structure)
│   └── README.md
└── lib/
    └── utils.ts        # cn() helper (clsx + tailwind-merge)
```

---

## Folder Structure — Logic Isolation

### The Core Rule

> **Components are dumb. Hooks are smart.**

| Where | What goes here |
|-------|---------------|
| `components/*.tsx` | JSX only — markup, props, calls to hooks |
| `hooks/use-*.ts` | Zod schemas, `useForm`, `fetch`, derived state |
| `app/**/page.tsx` | Layout assembly — calls Server Components or hooks only |

### ✅ Correct Pattern

```tsx
// hooks/use-login-form.ts  ← logic here
export function useLoginForm() {
  const form = useForm<z.infer<typeof emailSchema>>({
    resolver: zodResolver(emailSchema),
  })
  const onSubmit = async (data: z.infer<typeof emailSchema>) => {
    const res = await fetch("/api/v1/auth/send-code", { ... })
  }
  return { form, onSubmit }
}

// components/login-form.tsx  ← UI only
export function LoginForm() {
  const { form, onSubmit } = useLoginForm()  // consumed here
  return <Form {...form}><Input .../><Button .../></Form>
}
```

### ❌ Forbidden Pattern

```tsx
// ❌ Logic inside a page or component
export default function LoginPage() {
  const form = useForm(...)    // ← forbidden
  const onSubmit = async () => {
    await fetch(...)           // ← forbidden
  }
}
```

### Colocation Rule (Go-style)

Test file lives **in the same directory** as its source:

```
hooks/
├── use-auth-login.ts
├── use-auth-login.test.ts    ← here, not in __tests__/

components/
├── login-form.tsx
└── login-form.test.tsx       ← here, not in __tests__/
```

---

## Stitch & UI Prototyping Workflow

### Task Split

UI design and logic implementation are **two separate, independent tasks**. They must never be mixed.

| Task type | Who/What | Output |
|-----------|---------|--------|
| UI Design | Stitch (https://stitch.withgoogle.com/) | Clean JSX, no logic |
| Logic | Hook in `hooks/` | Zod, fetch, state |

### When to Use Stitch

Use Stitch for **every new form or screen**. Required cases:
- New page layout
- New form (any ERP screen from `docs/frontend/forms/`)
- New complex UI block (multi-step wizard, data table, card grid)

### Stitch Prompt Template

Always include this context when prompting Stitch:

```
Use Tailwind v4 with semantic CSS variables from globals.css.
Use shadcn/ui components. Create clean JSX without business logic,
fetch calls, or state management.
Variables available: bg-primary, text-foreground, bg-muted,
text-destructive, border-border, bg-secondary, rounded-xl, etc.
```

### Implementation Decision

| Stitch result | Action |
|---------------|--------|
| Matches project color palette and structure | Integrate immediately |
| Changes layout structure or adds new elements | Request approval first |
| Uses raw Tailwind colors (bg-gray-*, etc.) | Replace with semantic variables before integrating |

---

## Frontend TDD Rule

> **Write the test before writing the hook. No exceptions for logic with conditions.**

TDD cycle for every hook with non-trivial logic:

```
1. Describe scenario (Given/When/Then comment in test file)
2. Write the test  →  hooks/use-xxx.test.ts
3. Present test for approval
4. Implement the hook  →  hooks/use-xxx.ts
5. Confirm: `pnpm test -- use-xxx` passes
```

**Exempt from mandatory tests** (same rule as backend):
- Simple data fetching with no conditions or transformations
- Static display components

**Requires tests:**
- Any hook with `if` / branching logic
- Validation hooks (Zod refinements)
- Hooks with redirects or role-based behavior
- Hooks with timers or side effects

---

## Navigation Directory

To prevent reading unnecessary context, only look into these files if your current task requires them:

→ **Styling, HTTP, UI Components, Icons:** [`conventions.md`](./conventions.md)
→ **Hooks & Business Logic Rules:** [`hooks.md`](./hooks.md)
→ **Form Layouts & UI Specs:** [`forms/`](./forms/)
→ **Form API Logic & Validators:** [`forms-logic.md`](./forms-logic.md)
→ **Backend SOP (if full-stack feature):** [`docs/backend/sop.md`](../backend/sop.md)


