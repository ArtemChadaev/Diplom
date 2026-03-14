package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/ima/diplom-backend/internal/domain"
	"github.com/ima/diplom-backend/internal/handler/dto"
)

func (h *Handler) register(w http.ResponseWriter, r *http.Request) {
	var req dto.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	// Validate login, email, password...
	if req.Login == "" || req.Password == "" {
		http.Error(w, "login and password required", http.StatusBadRequest)
		return
	}

	input := domain.RegisterInput{
		Login:    req.Login,
		Email:    req.Email,
		Password: req.Password,
	}

	_, err := h.service.Auth.Register(r.Context(), input)
	if err != nil {
		if err == domain.ErrLoginTaken || err == domain.ErrEmailTaken {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}
		http.Error(w, "registration failed", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"message": "account pending admin approval"}`))
}

func (h *Handler) login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	userAgent := r.UserAgent()
	ip := r.RemoteAddr // For proxies, you'd extract from X-Forwarded-For

	pair, err := h.service.Auth.LoginWithPassword(r.Context(), req.Login, req.Password, userAgent, ip)
	if err != nil {
		if err == domain.ErrUserUnverified || err == domain.ErrUserBlocked {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	// Set HttpOnly refresh token cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    pair.RefreshToken,
		Path:     "/",
		Expires:  time.Now().Add(15 * 24 * time.Hour), // 15d
		HttpOnly: true,
		Secure:   false, // Set true in production with HTTPS
		SameSite: http.SameSiteStrictMode,
	})

	resp := dto.TokenResponse{
		AccessToken: pair.AccessToken,
		ExpiresIn:   pair.ExpiresIn,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *Handler) refresh(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		http.Error(w, "refresh token missing", http.StatusUnauthorized)
		return
	}

	meta := domain.SessionMeta{
		UserAgent: r.UserAgent(),
		IPAddress: r.RemoteAddr,
	}

	pair, err := h.service.Auth.RefreshTokens(r.Context(), cookie.Value, meta)
	if err != nil {
		// e.g. token expired, revoked, IP mismatch
		http.SetCookie(w, &http.Cookie{
			Name:   "refresh_token",
			Value:  "",
			Path:   "/",
			MaxAge: -1, // delete cookie
		})
		http.Error(w, "invalid session, please login again", http.StatusUnauthorized)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    pair.RefreshToken,
		Path:     "/",
		MaxAge:   15 * 24 * 3600,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteStrictMode,
	})

	resp := dto.TokenResponse{
		AccessToken: pair.AccessToken,
		ExpiresIn:   pair.ExpiresIn,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *Handler) logout(w http.ResponseWriter, r *http.Request) {
	// Not fully implemented revocation in the handler context yet,
	// but normally you'd read the access token to get userID and session ID
	// Or just do a client-side logout by invalidating cookie.
	http.SetCookie(w, &http.Cookie{
		Name:   "refresh_token",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})
	w.WriteHeader(http.StatusNoContent)
}
