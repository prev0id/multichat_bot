package bootstrap

import (
	"context"
	"log/slog"

	"golang.org/x/oauth2"
	ouath2_twitch "golang.org/x/oauth2/twitch"

	"multichat_bot/internal/config"
	"multichat_bot/internal/twitch/client/irc"
	"multichat_bot/internal/twitch/processor"
	twitch "multichat_bot/internal/twitch/service"
	"multichat_bot/internal/twitch/service/message_manager"
)

func Twitch(ctx context.Context, cfg config.Twitch) (*twitch.Service, error) {
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

	ircClient.WithMessageProcessor(
		processor.New(twitchService),
	)

	if err := ircClient.Connect(ctx, cfg.IRCServer); err != nil {
		slog.Error("1", slog.StringValue(err.Error()))
		return nil, err
	}

	if err := twitchService.Connect(ctx, cfg); err != nil {
		slog.Error("2", slog.StringValue(err.Error()))
		return nil, err
	}

	return twitchService, nil
}
