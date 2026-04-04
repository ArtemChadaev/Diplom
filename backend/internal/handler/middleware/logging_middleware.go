package middleware

import (
	"context"
	"net/http"
	"time"

	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/ima/diplom-backend/internal/pkg/logger"
)

type LogData struct {
	UserID int
}

const CtxLogData contextKey = "log_data"

// RequestLogger enriches the context with request_id (from chi's
// RequestID middleware which must run first), then logs start/end.
func RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// chi's RequestID middleware sets this header automatically.
		requestID := chimiddleware.GetReqID(r.Context())
		
		logData := &LogData{}
		ctx := context.WithValue(r.Context(), CtxLogData, logData)
		ctx = logger.WithRequestID(ctx, requestID)
		r = r.WithContext(ctx)

		log := logger.FromContext(ctx)
		log.Info("request started",
			"method", r.Method,
			"path", r.URL.Path,
			"ip", r.RemoteAddr,
		)

		ww := chimiddleware.NewWrapResponseWriter(w, r.ProtoMajor)
		next.ServeHTTP(ww, r.WithContext(ctx))

		logFinal := log
		if logData.UserID != 0 {
			logFinal = logFinal.With("user_id", logData.UserID)
		}

		logFinal.Info("request completed",
			"method", r.Method,
			"path", r.URL.Path,
			"status", ww.Status(),
			"duration_ms", time.Since(start).Milliseconds(),
		)
	})
}
