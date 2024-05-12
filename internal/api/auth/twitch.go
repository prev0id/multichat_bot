package auth

import (
	"net/http"

	"github.com/dghubble/gologin/v2"
	"golang.org/x/oauth2"
	twitchOAuth2 "golang.org/x/oauth2/twitch"

	"multichat_bot/internal/api/auth/twitch"
	"multichat_bot/internal/common/cookie"
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

	s.callBack[domain.Twitch.String()] = twitch.StateHandler(stateConfig, twitch.CallbackHandler(twitchConfig, s.issueTwitchSession(), nil))
	s.login[domain.Twitch.String()] = twitch.StateHandler(stateConfig, twitch.LoginHandler(twitchConfig, nil))
	s.logout[domain.Twitch.String()] = http.HandlerFunc(s.twitchLogOut)
}

func (s *Service) issueTwitchSession() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		user, err := twitch.UserFromContext(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		session := s.cookieStore.New(cookie.TwitchSession)
		session.Set(cookie.IDKey, user.ID)
		session.Set(cookie.UsernameKey, user.DisplayName)
		session.Set(cookie.EmailKey, user.Email)
		//session.Set(cookie.AccessTokenKey, token.AccessToken)

		if err := s.cookieStore.Save(w, session); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/", http.StatusFound)
	})
}

func (s *Service) twitchLogOut(w http.ResponseWriter, r *http.Request) {
	s.cookieStore.Destroy(w, cookie.TwitchSession)
	http.Redirect(w, r, "/", http.StatusFound)
}
