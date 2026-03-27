# Форма 5.1: Карточка товара (Медикамента)

## Требования (пресловой-чеклист)

- [ ] Таблица `products`: `id`, `mnn`, `trade_name`, `sku`, `barcode`, `datamatrix_gtin`, `ru_number`, `atc_code`, `dosage_form`, `dosage`, `package_multiplicity`, `weight_g`, `width_cm`, `height_cm`, `depth_cm`
- [ ] Поля-флаги: `is_jnvlp: bool`, `is_mdlp: bool`, `is_ns_pv: bool`, `cold_chain: bool`
- [ ] Поля хранения: `temp_min`, `temp_max`, `humidity_max`
- [ ] Таблица `product_photos`: `id`, `product_id`, `url`, `is_primary: bool`
- [ ] Бэкенд: `GET /api/v1/products` — список товаров (поиск + пагинация)
- [ ] Бэкенд: `GET /api/v1/products/:id` — карточка товара
- [ ] Бэкенд: `POST /api/v1/products` — (admin) создание товара
- [ ] Бэкенд: `PUT /api/v1/products/:id` — (admin) редактирование
- [ ] Бэкенд: `DELETE /api/v1/products/:id` — (admin) удаление (мягкое: `deleted_at`)
- [ ] Бэкенд: `POST /api/v1/products/:id/photos` — загрузка фото
- [ ] Справочник ATC-классификации (`GET /api/v1/ref/atc`)
- [ ] shadcn/ui: `npx shadcn@latest add tooltip separator`

---

## Интерфейс

### Лейаут страницы

- `/products` — поисковая таблица товаров (доступна всем)
- `/products/:id` — карточка товара (доступна всем)
- `/products/new` — (admin) создание нового товара
- `/products/:id/edit` — (admin) редактирование

### Компоненты shadcn/ui

| Компонент | Назначение |
|-----------|-----------|
| `<Card>` | Основной контейнер карточки |
| `<Badge>` | Флаги (ЖНВЛП, МДЛП, НС/ПВ, ❄️ Холодовая цепь) |
| `<Tabs>` | Разделы: Основные / Хранение / Остатки / История |
| `<Table>` | Текущие серии с остатками |
| `<Tooltip>` | Подсказки к флагам и кодам |
| `<Separator>` | Разделители секций |
| `<Form>`, `<Input>`, `<Select>`, `<Textarea>` | Форма создания/редактирования |
| `<Button>` | Редактировать, Загрузить фото |
| `<MultiSelect>` | Выбор нескольких ATC-кодов (уже есть в проекте) |

### Структура UI — Карточка

```
┌──────────────────────────────────────────────────────────┐
│  [Фото]   Амоксициллин 500мг                             │
│           [ЖНВЛП] [МДЛП] [❄️ Холодовая цепь]           │
│           МНН: Амоксициллин                              │
│           РУ: ЛП-002712  |  SKU: AM500CAP               │
│           ATC: J01CA04                                   │
├──────────────────────────────────────────────────────────┤
│  [Основные] [Условия хранения] [Остатки] [История]       │
├──────────────────────────────────────────────────────────┤
│  Вкладка "Основные":                                     │
│  Форма выпуска: Капсулы 500 мг                          │
│  Кратность упаковки: 16 штук                            │
│  Штрих-код: 4607053860207                               │
│  GTIN (DataMatrix): 04607053860207                      │
│  ВГХ: 120×80×30 мм | 45 г                              │
├──────────────────────────────────────────────────────────┤
│  Вкладка "Условия хранения":                             │
│  Температура: 2–8°C                                      │
│  Влажность: не более 60%                                 │
├──────────────────────────────────────────────────────────┤
│  Вкладка "Остатки" (текущие серии):                      │
│  │ Серия │ Партия │ Годен до │ Кол-во │ Зона │ Статус │  │
└──────────────────────────────────────────────────────────┘
```

### Форма создания/редактирования (admin)

