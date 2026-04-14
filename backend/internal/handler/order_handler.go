package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/ima/diplom-backend/internal/domain"
	"github.com/ima/diplom-backend/internal/handler/dto"
	"github.com/ima/diplom-backend/internal/handler/middleware"
)

// listOrders godoc
// @Summary      List shipping orders
// @Description  Returns a paginated list of orders, sorted by priority and date
// @Tags         Orders
// @Produce      json
// @Param        limit   query     int  false  "Page size"
// @Param        offset  query     int  false  "Offset"
// @Success      200  {object}  dto.OrderListResponse
// @Router       /api/v1/orders [get]
func (h *Handler) listOrders(w http.ResponseWriter, r *http.Request) {
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

	orders, total, err := h.service.Order.ListOrders(r.Context(), limit, offset)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to fetch orders")
		return
	}

	resp := dto.OrderListResponse{
		Total:  total,
		Orders: make([]dto.OrderResponse, len(orders)),
	}
	for i, o := range orders {
		resp.Orders[i] = dto.ToOrderResponse(o)
	}

	writeJSON(w, http.StatusOK, resp)
}

// getOrder godoc
// @Summary      Get order by ID
// @Description  Returns detailed info about a single order including items
// @Tags         Orders
// @Produce      json
// @Param        id   path      string  true  "Order UUID"
// @Success      200  {object}  dto.OrderResponse
// @Router       /api/v1/orders/{id} [get]
func (h *Handler) getOrder(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	o, err := h.service.Order.GetOrder(r.Context(), id)
	if err != nil {
		if err == domain.ErrOrderNotFound {
			writeError(w, http.StatusNotFound, "Order not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "Failed to fetch order")
		return
	}

	writeJSON(w, http.StatusOK, dto.ToOrderResponse(*o))
}

// createOrder godoc
// @Summary      Create shipping order
// @Description  Registers a new order for subsequent assembly and shipment
// @Tags         Orders
// @Accept       json
// @Produce      json
// @Param        request body dto.CreateOrderRequest true "Order data"
// @Success      201  {object}  dto.OrderResponse
// @Router       /api/v1/orders [post]
func (h *Handler) createOrder(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	userID, _ := r.Context().Value(middleware.CtxUserID).(int)
	
	o := req.ToDomain()
	o.CreatedBy = userID
	o.OrderNumber = fmt.Sprintf("ORD-%d", time.Now().Unix()) // Simple generator

	created, err := h.service.Order.CreateOrder(r.Context(), &o)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to create order")
		return
	}

	writeJSON(w, http.StatusCreated, dto.ToOrderResponse(*created))
}

// updateOrderStatus godoc
// @Summary      Update order status
// @Description  Changes the status of an order (e.g. Assembling -> Assembled)
// @Tags         Orders
// @Accept       json
// @Produce      json
// @Param        id      path      string  true  "Order UUID"
// @Param        request body dto.UpdateOrderStatusRequest true "New status"
// @Success      204  "No Content"
// @Router       /api/v1/orders/{id}/status [patch]
func (h *Handler) updateOrderStatus(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req dto.UpdateOrderStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	role, _ := r.Context().Value(middleware.CtxRole).(domain.UserRole)

	if err := h.service.Order.UpdateStatus(r.Context(), role, id, req.Status); err != nil {
		if err == domain.ErrOrderNotFound {
			writeError(w, http.StatusNotFound, "Order not found")
			return
		}
		if err == domain.ErrInsufficientPerms {
			writeError(w, http.StatusForbidden, "No permission to update status")
			return
		}
		writeError(w, http.StatusInternalServerError, "Failed to update order status")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
