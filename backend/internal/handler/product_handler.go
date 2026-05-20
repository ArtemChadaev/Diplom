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

// listProducts godoc
// @Summary      List products
// @Description  Returns a paginated list of products with optional filtering
// @Tags         Products
// @Produce      json
// @Param        q          query     string  false  "Search by name, generic name or SKU"
// @Param        is_jnvlp   query     bool    false  "Filter by JNVLP status"
// @Param        atc_code   query     string  false  "Filter by ATC code"
// @Param        limit      query     int     false  "Page size"
// @Param        offset     query     int     false  "Offset"
// @Success      200  {object}  dto.ProductListResponse
// @Router       /api/v1/products [get]
func (h *Handler) listProducts(w http.ResponseWriter, r *http.Request) {
	filter := domain.ProductFilter{
		Query:    r.URL.Query().Get("q"),
		ATCCode:  r.URL.Query().Get("atc_code"),
		Limit:    10,
		Offset:   0,
	}

	if jnvlpStr := r.URL.Query().Get("is_jnvlp"); jnvlpStr != "" {
		val, _ := strconv.ParseBool(jnvlpStr)
		filter.IsJNVLP = &val
	}

	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if val, err := strconv.Atoi(limitStr); err == nil {
			filter.Limit = val
		}
	}

	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		if val, err := strconv.Atoi(offsetStr); err == nil {
			filter.Offset = val
		}
	}

	products, total, err := h.service.Product.ListProducts(r.Context(), filter)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to fetch products")
		return
	}

	resp := dto.ProductListResponse{
		Total:    total,
		Products: make([]dto.ProductResponse, len(products)),
	}
	for i, p := range products {
		resp.Products[i] = dto.ToProductResponse(p)
	}

	writeJSON(w, http.StatusOK, resp)
}

// getProduct godoc
// @Summary      Get product by ID
// @Description  Returns detailed info about a single product
// @Tags         Products
// @Produce      json
// @Param        id   path      string  true  "Product UUID"
// @Success      200  {object}  dto.ProductResponse
// @Router       /api/v1/products/{id} [get]
func (h *Handler) getProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	p, err := h.service.Product.GetProduct(r.Context(), id)
	if err != nil {
		if err == domain.ErrProductNotFound {
			writeError(w, http.StatusNotFound, "Product not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "Failed to fetch product")
		return
	}

	writeJSON(w, http.StatusOK, dto.ToProductResponse(*p))
}

// createProduct godoc
// @Summary      Create product
// @Description  Creates a new product (Admin or Warehouse Manager only)
// @Tags         Products
// @Accept       json
// @Produce      json
// @Param        request body dto.CreateProductRequest true "Product data"
// @Success      201  {object}  dto.ProductResponse
// @Router       /api/v1/products [post]
func (h *Handler) createProduct(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	role, _ := r.Context().Value(middleware.CtxRole).(domain.UserRole)
	p := req.ToDomain()

	created, err := h.service.Product.CreateProduct(r.Context(), role, &p)
	if err != nil {
		if err == domain.ErrInsufficientPerms {
			writeError(w, http.StatusForbidden, "No permission to create products")
			return
		}
		if err == domain.ErrConflict {
			writeError(w, http.StatusConflict, "Product with this SKU already exists")
			return
		}
		writeError(w, http.StatusInternalServerError, "Failed to create product")
		return
	}

	writeJSON(w, http.StatusCreated, dto.ToProductResponse(*created))
}

// patchProduct godoc
// @Summary      Partial update product
// @Description  Updates specific fields of an existing product (Admin or Warehouse Manager only)
// @Tags         Products
// @Accept       json
// @Produce      json
// @Param        id      path      string  true  "Product UUID"
// @Param        request body dto.UpdateProductRequest true "Partial update data"
// @Success      200  {object}  dto.ProductResponse
// @Router       /api/v1/products/{id} [patch]
func (h *Handler) patchProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req dto.UpdateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	role, _ := r.Context().Value(middleware.CtxRole).(domain.UserRole)
	
	p, err := h.service.Product.GetProduct(r.Context(), id)
	if err != nil {
		if err == domain.ErrProductNotFound {
			writeError(w, http.StatusNotFound, "Product not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "Failed to fetch product")
		return
	}

	req.ApplyTo(p)

	updated, err := h.service.Product.UpdateProduct(r.Context(), role, p)
	if err != nil {
		if err == domain.ErrInsufficientPerms {
			writeError(w, http.StatusForbidden, "No permission to update products")
			return
		}
		writeError(w, http.StatusInternalServerError, "Failed to update product")
		return
	}

	writeJSON(w, http.StatusOK, dto.ToProductResponse(*updated))
}

// deleteProduct godoc
// @Summary      Delete (soft) product
// @Description  Marks a product as deleted (Admin only)
// @Tags         Products
// @Produce      json
// @Param        id   path      string  true  "Product UUID"
// @Success      204  "No Content"
// @Router       /api/v1/products/{id} [delete]
func (h *Handler) deleteProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	role, _ := r.Context().Value(middleware.CtxRole).(domain.UserRole)

	if err := h.service.Product.DeleteProduct(r.Context(), role, id); err != nil {
		if err == domain.ErrProductNotFound {
			writeError(w, http.StatusNotFound, "Product not found")
			return
		}
		if err == domain.ErrInsufficientPerms {
			writeError(w, http.StatusForbidden, "No permission to delete products")
			return
		}
		writeError(w, http.StatusInternalServerError, "Failed to delete product")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// runReorderCheck godoc
// @Summary      Run reorder check
// @Description  Calculates reorder point (ROP) and reorder quantities for all products (Admin or Warehouse Manager only)
// @Tags         Products
// @Produce      json
// @Success      200  {array}   domain.ROPResult
// @Router       /api/v1/products/reorder-check [get]
func (h *Handler) runReorderCheck(w http.ResponseWriter, r *http.Request) {
	role, _ := r.Context().Value(middleware.CtxRole).(domain.UserRole)
	if role != domain.RoleAdmin && role != domain.RoleWarehouseManager {
		writeError(w, http.StatusForbidden, "No permission to perform reorder check")
		return
	}

	results, err := h.service.Product.RunReorderCheckAll(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to run reorder check: "+err.Error())
		return
	}

	writeJSON(w, http.StatusOK, results)
}
