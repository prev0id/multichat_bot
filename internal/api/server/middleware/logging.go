package middleware

import (
	"log/slog"
	"net/http"
	"time"
)

func WithLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		next.ServeHTTP(w, r)

		slog.Info(
			"request processed",
			slog.String("method", r.Method),
			slog.String("uri", r.RequestURI),
			slog.String("time", time.Since(start).String()),
		)
	})
}
