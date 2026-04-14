package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/ima/diplom-backend/internal/domain"
	"github.com/ima/diplom-backend/internal/handler/dto"
	"github.com/ima/diplom-backend/internal/handler/middleware"
)

// listEnvLogs godoc
// @Summary      List environment logs
// @Description  Returns a paginated list of temperature/humidity logs
// @Tags         Environment
// @Produce      json
// @Param        zone_id  query     string  false  "Filter by zone UUID"
// @Param        limit    query     int     false  "Page size"
// @Param        offset   query     int     false  "Offset"
// @Success      200  {object}  dto.EnvLogListResponse
// @Router       /api/v1/env/logs [get]
func (h *Handler) listEnvLogs(w http.ResponseWriter, r *http.Request) {
	zoneID := r.URL.Query().Get("zone_id")
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

	logs, total, err := h.service.EnvironmentLog.ListLogs(r.Context(), zoneID, limit, offset)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to fetch environment logs")
		return
	}

	resp := dto.EnvLogListResponse{
		Total: total,
		Logs:  make([]dto.EnvLogResponse, len(logs)),
	}
	for i, l := range logs {
		resp.Logs[i] = dto.ToEnvLogResponse(l)
	}

	writeJSON(w, http.StatusOK, resp)
}

// recordEnvLogs godoc
// @Summary      Batch record environment logs
// @Description  Records temperature and humidity for one or more zones
// @Tags         Environment
// @Accept       json
// @Produce      json
// @Param        request body []dto.RecordEnvLogRequest true "List of logs"
// @Success      204  "No Content"
// @Router       /api/v1/env/logs [post]
func (h *Handler) recordEnvLogs(w http.ResponseWriter, r *http.Request) {
	var req []dto.RecordEnvLogRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	userID, _ := r.Context().Value(middleware.CtxUserID).(int)
	
	domainLogs := make([]domain.EnvironmentLog, len(req))
	for i, r := range req {
		domainLogs[i] = r.ToDomain()
	}

	if err := h.service.EnvironmentLog.RecordLogs(r.Context(), userID, domainLogs); err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to record environment logs")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
