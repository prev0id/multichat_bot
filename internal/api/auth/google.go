package auth

import (
	"fmt"
	"net/http"

	"github.com/dghubble/gologin/v2"
	"github.com/dghubble/gologin/v2/google"
	oauth2Login "github.com/dghubble/gologin/v2/oauth2"
	"golang.org/x/oauth2"
	googleOAuth2 "golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"

	"multichat_bot/internal/common/cookie"
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

	s.callBack[domain.YouTube.String()] = google.StateHandler(stateConfig, google.CallbackHandler(googleConfig, s.googleSuccessCallBack(googleConfig), nil))
	s.login[domain.YouTube.String()] = google.StateHandler(stateConfig, google.LoginHandler(googleConfig, nil))
	s.logout[domain.YouTube.String()] = http.HandlerFunc(s.googleLogOut)
}

func (s *Service) googleSuccessCallBack(cfg *oauth2.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		token, err := oauth2Login.TokenFromContext(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		youtubeService, err := youtube.NewService(ctx, option.WithHTTPClient(cfg.Client(ctx, token)))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		resp, err := youtubeService.Channels.List([]string{"id", "snippet"}).Mine(true).Do()
		fmt.Println(err)
		fmt.Println(resp.Items[0].Id)
		fmt.Printf("%+v\n", resp.Items[0].Snippet)

		googleUser, err := google.UserFromContext(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		session := s.cookieStore.New(cookie.GoogleSession)
		session.Set(cookie.IDKey, googleUser.Id)
		session.Set(cookie.UsernameKey, googleUser.Name)
		session.Set(cookie.EmailKey, googleUser.Email)
		//session.Set(cookie.AccessTokenKey, token.AccessToken)

		if err := s.cookieStore.Save(w, session); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/", http.StatusFound)
	})
}

func (s *Service) googleLogOut(w http.ResponseWriter, r *http.Request) {
	s.cookieStore.Destroy(w, cookie.GoogleSession)
	http.Redirect(w, r, "/", http.StatusFound)
}
