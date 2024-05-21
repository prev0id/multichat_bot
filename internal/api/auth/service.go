package auth

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptrace"
	"time"

	"github.com/dghubble/gologin/v2"
	"github.com/go-chi/chi/v5"
	"golang.org/x/oauth2"

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

	if platform == "twitch" {
		r = r.WithContext(prepareTwitchHTTP(r.Context()))
	}

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
	s.auth.HandleLogout(w)
	http.Redirect(w, r, "/", http.StatusFound)
}

func (s *Service) Delete(w http.ResponseWriter, r *http.Request) {
	user, ok := s.auth.IsLoggedIn(r)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if err := s.auth.HandleDelete(w, user); err != nil {
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

func prepareTwitchHTTP(ctx context.Context) context.Context {
	client := http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	ctx = context.WithValue(ctx, oauth2.HTTPClient, client)

	trace := &httptrace.ClientTrace{
		GetConn: func(hostPort string) {
			fmt.Printf("GetConn: %s", hostPort)
		},
		GotConn: func(connInfo httptrace.GotConnInfo) {
			fmt.Printf("Got Conn: %+v\n", connInfo)
		},
		DNSStart: func(info httptrace.DNSStartInfo) {
			fmt.Printf("DNS Start: %+v\n", info)
		},
		DNSDone: func(dnsInfo httptrace.DNSDoneInfo) {
			fmt.Printf("DNS Info: %+v\n", dnsInfo)
		},
		ConnectStart: func(network, addr string) {
			fmt.Printf("Connect Start: %s %s\n", network, addr)
		},
		ConnectDone: func(network, addr string, err error) {
			fmt.Printf("Connect Done: %s %s %v\n", network, addr, err)
		},
		TLSHandshakeStart: func() {
			fmt.Printf("TLS Handshake Start\n")
		},
		TLSHandshakeDone: func(state tls.ConnectionState, err error) {
			fmt.Printf("TLS Handshake Done: %v\n", err)
		},
		WroteHeaderField: func(key string, value []string) {
			fmt.Printf("Wrote Header Field: %s %v\n", key, value)
		},
		WroteRequest: func(info httptrace.WroteRequestInfo) {
			fmt.Printf("Wrote Request Info: %+v\n", info)
		},
	}
	ctx = httptrace.WithClientTrace(ctx, trace)

	ctx, _ = context.WithDeadline(ctx, time.Now().Add(time.Minute))

	return ctx
}
