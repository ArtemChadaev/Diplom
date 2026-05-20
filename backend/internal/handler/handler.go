package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/ima/diplom-backend/internal/config"
	"github.com/ima/diplom-backend/internal/domain"
	h_middleware "github.com/ima/diplom-backend/internal/handler/middleware"
	"github.com/ima/diplom-backend/internal/service"
)

type Handler struct {
	service             service.Service
	tokenSvc            domain.TokenService
	cfg                 *config.Config
	userRepo            domain.UserRepository
	employeeProfileRepo domain.EmployeeProfileRepository
}

func NewHandler(service *service.Service, tokenSvc domain.TokenService, cfg *config.Config, userRepo domain.UserRepository, employeeProfileRepo domain.EmployeeProfileRepository) *Handler {
	return &Handler{
		service:             *service,
		tokenSvc:            tokenSvc,
		cfg:                 cfg,
		userRepo:            userRepo,
		employeeProfileRepo: employeeProfileRepo,
	}
}

func (h *Handler) Router() chi.Router {
	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   h.cfg.AllowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Use(middleware.RequestID)
	r.Use(middleware.CleanPath)       // Убирает двойные слэши //
	r.Use(middleware.RedirectSlashes) // Направляет /path/ на /path
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(middleware.Heartbeat("/ping"))
	r.Use(h_middleware.InjectIPAddress) // Записывает IP клиента в контекст (используется аудит-логом)
	r.Use(h_middleware.RequestLogger)

	// Healthcheck
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	// Public Auth Routes (OAuth only + OTP)
	r.Route("/auth", func(r chi.Router) {
		r.Post("/google", h.googleLogin)
		r.Post("/register", h.register)
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
		r.Get("/users/me/profile", h.getMeProfile)
		r.Patch("/users/me", h.patchMe)

		// Session management (own)
		r.Delete("/sessions/{sessionID}", h.revokeSession)

		// Admin routes
		r.Route("/admin", func(r chi.Router) {
			r.Use(h_middleware.RequireRole(domain.RoleAdmin))

			// Caddy forward_auth gate: validates JWT + admin role before proxying to Dozzle
			r.Get("/auth/check-logs-auth", h.checkLogsAuth)

			r.Patch("/users/{id}/role", h.adminAssignRole)
			r.Patch("/users/{id}/blocked", h.adminSetBlocked)
			r.Delete("/sessions/{sessionID}", h.adminRevokeSession)

			// Admin Users list/edit
			r.Get("/users", h.listUsers)
			r.Get("/users/{id}", h.getUserByID)
			r.Patch("/users/{id}", h.patchUser)

			// Admin Employee Profile routes
			r.Route("/employees", func(r chi.Router) {
				r.Post("/", h.adminCreateEmployeeProfile)
				r.Get("/", h.adminListEmployeeProfiles)
				r.Get("/{userID}", h.adminGetEmployeeProfile)
				r.Patch("/{userID}", h.adminPatchEmployeeProfile)
			})
		})

		// Reference data
		r.Route("/ref", func(r chi.Router) {
			r.Get("/countries", h.listCountries)
			r.Get("/atc", h.searchATC)
		})

		// Product management
		r.Route("/products", func(r chi.Router) {
			r.Get("/", h.listProducts)
			r.Get("/reorder-check", h.runReorderCheck)
			r.Get("/{id}", h.getProduct)
			r.Post("/", h.createProduct)
			r.Patch("/{id}", h.patchProduct)
			r.Delete("/{id}", h.deleteProduct)
		})

		// Supplier management
		r.Route("/suppliers", func(r chi.Router) {
			r.Get("/", h.listSuppliers)
			r.Get("/{id}", h.getSupplier)
			r.Post("/", h.createSupplier)
			r.Patch("/{id}", h.patchSupplier)
			r.Delete("/{id}", h.deleteSupplier)
		})

		// Warehouse zoning
		r.Route("/zones", func(r chi.Router) {
			r.Get("/", h.listZones)
			r.Get("/{id}", h.getZone)
			r.Post("/", h.createZone)
			r.Patch("/{id}", h.patchZone)
			r.Delete("/{id}", h.deleteZone)
		})

		// Inbound receiving
		r.Route("/inbound", func(r chi.Router) {
			r.Get("/", h.listInbounds)
			r.Get("/{id}", h.getInbound)
			r.Post("/", h.createInbound)
			r.Patch("/{id}/status", h.updateInboundStatus)
			r.Delete("/{id}", h.deleteInbound)
		})

		// Environment logs
		r.Route("/env", func(r chi.Router) {
			r.Get("/logs", h.listEnvLogs)
			r.Get("/logs/export", h.exportEnvLogs)
			r.Post("/logs", h.recordEnvLogs)
		})

		// Orders
		r.Route("/orders", func(r chi.Router) {
			r.Get("/", h.listOrders)
			r.Get("/{id}", h.getOrder)
			r.Post("/", h.createOrder)
			r.Patch("/{id}/status", h.updateOrderStatus)
			r.Get("/{id}/pdf/ttn", h.getOrderTTN)
			r.Get("/{id}/pdf/quality-registry", h.getOrderQualityRegistry)
		})

		// Inventory
		r.Route("/inventory", func(r chi.Router) {
			r.Get("/", h.listInventorySessions)
			r.Get("/{id}", h.getInventorySession)
			r.Post("/", h.startInventorySession)
			r.Post("/{id}/count", h.submitInventoryCount)
			r.Post("/{id}/finish", h.finishInventorySession)
			r.With(h_middleware.RequireRole(domain.RoleAdmin, domain.RoleWarehouseManager)).Get("/{id}/netting", h.getInventoryNetting)
		})

		// Claims
		r.Route("/claims", func(r chi.Router) {
			r.Get("/", h.listClaims)
			r.Get("/{id}", h.getClaim)
			r.Post("/", h.createClaim)
			r.Patch("/{id}/status", h.updateClaimStatus)
		})

		// System Settings
		r.Route("/settings", func(r chi.Router) {
			r.Get("/", h.listSettings)
			r.Patch("/{key}", h.updateSetting)
		})

		// Batches (Standing stock)
		r.Route("/batches", func(r chi.Router) {
			r.Get("/", h.listBatches)
			r.Get("/{id}", h.getBatch)
			r.Patch("/{id}/status", h.updateBatchStatus)
			r.Post("/{id}/transfer", h.transferBatch)
		})

		// Recalled (Blocked series)
		r.Route("/recalled", func(r chi.Router) {
			r.Get("/", h.listRecalled)
			r.Get("/check/{serial}", h.checkBatch)
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
