package main

import (
	"context"
	"errors"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/ima/diplom-backend/internal/bootstrap"
	"github.com/ima/diplom-backend/internal/config"
	"github.com/ima/diplom-backend/internal/domain"
	"github.com/ima/diplom-backend/internal/handler"
	"github.com/ima/diplom-backend/internal/pkg/logger"
	"github.com/ima/diplom-backend/internal/pkg/mailer"
	"github.com/ima/diplom-backend/internal/repository"
	"github.com/ima/diplom-backend/internal/service"
	"github.com/valkey-io/valkey-go"
)

func main() {
	// 1. Инициализация логгера
	logger.Setup(os.Getenv("APP_ENV"))

	// 2. Graceful shutdown контекст
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	log := logger.FromContext(ctx)

	// 2. Загрузка конфигурации из .env
	cfg, err := config.Load()
	if err != nil {
		log.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	log.Info("config loaded", "port", cfg.Port)

	// 3. Подключение к БД
	db, err := repository.NewPostgresDB(repository.PostgresConfig{
		Host:     cfg.DBHost,
		Port:     cfg.DBPort,
		Username: cfg.DBUser,
		Database: cfg.DBName,
		Password: cfg.DBPassword,
		SSLMode:  "disable",
	})
	if err != nil {
		log.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	sqlDB, err := db.DB()
	if err != nil {
		log.Error("failed to get sql.DB from gorm explorer", "error", err)
		os.Exit(1)
	}
	defer sqlDB.Close()

	log.Info("database connection established")

	// 4. Подключение к Valkey (Redis)
	valkeyClient, err := valkey.NewClient(valkey.ClientOption{
		InitAddress: []string{net.JoinHostPort(cfg.ValkeyHost, cfg.ValkeyPort)},
		Password:    cfg.ValkeyPassword,
	})
	if err != nil {
		log.Error("failed to connect to valkey", "error", err)
		os.Exit(1)
	}
	defer valkeyClient.Close()
	log.Info("valkey connection established")

	// 5. Инициализация Mailer
	m := mailer.New(mailer.Config{
		SMTPServer: cfg.MAILER_SMTP_SERVER,
		SMTPPort:   cfg.MAILER_SMTP_PORT,
		Username:   cfg.MAILER_USERNAME,
		Password:   cfg.MAILER_PASSWORD,
	})

	// 6. Инициализация слоёв: Repository → Service → Handler
	repos := repository.NewRepository(db, valkeyClient)

	services := service.NewService(repos, cfg.JWTSecret, cfg.GoogleClientID, cfg.OTPHMACSecret, m)
	handlers := handler.NewHandler(services, services.Token, cfg, repos.User, repos.EmployeeProfile)

	router := handlers.Router()

	// 7. Запуск HTTP-сервера
	srv := new(domain.Server)
	go func() {
		log.Info("server started", "addr", ":"+cfg.Port)
		if err := srv.Run(cfg.Port, router); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("server error", "error", err)
			os.Exit(1)
		}
	}()

	// Admin Seeding
	if err := bootstrap.SeedAdmin(ctx, cfg, repos.User); err != nil {
		log.Error("admin bootstrap failed", "error", err)
		os.Exit(1)
	}

	// 8. Ожидание сигнала завершения
	<-ctx.Done()
	log.Info("shutting down server...")

	if err := srv.Shutdown(context.Background()); err != nil {
		log.Error("shutdown error", "error", err)
	}
}
