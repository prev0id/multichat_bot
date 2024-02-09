package middleware

import (
	"log/slog"
	"net/http"
	"runtime/debug"
)

func WithPanicRecovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

				slog.Error(
					"recovered",
					slog.String("stack", string(debug.Stack())),
				)
			}
		}()

		next.ServeHTTP(w, r)
	})
}
