# Форма 4.3: Рекламации и Брак

## Требования (пресловой-чеклист)

- [ ] Таблица `claims`: `id`, `type: enum(recall, return_from_pharmacy, return_to_supplier, defect)`, `batch_id`, `product_id`, `status: enum(open, blocked, closed)`, `source`, `notes`, `created_by`, `created_at`
- [ ] Таблица `claim_photos`: `id`, `claim_id`, `url`, `uploaded_at`
- [ ] Интеграция с Росздравнадзором: таблица `recalled_batches` (синхронизируется фоновым воркером)
- [ ] Бэкенд: `GET /api/v1/claims` — список рекламаций
- [ ] Бэкенд: `POST /api/v1/claims` — создание рекламации
- [ ] Бэкенд: `GET /api/v1/claims/:id` — детали
- [ ] Бэкенд: `POST /api/v1/claims/:id/photos` — загрузка фото (multipart)
- [ ] Бэкенд: `POST /api/v1/claims/:id/close` — (admin/qp) закрытие рекламации
- [ ] Бэкенд: `GET /api/v1/recalled-batches/sync` — (admin) ручной запуск синхронизации с Росздравнадзором
- [ ] STOP-сигнал: при импорте изъятой серии `batch.status → "blocked"` — автоматически
- [ ] Хранилище файлов (S3/MinIO) для фотографий брака
- [ ] shadcn/ui: `npx shadcn@latest add alert-dialog`

---

## Интерфейс

### Лейаут страницы

- `/claims` — список рекламаций
- `/claims/new` — форма создания рекламации
- `/claims/:id` — детальный просмотр

### Компоненты shadcn/ui

| Компонент | Назначение |
|-----------|-----------|
| `<Table>` | Список рекламаций |
| `<Badge>` | Тип и статус рекламации |
| `<Form>`, `<FormField>` | Форма создания |
| `<Select>` | Тип рекламации |
| `<Combobox>` | Поиск товара/серии |
| `<Input>` | Источник, описание |
| `<Textarea>` | Примечания |
| `<Alert variant="destructive">` | STOP-сигнал, блокировка |
| `<AlertDialog>` | Подтверждение закрытия рекламации |
| `<Card>` | Блок загруженных фото |
| `<Button>` | Загрузить фото, Сохранить, Закрыть |

### Структура UI — Список

```
┌──────────────────────────────────────────────────────┐
│  Рекламации и Брак                     [+ Новая]    │
│  [Все] [Открытые] [Заблокированные] [Закрытые]     │
├────┬─────────────┬──────────┬──────────┬────────────┤
│ №  │  Товар      │  Серия   │  Тип     │ Статус     │
├────┼─────────────┼──────────┼──────────┼────────────┤
│ 12 │ Аспирин... │ RZ2024A  │ Изъятие  │🔴 ЗАБЛОК. │
│ 11 │ Омепразол  │ PK2025C  │ Возврат  │🟡 Открыта  │
└────┴─────────────┴──────────┴──────────┴────────────┘
```

### STOP-сигнал (глобальный алерт)

При совпадении серии с реестром Росздравнадзора — **во вверху любой страницы** показывается фиксированный `<Alert>`:
```
🛑 STOP-СИГНАЛ: Серия RZ2024A (Аспирин 500мг) включена в 
   реестр изъятых Росздравнадзором. Остатки заблокированы.
   [Перейти к рекламации]
```

### Структура UI — Форма создания рекламации

```
Тип рекламации:   [Select: Изъятие Росздр. / Возврат из аптеки /
                           Возврат поставщику / Брак]
Товар:            [Combobox поиск по МНН/SKU]
Серия:            [Combobox — только серии выбранного товара]
Источник:         [Input — аптека-поставщик / № предписания РЗН]
Описание дефекта: [Textarea]
Фотофиксация:     [Зона drag-and-drop или Input type=file]
                  [Превью загруженных фото — grid картинок]
```

---

## Логика

### Zod-схема

```ts
const claimSchema = z.object({
  type: z.enum(["recall", "return_from_pharmacy", "return_to_supplier", "defect"]),
  product_id: z.string().min(1, "Выберите товар"),
  batch_id: z.string().min(1, "Выберите серию"),
  source: z.string().optional(),
  notes: z.string().optional(),
  photos: z.array(z.instanceof(File)).optional(),
})
```

### Данные с сервера

```
GET /api/v1/claims?status=open&page=1&limit=20
Response: Array<{
  id, type, status, created_at,
  product: { id, mnn, sku },
  batch: { id, serial_number, expiry_date },
  source: string | null,
  photos_count: number,
  created_by: { full_name }
}>
```

### Парсинг статусов

```ts
const statusConfig = {
  open: { label: "Открыта", variant: "warning" },
  blocked: { label: "ЗАБЛОКИРОВАНО", variant: "destructive" },
  closed: { label: "Закрыта", variant: "secondary" },
}
```

### Автоматическая блокировка (STOP-сигнал)

Фоновый воркер (Go) периодически синхронизирует `recalled_batches`.  
На клиенте — WebSocket или polling `GET /api/v1/stop-signals` проверяет наличие активных STOP-сигналов.  
При наличии `active: true` — рендерится глобальный `<Alert>` вне основного лейаута.

### Загрузка фото

```ts
const uploadPhotos = async (claimId: string, files: File[]) => {
  const formData = new FormData()
  files.forEach(f => formData.append("photos", f))
  await fetch(`/api/v1/claims/${claimId}/photos`, { method: "POST", body: formData })
}
```

---

## 🔐 Блок для администратора

| Действие | Кто |
|---------|-----|
| Закрыть рекламацию (`POST /api/v1/claims/:id/close`) | `admin`, `qp` |
| Ручная синхронизация с Росздравнадзором | `admin` — кнопка «Синхронизировать» в `/admin/settings` |
| Снять блокировку серии (`PATCH /api/v1/batches/:id/unblock`) | `admin` с обязательным указанием причины |

Кнопка «Снять блокировку» показывается только `admin` и только для серий со статусом `blocked`, у которых источник — **не Росздравнадзор** (например, ошибочная блокировка).

---

## Ссылки на спецификацию

→ [next-form-ai.md — Раздел 4.3 Рекламации и Брак](../next-form-ai.md#43-рекламации-и-брак)
