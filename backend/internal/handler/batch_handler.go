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

// listBatches godoc
// @Summary      List product batches (Inventory)
// @Description  Returns a paginated list of batches across all or specific zones
// @Tags         Batches
// @Produce      json
// @Param        product_id  query     string  false  "Filter by product UUID"
// @Param        zone_id     query     string  false  "Filter by zone UUID"
// @Param        status      query     string  false  "Filter by status"
// @Param        limit       query     int     false  "Page size"
// @Param        offset      query     int     false  "Offset"
// @Success      200  {object}  dto.BatchListResponse
// @Router       /api/v1/batches [get]
func (h *Handler) listBatches(w http.ResponseWriter, r *http.Request) {
	filter := domain.BatchFilter{
		Limit:  10,
		Offset: 0,
	}

	q := r.URL.Query()
	filter.ProductID = q.Get("product_id")
	filter.ZoneID = q.Get("zone_id")
	filter.Status = domain.BatchStatus(q.Get("status"))

	if val, err := strconv.Atoi(q.Get("limit")); err == nil {
		filter.Limit = val
	}
	if val, err := strconv.Atoi(q.Get("offset")); err == nil {
		filter.Offset = val
	}

	batches, total, err := h.service.Batch.ListBatches(r.Context(), filter)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to fetch batches")
		return
	}

	resp := dto.BatchListResponse{
		Total:   total,
		Batches: make([]dto.BatchResponse, len(batches)),
	}
	for i, b := range batches {
		resp.Batches[i] = dto.ToBatchResponse(b)
	}

	writeJSON(w, http.StatusOK, resp)
}

// getBatch godoc
// @Summary      Get batch by ID
// @Description  Returns detailed info about a specific batch including quantity and expiry
// @Tags         Batches
// @Produce      json
// @Param        id   path      string  true  "Batch UUID"
// @Success      200  {object}  dto.BatchResponse
// @Router       /api/v1/batches/{id} [get]
func (h *Handler) getBatch(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	b, err := h.service.Batch.GetBatch(r.Context(), id)
	if err != nil {
		if err == domain.ErrBatchNotFound {
			writeError(w, http.StatusNotFound, "Batch not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "Failed to fetch batch")
		return
	}

	writeJSON(w, http.StatusOK, dto.ToBatchResponse(*b))
}

// updateBatchStatus godoc
// @Summary      Update batch status
// @Description  Manually changes batch status (e.g. Quarantine -> Available)
// @Tags         Batches
// @Accept       json
// @Produce      json
// @Param        id      path      string  true  "Batch UUID"
// @Param        request body dto.UpdateBatchStatusRequest true "New status"
// @Success      204  "No Content"
// @Router       /api/v1/batches/{id}/status [patch]
func (h *Handler) updateBatchStatus(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req dto.UpdateBatchStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	role, _ := r.Context().Value(middleware.CtxRole).(domain.UserRole)

	if err := h.service.Batch.UpdateStatus(r.Context(), role, id, req.Status); err != nil {
		if err == domain.ErrBatchNotFound {
			writeError(w, http.StatusNotFound, "Batch not found")
			return
		}
		if err == domain.ErrInsufficientPerms {
			writeError(w, http.StatusForbidden, "No permission to update batch status")
			return
		}
		writeError(w, http.StatusInternalServerError, "Failed to update batch status")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// transferBatch godoc
// @Summary      Transfer batch to another zone
// @Description  Moves stock series from one warehouse zone to another
// @Tags         Batches
// @Accept       json
// @Produce      json
// @Param        id      path      string  true  "Batch UUID"
// @Param        request body dto.TransferBatchRequest true "Target zone info"
// @Success      204  "No Content"
// @Router       /api/v1/batches/{id}/transfer [post]
func (h *Handler) transferBatch(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req dto.TransferBatchRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	role, _ := r.Context().Value(middleware.CtxRole).(domain.UserRole)

	if err := h.service.Batch.TransferBatch(r.Context(), role, id, req.TargetZoneID); err != nil {
		if err == domain.ErrBatchNotFound {
			writeError(w, http.StatusNotFound, "Batch not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "Failed to transfer batch")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
