# Форма 2.2: Профиль сотрудника

## Требования (пресловой-чеклист)

- [ ] Таблица `users` с полями: `full_name`, `email`, `role`, `ukep_bound`, `ns_pv_access`
- [ ] Таблица `employee_profiles`: `medical_book_scan_url`, `gdp_training_history (jsonb)`, `special_zone_access: bool`
- [ ] Бэкенд: `GET /api/v1/users/me` — данные текущего пользователя
- [ ] Бэкенд: `PUT /api/v1/users/me` — обновление профиля
- [ ] Бэкенд: `POST /api/v1/users/me/medical-book` — загрузка скана медкнижки (multipart)
- [ ] Бэкенд: `GET /api/v1/users/:id` — (admin) просмотр любого сотрудника
- [ ] Хранилище файлов (S3/MinIO) для сканов медкнижек
- [ ] shadcn/ui доп. компоненты: `npx shadcn@latest add avatar tabs toast`

---

## Интерфейс

### Лейаут страницы

- Страница `/profile` — доступна всем авторизованным пользователям
- Страница `/admin/users/:id` — только для `admin`
- Разбита на вкладки через `<Tabs>`

### Компоненты shadcn/ui

| Компонент | Назначение |
|-----------|-----------|
| `<Avatar>`, `<AvatarImage>`, `<AvatarFallback>` | Аватар сотрудника |
| `<Tabs>`, `<TabsList>`, `<TabsTrigger>`, `<TabsContent>` | Разделение секций |
| `<Card>` | Контейнеры секций |
| `<Form>`, `<FormField>`, etc. | Редактируемые поля |
| `<Input>` | ФИО, email |
| `<Badge>` | Роль, статус допусков |
| `<Button>` | Сохранение, загрузка файла |
| `<Toast>` | Уведомление об успешном сохранении |
| `<Separator>` | Разделители секций |

---

## Логика

### Данные с сервера

```
GET /api/v1/users/me
Response: {
  id: string
  full_name: string
  email: string
  role: RoleEnum
  ukep_bound: boolean
  ns_pv_access: boolean
  profile: {
    medical_book_scan_url: string | null
    special_zone_access: boolean
    gdp_training_history: Array<{
      date: string        // ISO 8601
      course_name: string
      result: "pass" | "fail"
      certificate_url: string | null
    }>
  }
}
```

### Парсинг

- `gdp_training_history` — отображается как `<Table>` с колонками: Дата, Курс, Результат, Сертификат
- `medical_book_scan_url` — если `null`, показать кнопку «Загрузить скан»; иначе — «Просмотреть» (открывает в новой вкладке) + «Обновить»
- `role` — маппинг на читаемые названия:

```ts
const roleLabels: Record<string, string> = {
  admin: "Администратор",
  qp: "Уполномоченное лицо (QP)",
  warehouse_manager: "Заведующий складом",
  storekeeper: "Кладовщик",
  pharmacist: "Фармацевт",
}
```

### Загрузка медкнижки

```ts
const onUpload = async (file: File) => {
  const formData = new FormData()
  formData.append("file", file)
  await fetch("/api/v1/users/me/medical-book", {
    method: "POST",
    body: formData,
  })
}
```

### Readonly-поля

Поля `full_name`, `email`, `role` — **только для чтения** (редактируются только через `admin`-форму).

---

## 🔐 Блок для администратора

На странице `/admin/users/:id` `admin` видит дополнительно:

| Элемент | Действие |
|---------|---------|
| Переключатель **«Допуск к НС/ПВ»** | `PATCH /api/v1/admin/users/:id` `{ ns_pv_access: bool }` |
| Переключатель **«Допуск к спецзонам»** | то же поле `special_zone_access` |
| Выбор **роли** (`<Select>`) | `PATCH /api/v1/admin/users/:id` `{ role: RoleEnum }` |
| Кнопка **«Выслать новый код для входа»** | `POST /api/v1/admin/users/:id/send-login-link` |

Переключатели реализуются через shadcn/ui `<Switch>` (установить: `npx shadcn@latest add switch`).

---

## Ссылки на спецификацию

→ [next-form-ai.md — Раздел 2.2 Профиль сотрудника](../next-form-ai.md#22-профиль-сотрудника)
