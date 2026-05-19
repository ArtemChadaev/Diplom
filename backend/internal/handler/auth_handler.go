package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/ima/diplom-backend/internal/domain"
	"github.com/ima/diplom-backend/internal/handler/dto"
	"github.com/ima/diplom-backend/internal/pkg/logger"
)

//TODO: почему нету описания путей тут всех, а только часть
// NOTE: register and login endpoints removed — authentication is OAuth-only.
// Use POST /auth/google for Google OAuth login.

// refresh godoc
// @Summary      Refresh tokens
// @Description  Reissues access token and refresh token using refresh_token cookie
// @Tags         auth
// @Produce      json
// @Success      200  {object}  dto.TokenResponse
// @Failure      401  {object}  dto.ErrorResponse  "refresh token missing or invalid session"
// @Router       /auth/refresh [post]
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
		// token expired, revoked, or invalid
		http.SetCookie(w, &http.Cookie{
			Name:   "refresh_token",
			Value:  "",
			Path:   "/",
			MaxAge: -1, // delete cookie
		})
		http.Error(w, "invalid session, please login again", http.StatusUnauthorized)
		return
	}

	sameSite := http.SameSiteLaxMode
	secure := false
	if h.cfg.Env == "production" {
		sameSite = http.SameSiteNoneMode
		secure = true
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    pair.RefreshToken,
		Path:     "/auth",
		MaxAge:   15 * 24 * 3600,
		HttpOnly: true,
		Secure:   secure,
		SameSite: sameSite,
	})

	resp := dto.TokenResponse{
		AccessToken: pair.AccessToken,
		ExpiresIn:   pair.ExpiresIn,
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}

// logout godoc
// @Summary      Log out
// @Description  Clears the refresh token cookie and revokes the active session
// @Tags         auth
// @Security     BearerAuth
// @Success      204  "No Content"
// @Router       /auth/logout [post]
func (h *Handler) logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		h.clearRefreshCookie(w)
		w.WriteHeader(http.StatusNoContent)
		return
	}

	// Hash token and find session
	hash := h.tokenSvc.HashToken(cookie.Value)
	_ = hash // session revocation happens via token hash lookup in service

	authHeader := r.Header.Get("Authorization")
	if authHeader != "" {
		const prefix = "Bearer "
		if len(authHeader) > len(prefix) {
			token := authHeader[len(prefix):]
			claims, parseErr := h.tokenSvc.ParseAccessToken(token)
			if parseErr == nil {
				revokeErr := h.service.Auth.RevokeSession(r.Context(), claims.SessionID, claims.UserID, claims.Role)
				if revokeErr != nil {
					logger.FromContext(r.Context()).Warn("failed to revoke session during logout",
						"sessionID", claims.SessionID.String(),
						"error", revokeErr.Error(),
					)
				}
			}
		}
	}

	h.clearRefreshCookie(w)
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) clearRefreshCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/auth",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		HttpOnly: true,
	})
}

// sendCode godoc
// @Summary      Send OTP code
// @Description  Sends a 6-digit OTP code to the user's email for passwordless login
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body      dto.SendCodeRequest  true  "Email address to send OTP code to"
// @Success      200   {object}  dto.SendCodeResponse
// @Failure      400   {object}  dto.ErrorResponse  "invalid request body or email is required"
// @Failure      404   {object}  dto.ErrorResponse  "user with this email not found"
// @Failure      429   {object}  dto.ErrorResponse  "too many requests, try again later"
// @Failure      500   {object}  dto.ErrorResponse  "failed to send code"
// @Router       /auth/send-code [post]
func (h *Handler) sendCode(w http.ResponseWriter, r *http.Request) {
	var req dto.SendCodeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Email == "" {
		writeError(w, http.StatusBadRequest, "email is required")
		return
	}

	if err := h.service.Auth.SendOTPCode(r.Context(), req.Email); err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			writeError(w, http.StatusNotFound, "user with this email not found")
			return
		}
		if errors.Is(err, domain.ErrOTPMaxAttempts) {
			writeError(w, http.StatusTooManyRequests, "too many requests, try again later")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to send code")
		return
	}

	resp := dto.SendCodeResponse{
		Message:   "OTP code sent successfully",
		ExpiresIn: 600, // 10 minutes
	}
	writeJSON(w, http.StatusOK, resp)
}

