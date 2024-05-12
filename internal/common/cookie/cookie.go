package cookie

import (
	"net/http"

	"github.com/dghubble/sessions"

	"multichat_bot/internal/config"
)

type SessionName string

const (
	TwitchSession SessionName = "twitch-session"
	GoogleSession SessionName = "google-session"

	IDKey          = "id"
	UsernameKey    = "username"
	EmailKey       = "email"
	AccessTokenKey = "access_token"
)

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

func (s *Store) New(name SessionName) *sessions.Session[string] {
	return s.cookieStore.New(string(name))
}

func (s *Store) Get(req *http.Request, name SessionName) (*sessions.Session[string], error) {
	session, err := s.cookieStore.Get(req, string(name))
	if err != nil {
		return nil, err
	}

	return session, nil
}

func (s *Store) Save(w http.ResponseWriter, session *sessions.Session[string]) error {
	return s.cookieStore.Save(w, session)
}

func (s *Store) Destroy(w http.ResponseWriter, name SessionName) {
	s.cookieStore.Destroy(w, string(name))
}