```
Торговое название: [Input]
МНН:               [Input]
SKU:               [Input] (авто-генерация или ручной ввод)
Штрих-код:         [Input]
GTIN:              [Input]
№ РУ:              [Input] ← Обязательно
ATC-код:           [MultiSelect — из справочника]
Форма выпуска:     [Select: Таблетки / Капсулы / Раствор / ...]
Дозировка:         [Input]
Кратность уп.:     [Input type=number]

Флаги:
  [Switch] ЖНВЛП
  [Switch] Подлежит МДЛП (маркировка)
  [Switch] НС/ПВ (ПКУ)
  [Switch] Холодовая цепь

Условия хранения:
  Температура от: [Input] °C
  Температура до: [Input] °C
  Влажность макс: [Input] %

ВГХ:
  Вес: [Input] г
  Ш × В × Г: [Input] × [Input] × [Input] мм

Фото:
  [Загрузить основное фото] ← drag-and-drop
```

---

## Логика

### Zod-схема

```ts
const productSchema = z.object({
  trade_name: z.string().min(2),
  mnn: z.string().min(2),
  sku: z.string().optional(),
  barcode: z.string().optional(),
  datamatrix_gtin: z.string().optional(),
  ru_number: z.string().min(1, "РУ обязателен"),
  atc_codes: z.array(z.string()).min(1, "Минимум один ATC-код"),
  dosage_form: z.string().min(1),
  dosage: z.string().min(1),
  package_multiplicity: z.number().int().min(1),
  is_jnvlp: z.boolean().default(false),
  is_mdlp: z.boolean().default(false),
  is_ns_pv: z.boolean().default(false),
  cold_chain: z.boolean().default(false),
  temp_min: z.number().nullable(),
  temp_max: z.number().nullable(),
  humidity_max: z.number().min(0).max(100).nullable(),
  weight_g: z.number().min(0).optional(),
  width_cm: z.number().min(0).optional(),
  height_cm: z.number().min(0).optional(),
  depth_cm: z.number().min(0).optional(),
})
.refine(
  (d) => !d.cold_chain || (d.temp_min !== null && d.temp_max !== null),
  { message: "Для товара с холодовой цепью укажите диапазон температур" }
)
```

### Данные с сервера

```
GET /api/v1/products/:id
Response: {
  id, trade_name, mnn, sku, barcode, datamatrix_gtin,
  ru_number, atc_codes: string[],
  dosage_form, dosage, package_multiplicity,
  is_jnvlp, is_mdlp, is_ns_pv, cold_chain,
  temp_min, temp_max, humidity_max,
  weight_g, width_cm, height_cm, depth_cm,
  photos: Array<{ id, url, is_primary }>,
  batches: Array<{
    id, serial_number, manufacture_date, expiry_date,
    quantity, zone_name, status
  }>
}
```

### Парсинг флагов для отображения

```ts
const flags = [
  { key: "is_jnvlp", label: "ЖНВЛП", variant: "outline" },
  { key: "is_mdlp", label: "МДЛП", variant: "secondary" },
  { key: "is_ns_pv", label: "НС/ПВ", variant: "destructive" },
  { key: "cold_chain", label: "❄️ Холодовая цепь", variant: "blue" },
]
// Отображать только флаги с true
product.flags = flags.filter(f => product[f.key])
```

### Поиск товаров

```
GET /api/v1/products?q={query}&is_jnvlp=true&page=1&limit=20
```

Поиск поддерживает МНН, торговое название, SKU, штрих-код, РУ, DataMatrix (GTIN).

---

## 🔐 Блок для администратора

- Кнопка «Редактировать» — только `admin`
- Кнопка «Создать товар» — только `admin`
- Карточки товаров с `is_ns_pv = true` видят **все пользователи**, но вкладку «Остатки» для таких товаров — только `admin` и пользователи с `ns_pv_access: true`
- При `is_ns_pv = true` на форме редактирования появляется доп. секция «ПКУ-требования» (подучётно-количественные) — только `admin`

---

## Ссылки на спецификацию

→ [next-form-ai.md — Раздел 5.1 Карточка товара (Медикамента)](../next-form-ai.md#51-карточка-товара-медикамента)
