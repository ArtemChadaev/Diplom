package main

import (
	"context"
	"errors"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/ima/diplom-backend/internal/bootstrap"
	"github.com/ima/diplom-backend/internal/config"
	"github.com/ima/diplom-backend/internal/domain"
	"github.com/ima/diplom-backend/internal/handler"
	"github.com/ima/diplom-backend/internal/repository"
	"github.com/ima/diplom-backend/internal/service"
)

func main() {
	// 1. Загрузка конфигурации из .env
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("ошибка загрузки конфига: %s", err)
	}

	slog.Info("конфигурация загружена", slog.String("port", cfg.Port))

	// 2. Graceful shutdown контекст
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

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
		log.Fatalf("ошибка подключения к БД: %s", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("ошибка получения sql.DB: %s", err)
	}
	defer sqlDB.Close()

	slog.Info("подключение к БД установлено")

	// 4. Инициализация слоёв: Repository → Service → Handler
	repos := repository.NewRepository(db)

	// Admin Seeding
	if err := bootstrap.SeedAdmin(ctx, cfg, repos.User); err != nil {
		log.Fatalf("ошибка инициализации админа: %s", err)
	}

	services := service.NewService(repos, cfg.JWTSecret)
	handlers := handler.NewHandler(services)

	router := handlers.Router()

	// 5. Запуск HTTP-сервера
	srv := new(domain.Server)
	go func() {
		slog.Info("сервер запущен", slog.String("addr", ":"+cfg.Port))
		if err := srv.Run(cfg.Port, router); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("ошибка сервера: %s", err)
		}
	}()

	// 6. Ожидание сигнала завершения
	<-ctx.Done()
	slog.Info("завершение работы сервера...")

	if err := srv.Shutdown(context.Background()); err != nil {
		log.Printf("ошибка при shutdown: %s", err)
	}
}
