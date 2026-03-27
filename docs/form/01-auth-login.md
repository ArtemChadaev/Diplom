# Форма 2.1: Авторизация

## Требования (пресловой-чеклист)

- [ ] Таблица `users`, поля: `email`, `role`, `ns_pv_access`, `ukep_bound`
- [ ] Valkey: OTP-код в `otp:user:<id>` (Hash), TTL 600 сек — см. [valkey-cache.md](../valkey-cache.md)
- [ ] Бэкенд: `POST /api/v1/auth/send-code` — отправка 6-значного кода на email
- [ ] Бэкенд: `POST /api/v1/auth/verify-code` — проверка кода, выдача токенов
- [ ] Бэкенд: `GET /api/v1/auth/google` — OAuth-редирект
- [ ] Бэкенд: `GET /api/v1/auth/google/callback` — получение токена после Google
- [ ] Коды: 6 символов (цифры), TTL 10 минут, одноразовые
- [ ] Ограничение: не более 5 попыток отправки кода в час (rate limit)
- [ ] shadcn/ui: уже есть `Input`, `Button`, `Card`, `Label`; добавить: `npx shadcn@latest add form alert`

---

## Интерфейс

### Компоненты shadcn/ui

| Компонент | Назначение |
|-----------|-----------|
| `<Card>`, `<CardHeader>`, `<CardContent>` | Контейнер формы |
| `<Form>` | Обёртка (react-hook-form + zod) |
| `<FormField>`, `<FormItem>`, `<FormLabel>`, `<FormControl>`, `<FormMessage>` | Email-поле |
| `<Input>` | Email-адрес |
| `<Button>` | «Получить код», «Войти через Google» |
| `<Alert>`, `<AlertDescription>` | Ошибки (неверный email, rate limit) |
| `<Separator>` | Разделитель между OAuth и email |

---

## Логика

### Zod-схема

```ts
const emailSchema = z.object({
  email: z.string().email("Некорректный email"),
})
```

### API-запросы

**Шаг 1 — Запрос кода:**
```
POST /api/v1/auth/send-code
Body: { email: string }
Response 200: { message: "Code sent", expires_in: 600 }
Response 429: { message: "Too many requests" }
Response 404: { message: "User not found" }
```

**Google OAuth:**
```
GET /api/v1/auth/google  → редирект на Google
GET /api/v1/auth/google/callback → { access_token, user: UserDTO }
                                    + Set-Cookie: refresh_token=...; HttpOnly; Secure
```

### UserDTO

```ts
type UserDTO = {
  id: string
  email: string
  full_name: string
  role: "admin" | "qp" | "warehouse_manager" | "storekeeper" | "pharmacist"
  ns_pv_access: boolean
  ukep_bound: boolean
}
```

### Редиректы после успешного входа

| Роль | Редирект |
|------|---------|
| `admin` | `/admin/users` |
| `qp` | `/receiving` |
| `warehouse_manager` | `/warehouse` |
| `storekeeper`, `pharmacist` | `/search` |

### Обработка ошибок

- `404` → «Пользователь с таким email не найден»
- `429` → «Слишком много запросов. Повторите через X минут»

---

## 🔐 Блок для администратора

Нет специфичных admin-элементов на этой форме.  
После входа с ролью `admin` — редирект на `/admin`.

---

## Ссылки на спецификацию

→ [next-form-ai.md — Раздел 2.1 Авторизация](../next-form-ai.md#21-auth--login)
