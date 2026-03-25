package handler

import (
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/ima/diplom-backend/internal/domain"
	h_middleware "github.com/ima/diplom-backend/internal/handler/middleware"
	"github.com/ima/diplom-backend/internal/service"
)

type Handler struct {
	service  service.Service
	tokenSvc domain.TokenService
}

func NewHandler(service *service.Service, tokenSvc domain.TokenService) *Handler {
	return &Handler{
		service:  *service,
		tokenSvc: tokenSvc,
	}
}

func (h *Handler) Router() chi.Router {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(middleware.Heartbeat("/ping"))
	r.Use(h.loggingMiddleware)


	// Public Auth Routes
	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", h.register)
		r.Post("/login", h.login)
		r.Post("/google", h.googleLogin)
		r.Post("/refresh", h.refresh)
		r.Post("/logout", h.logout) // Now requires Authorization header inside handler
	})

	// Protected API Routes
	r.Route("/api/v1", func(r chi.Router) {
		r.Use(h_middleware.AuthRequired(h.tokenSvc))

		// Session management (own)
		r.Delete("/sessions/{sessionID}", h.revokeSession)

		// Admin routes
		r.Route("/admin", func(r chi.Router) {
			r.Use(h_middleware.RequireRole(domain.RoleAdmin))

			r.Patch("/users/{id}/verify", h.adminVerifyUser)
			r.Patch("/users/{id}/role", h.adminAssignRole)
			r.Delete("/sessions/{sessionID}", h.adminRevokeSession)

			// Admin Employee Profile routes
			r.Route("/employees", func(r chi.Router) {
				r.Get("/", h.adminListEmployeeProfiles)
				r.Get("/{userID}", h.adminGetEmployeeProfile)
				r.Patch("/{userID}", h.adminPatchEmployeeProfile)
			})
		})

	})

	return r
}
