package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/ima/diplom-backend/internal/domain"
	h_middleware "github.com/ima/diplom-backend/internal/handler/middleware"
	"github.com/ima/diplom-backend/internal/config"
	"github.com/ima/diplom-backend/internal/service"
)

type Handler struct {
	service             service.Service
	tokenSvc            domain.TokenService
	cfg                 *config.Config
	userRepo           domain.UserRepository
	employeeProfileRepo domain.EmployeeProfileRepository
}

func NewHandler(service *service.Service, tokenSvc domain.TokenService, cfg *config.Config, userRepo domain.UserRepository, employeeProfileRepo domain.EmployeeProfileRepository) *Handler {
	return &Handler{
		service:             *service,
		tokenSvc:            tokenSvc,
		cfg:                 cfg,
		userRepo:           userRepo,
		employeeProfileRepo: employeeProfileRepo,
	}
}

func (h *Handler) Router() chi.Router {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(middleware.Heartbeat("/ping"))
	r.Use(h_middleware.RequestLogger)


	// Public Auth Routes (OAuth only + OTP)
	r.Route("/auth", func(r chi.Router) {
		r.Post("/google", h.googleLogin)
		r.Post("/refresh", h.refresh)
		r.Post("/logout", h.logout)
		r.Post("/send-code", h.sendCode)
		r.Post("/verify-code", h.verifyCode)
	})

	// Protected API Routes
	r.Route("/api/v1", func(r chi.Router) {
		r.Use(h_middleware.AuthRequired(h.tokenSvc))

		// User profile (own)
		r.Get("/users/me", h.getMe)

		// Session management (own)
		r.Delete("/sessions/{sessionID}", h.revokeSession)

		// Admin routes
		r.Route("/admin", func(r chi.Router) {
			r.Use(h_middleware.RequireRole(domain.RoleAdmin))

			r.Patch("/users/{id}/role", h.adminAssignRole)
			r.Patch("/users/{id}/blocked", h.adminSetBlocked)
			r.Delete("/sessions/{sessionID}", h.adminRevokeSession)

			// Admin Users list/edit
			r.Get("/users", h.listUsers)
			r.Get("/users/{id}", h.getUserByID)
			r.Patch("/users/{id}", h.patchUser)

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
func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}
