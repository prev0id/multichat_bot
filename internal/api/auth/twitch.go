package auth

import (
	"github.com/dghubble/gologin/v2"
	"golang.org/x/oauth2"
	twitchOAuth2 "golang.org/x/oauth2/twitch"

	"multichat_bot/internal/api/auth/twitch"
	"multichat_bot/internal/config"
	"multichat_bot/internal/domain"
)

func (s *Service) initTwitch(cfg config.Auth, stateConfig gologin.CookieConfig) {
	twitchConfig := &oauth2.Config{
		ClientID:     cfg.Twitch.ClientKey,
		ClientSecret: cfg.Twitch.ClientSecret,
		RedirectURL:  cfg.Twitch.CallbackURL,
		Scopes:       cfg.Twitch.Scopes,
		Endpoint:     twitchOAuth2.Endpoint,
	}

	s.callBack[domain.Twitch.String()] = twitch.StateHandler(stateConfig, twitch.CallbackHandler(twitchConfig, s.issueNewSession(domain.Twitch), nil))
	s.login[domain.Twitch.String()] = twitch.StateHandler(stateConfig, twitch.LoginHandler(twitchConfig, nil))
}
