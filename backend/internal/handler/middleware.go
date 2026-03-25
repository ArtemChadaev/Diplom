package handler

import (
	"net/http"
	"time"

	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/ima/diplom-backend/internal/pkg/logger"
)

// loggingMiddleware enriches the context with request_id (from chi's
// RequestID middleware which must run first), then logs start/end.
func (h *Handler) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// chi's RequestID middleware sets this header automatically.
		requestID := chimiddleware.GetReqID(r.Context())
		ctx := logger.WithRequestID(r.Context(), requestID)
		r = r.WithContext(ctx)

		log := logger.FromContext(ctx)
		log.Info("request started",
			"method", r.Method,
			"path", r.URL.Path,
			"ip", r.RemoteAddr,
		)

		ww := chimiddleware.NewWrapResponseWriter(w, r.ProtoMajor)
		next.ServeHTTP(ww, r.WithContext(ctx))

		log.Info("request completed",
			"method", r.Method,
			"path", r.URL.Path,
			"status", ww.Status(),
			"duration_ms", time.Since(start).Milliseconds(),
		)
	})
}

