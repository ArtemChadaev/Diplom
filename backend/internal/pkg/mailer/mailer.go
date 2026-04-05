package mailer

import (
	"context"
	"crypto/tls"
	"log/slog"
	"net/smtp"

	"github.com/jordan-wright/email"
	"github.com/ima/diplom-backend/internal/domain" // Импорт твоих ошибок
)

type Mailer interface {
	SendOTP(ctx context.Context, toEmail, code string) error
}

type Config struct {
	SMTPServer string
	SMTPPort   string
	Username   string
	Password   string
}

type smtpMailer struct {
	cfg Config
}

func New(cfg Config) Mailer {
	return &smtpMailer{cfg: cfg}
}

func (m *smtpMailer) SendOTP(ctx context.Context, toEmail, code string) error {
	// 1. Сборка письма через библиотеку
	e := email.NewEmail()
	e.From = "Diplom Admin <" + m.cfg.Username + ">" // Простая конкатенация вместо fmt
	e.To = []string{toEmail}
	e.Subject = "Verification Code"
	e.HTML = []byte("<h2>Your code: " + code + "</h2>") // Конкатенация для скорости

	auth := smtp.PlainAuth("", m.cfg.Username, m.cfg.Password, m.cfg.SMTPServer)
	addr := m.cfg.SMTPServer + ":" + m.cfg.SMTPPort

	// 2. Отправка
	err := e.SendWithTLS(
		addr,
		auth,
		&tls.Config{ServerName: m.cfg.SMTPServer},
	)

	if err != nil {
		// Создаем структурированную ошибку приложения
		appErr := domain.NewAppError(
			"mailer_send_error",
			"failed to physically send email via SMTP",
			err, // Оборачиваем реальную ошибку библиотеки
			slog.String("target_email", toEmail),
			slog.String("smtp_server", m.cfg.SMTPServer),
		)
		
		// Логируем через встроенный метод AppError
		appErr.LogError(ctx)
		return appErr
	}

	return nil
}