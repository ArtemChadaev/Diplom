# Форма 3.1: Приход товара (Inbound)

## Требования (пресловой-чеклист)

- [ ] Таблица `suppliers` (поставщики) с полями: `id`, `name`, `inn`, `license_number`
- [ ] Таблица `products` (товары/медикаменты) с полями: `id`, `mnn`, `ru_number`, `sku`, `atc_code`, `is_jnvlp`, `is_mdlp`, `is_ns_pv`, `cold_chain`
- [ ] Таблица `batches` (серии): `id`, `product_id`, `serial_number`, `manufacture_date`, `expiry_date`, `quantity`, `status: enum(quarantine, available, rejected)`
- [ ] Таблица `inbound_receipts` (накладные): `id`, `supplier_id`, `invoice_number`, `country_of_origin`, `manufacturer`, `vat_rate`, `markup_jnvlp`, `created_by`, `created_at`
- [ ] Справочник стран (`GET /api/v1/ref/countries`)
- [ ] Справочник поставщиков (`GET /api/v1/suppliers`)
- [ ] Справочник товаров (`GET /api/v1/products/search?q=...`)
- [ ] Бэкенд: `POST /api/v1/inbound` — создание прихода
- [ ] Бэкенд: `GET /api/v1/inbound/:id` — просмотр прихода
- [ ] Бэкенд: `POST /api/v1/inbound/:id/quarantine-release` — (admin/QP) выпуск из карантина
- [ ] shadcn/ui доп. компоненты: `npx shadcn@latest add form select combobox date-picker textarea toast`
- [ ] `<DatePicker>` уже есть в проекте `components/ui/date-picker.tsx`
- [ ] `<MultiSelect>` уже есть в проекте `components/ui/multi-select.tsx`

---

## Интерфейс

### Лейаут страницы

- Страница `/receiving/new` — мастер-форма в несколько шагов (Stepper)
- Степпер: **Шаг 1 → Шаг 2 → Шаг 3 → Подтверждение**
- Адаптивная: на мобильных — вертикальный вид

### Компоненты shadcn/ui

| Компонент | Назначение |
|-----------|-----------|
| `<Form>` | Корневая обёртка (react-hook-form) |
| `<Select>`, `<SelectContent>`, `<SelectItem>` | Поставщик, Тип закупки, Страна, НДС |
| `<Input>` | № Накладной, Завод-изготовитель |
| `<Combobox>` | Поиск товара по МНН/SKU (autocomplete) |
| `<DatePicker>` | Дата производства, Дата истечения срока |
| `<Textarea>` | Примечания |
| `<Table>` | Список позиций в накладной |
| `<Badge>` | Статус серии (Карантин / Доступен) |
| `<Button>` | Следующий шаг, Сохранить, Добавить позицию |
| `<Dialog>` | Подтверждение выпуска из карантина |
| `<Toast>` | Уведомления |
| `<Alert>` | Предупреждение о карантине |

### Структура UI (по шагам)

**Шаг 1 — Реквизиты накладной:**
```
Поставщик:         [Combobox — поиск по названию/ИНН]
Тип закупки:       [Select: Прямая / Тендер / Гос.закупка]
№ Накладной:       [Input]
Страна происх.:    [Select — справочник стран]
Завод-изготовит.:  [Input]
```

**Шаг 2 — Финансы:**
```
Ставка НДС:        [Select: 0% / 10% / 20%]
Контроль наценки ЖНВЛП: [Switch — активен если товар is_jnvlp]
Наценка (%):       [Input type=number — если ЖНВЛП включён]
```

**Шаг 3 — Позиции накладной (повторяющийся блок):**
```
Товар (МНН/SKU):   [Combobox с поиском]
№ Серии:           [Input]
Дата производства: [DatePicker]
Срок годности:     [DatePicker] ← КРИТИЧНО
Количество:        [Input type=number]
[+ Добавить позицию]

─────────────────────────────────────────
Таблица позиций:
│ МНН │ Серия │ Произв. │ Годен до │ Кол-во │ Удалить │
```

**Шаг 4 — Подтверждение:**
```
⚠️  Товар будет помещён в КАРАНТИН до подписания протокола приёмки.
[Итоговый просмотр всех данных]
[Сохранить и отправить в карантин]
```

---

## Логика

### Zod-схема

```ts
const inboundSchema = z.object({
  supplier_id: z.string().min(1),
  purchase_type: z.enum(["direct", "tender", "state"]),
  invoice_number: z.string().min(1),
  country_of_origin: z.string().min(1),
  manufacturer: z.string().min(1),
  vat_rate: z.enum(["0", "10", "20"]),
  is_jnvlp_controlled: z.boolean(),
  jnvlp_markup: z.number().min(0).max(100).optional(),
  positions: z.array(z.object({
    product_id: z.string().min(1),
    serial_number: z.string().min(1),
    manufacture_date: z.date(),
    expiry_date: z.date(),
    quantity: z.number().int().min(1),
  })).min(1, "Добавьте хотя бы одну позицию"),
})
  .refine((d) => d.positions.every(p => p.expiry_date > p.manufacture_date), {
    message: "Дата истечения должна быть позже даты производства",
    path: ["positions"],
  })
```

### Данные с сервера

```
GET /api/v1/suppliers?q={query}
Response: Array<{ id, name, inn }>

GET /api/v1/ref/countries
Response: Array<{ code: string, name_ru: string }>

GET /api/v1/products/search?q={mnn_or_sku}
Response: Array<{
  id, mnn, sku, ru_number,
  is_jnvlp: boolean,
  is_mdlp: boolean,
  is_ns_pv: boolean,
  cold_chain: boolean,
  temp_min: number, temp_max: number
}>
```

### Парсинг и логика после выбора товара

При выборе товара из Combobox:
1. Если `is_jnvlp = true` → автоматически включить `is_jnvlp_controlled` Switch и показать поле наценки
2. Если `is_ns_pv = true` → показать предупреждение `<Alert>` «Товар относится к НС/ПВ. Убедитесь в наличии специальных разрешений»
3. Если `cold_chain = true` → показать `<Badge variant="outline">❄️ Холодовая цепь</Badge>` рядом с позицией

### Статус после сохранения

После `POST /api/v1/inbound` все серии получают `status: "quarantine"` — они **заблокированы** для любых перемещений.  
В списке прихода показывается `<Badge variant="warning">В карантине</Badge>`.

---

## 🔐 Блок для администратора (admin / QP)

### Протокол приёмки и выпуск из карантина

Блок видим **только** пользователям с ролью `admin` или `qp`:

```
┌──────────────────────────────────────────┐
│  📋 Протокол приёмки                    │
│                                          │
│  Уполномоченное лицо: [Автозаполнение]  │
│  Дата проверки: [DatePicker]            │
│  Результат: [✓ Соответствует / ✗ Брак] │
│  Замечания: [Textarea]                  │
│                                          │
│  [🔐 Подписать ЭЦП и выпустить]         │
└──────────────────────────────────────────┘
```

Кнопка «Подписать ЭЦП и выпустить» → `POST /api/v1/inbound/:id/quarantine-release`  
Тело: `{ qp_user_id, inspection_date, result: "approved"|"rejected", notes }`  
После успеха: серии переходят в `status: "available"`.

**Если `result = "rejected"`** → серии переходят в `status: "rejected"` и появляется кнопка «Оформить возврат поставщику».

---

## Ссылки на спецификацию

→ [next-form-ai.md — Раздел 3.1 Форма прихода (Inbound)](../next-form-ai.md#31-форма-прихода-inbound)
