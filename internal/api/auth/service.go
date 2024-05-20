package auth

import (
	"context"
	"errors"
	"net/http"

	"github.com/dghubble/gologin/v2"
	"github.com/go-chi/chi/v5"

	"multichat_bot/internal/api/auth/google"
	"multichat_bot/internal/api/auth/twitch"
	"multichat_bot/internal/common/auth"
	"multichat_bot/internal/config"
	"multichat_bot/internal/database"
	"multichat_bot/internal/domain"
)

type Service struct {
	db   *database.Manager
	auth *auth.Auth

	callBack map[string]http.Handler
	login    map[string]http.Handler
}

func NewService(cfg config.Auth, db *database.Manager, authService *auth.Auth) *Service {
	s := &Service{
		db:       db,
		auth:     authService,
		callBack: make(map[string]http.Handler),
		login:    make(map[string]http.Handler),
	}

	stateConfig := gologin.DefaultCookieConfig
	if !cfg.IsProd {
		stateConfig = gologin.DebugOnlyCookieConfig
	}

	s.initGoogle(cfg, stateConfig)
	s.initTwitch(cfg, stateConfig)

	return s
}

func (s *Service) CallBack(w http.ResponseWriter, r *http.Request) {
	platform := chi.URLParam(r, domain.URLParamPlatform)

	callback, ok := s.callBack[platform]
	if !ok {
		w.WriteHeader(http.StatusNotImplemented)
		return
	}

	callback.ServeHTTP(w, r)
}

func (s *Service) Login(w http.ResponseWriter, r *http.Request) {
	platform := chi.URLParam(r, domain.URLParamPlatform)

	callback, ok := s.login[platform]
	if !ok {
		w.WriteHeader(http.StatusNotImplemented)
		return
	}

	callback.ServeHTTP(w, r)
}

func (s *Service) Logout(w http.ResponseWriter, r *http.Request) {
	_, ok := s.auth.IsLoggedIn(r)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	s.auth.HandleLogout(w)

	http.Redirect(w, r, "/", http.StatusFound)
}

func (s *Service) DeletePlatform(w http.ResponseWriter, r *http.Request) {
	rawPlatform := chi.URLParam(r, domain.URLParamPlatform)
	platform, ok := domain.StringToPlatform[rawPlatform]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user, ok := s.auth.IsLoggedIn(r)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if err := s.auth.HandleDelete(w, user, platform); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusFound)
}

func (s *Service) issueNewSession(platform domain.Platform) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		newConfig, err := getConfigFromContext(ctx, platform)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := s.auth.HandleLogin(w, r, platform, newConfig); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/", http.StatusFound)
	})
}

func getConfigFromContext(ctx context.Context, platform domain.Platform) (*domain.PlatformConfig, error) {
	switch platform {
	case domain.Twitch:
		cfg, err := twitch.PlatformInfoFromContext(ctx)
		return cfg, err
	case domain.YouTube:
		cfg, err := google.PlatformInfoFromContext(ctx)
		return cfg, err
	default:
		return nil, errors.New("Invalid platform")
	}
}
