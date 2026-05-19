import js from '@eslint/js'
import globals from 'globals'
import reactPlugin from 'eslint-plugin-react'
import reactHooks from 'eslint-plugin-react-hooks'
import reactRefresh from 'eslint-plugin-react-refresh'
import tseslint from 'typescript-eslint'
import boundaries from 'eslint-plugin-boundaries'
import importX from 'eslint-plugin-import-x'
import { createTypeScriptImportResolver } from 'eslint-import-resolver-typescript'
import { defineConfig, globalIgnores } from 'eslint/config'

// ─── FSD layer order (high → low) ────────────────────────────────────────────
// Each layer may only import from layers listed AFTER it in this array.
const FSD_LAYERS = ['app', 'pages', 'widgets', 'features', 'entities', 'shared']

export default defineConfig([
  // ── Ignore built artifacts ─────────────────────────────────────────────────
  globalIgnores(['dist', 'node_modules', 'public']),

  // ── Base config for all TS/TSX files ──────────────────────────────────────
  {
    files: ['**/*.{ts,tsx}'],

    extends: [
      js.configs.recommended,
      tseslint.configs.strictTypeChecked,
      tseslint.configs.stylisticTypeChecked,
    ],

    plugins: {
      react: reactPlugin,
      'react-hooks': reactHooks,
      'react-refresh': reactRefresh,
      'import-x': importX,
      boundaries,
    },

    languageOptions: {
      globals: globals.browser,
      parserOptions: {
        // Required for type-aware rules (strictTypeChecked)
        projectService: true,
        tsconfigRootDir: import.meta.dirname,
      },
    },

    settings: {
      // ── React ──────────────────────────────────────────────────────────────
      react: {
        version: '19',
      },

      // ── eslint-plugin-import-x: new resolver-next API (ESLint 10 compatible)
      'import-x/resolver-next': [
        createTypeScriptImportResolver({
          alwaysTryTypes: true,
          project: './tsconfig.app.json',
        }),
      ],

      // ── FSD boundaries elements ────────────────────────────────────────────
      'boundaries/elements': [
        { type: 'app',      pattern: 'src/app/**' },
        { type: 'pages',    pattern: 'src/pages/**' },
        { type: 'widgets',  pattern: 'src/widgets/**' },
        { type: 'features', pattern: 'src/features/**' },
        { type: 'entities', pattern: 'src/entities/**' },
        { type: 'shared',   pattern: 'src/shared/**' },
      ],
      'boundaries/ignore': ['**/*.test.*', '**/*.spec.*'],
    },

    rules: {
      // ── React ───────────────────────────────────────────────────────────────
      ...reactPlugin.configs.recommended.rules,
      ...reactPlugin.configs['jsx-runtime'].rules, // React 19: no need to import React
      ...reactHooks.configs.recommended.rules,
      'react-refresh/only-export-components': ['warn', { allowConstantExport: true }],

      'react/display-name': 'error',
      'react/prop-types': 'off',           // covered by TypeScript
      'react/no-deprecated': 'error',
      'react/no-unknown-property': 'error',
      'react/jsx-key': ['error', { checkFragmentShorthand: true }],
      'react/self-closing-comp': 'warn',

      // ── TypeScript ──────────────────────────────────────────────────────────
      '@typescript-eslint/consistent-type-imports': [
        'error',
        { prefer: 'type-imports', fixStyle: 'inline-type-imports' },
      ],
      '@typescript-eslint/no-import-type-side-effects': 'error',
      '@typescript-eslint/no-unused-vars': [
        'error',
        { argsIgnorePattern: '^_', varsIgnorePattern: '^_' },
      ],
      '@typescript-eslint/no-explicit-any': 'warn',
      '@typescript-eslint/no-floating-promises': 'error',
      '@typescript-eslint/no-misused-promises': [
        'error',
        { checksVoidReturn: { attributes: false } }, // allow async onClick handlers
      ],
      '@typescript-eslint/prefer-nullish-coalescing': 'warn',
      '@typescript-eslint/prefer-optional-chain': 'warn',

      // ── Import order (import-x, ESLint 10 compatible) ───────────────────────
      'import-x/order': [
        'warn',
        {
          groups: [
            'builtin',
            'external',
            'internal',
            ['parent', 'sibling', 'index'],
            'type',
          ],
          pathGroups: FSD_LAYERS.map((layer) => ({
            pattern: `@/${layer}/**`,
            group: 'internal',
            position: 'after',
          })),
          pathGroupsExcludedImportTypes: ['type'],
          'newlines-between': 'always',
          alphabetize: { order: 'asc', caseInsensitive: true },
        },
      ],
      'import-x/no-cycle': 'error',
      'import-x/no-self-import': 'error',
      'import-x/no-duplicates': 'error',

      // ── FSD Layer Boundaries (eslint-plugin-boundaries v6) ─────────────────
      // Rule: a layer can only import from layers BELOW it.
      // Stack: app > pages > widgets > features > entities > shared
      // Note: `allow` uses v6 DependencySelector format: { to: objectElementMatcher }
      'boundaries/dependencies': [
        'error',
        {
          default: 'disallow',
          rules: [
            // app — can import anything
            {
              from: { type: 'app' },
              allow: FSD_LAYERS.map((l) => ({ to: { type: l } })),
            },
            // pages — can import widgets, features, entities, shared
            {
              from: { type: 'pages' },
              allow: ['widgets', 'features', 'entities', 'shared'].map((l) => ({ to: { type: l } })),
            },
            // widgets — can import features, entities, shared
            {
              from: { type: 'widgets' },
              allow: ['features', 'entities', 'shared'].map((l) => ({ to: { type: l } })),
            },
            // features — can import entities, shared
            {
              from: { type: 'features' },
              allow: ['entities', 'shared'].map((l) => ({ to: { type: l } })),
            },
            // entities — can import only shared
            {
              from: { type: 'entities' },
              allow: [{ to: { type: 'shared' } }],
            },
            // shared — no FSD imports (pure utilities)
            {
              from: { type: 'shared' },
              allow: [],
            },
          ],
        },
      ],

      // ── General code quality ────────────────────────────────────────────────
      'no-console': ['warn', { allow: ['warn', 'error'] }],
      'prefer-const': 'error',
      'no-var': 'error',
      eqeqeq: ['error', 'always'],
    },
  },

  // ── Relax strict rules for vite/eslint config files themselves ─────────────
  {
    files: ['eslint.config.js', 'vite.config.ts'],
    rules: {
      '@typescript-eslint/no-require-imports': 'off',
    },
  },

  // ── Relax strict TS rules for auto-generated shadcn/ui components ──────────
  // These files are generated by `npx shadcn@latest add` and are not manually
  // maintained. Strict type rules like no-non-null-assertion and
  // no-unnecessary-condition would produce noise without value here.
  {
    files: ['src/shared/ui/**/*.{ts,tsx}'],
    rules: {
      '@typescript-eslint/no-non-null-assertion': 'off',
      '@typescript-eslint/no-unnecessary-condition': 'off',
      '@typescript-eslint/prefer-promise-reject-errors': 'off',
      '@typescript-eslint/use-unknown-in-catch-callback-variable': 'off',
      '@typescript-eslint/return-await': 'off',
      '@typescript-eslint/prefer-nullish-coalescing': 'off',
      '@typescript-eslint/consistent-type-definitions': 'off',
      '@typescript-eslint/no-import-type-side-effects': 'off',
      'react-refresh/only-export-components': 'off',
    },
  },
  // ── Same relaxations for legacy src/components/ui (old shadcn location) ───
  {
    files: ['src/components/ui/**/*.{ts,tsx}'],
    rules: {
      '@typescript-eslint/no-non-null-assertion': 'off',
      '@typescript-eslint/no-unnecessary-condition': 'off',
      '@typescript-eslint/prefer-promise-reject-errors': 'off',
      '@typescript-eslint/use-unknown-in-catch-callback-variable': 'off',
      '@typescript-eslint/return-await': 'off',
      '@typescript-eslint/prefer-nullish-coalescing': 'off',
      '@typescript-eslint/consistent-type-definitions': 'off',
      '@typescript-eslint/no-import-type-side-effects': 'off',
      'react-refresh/only-export-components': 'off',
    },
  },
])
