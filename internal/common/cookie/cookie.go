package cookie

import (
	"fmt"
	"net/http"

	"github.com/dghubble/sessions"

	"multichat_bot/internal/config"
)

const (
	appSession     = "access"
	accessTokenKey = "access_token"
)

type PlatformInfo struct {
	ChannelID   string
	Username    string
	AccessToken string
}

type Store struct {
	config      *sessions.CookieConfig
	cookieStore sessions.Store[string]
}

func NewStore(cfg config.CookieStore) *Store {
	cookieConfig := sessions.DefaultCookieConfig
	if !cfg.IsProd {
		cookieConfig = sessions.DefaultCookieConfig
	}

	return &Store{
		config:      cookieConfig,
		cookieStore: sessions.NewCookieStore[string](cookieConfig, []byte(cfg.SessionKey)),
	}
}

func (s *Store) GetAccessToken(req *http.Request) string {
	app, err := s.cookieStore.Get(req, appSession)
	if err != nil {
		return ""
	}

	return app.Get(accessTokenKey)
}

func (s *Store) SaveAccessToken(w http.ResponseWriter, token string) error {
	session := s.cookieStore.New(appSession)
	session.Set(accessTokenKey, token)

	if err := s.cookieStore.Save(w, session); err != nil {
		return fmt.Errorf("save cookie: %w", err)
	}

	return nil
}

func (s *Store) DestroyAppSession(w http.ResponseWriter) {
	s.cookieStore.Destroy(w, appSession)
}
