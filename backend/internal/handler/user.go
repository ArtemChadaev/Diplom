package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/ima/diplom-backend/internal/handler/dto"
)

// isLoginTaken godoc
// GET /user/is-login-taken?login=employee1
// Response: {"taken": bool}
func (h *Handler) isLoginTaken(w http.ResponseWriter, r *http.Request) {
	login := r.URL.Query().Get("login")

	// Используем DTO для валидации входящих данных (можно добавить go-playground/validator позже)
	req := dto.IsLoginTakenRequest{Login: login}
	if req.Login == "" {
		http.Error(w, `{"error":"параметр login обязателен"}`, http.StatusBadRequest)
		return
	}

	taken, err := h.service.IsLoginTaken(r.Context(), req.Login)
	if err != nil {
		slog.Error("isLoginTaken: ошибка при проверке login", slog.String("login", req.Login), slog.Any("error", err))
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	resp := dto.IsLoginTakenResponse{Taken: taken}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}
