package bootstrap

import (
	"log/slog"
	"net"
	"net/http"

	"multichat_bot/internal/api/server"
	"multichat_bot/internal/api/server/middleware"
	"multichat_bot/internal/api/service"
	"multichat_bot/internal/config"
	twitch "multichat_bot/internal/twitch/service"
)

func API(cfg config.Api, twitchService *twitch.Service) error {
	httpServer := server.New()
	httpServer.WithMiddleware(
		middleware.WithPanicRecovery,
		middleware.WithLogging,
	)

	apiService := service.New(twitchService)
	httpServer.RegisterHandler("/", apiService.Default)
	httpServer.RegisterHandler("/twitch/join", apiService.TwitchJoin)
	httpServer.RegisterHandler("/twitch/leave", apiService.TwitchLeave)

	address := net.JoinHostPort(cfg.Host, cfg.Port)
	slog.Info("starting server", slog.String("address", address))

	return http.ListenAndServe(address, httpServer.GetHandler())
}
