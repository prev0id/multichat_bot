package main

import (
	"context"
	"flag"
	"log"
	"log/slog"
	"sync"

	"multichat_bot/internal/client/ws"
	"multichat_bot/internal/config"
)

const (
	twitchWSSAddress = "ws://irc-ws.chat.twitch.tv:80"
	localHost        = "http://localhost:3000"
	protocol         = ""

	accessToken  = "at1om39cl148k7b6vuahxj9xjri6lg"
	clientSecret = "qdiydh8hjj6v2kpcr9ol038t78493h"
	clientID     = "v59a8wp2evo747lhge8ynzrvhdnbx1"

	userNick     = "pre_void"
	botName      = "multichat_bot"
	userPassword = "deevsemen2YANDEX"
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

	twitchConfig, isExist := cfg.Clients[config.TwitchClient]
	if !isExist {
		log.Fatalf("%s config should be specified", config.TwitchClient)
	}

	twitchClient, err := ws.New(appCtx, config.TwitchClient, twitchConfig)
	if err != nil {
		log.Fatalf("error while creating twitch client: %s", err.Error())
	}

	wg := &sync.WaitGroup{}
	go func() {
		receiveMessage(appCtx, wg, twitchClient.GetMessageChannel())
	}()

	err = twitchClient.Send("CAP REQ :twitch.tv/membership twitch.tv/tags twitch.tv/commands")
	if err != nil {
		slog.Error("CAP REQ " + err.Error())
	}

	err = twitchClient.Send("PASS oauth:" + cfg.Credentials.Token)
	if err != nil {
		slog.Error("PASS " + err.Error())
	}

	err = twitchClient.Send("NICK " + cfg.Credentials.UserName)
	if err != nil {
		slog.Error("NICK " + err.Error())
	}

	err = twitchClient.Send("JOIN #pre_void")

	wg.Wait()
}

func receiveMessage(ctx context.Context, wg *sync.WaitGroup, ch <-chan string) {
	wg.Add(1)
	for {
		select {
		case <-ctx.Done():
			slog.Info("shutdown")
			wg.Done()
			return
		case <-ch:
			slog.Info("received in main")
		}
	}
}
