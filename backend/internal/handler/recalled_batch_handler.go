package handler

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/ima/diplom-backend/internal/handler/dto"
)

// listRecalled godoc
// @Summary      List recalled batches
// @Description  Returns a paginated list of blocked pharmaceutical series from the central register
// @Tags         Recalled
// @Produce      json
// @Param        limit   query     int  false  "Page size"
// @Param        offset  query     int  false  "Offset"
// @Success      200  {array}  dto.RecalledBatchResponse
// @Router       /api/v1/recalled [get]
func (h *Handler) listRecalled(w http.ResponseWriter, r *http.Request) {
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

	batches, _, err := h.service.RecalledBatch.ListRecalled(r.Context(), limit, offset)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to fetch recalled batches")
		return
	}

	resp := make([]dto.RecalledBatchResponse, len(batches))
	for i, b := range batches {
		resp[i] = dto.ToRecalledBatchResponse(b)
	}

	writeJSON(w, http.StatusOK, resp)
}

// checkBatch godoc
// @Summary      Check series for recall
// @Description  Verifies if a specific serial number is blocked for pharmaceutical circulation
// @Tags         Recalled
// @Produce      json
// @Param        serial   path      string  true  "Series Number"
// @Success      200  {object}  dto.RecalledCheckResponse
// @Router       /api/v1/recalled/check/{serial} [get]
func (h *Handler) checkBatch(w http.ResponseWriter, r *http.Request) {
	serial := chi.URLParam(r, "serial")
	isRecalled, details, err := h.service.RecalledBatch.CheckBatch(r.Context(), serial)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to check batch")
		return
	}

	resp := dto.RecalledCheckResponse{
		IsRecalled: isRecalled,
	}
	if details != nil {
		d := dto.ToRecalledBatchResponse(*details)
		resp.Details = &d
	}

	writeJSON(w, http.StatusOK, resp)
}
