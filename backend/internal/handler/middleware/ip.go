package middleware

import (
	"net"
	"net/http"

	"github.com/ima/diplom-backend/internal/pkg/logger"
)

// InjectIPAddress is an HTTP middleware that extracts the real client IP
// from the request and stores it in the context via logger.WithIPAddress.
//
// Priority order:
//  1. X-Forwarded-For (first entry — set by reverse proxies/load balancers)
//  2. X-Real-IP (Nginx convention)
//  3. RemoteAddr (direct connection fallback)
func InjectIPAddress(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := extractIP(r)
		ctx := logger.WithIPAddress(r.Context(), ip)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func extractIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		// X-Forwarded-For may contain a comma-separated list; take the first.
		if idx := len(xff); idx > 0 {
			for i, c := range xff {
				if c == ',' {
					idx = i
					break
				}
			}
			if ip := net.ParseIP(trimSpace(xff[:idx])); ip != nil {
				return ip.String()
			}
		}
	}

	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		if ip := net.ParseIP(trimSpace(xri)); ip != nil {
			return ip.String()
		}
	}

	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}

func trimSpace(s string) string {
	start, end := 0, len(s)
	for start < end && (s[start] == ' ' || s[start] == '\t') {
		start++
	}
	for end > start && (s[end-1] == ' ' || s[end-1] == '\t') {
		end--
	}
	return s[start:end]
}
