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

// listSuppliers godoc
// @Summary      List suppliers
// @Description  Returns a paginated list of suppliers
// @Tags         Suppliers
// @Produce      json
// @Param        limit      query     int     false  "Page size"
// @Param        offset     query     int     false  "Offset"
// @Success      200  {object}  dto.SupplierListResponse
// @Router       /api/v1/suppliers [get]
func (h *Handler) listSuppliers(w http.ResponseWriter, r *http.Request) {
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

	suppliers, total, err := h.service.Supplier.ListSuppliers(r.Context(), limit, offset)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to fetch suppliers")
		return
	}

	resp := dto.SupplierListResponse{
		Total:    total,
		Suppliers: make([]dto.SupplierResponse, len(suppliers)),
	}
	for i, s := range suppliers {
		resp.Suppliers[i] = dto.ToSupplierResponse(s)
	}

	writeJSON(w, http.StatusOK, resp)
}

// getSupplier godoc
// @Summary      Get supplier by ID
// @Description  Returns detailed info about a single supplier
// @Tags         Suppliers
// @Produce      json
// @Param        id   path      string  true  "Supplier UUID"
// @Success      200  {object}  dto.SupplierResponse
// @Router       /api/v1/suppliers/{id} [get]
func (h *Handler) getSupplier(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	s, err := h.service.Supplier.GetSupplier(r.Context(), id)
	if err != nil {
		if err == domain.ErrSupplierNotFound {
			writeError(w, http.StatusNotFound, "Supplier not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "Failed to fetch supplier")
		return
	}

	writeJSON(w, http.StatusOK, dto.ToSupplierResponse(*s))
}

// createSupplier godoc
// @Summary      Create supplier
// @Description  Creates a new supplier (Admin or Warehouse Manager only)
// @Tags         Suppliers
// @Accept       json
// @Produce      json
// @Param        request body dto.CreateSupplierRequest true "Supplier data"
// @Success      201  {object}  dto.SupplierResponse
// @Router       /api/v1/suppliers [post]
func (h *Handler) createSupplier(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateSupplierRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	role, _ := r.Context().Value(middleware.CtxRole).(domain.UserRole)
	s := req.ToDomain()

	created, err := h.service.Supplier.CreateSupplier(r.Context(), role, &s)
	if err != nil {
		if err == domain.ErrInsufficientPerms {
			writeError(w, http.StatusForbidden, "No permission to create suppliers")
			return
		}
		if err == domain.ErrConflict {
			writeError(w, http.StatusConflict, "Supplier with this INN already exists")
			return
		}
		writeError(w, http.StatusInternalServerError, "Failed to create supplier")
		return
	}

	writeJSON(w, http.StatusCreated, dto.ToSupplierResponse(*created))
}

// patchSupplier godoc
// @Summary      Partial update supplier
// @Description  Updates specific fields of an existing supplier (Admin or Warehouse Manager only)
// @Tags         Suppliers
// @Accept       json
// @Produce      json
// @Param        id      path      string  true  "Supplier UUID"
// @Param        request body dto.UpdateSupplierRequest true "Partial update data"
// @Success      200  {object}  dto.SupplierResponse
// @Router       /api/v1/suppliers/{id} [patch]
func (h *Handler) patchSupplier(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req dto.UpdateSupplierRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	role, _ := r.Context().Value(middleware.CtxRole).(domain.UserRole)
	
	s, err := h.service.Supplier.GetSupplier(r.Context(), id)
	if err != nil {
		if err == domain.ErrSupplierNotFound {
			writeError(w, http.StatusNotFound, "Supplier not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "Failed to fetch supplier")
		return
	}

	req.ApplyTo(s)

	updated, err := h.service.Supplier.UpdateSupplier(r.Context(), role, s)
	if err != nil {
		if err == domain.ErrInsufficientPerms {
			writeError(w, http.StatusForbidden, "No permission to update suppliers")
			return
		}
		writeError(w, http.StatusInternalServerError, "Failed to update supplier")
		return
	}

	writeJSON(w, http.StatusOK, dto.ToSupplierResponse(*updated))
}

// deleteSupplier godoc
// @Summary      Delete supplier
// @Description  Permanently deletes a supplier (Admin only)
// @Tags         Suppliers
// @Produce      json
// @Param        id   path      string  true  "Supplier UUID"
// @Success      204  "No Content"
// @Router       /api/v1/suppliers/{id} [delete]
func (h *Handler) deleteSupplier(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	role, _ := r.Context().Value(middleware.CtxRole).(domain.UserRole)

	if err := h.service.Supplier.DeleteSupplier(r.Context(), role, id); err != nil {
		if err == domain.ErrSupplierNotFound {
			writeError(w, http.StatusNotFound, "Supplier not found")
			return
		}
		if err == domain.ErrInsufficientPerms {
			writeError(w, http.StatusForbidden, "No permission to delete suppliers")
			return
		}
		writeError(w, http.StatusInternalServerError, "Failed to delete supplier")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
