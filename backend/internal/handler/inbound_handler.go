package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/ima/diplom-backend/internal/domain"
	"github.com/ima/diplom-backend/internal/handler/dto"
	"github.com/ima/diplom-backend/internal/handler/middleware"
)

// listInbounds godoc
// @Summary      List inbound receivings
// @Description  Returns a paginated list of inbound records
// @Tags         Inbound
// @Produce      json
// @Param        limit   query     int  false  "Page size"
// @Param        offset  query     int  false  "Offset"
// @Success      200  {object}  dto.InboundListResponse
// @Router       /api/v1/inbound [get]
func (h *Handler) listInbounds(w http.ResponseWriter, r *http.Request) {
	limit := 10
	offset := 0

	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if val, err := strconv.Atoi(limitStr); err == nil {
			limit = val
		}
	}
	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		if val, err := strconv.Atoi(offsetStr); err == nil {
			offset = val
		}
	}

	inbounds, total, err := h.service.Inbound.ListInbounds(r.Context(), limit, offset)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to fetch inbounds")
		return
	}

	resp := dto.InboundListResponse{
		Total:    total,
		Inbounds: make([]dto.InboundResponse, len(inbounds)),
	}
	for i, in := range inbounds {
		resp.Inbounds[i] = dto.ToInboundResponse(in)
	}

	writeJSON(w, http.StatusOK, resp)
}

// getInbound godoc
// @Summary      Get inbound receiving by ID
// @Description  Returns detailed info about a single inbound record including items
// @Tags         Inbound
// @Produce      json
// @Param        id   path      string  true  "Inbound UUID"
// @Success      200  {object}  dto.InboundResponse
// @Router       /api/v1/inbound/{id} [get]
func (h *Handler) getInbound(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	in, err := h.service.Inbound.GetInbound(r.Context(), id)
	if err != nil {
		if err == domain.ErrInboundNotFound {
			writeError(w, http.StatusNotFound, "Inbound record not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "Failed to fetch inbound record")
		return
	}

	writeJSON(w, http.StatusOK, dto.ToInboundResponse(*in))
}

// createInbound godoc
// @Summary      Create inbound receiving
// @Description  Creates a new inbound record (Draft status)
// @Tags         Inbound
// @Accept       json
// @Produce      json
// @Param        request body dto.CreateInboundRequest true "Inbound data"
// @Success      201  {object}  dto.InboundResponse
// @Router       /api/v1/inbound [post]
func (h *Handler) createInbound(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateInboundRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	role, _ := r.Context().Value(middleware.CtxRole).(domain.UserRole)
	userID, _ := r.Context().Value(middleware.CtxUserID).(int)
	
	in := req.ToDomain()
	in.ReceivedBy = userID

	created, err := h.service.Inbound.CreateInbound(r.Context(), role, &in)
	if err != nil {
		if err == domain.ErrInsufficientPerms {
			writeError(w, http.StatusForbidden, "No permission to create inbound records")
			return
		}
		writeError(w, http.StatusInternalServerError, "Failed to create inbound record")
		return
	}

	writeJSON(w, http.StatusCreated, dto.ToInboundResponse(*created))
}

// updateInboundStatus godoc
// @Summary      Update inbound status
// @Description  Changes the status of an inbound record (e.g. Received -> Completed)
// @Tags         Inbound
// @Accept       json
// @Produce      json
// @Param        id      path      string  true  "Inbound UUID"
// @Param        request body dto.UpdateInboundStatusRequest true "New status"
// @Success      204  "No Content"
// @Router       /api/v1/inbound/{id}/status [patch]
func (h *Handler) updateInboundStatus(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req dto.UpdateInboundStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	role, _ := r.Context().Value(middleware.CtxRole).(domain.UserRole)

	if err := h.service.Inbound.UpdateStatus(r.Context(), role, id, req.Status); err != nil {
		if err == domain.ErrInboundNotFound {
			writeError(w, http.StatusNotFound, "Inbound record not found")
			return
		}
		if err == domain.ErrInsufficientPerms {
			writeError(w, http.StatusForbidden, "No permission to update status")
			return
		}
		if err == domain.ErrConflict {
			writeError(w, http.StatusConflict, "Cannot update status of completed record")
			return
		}
		writeError(w, http.StatusInternalServerError, "Failed to update inbound status")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// deleteInbound godoc
// @Summary      Delete inbound receiving
// @Description  Permanently deletes an inbound record (Draft status only, Admin only)
// @Tags         Inbound
// @Produce      json
// @Param        id   path      string  true  "Inbound UUID"
// @Success      204  "No Content"
// @Router       /api/v1/inbound/{id} [delete]
func (h *Handler) deleteInbound(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	role, _ := r.Context().Value(middleware.CtxRole).(domain.UserRole)

	if err := h.service.Inbound.DeleteInbound(r.Context(), role, id); err != nil {
		if err == domain.ErrInboundNotFound {
			writeError(w, http.StatusNotFound, "Inbound record not found")
			return
		}
		if err == domain.ErrInsufficientPerms {
			writeError(w, http.StatusForbidden, "No permission to delete inbound records")
			return
		}
		if err == domain.ErrConflict {
			writeError(w, http.StatusConflict, "Cannot delete completed record")
			return
		}
		writeError(w, http.StatusInternalServerError, "Failed to delete inbound record")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
