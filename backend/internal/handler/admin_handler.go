package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/ima/diplom-backend/internal/domain"
	"github.com/ima/diplom-backend/internal/handler/dto"
	"github.com/ima/diplom-backend/internal/handler/middleware"
)

func mapAdminError(w http.ResponseWriter, err error) {
	if errors.Is(err, domain.ErrUserNotFound) || errors.Is(err, domain.ErrSessionNotFound) {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if errors.Is(err, domain.ErrInsufficientPerms) {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}
	http.Error(w, err.Error(), http.StatusInternalServerError)
}

// adminSetBlocked godoc
// @Summary      Block/Unblock user (Admin)
// @Description  Toggles the blocked status of a specific user. Admin role required.
// @Tags         Admin
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id    path      int  true  "User ID"
// @Param        body  body      dto.SetBlockedRequest  true  "Block status"
// @Success      200   {object}  map[string]bool  "blocked status response"
// @Failure      400   {object}  dto.ErrorResponse  "invalid user ID or invalid JSON"
// @Failure      403   {object}  dto.ErrorResponse  "insufficient permissions"
// @Failure      404   {object}  dto.ErrorResponse  "user not found"
// @Router       /api/v1/admin/users/{id}/blocked [patch]
func (h *Handler) adminSetBlocked(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	userID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid user id", http.StatusBadRequest)
		return
	}

	var req dto.SetBlockedRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	callerID := r.Context().Value(middleware.CtxUserID).(int)
	callerRole := r.Context().Value(middleware.CtxRole).(domain.UserRole)

	if err = h.service.Auth.SetBlocked(r.Context(), callerID, callerRole, userID, req.Blocked); err != nil {
		mapAdminError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]bool{"blocked": req.Blocked})
}

// adminAssignRole godoc
// @Summary      Assign user role (Admin)
// @Description  Updates the role of a user. Admin role required.
// @Tags         Admin
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id    path      int  true  "User ID"
// @Param        body  body      dto.AssignRoleRequest  true  "New role details"
// @Success      200   {object}  map[string]string  "assigned role response"
// @Failure      400   {object}  dto.ErrorResponse  "invalid user ID or invalid JSON"
// @Failure      403   {object}  dto.ErrorResponse  "insufficient permissions"
// @Failure      404   {object}  dto.ErrorResponse  "user not found"
// @Router       /api/v1/admin/users/{id}/role [patch]
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
		mapAdminError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"role": req.Role})
}

// adminRevokeSession godoc
// @Summary      Revoke user session (Admin)
// @Description  Revokes a specific active user session. Admin role required.
// @Tags         Admin
// @Security     BearerAuth
// @Param        sessionID  path      string  true  "Session UUID"
// @Success      204        "No Content"
// @Failure      400        {object}  dto.ErrorResponse  "invalid session ID"
// @Failure      403        {object}  dto.ErrorResponse  "insufficient permissions"
// @Failure      404        {object}  dto.ErrorResponse  "session not found"
// @Router       /api/v1/admin/sessions/{sessionID} [delete]
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
		mapAdminError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// revokeSession godoc
// @Summary      Revoke own session
// @Description  Revokes a specific active session belonging to the caller.
// @Tags         Sessions
// @Security     BearerAuth
// @Param        sessionID  path      string  true  "Session UUID"
// @Success      204        "No Content"
// @Failure      400        {object}  dto.ErrorResponse  "invalid session ID"
// @Failure      403        {object}  dto.ErrorResponse  "insufficient permissions (not owner and not admin)"
// @Failure      404        {object}  dto.ErrorResponse  "session not found"
// @Router       /api/v1/sessions/{sessionID} [delete]
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
		mapAdminError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
