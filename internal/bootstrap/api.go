package bootstrap

import (
	"net"
	"net/http"

	"multichat_bot/internal/api/server"
	"multichat_bot/internal/api/server/middleware"
	"multichat_bot/internal/api/service"
	"multichat_bot/internal/config"
	twitch "multichat_bot/internal/platforms/twitch/service"
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

	return http.ListenAndServe(address, httpServer.GetHandler())
}
