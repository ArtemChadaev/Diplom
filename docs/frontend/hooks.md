# hooks/ — Business Logic Directory

← [Back to Frontend Overview](./README.md) | [Forms Logic →](./forms-logic.md)

This directory (`frontend/src/hooks/`) is the **only** place for frontend business logic.

## Rule: No Logic in Components

All of the following **must** live in a hook in this directory — never inline in a page or component:

- Form state management (`react-hook-form` + `useForm`)
- Zod schema instantiation and validation
- API calls (`fetch`)
- Derived state or computed values
- Side effects tied to business actions

## Colocation Rule (Go-style)

Test files live **next to** their source:

```
hooks/
├── use-auth-login.ts
├── use-auth-login.test.ts    ← same folder, not __tests__/
├── use-otp-verify.ts
└── use-otp-verify.test.ts
```

## TDD Rule

**Write the test before writing the hook.**

```
Step 1: Write use-xxx.test.ts with scenario (Given/When/Then in comments)
Step 2: Present test for review
Step 3: Implement use-xxx.ts
Step 4: Confirm tests pass
```

A hook is not considered complete until its tests pass.

## Naming Convention

```
use-{feature-name}.ts
use-{feature-name}.test.ts
```

Examples:
- `use-auth-login.ts` — login form logic
- `use-otp-verify.ts` — OTP verification logic
- `use-inbound-receiving.ts` — inbound goods form logic

## Component Rule

Components in `components/` must be "dumb":
- Accept props or call hooks
- Handle JSX only
- Zero `fetch` calls — all API calls go through hooks
- Zero Zod schemas — all validation goes through hooks

## Logic Specs

For the full list of required fields, Zod schemas, and API contracts per feature,
see [`forms-logic.md`](./forms-logic.md).
