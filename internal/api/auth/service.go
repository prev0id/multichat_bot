package auth

import (
	"net/http"

	"github.com/dghubble/gologin/v2"
	"github.com/go-chi/chi/v5"

	"multichat_bot/internal/common/cookie"
	"multichat_bot/internal/config"
	"multichat_bot/internal/domain"
)

type Service struct {
	cookieStore *cookie.Store

	callBack map[string]http.Handler
	login    map[string]http.Handler
	logout   map[string]http.Handler
}

func NewService(cfg config.Auth, store *cookie.Store) *Service {
	s := &Service{
		cookieStore: store,
		callBack:    make(map[string]http.Handler),
		login:       make(map[string]http.Handler),
		logout:      make(map[string]http.Handler),
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
	platform := chi.URLParam(r, domain.URLParamPlatform)

	callback, ok := s.logout[platform]
	if !ok {
		w.WriteHeader(http.StatusNotImplemented)
		return
	}

	callback.ServeHTTP(w, r)
}
