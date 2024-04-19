package middleware

import (
	"log/slog"
	"net/http"
	"runtime/debug"

	"multichat_bot/internal/domain/logger"
)

func WithPanicRecovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

				slog.Error(
					"recovered",
					slog.String(logger.Stack, string(debug.Stack())),
				)
			}
		}()

		next.ServeHTTP(w, r)
	})
}
