package bootstrap

import (
	"context"
	"log/slog"

	"golang.org/x/oauth2"
	ouath2_twitch "golang.org/x/oauth2/twitch"

	"multichat_bot/internal/config"
	"multichat_bot/internal/domain"
	"multichat_bot/internal/domain/logger"
	"multichat_bot/internal/platforms/twitch/client/irc"
	"multichat_bot/internal/platforms/twitch/processor"
	twitch "multichat_bot/internal/platforms/twitch/service"
	"multichat_bot/internal/platforms/twitch/service/message_manager"
)

func Twitch(ctx context.Context, cfg config.Twitch, broadcast chan<- *domain.Message) (*twitch.Service, error) {
	token := &oauth2.Token{
		RefreshToken: cfg.Oauth.RefreshToken,
	}
	oauthConfig := &oauth2.Config{
		ClientID:     cfg.Oauth.ClientID,
		ClientSecret: cfg.Oauth.ClientSecret,
		RedirectURL:  cfg.Oauth.RedirectURL,
		Scopes:       cfg.Oauth.Scopes,
		Endpoint:     ouath2_twitch.Endpoint,
	}

	ircClient := irc.NewClient()

	twitchService := twitch.New(
		message_manager.New(ircClient),
		oauthConfig.TokenSource(ctx, token),
	)

	ircClient.WithMessageProcessor(processor.New(twitchService, broadcast))

	if err := ircClient.StartWorker(ctx, cfg.IRCServer); err != nil {
		slog.Error("unable to connect to twitch irc server", slog.String(logger.Error, err.Error()))
		return nil, err
	}

	if err := twitchService.Connect(cfg); err != nil {
		slog.Error("unable to bootstrap twitch service", slog.String(logger.Error, err.Error()))
		return nil, err
	}

	return twitchService, nil
}
