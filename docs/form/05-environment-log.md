# Форма 3.2б: Журнал параметров среды

## Требования (пресловой-чеклист)

- [ ] Таблица `environment_logs`: `id`, `zone_id`, `recorded_by`, `recorded_at`, `temperature`, `humidity`, `shift: enum(morning, evening)`, `notes`
- [ ] Бэкенд: `POST /api/v1/zones/:id/environment-log` — создание записи
- [ ] Бэкенд: `GET /api/v1/zones/:id/environment-log` — история записей (с пагинацией и фильтром по дате)
- [ ] Бэкенд: `GET /api/v1/environment-log/today` — все записи за сегодня (агрегация по зонам)
- [ ] Правило: **2 записи в день** — утренняя (`morning`) и вечерняя (`evening`) для каждой зоны
- [ ] Алерт: если запись за текущую смену отсутствует → уведомление в Header-центре уведомлений
- [ ] Зоны с `temp_min`/`temp_max` — валидация допустимого диапазона на сервере и клиенте
- [ ] shadcn/ui: `npx shadcn@latest add form input select textarea toast alert`

---

## Интерфейс

### Лейаут страницы

- Страница `/warehouse/environment-log` — доступна `warehouse_manager`, `storekeeper`, `qp`, `admin`
- Состоит из двух частей:
  1. **Панель быстрого ввода** — карточки зон с кнопкой «Внести данные»
  2. **История** — таблица прошлых записей с фильтром по дате/зоне

### Компоненты shadcn/ui

| Компонент | Назначение |
|-----------|-----------|
| `<Card>` | Карточка зоны в панели ввода |
| `<Badge>` | Статус записи (✓ внесена / ⚠ ожидает) |
| `<Dialog>` | Форма ввода данных |
| `<Form>`, `<FormField>` | Поля ввода |
| `<Input type="number">` | Температура, влажность |
| `<Select>` | Выбор смены (утренняя/вечерняя) |
| `<Textarea>` | Примечания |
| `<Table>` | История записей |
| `<Alert>` | Предупреждение о выходе значений за пределы нормы |
| `<Toast>` | Подтверждение сохранения |

### Структура UI

**Панель быстрого ввода:**
```
┌──────────────────────────────────────────────────┐
│  Журнал температуры и влажности — Сегодня        │
│  27.03.2026 | Смена: [Утренняя ▼]               │
├──────────┬──────────┬──────────┬────────────────┤
│ 📦 Общая │ ❄️ Холодов│ 🔥 Огнеоп│  🔒 НС/ПВ     │
│  ✓ Внес. │  ⚠ Ожид. │  ✓ Внес. │  ⚠ Ожид.     │
│ [Просм.] │ [Внести] │ [Просм.] │  [Внести]     │
└──────────┴──────────┴──────────┴────────────────┘
```

**Диалог ввода данных:**
```
┌──────────────────────────────────────┐
│  ❄️ Холодовая цепь — Утренняя смена  │
│  Норма: 2–8°C | Влажность ≤60%      │
├──────────────────────────────────────┤
│  Температура (°C): [____] ← input   │
│  Влажность (%):    [____] ← input   │
│  Примечания:       [_____________]  │
│                                      │
│  ⚠ 9.2°C — ВЫШЕ НОРМЫ! (2–8°C)    │
│                                      │
│  [Отмена]          [Сохранить]      │
└──────────────────────────────────────┘
```

**История (таблица):**
```
│ Зона │ Дата │ Смена │ Темп. │ Влаж. │ Сотрудник │ Примечания │
```

---

## Логика

### Zod-схема ввода

```ts
const envLogSchema = (zone: ZoneDTO) => z.object({
  shift: z.enum(["morning", "evening"]),
  temperature: z.number()
    .min(-50).max(100)
    .refine(
      (v) => zone.temp_min == null || v >= zone.temp_min,
      `Ниже минимума (${zone.temp_min}°C)`
    )
    .refine(
      (v) => zone.temp_max == null || v <= zone.temp_max,
      `Выше максимума (${zone.temp_max}°C)`
    ),
  humidity: z.number().min(0).max(100)
    .refine(
      (v) => zone.humidity_max == null || v <= zone.humidity_max,
      `Влажность выше нормы (${zone.humidity_max}%)`
    ),
  notes: z.string().optional(),
})
```

### Данные с сервера

```
GET /api/v1/environment-log/today
Response: Array<{
  zone_id: string
  zone_name: string
  zone_type: ZoneTypeEnum
  temp_min: number | null
  temp_max: number | null
  humidity_max: number | null
  morning_log: EnvLogDTO | null    // null = не внесена
  evening_log: EnvLogDTO | null
}>

// EnvLogDTO
{
  id, temperature, humidity, shift,
  recorded_by: { id, full_name },
  recorded_at: string  // ISO 8601
  notes: string | null
}
```

### Парсинг

- Если `morning_log = null` и текущее время < 14:00 → Badge «Ожидает ввода» (yellow)
- Если значение выходит за пределы `temp_min/temp_max` или `humidity_max` → `<Alert variant="destructive">` в форме
- После успешного сохранения → `<Toast>` «Данные сохранены. Сотрудник: [ФИО]»

### История

```
GET /api/v1/zones/:id/environment-log?from=2026-01-01&to=2026-03-31&page=1&limit=50
```

Данные таблицы отображаются с подсветкой строк где значение вышло за пределы нормы (`text-destructive`).

---

## 🔐 Блок для администратора

- `admin` и `qp` могут **редактировать** уже внесённую запись (кнопка «Исправить»)
- `admin` видит журнал **всех зон** включая `safe_strong` (НС/ПВ)
- Экспорт истории журнала в Excel (`GET /api/v1/environment-log/export?from=...&to=...`) — кнопка «Выгрузить» видна только `admin`/`qp`

---

## Ссылки на спецификацию

→ [next-form-ai.md — Раздел 3.2 Журнал параметров среды](../next-form-ai.md#32-зонирование-и-хранение)
