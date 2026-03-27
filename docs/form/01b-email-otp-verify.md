# Форма 2.1б: Подтверждение email (OTP-код)

## Требования (пресловой-чеклист)

- [ ] Зависит от формы 01-авторизация: пользователь уже ввёл email и нажал «Получить код»
- [ ] Valkey: ключ `otp:user:<user_id>` (Hash), TTL 600 сек — см. [valkey-cache.md](../valkey-cache.md#1-otp-коды-otp)
- [ ] Бэкенд: `POST /api/v1/auth/verify-code` — проверка 6-значного кода
- [ ] Бэкенд: `POST /api/v1/auth/send-code` — повторная отправка кода (из той же формы)
- [ ] Код: 6 цифр, TTL 10 минут, одноразовый; после использования — `DEL otp:user:<id>` в Valkey
- [ ] Максимум 3 попытки ввода неверного кода (поле `attempts` в Hash) → `DEL`, кнопка «Выслать новый»
- [ ] Таймер обратного отсчёта (10:00 → 0:00) до истечения кода
- [ ] shadcn/ui доп. компоненты: `npx shadcn@latest add input-otp`

---

## Интерфейс

### Компоненты shadcn/ui

| Компонент | Назначение |
|-----------|-----------|
| `<InputOTP>`, `<InputOTPGroup>`, `<InputOTPSlot>` | 6 отдельных ячеек для кода |
| `<Form>`, `<FormField>` | Обёртка с валидацией |
| `<Button>` | «Подтвердить», «Выслать код повторно» |
| `<Alert>` | Ошибки (неверный код, истёк срок) |

---

## Логика

### Zod-схема

```ts
const otpSchema = z.object({
  code: z.string().length(6, "Код должен содержать 6 цифр").regex(/^\d{6}$/, "Только цифры"),
})
```

### API-запросы

**Проверка кода:**
```
POST /api/v1/auth/verify-code
Body: { email: string, code: string }

Response 200: {
  access_token: string,
  refresh_token: string,
  user: UserDTO
}
Response 400: { message: "Invalid code", attempts_left: number }
Response 410: { message: "Code expired" }
Response 429: { message: "Max attempts reached" }
```

**Повторная отправка кода:**
```
POST /api/v1/auth/send-code
Body: { email: string }
```

### Поведение таймера

```ts
// Таймер запускается при загрузке страницы/компонента
// expires_at приходит из предыдущего ответа send-code
const [secondsLeft, setSecondsLeft] = useState(
  Math.max(0, Math.floor((expiresAt - Date.now()) / 1000))
)

useEffect(() => {
  const interval = setInterval(() => {
    setSecondsLeft(s => Math.max(0, s - 1))
  }, 1000)
  return () => clearInterval(interval)
}, [])
```

Кнопка «Выслать повторно» — задизейблена пока `secondsLeft > 0`.  
После `secondsLeft === 0` → показать текст «Код истёк» и активировать кнопку.

### Авто-сабмит

При заполнении всех 6 ячеек — форма отправляется **автоматически** без нажатия кнопки:

```ts
<InputOTP
  maxLength={6}
  onComplete={(code) => form.handleSubmit(onSubmit)()}
>
```

### Парсинг ответа

- `200` → сохранить токены → редирект по роли (см. [01-авторизация.md](./01-авторизация.md))
- `400` `attempts_left` → показать `<Alert>` «Неверный код. Осталось попыток: N»
- `429` → показать «Превышено число попыток. Запросите новый код» + активировать кнопку «Выслать повторно»
- `410` → показать «Код истёк. Запросите новый»

### Передача email между формами

Email передаётся через URL `?email=...` или через состояние роутера (Next.js `router.push` со `state`):

```ts
// На форме авторизации после успешного send-code:
router.push(`/auth/verify?email=${encodeURIComponent(email)}`)
```

---

## 🔐 Блок для администратора

Нет специфичных admin-элементов.

---

## Ссылки на спецификацию

→ [next-form-ai.md — Раздел 2.1 Авторизация](../next-form-ai.md#21-auth--login)
