package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/ima/diplom-backend/internal/domain"
	"github.com/ima/diplom-backend/internal/handler/dto"
	"github.com/ima/diplom-backend/internal/handler/middleware"
)

func (h *Handler) adminVerifyUser(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	userID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid user id", http.StatusBadRequest)
		return
	}

	callerID := r.Context().Value(middleware.CtxUserID).(int)
	callerRole := r.Context().Value(middleware.CtxRole).(domain.UserRole)

	err = h.service.Auth.VerifyUser(r.Context(), callerID, callerRole, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "active"})
}

func (h *Handler) adminAssignRole(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	userID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid user id", http.StatusBadRequest)
		return
	}

	var req dto.AssignRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	callerID := r.Context().Value(middleware.CtxUserID).(int)
	callerRole := r.Context().Value(middleware.CtxRole).(domain.UserRole)

	err = h.service.Auth.AssignRole(r.Context(), callerID, callerRole, userID, domain.UserRole(req.Role))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"role": req.Role})
}

func (h *Handler) adminRevokeSession(w http.ResponseWriter, r *http.Request) {
	sessionIDStr := chi.URLParam(r, "sessionID")
	sessionID, err := uuid.Parse(sessionIDStr)
	if err != nil {
		http.Error(w, "invalid session id", http.StatusBadRequest)
		return
	}

	callerID := r.Context().Value(middleware.CtxUserID).(int)
	callerRole := r.Context().Value(middleware.CtxRole).(domain.UserRole)

	err = h.service.Auth.RevokeSession(r.Context(), sessionID, callerID, callerRole)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) revokeSession(w http.ResponseWriter, r *http.Request) {
	sessionIDStr := chi.URLParam(r, "sessionID")
	sessionID, err := uuid.Parse(sessionIDStr)
	if err != nil {
		http.Error(w, "invalid session id", http.StatusBadRequest)
		return
	}

	callerID := r.Context().Value(middleware.CtxUserID).(int)
	callerRole := r.Context().Value(middleware.CtxRole).(domain.UserRole)

	err = h.service.Auth.RevokeSession(r.Context(), sessionID, callerID, callerRole)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
