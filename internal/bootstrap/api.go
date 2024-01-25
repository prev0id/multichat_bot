package bootstrap

import (
	"log"
	"log/slog"
	"net"
	"net/http"

	"github.com/go-chi/chi/v5"

	"multichat_bot/internal/api"
	desc "multichat_bot/internal/api/gen"
	"multichat_bot/internal/config"
	twitch "multichat_bot/internal/twitch/service"
)

func API(cfg config.Api, twitchService *twitch.Service) {
	server := api.NewServer(twitchService)

	handler := desc.NewStrictHandler(server, nil)

	router := chi.NewRouter()

	s := &http.Server{
		Handler: desc.HandlerFromMux(handler, router),
		Addr:    net.JoinHostPort(cfg.Host, cfg.Port),
	}

	slog.Info("starting http server", slog.String("PATH", s.Addr))
	if err := s.ListenAndServe(); err != nil {
		log.Fatalf("error while serving api: %s", err.Error())
	}
}
