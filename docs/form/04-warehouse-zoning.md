# Форма 3.2а: Зонирование склада

## Требования (пресловой-чеклист)

- [ ] Таблица `warehouse_zones`: `id`, `name`, `type: enum(general, cold_chain, flammable, safe_strong)`, `description`, `temp_min`, `temp_max`, `humidity_max`
- [ ] Связь `storage_locations`: `id`, `zone_id`, `rack`, `shelf`, `cell` (адресное хранение — Фаза 4)
- [ ] Бэкенд: `GET /api/v1/zones` — список зон
- [ ] Бэкенд: `POST /api/v1/zones` — (admin) создание зоны
- [ ] Бэкенд: `PUT /api/v1/zones/:id` — (admin) редактирование зоны
- [ ] Бэкенд: `GET /api/v1/zones/:id/stock` — остатки в зоне
- [ ] shadcn/ui компоненты: `npx shadcn@latest add card badge table tooltip`
- [ ] Тип `cold_chain` — в зону могут попадать только товары с флагом `cold_chain: true`
- [ ] Тип `safe_strong` — только пользователи с `ns_pv_access: true` видят остатки

---

## Интерфейс

### Лейаут страницы

- Страница `/warehouse/zones` — доступна `admin`, `qp`, `warehouse_manager`
- Отображение в виде сетки карточек зон (`grid 2x2` или `3x...`)
- Каждая карточка — краткая информация о зоне + переход к детальному просмотру

### Компоненты shadcn/ui

| Компонент | Назначение |
|-----------|-----------|
| `<Card>`, `<CardHeader>`, `<CardContent>`, `<CardFooter>` | Карточка зоны |
| `<Badge>` | Тип зоны (цветовая кодировка) |
| `<Table>` | Остатки товаров в зоне |
| `<Tooltip>` | Подсказки о требованиях к зоне |
| `<Dialog>` | Форма создания/редактирования зоны (admin) |
| `<Form>`, `<Input>`, `<Select>` | Поля в диалоге редактирования |

### Цветовая кодировка типов зон

| Тип зоны | Badge color | Иконка |
|----------|-------------|--------|
| `general` | `default` | 📦 |
| `cold_chain` | `blue` | ❄️ |
| `flammable` | `orange` | 🔥 |
| `safe_strong` | `red` | 🔒 |

### Структура UI

```
┌──────────────────────────────────────────┐
│   Склад: Зонирование          [+ Зона]   │
├──────────────────────────────────────────┤
│ ┌──────────────┐  ┌──────────────┐       │
│ │ 📦 Общая    │  │ ❄️ Холодовая  │       │
│ │ зона        │  │ цепь          │       │
│ │             │  │ 2-8°C | ≤60% │       │
│ │ 142 позиции │  │ 38 позиций    │       │
│ │ [Просмотр]  │  │ [Просмотр]   │       │
│ └──────────────┘  └──────────────┘       │
│ ┌──────────────┐  ┌──────────────┐       │
│ │ 🔥 Огне-    │  │ 🔒 Сейфовая │       │
│ │ опасные     │  │ зона (НС/ПВ) │       │
│ │ 8 позиций   │  │ [Доступ огр.]│       │
│ └──────────────┘  └──────────────┘       │
└──────────────────────────────────────────┘
```

**Детальный просмотр зоны `/warehouse/zones/:id`:**
```
Зона: Холодовая цепь  [Редактировать]
Требования: 2–8°C | Влажность: ≤60%

Таблица остатков:
│ МНН │ Серия │ Годен до │ Кол-во │ Статус │
```

---

## Логика

### Данные с сервера

```
GET /api/v1/zones
Response: Array<{
  id: string
  name: string
  type: "general" | "cold_chain" | "flammable" | "safe_strong"
  description: string
  temp_min: number | null
  temp_max: number | null
  humidity_max: number | null
  stock_count: number    // кол-во позиций в зоне
}>

GET /api/v1/zones/:id/stock
Response: {
  zone: ZoneDTO,
  items: Array<{
    product_id, mnn, serial_number,
    expiry_date, quantity, status
  }>
}
```

### Парсинг

- `temp_min`/`temp_max` — если оба `null`, показывать «Без ограничений температуры»
- `humidity_max` — если `null`, показывать «Без ограничений влажности»
- Для зоны `safe_strong` — проверить `user.ns_pv_access`, если `false` — показать заглушку вместо остатков

### Zod для формы создания зоны (admin)

```ts
const zoneSchema = z.object({
  name: z.string().min(2),
  type: z.enum(["general", "cold_chain", "flammable", "safe_strong"]),
  description: z.string().optional(),
  temp_min: z.number().nullable(),
  temp_max: z.number().nullable(),
  humidity_max: z.number().min(0).max(100).nullable(),
})
```

---

## 🔐 Блок для администратора

- Кнопка `[+ Зона]` — видна только `admin`
- Диалог редактирования каждой карточки — только `admin`
- Карточка зоны `safe_strong` показывает остатки только `admin` и пользователям с `ns_pv_access: true`

---

## Ссылки на спецификацию

→ [next-form-ai.md — Раздел 3.2 Зонирование и Хранение](../next-form-ai.md#32-зонирование-и-хранение)
