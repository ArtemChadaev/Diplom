# Архитектура системы ERP (Фармацевтика)

Данный документ описывает программную архитектуру, взаимосвязь модулей и структуру данных системы.

---

## 1. Дерево функций (Functional Tree)

Система разделена на административный, складской и общепользовательский функционал.

*   **Модуль «Аутентификация и Доступы»**
    *   Вход через Google OAuth 2.0
    *   Вход по Email (OTP 6-значный код)
    *   Управление сессиями (Refresh Tokens)
    *   Разграничение ролей (Admin, QP, Manager, Storekeeper, Pharmacist)
*   **Модуль «Персонал»**
    *   Управление профилями сотрудников (Медкнижки, допуски)
    *   Учет обучения GDP (Good Distribution Practice)
*   **Модуль «Каталог товаров»**
    *   Реестр лекарственных средств (МНН, АТХ-коды, ЖНВЛП)
    *   Учет весогабаритных характеристик и температурных режимов
    *   Управление фото-контентом
*   **Модуль «Складской учет»**
    *   **Приемка (Inbound):** Создание накладных, перемещение в карантин.
    *   **Контроль качества (QP):** Акцепт или отбраковка серий.
    *   **Зонирование:** Размещение по спец-зонам (Холод, НС/ПВ, Сейф).
    *   **Журнал среды:** Фиксация температуры и влажности по сменам.
*   **Модуль «Заказы и Отгрузка»**
    *   Создание заказов (Regular/CITO)
    *   Автоподбор серий по **FEFO**
    *   Сборка и генерация ТТН (Товарно-транспортная накладная)
*   **Модуль «Рекламации и Регуляторика»**
    *   Учет брака и возвратов (Рекламации)
    *   Мониторинг STOP-сигналов Росздравнадзора (Recalls)
*   **Модуль «Инвентаризация»**
    *   Создание сессий инвентаризации (слепой пересчет)
    *   Анализ расхождений и акты списания/оприходования

---

## 2. Диаграмма взаимосвязи модулей (Next.js ↔ Go)

Система использует современный стек: **Go (Backend)** и **Next.js (Frontend)**.

```mermaid
graph LR
    subgraph "Frontend (Next.js 16)"
        A[Client Components] -- "Interaction" --> B[React State/Zod]
        C[Server Components] -- "Initial Data" --> D[Next.js Fetch Cache]
    end

    subgraph "Communication"
        API[REST API / JSON]
        Auth[JWT + HttpOnly Cookie]
    end

    subgraph "Backend (Go 1.22+)"
        H[Chi Router / Middleware] --> S[Service Layer (Business Logic)]
        S --> R[Repository Layer (GORM/PostgreSQL)]
        S --> V[Valkey Cache (OTP/Rate Limit)]
    end

    A <--> API
    C <--> API
    B <--> Auth
```

### Ключевые принципы взаимодействия:
1.  **Stateless API:** Бэкенд не хранит состояние сессии в памяти, используя JWT и PostgreSQL/Valkey.
2.  **Onion Architecture (Go):** Четкое разделение на Domain, Service, Repository и Handler.
3.  **App Router (Next.js):** Использование серверных компонентов для ускорения рендеринга и SEO, и клиентских для интерактивных форм.
4.  **Shared DTO:** Структуры данных на Go (Backend) соответствуют TypeScript-интерфейсам (Frontend).

---

## 3. Схема данных и взаимосвязь классов (Entities)

Ниже представлена детальная структура основных сущностей системы и их атрибутов.

```mermaid
classDiagram
    class User {
        +int ID
        +string Email
        +string GoogleID
        +string TelegramID
        +user_role Role
        +bool NsPvAccess
        +bool UkepBound
        +bool IsBlocked
    }
    class EmployeeProfile {
        +int ID
        +int UserID
        +string FullName
        +string EmployeeCode
        +string Position
        +string Department
        +date HireDate
        +string MedicalBookUrl
        +json GDPHistory
    }
    class Product {
        +uuid ID
        +string TradeName
        +string mnn
        +string Sku
        +string RuNumber
        +bool IsJNVLP
        +bool IsNS_PV
        +bool ColdChain
        +float TempMin
        +float TempMax
    }
    class Batch {
        +uuid ID
        +uuid ProductID
        +uuid ZoneID
        +string SerialNumber
        +date ExpiryDate
        +int Quantity
        +batch_status Status
    }
    class InboundReceipt {
        +uuid ID
        +uuid SupplierID
        +string InvoiceNumber
        +purchase_type Type
        +string InspectionResult
        +int QpUserID
        +date InspectionDate
    }
    class Order {
        +uuid ID
        +order_type Type
        +order_status Status
        +uuid DestinationID
        +int CreatedBy
        +timestamp ShippedAt
    }
    class WarehouseZone {
        +uuid ID
        +string Name
        +zone_type Type
        +float TempMin
        +float TempMax
        +float HumidityMax
    }

    User "1" -- "1" EmployeeProfile : has
    User "1" -- "n" InboundReceipt : approves (QP)
    Product "1" -- "n" Batch : comprises
    WarehouseZone "1" -- "n" Batch : stores
    InboundReceipt "1" -- "n" InboundPosition : defines
    InboundPosition "n" -- "1" Batch : links
    Order "1" -- "n" OrderItem : contains
    OrderItem "n" -- "1" Batch : allocates
```

### Описание ключевых атрибутов сущностей:

*   **User (Пользователь)**: Базовая сущность для аутентификации. `NsPvAccess` определяет допуск к сильнодействующим веществам, `UkepBound` — привязку электронной подписи.
*   **EmployeeProfile (Профиль сотрудника)**: Хранит кадровые данные. `GDPHistory` — JSON-массив с историей ежегодных обучений надлежащей практике распределения.
*   **Product (Товар)**: Карточка лекарственного средства. Включает требования по хранению (`TempMin`/`Max`) и регуляторные флаги (`IsJNVLP` — ЖНВЛП, `IsNS_PV` — ПКУ).
*   **Batch (Серия)**: Партия конкретного товара на складе. Ключевые поля: `ExpiryDate` (срок годности для FEFO) и `Status` (контроль доступности).
*   **InboundReceipt (Приемка)**: Документ поступления. Связывает поставщика и принятые серии. Содержит результаты контроля `QP` (Уполномоченного лица).
*   **Order (Заказ)**: Документ отгрузки. Содержит статус жизненного цикла заказа и ссылки на ответственных за сборку и отгрузку.
*   **WarehouseZone (Зона склада)**: Физическое или логическое место хранения с заданными параметрами микроклимата (`TempMin`, `TempMax`, `HumidityMax`).

---

## 4. Схема инфраструктуры (Deployment)

*   **Reverse Proxy:** Nginx (Proxy pass к Next.js и Go API)
*   **Database:** PostgreSQL 16
*   **Caching:** Valkey 7.x (for OTP and Rate-limiting)
*   **Storage:** S3-compatible (for Photos and PDF Scans)
