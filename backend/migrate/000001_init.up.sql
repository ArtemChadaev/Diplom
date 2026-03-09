-- Создание таблицы пользователей
CREATE TABLE users (
    id            SERIAL PRIMARY KEY,
    full_name     VARCHAR(255) NOT NULL,                         -- ФИО
    password_hash TEXT         NOT NULL,                         -- Хеш пароля
    email         VARCHAR(255) NOT NULL UNIQUE,                  -- E-mail
    position      VARCHAR(255) NOT NULL,                         -- Должность
    role          VARCHAR(100) NOT NULL DEFAULT 'employee',      -- Полномочия / роль
    created_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW(),           -- Дата создания
    is_active     BOOLEAN      NOT NULL DEFAULT TRUE             -- true — работает, false — уволен
);
