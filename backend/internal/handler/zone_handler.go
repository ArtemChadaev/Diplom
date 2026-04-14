package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/ima/diplom-backend/internal/domain"
	"github.com/ima/diplom-backend/internal/handler/dto"
	"github.com/ima/diplom-backend/internal/handler/middleware"
)

// listZones godoc
// @Summary      List warehouse zones
// @Description  Returns all warehouse zones
// @Tags         Zones
// @Produce      json
// @Success      200  {array}  dto.ZoneResponse
// @Router       /api/v1/zones [get]
func (h *Handler) listZones(w http.ResponseWriter, r *http.Request) {
	zones, err := h.service.Zone.ListZones(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to fetch zones")
		return
	}

	resp := make([]dto.ZoneResponse, len(zones))
	for i, z := range zones {
		resp[i] = dto.ToZoneResponse(z)
	}

	writeJSON(w, http.StatusOK, resp)
}

// getZone godoc
// @Summary      Get warehouse zone by ID
// @Description  Returns detailed info about a single zone
// @Tags         Zones
// @Produce      json
// @Param        id   path      string  true  "Zone UUID"
// @Success      200  {object}  dto.ZoneResponse
// @Router       /api/v1/zones/{id} [get]
func (h *Handler) getZone(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	z, err := h.service.Zone.GetZone(r.Context(), id)
	if err != nil {
		if err == domain.ErrZoneNotFound {
			writeError(w, http.StatusNotFound, "Zone not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "Failed to fetch zone")
		return
	}

	writeJSON(w, http.StatusOK, dto.ToZoneResponse(*z))
}

// createZone godoc
// @Summary      Create warehouse zone
// @Description  Creates a new zone (Admin or Warehouse Manager only)
// @Tags         Zones
// @Accept       json
// @Produce      json
// @Param        request body dto.CreateZoneRequest true "Zone data"
// @Success      201  {object}  dto.ZoneResponse
// @Router       /api/v1/zones [post]
func (h *Handler) createZone(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateZoneRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	role, _ := r.Context().Value(middleware.CtxRole).(domain.UserRole)
	z := req.ToDomain()

	created, err := h.service.Zone.CreateZone(r.Context(), role, &z)
	if err != nil {
		if err == domain.ErrInsufficientPerms {
			writeError(w, http.StatusForbidden, "No permission to create zones")
			return
		}
		writeError(w, http.StatusInternalServerError, "Failed to create zone")
		return
	}

	writeJSON(w, http.StatusCreated, dto.ToZoneResponse(*created))
}

// patchZone godoc
// @Summary      Partial update warehouse zone
// @Description  Updates specific fields of an existing zone (Admin or Warehouse Manager only)
// @Tags         Zones
// @Accept       json
// @Produce      json
// @Param        id      path      string  true  "Zone UUID"
// @Param        request body dto.UpdateZoneRequest true "Partial update data"
// @Success      200  {object}  dto.ZoneResponse
// @Router       /api/v1/zones/{id} [patch]
func (h *Handler) patchZone(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req dto.UpdateZoneRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	role, _ := r.Context().Value(middleware.CtxRole).(domain.UserRole)
	
	z, err := h.service.Zone.GetZone(r.Context(), id)
	if err != nil {
		if err == domain.ErrZoneNotFound {
			writeError(w, http.StatusNotFound, "Zone not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "Failed to fetch zone")
		return
	}

	req.ApplyTo(z)

	updated, err := h.service.Zone.UpdateZone(r.Context(), role, z)
	if err != nil {
		if err == domain.ErrInsufficientPerms {
			writeError(w, http.StatusForbidden, "No permission to update zones")
			return
		}
		writeError(w, http.StatusInternalServerError, "Failed to update zone")
		return
	}

	writeJSON(w, http.StatusOK, dto.ToZoneResponse(*updated))
}

// deleteZone godoc
// @Summary      Delete warehouse zone
// @Description  Permanently deletes a zone (Admin only)
// @Tags         Zones
// @Produce      json
// @Param        id   path      string  true  "Zone UUID"
// @Success      204  "No Content"
// @Router       /api/v1/zones/{id} [delete]
func (h *Handler) deleteZone(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	role, _ := r.Context().Value(middleware.CtxRole).(domain.UserRole)

	if err := h.service.Zone.DeleteZone(r.Context(), role, id); err != nil {
		if err == domain.ErrZoneNotFound {
			writeError(w, http.StatusNotFound, "Zone not found")
			return
		}
		if err == domain.ErrInsufficientPerms {
			writeError(w, http.StatusForbidden, "No permission to delete zones")
			return
		}
		writeError(w, http.StatusInternalServerError, "Failed to delete zone")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
