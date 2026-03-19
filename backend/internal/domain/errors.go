package domain

import (
	"errors"
)

var (
	ErrUserNotFound      = errors.New("пользователь не найден")
	ErrLoginTaken        = errors.New("логин уже занят")
	ErrEmailTaken        = errors.New("email уже занят")
	ErrInvalidCreds      = errors.New("неверный логин или пароль")
	ErrUserUnverified    = errors.New("аккаунт ожидает подтверждения администратором")
	ErrUserBlocked       = errors.New("аккаунт заблокирован")
	ErrTokenExpired      = errors.New("срок действия токена истёк")
	ErrInvalidToken      = errors.New("недействительный токен")
	ErrSessionNotFound   = errors.New("сессия не найдена или была завершена")
	ErrInvalidTelegram   = errors.New("неверные данные авторизации Telegram")
	ErrInsufficientPerms = errors.New("недостаточно прав для выполнения операции")
)
