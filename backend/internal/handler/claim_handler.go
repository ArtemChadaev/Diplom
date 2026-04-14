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

// listClaims godoc
// @Summary      List claims and defects
// @Description  Returns a paginated list of claims
// @Tags         Claims
// @Produce      json
// @Param        limit   query     int  false  "Page size"
// @Param        offset  query     int  false  "Offset"
// @Success      200  {object}  dto.ClaimListResponse
// @Router       /api/v1/claims [get]
func (h *Handler) listClaims(w http.ResponseWriter, r *http.Request) {
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

	claims, total, err := h.service.Claim.ListClaims(r.Context(), limit, offset)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to fetch claims")
		return
	}

	resp := dto.ClaimListResponse{
		Total:  total,
		Claims: make([]dto.ClaimResponse, len(claims)),
	}
	for i, c := range claims {
		resp.Claims[i] = dto.ToClaimResponse(c)
	}

	writeJSON(w, http.StatusOK, resp)
}

// getClaim godoc
// @Summary      Get claim by ID
// @Description  Returns detailed info about a single claim
// @Tags         Claims
// @Produce      json
// @Param        id   path      string  true  "Claim UUID"
// @Success      200  {object}  dto.ClaimResponse
// @Router       /api/v1/claims/{id} [get]
func (h *Handler) getClaim(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	c, err := h.service.Claim.GetClaim(r.Context(), id)
	if err != nil {
		if err == domain.ErrClaimNotFound {
			writeError(w, http.StatusNotFound, "Claim not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "Failed to fetch claim")
		return
	}

	writeJSON(w, http.StatusOK, dto.ToClaimResponse(*c))
}

// createClaim godoc
// @Summary      Create claim or defect
// @Description  Registers a new quality control issue or delivery defect
// @Tags         Claims
// @Accept       json
// @Produce      json
// @Param        request body dto.CreateClaimRequest true "Claim data"
// @Success      201  {object}  dto.ClaimResponse
// @Router       /api/v1/claims [post]
func (h *Handler) createClaim(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateClaimRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	userID, _ := r.Context().Value(middleware.CtxUserID).(int)
	
	c := req.ToDomain()
	c.CreatedBy = userID

	created, err := h.service.Claim.CreateClaim(r.Context(), &c)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to create claim")
		return
	}

	writeJSON(w, http.StatusCreated, dto.ToClaimResponse(*created))
}

// updateClaimStatus godoc
// @Summary      Update claim status
// @Description  Changes the resolution status of a claim (e.g. New -> Resolved)
// @Tags         Claims
// @Accept       json
// @Produce      json
// @Param        id      path      string  true  "Claim UUID"
// @Param        request body dto.UpdateClaimStatusRequest true "New status"
// @Success      204  "No Content"
// @Router       /api/v1/claims/{id}/status [patch]
func (h *Handler) updateClaimStatus(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req dto.UpdateClaimStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	role, _ := r.Context().Value(middleware.CtxRole).(domain.UserRole)

	if err := h.service.Claim.UpdateStatus(r.Context(), role, id, req.Status); err != nil {
		if err == domain.ErrClaimNotFound {
			writeError(w, http.StatusNotFound, "Claim not found")
			return
		}
		if err == domain.ErrInsufficientPerms {
			writeError(w, http.StatusForbidden, "No permission to resolve claims")
			return
		}
		writeError(w, http.StatusInternalServerError, "Failed to update claim status")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
