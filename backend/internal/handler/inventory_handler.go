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

// listInventorySessions godoc
// @Summary      List inventory sessions
// @Description  Returns a paginated list of inventory audit sessions
// @Tags         Inventory
// @Produce      json
// @Param        limit   query     int  false  "Page size"
// @Param        offset  query     int  false  "Offset"
// @Success      200  {array}  dto.InventorySessionResponse
// @Router       /api/v1/inventory [get]
func (h *Handler) listInventorySessions(w http.ResponseWriter, r *http.Request) {
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

	sessions, _, err := h.service.Inventory.ListSessions(r.Context(), limit, offset)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to fetch inventory sessions")
		return
	}

	resp := make([]dto.InventorySessionResponse, len(sessions))
	for i, s := range sessions {
		resp[i] = dto.ToInventorySessionResponse(s)
	}

	writeJSON(w, http.StatusOK, resp)
}

// getInventorySession godoc
// @Summary      Get inventory session by ID
// @Description  Returns detailed info about a single audit session including results
// @Tags         Inventory
// @Produce      json
// @Param        id   path      string  true  "Session UUID"
// @Success      200  {object}  dto.InventorySessionResponse
// @Router       /api/v1/inventory/{id} [get]
func (h *Handler) getInventorySession(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	s, err := h.service.Inventory.GetSession(r.Context(), id)
	if err != nil {
		if err == domain.ErrInventorySessionNotFound {
			writeError(w, http.StatusNotFound, "Inventory session not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "Failed to fetch session")
		return
	}

	writeJSON(w, http.StatusOK, dto.ToInventorySessionResponse(*s))
}

// startInventorySession godoc
// @Summary      Start inventory session
// @Description  Initiates a new audit session for a specific zone
// @Tags         Inventory
// @Accept       json
// @Produce      json
// @Param        request body dto.StartInventoryRequest true "Zone info"
// @Success      201  {object}  dto.InventorySessionResponse
// @Router       /api/v1/inventory [post]
func (h *Handler) startInventorySession(w http.ResponseWriter, r *http.Request) {
	var req dto.StartInventoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	userID, _ := r.Context().Value(middleware.CtxUserID).(int)
	
	created, err := h.service.Inventory.StartSession(r.Context(), userID, req.ZoneID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to start session")
		return
	}

	writeJSON(w, http.StatusCreated, dto.ToInventorySessionResponse(*created))
}

// submitInventoryCount godoc
// @Summary      Submit inventory counts
// @Description  Records physical quantities for products within an active session
// @Tags         Inventory
// @Accept       json
// @Produce      json
// @Param        id      path      string  true  "Session UUID"
// @Param        request body dto.SubmitCountRequest true "Counts data"
// @Success      204  "No Content"
// @Router       /api/v1/inventory/{id}/count [post]
func (h *Handler) submitInventoryCount(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req dto.SubmitCountRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	domainItems := make([]domain.InventoryItem, len(req.Items))
	for i, item := range req.Items {
		domainItems[i] = item.ToDomain()
	}

	if err := h.service.Inventory.SubmitCount(r.Context(), id, domainItems); err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to submit counts")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// finishInventorySession godoc
// @Summary      Finish inventory session
// @Description  Closes an active session, marking it as completed
// @Tags         Inventory
// @Produce      json
// @Param        id   path      string  true  "Session UUID"
// @Success      204  "No Content"
// @Router       /api/v1/inventory/{id}/finish [post]
func (h *Handler) finishInventorySession(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := h.service.Inventory.FinishSession(r.Context(), id); err != nil {
		if err == domain.ErrInventorySessionNotFound {
			writeError(w, http.StatusNotFound, "Inventory session not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "Failed to finish session")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
