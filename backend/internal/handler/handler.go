package handler

import (
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/ima/diplom-backend/internal/service"
)

type Handler struct {
	service service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{
		service: *service,
	}
}

func (h *Handler) Router() chi.Router {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(middleware.Heartbeat("/ping"))
	r.Use(h.loggingMiddleware)

	// TODO: группировать под /api/v1 по мере роста
	// TODO: Сделать проверку в middleware на то какой пользователь, рандомный, стандартный или админ
	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", h.register)
		r.Post("/login", h.login)
		r.Post("/refresh", h.refresh)
		r.Post("/logout", h.logout)
	})

	return r
}
