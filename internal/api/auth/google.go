package auth

import (
	"github.com/dghubble/gologin/v2"
	"golang.org/x/oauth2"
	googleOAuth2 "golang.org/x/oauth2/google"

	"multichat_bot/internal/api/auth/google"
	"multichat_bot/internal/config"
	"multichat_bot/internal/domain"
)

func (s *Service) initGoogle(cfg config.Auth, stateConfig gologin.CookieConfig) {
	googleConfig := &oauth2.Config{
		ClientID:     cfg.Youtube.ClientKey,
		ClientSecret: cfg.Youtube.ClientSecret,
		RedirectURL:  cfg.Youtube.CallbackURL,
		Endpoint:     googleOAuth2.Endpoint,
		Scopes:       cfg.Youtube.Scopes,
	}

	s.callBack[domain.YouTube.String()] = google.StateHandler(stateConfig, google.CallbackHandler(googleConfig, s.issueNewSession(domain.YouTube), nil))
	s.login[domain.YouTube.String()] = google.StateHandler(stateConfig, google.LoginHandler(googleConfig, nil))
}
