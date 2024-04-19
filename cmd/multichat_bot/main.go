package main

import (
	"context"
	"flag"
	"log"

	"multichat_bot/internal/app/message_broadcaster"
	"multichat_bot/internal/bootstrap"
	"multichat_bot/internal/config"
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
		log.Fatalf("error parsing config: %s", err.Error())
	}

	messageManager := message_broadcaster.New()

	twitchService, err := bootstrap.Twitch(appCtx, cfg.Twitch, messageManager.GetMessageChannel())
	if err != nil {
		log.Fatalf("can not start twitch service: %s", err.Error())
	}

	messageManager.StartWorker(appCtx)

	if err := bootstrap.API(cfg.Api, twitchService); err != nil {
		log.Fatalf("can not bootstap api service: %s", err.Error())
	}
}
