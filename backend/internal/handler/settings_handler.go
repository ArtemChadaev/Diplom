package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/ima/diplom-backend/internal/domain"
	"github.com/ima/diplom-backend/internal/handler/dto"
	"github.com/ima/diplom-backend/internal/handler/middleware"
)

// listSettings godoc
// @Summary      List system settings
// @Description  Returns all global system parameters
// @Tags         Settings
// @Produce      json
// @Success      200  {array}  dto.SystemSettingResponse
// @Router       /api/v1/settings [get]
func (h *Handler) listSettings(w http.ResponseWriter, r *http.Request) {
	settings, err := h.service.Settings.ListSettings(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to fetch settings")
		return
	}

	resp := make([]dto.SystemSettingResponse, len(settings))
	for i, s := range settings {
		resp[i] = dto.SystemSettingResponse{
			Key:   s.Key,
			Value: s.Value,
		}
	}

	writeJSON(w, http.StatusOK, resp)
}

// updateSetting godoc
// @Summary      Update system setting
// @Description  Changes a specific global system parameter (Admin only)
// @Tags         Settings
// @Accept       json
// @Produce      json
// @Param        key     path      string  true  "Setting key"
// @Param        request body dto.UpdateSettingRequest true "New value"
// @Success      204  "No Content"
// @Router       /api/v1/settings/{key} [patch]
func (h *Handler) updateSetting(w http.ResponseWriter, r *http.Request) {
	key := chi.URLParam(r, "key")
	var req dto.UpdateSettingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	role, _ := r.Context().Value(middleware.CtxRole).(domain.UserRole)

	if err := h.service.Settings.UpdateSetting(r.Context(), role, key, req.Value); err != nil {
		if err == domain.ErrInsufficientPerms {
			writeError(w, http.StatusForbidden, "Only Admin can update settings")
			return
		}
		writeError(w, http.StatusInternalServerError, "Failed to update setting")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
