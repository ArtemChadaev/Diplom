package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

// isEmailTaken godoc
// GET /user/is-email-taken?email=user@example.com
// Response: {"taken": bool}
func (h *Handler) isEmailTaken(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	if email == "" {
		http.Error(w, `{"error":"параметр email обязателен"}`, http.StatusBadRequest)
		return
	}

	taken, err := h.service.IsEmailTaken(r.Context(), email)
	if err != nil {
		slog.Error("isEmailTaken: ошибка при проверке email", slog.String("email", email), slog.Any("error", err))
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]bool{"taken": taken})
}