// verifyCode godoc
// @Summary      Verify OTP code
// @Description  Verifies the OTP code sent to the user's email. Returns an access token and sets a HTTP-only refresh token cookie.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body      dto.VerifyCodeRequest  true  "Email and OTP code to verify"
// @Success      200   {object}  dto.TokenResponse
// @Failure      400   {object}  dto.ErrorResponse  "invalid request body or missing fields"
// @Failure      401   {object}  dto.ErrorResponse  "invalid or expired code"
// @Failure      500   {object}  dto.ErrorResponse  "failed to verify code"
// @Router       /auth/verify-code [post]
func (h *Handler) verifyCode(w http.ResponseWriter, r *http.Request) {
	var req dto.VerifyCodeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Email == "" || req.Code == "" {
		writeError(w, http.StatusBadRequest, "email and code are required")
		return
	}

	meta := domain.SessionMeta{
		UserAgent: r.UserAgent(),
		IPAddress: r.RemoteAddr,
	}

	pair, err := h.service.Auth.VerifyOTPCode(r.Context(), req.Email, req.Code, meta)
	if err != nil {
		if errors.Is(err, domain.ErrOTPNotFound) || errors.Is(err, domain.ErrOTPInvalid) || errors.Is(err, domain.ErrOTPMaxAttempts) { // using ErrOTPNotFound as "invalid or expired"
			writeError(w, http.StatusUnauthorized, "invalid or expired code")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to verify code")
		return
	}

	vcSameSite := http.SameSiteLaxMode
	vcSecure := false
	if h.cfg.Env == "production" {
		vcSameSite = http.SameSiteNoneMode
		vcSecure = true
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    pair.RefreshToken,
		Path:     "/auth",
		MaxAge:   15 * 24 * 3600,
		HttpOnly: true,
		Secure:   vcSecure,
		SameSite: vcSameSite,
	})

	resp := dto.TokenResponse{
		AccessToken: pair.AccessToken,
		ExpiresIn:   pair.ExpiresIn,
	}

	writeJSON(w, http.StatusOK, resp)
}

// @Summary     Регистрация нового пользователя по email
// @Description Создаёт нового пользователя с минимальной ролью (pharmacist).
//
//	Аккаунт остаётся неактивным до тех пор, пока администратор
//	не назначит пользователю рабочую роль. Сразу после создания
//	на указанный email отправляется одноразовый OTP-код (действует
//	10 минут), который нужно подтвердить через POST /auth/verify-code.
//
// @Tags        auth
// @Accept      json
// @Produce     json
// @Param       body  body      dto.RegisterRequest   true  "Email для регистрации"
// @Success     201   {object}  dto.RegisterResponse
// @Failure     400   {object}  dto.ErrorResponse     "email is required"
// @Failure     409   {object}  dto.ErrorResponse     "email already registered"
// @Failure     500   {object}  dto.ErrorResponse     "failed to register"
// @Router      /auth/register [post]
func (h *Handler) register(w http.ResponseWriter, r *http.Request) {
	var req dto.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Email == "" {
		writeError(w, http.StatusBadRequest, "email is required")
		return
	}

	if err := h.service.Auth.RegisterByEmail(r.Context(), req.Email); err != nil {
		if errors.Is(err, domain.ErrEmailTaken) {
			writeError(w, http.StatusConflict, "email already registered")
			return
		}
		logger.FromContext(r.Context()).Error("failed to register user", "email", req.Email, "error", err)
		writeError(w, http.StatusInternalServerError, "failed to register")
		return
	}

	resp := dto.RegisterResponse{
		Message:   "registered",
		ExpiresIn: 600, // 10 minutes
	}
	writeJSON(w, http.StatusCreated, resp)
}

// @Summary     Проверка доступа к логам (Caddy forward_auth)
// @Description Используется Caddy как точка forward_auth перед проксированием на Dozzle.
//
//	Caddy перенаправляет сюда оригинальный запрос с заголовком Authorization.
//	AuthRequired проверяет JWT, RequireRole проверяет роль admin.
//	Если оба middleware прошли — возвращает 200 OK и Caddy пускает к логам.
//
// @Tags        auth
// @Produce     json
// @Security    BearerAuth
// @Success     200  {object}  map[string]string  "ok"
// @Failure     401  "missing or invalid token"
// @Failure     403  "insufficient permissions (not admin)"
// @Router      /admin/auth/check-logs-auth [get]
func (h *Handler) checkLogsAuth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}
