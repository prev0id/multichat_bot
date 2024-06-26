package main

import (
	"context"
	"flag"
	"log"
	"log/slog"

	"multichat_bot/internal/api"
	"multichat_bot/internal/api/auth"
	"multichat_bot/internal/api/page"
	"multichat_bot/internal/api/user"
	"multichat_bot/internal/app/message_broadcaster"
	auth_service "multichat_bot/internal/common/auth"
	"multichat_bot/internal/common/cookie"
	"multichat_bot/internal/config"
	"multichat_bot/internal/database"
	"multichat_bot/internal/domain"
	"multichat_bot/internal/platform/twitch"
	"multichat_bot/internal/platform/youtube"
	"multichat_bot/internal/platform/youtube/client"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "config", "./configs/local.json", "sets the path to the config, by default ./configs/local.json")
}

func main() {
	flag.Parse()

	appCtx := context.Background()

	cfg, err := config.Parse(configPath)
	if err != nil {
		log.Fatalf("error parsing config: %v", err)
	}

	db, err := database.New(appCtx, cfg.DB)
	if err != nil {
		log.Fatalf("error connecting to database: %v", err)
	}

	broadcaster := message_broadcaster.New(db)

	ytClient, err := client.NewAdapter(appCtx, cfg.Youtube, cfg.Auth.Youtube, broadcaster.GetMessageChannel())
	if err != nil {
		log.Fatalf("error creating youtube client: %v", err)
	}

	ytService := youtube.NewService(ytClient)
	twitchService := twitch.NewService(cfg.Twitch, broadcaster.GetMessageChannel())

	broadcaster.AddPlatform(domain.Twitch, twitchService)
	broadcaster.AddPlatform(domain.YouTube, ytService)

	broadcaster.StartWorker(appCtx)

	ytClient.StartListening(appCtx)
	slog.Info("started youtube")

	if err = twitchService.Connect(); err != nil {
		log.Fatalf("error connecting to twitch: %v", err)
	}
	slog.Info("started twitch")

	cookieStore := cookie.NewStore(cfg.Cookie)
	authService := auth_service.NewAuthService(db, cookieStore)

	pageService, err := page.NewService(cfg.Auth.IsProd, authService)
	if err != nil {
		log.Fatalf("error creating page service: %v", err)
	}
	userAPI := user.NewService(db, authService).
		WithPlatformService(domain.Twitch, twitchService).
		WithPlatformService(domain.YouTube, ytService)
	authAPI := auth.NewService(cfg.Auth, db, authService)

	if err := api.Serve(cfg.API, userAPI, pageService, authAPI); err != nil {
		log.Fatalf("error starting api: %v", err)
	}
}
