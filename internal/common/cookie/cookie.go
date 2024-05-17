package cookie

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/dghubble/sessions"

	"multichat_bot/internal/config"
	"multichat_bot/internal/domain"
)

const (
	appSession    = "__Secure-access"
	twitchSession = "twitch-session"
	googleSession = "google-session"

	userIDKey      = "user_id"
	channelIDKey   = "channel_id"
	usernameKey    = "username"
	accessTokenKey = "access_token"
)

var (
	platforms = []string{appSession, twitchSession}

	EmptyUser = User{}
)

type User struct {
	Platforms map[domain.Platform]PlatformInfo
	ID        int64
	AuthToken string
}

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

func (s *Store) GetUser(req *http.Request) (User, bool) {
	app, err := s.cookieStore.Get(req, appSession)
	if err != nil {
		return EmptyUser, false
	}

	userID, err := strconv.ParseInt(app.Get(userIDKey), 10, 64)
	if err != nil {
		return EmptyUser, false
	}

	user := User{
		ID:        userID,
		Platforms: make(map[domain.Platform]PlatformInfo, len(platforms)),
	}

	for _, platformName := range platforms {
		session, err := s.cookieStore.Get(req, platformName)
		if err != nil {
			continue
		}

		user.Platforms[sessionToPlatform(platformName)] = convertSessionToPlatformInfo(session)
	}

	return user, true
}

func (s *Store) GetPlatformInfo(req *http.Request, platform domain.Platform) (PlatformInfo, bool) {
	session, err := s.cookieStore.Get(req, getPlatformSessionName(platform))
	if err != nil {
		return PlatformInfo{}, false
	}

	return convertSessionToPlatformInfo(session), true
}

func (s *Store) SaveDomainUser(w http.ResponseWriter, user domain.User) error {
	for platform, cfg := range user.Platforms {
		sessionName := getPlatformSessionName(platform)

		session := s.cookieStore.New(sessionName)

		session.Set(channelIDKey, cfg.ID)
		session.Set(usernameKey, cfg.Channel)
		session.Set(accessTokenKey, cfg.AccessToken)

		if err := s.cookieStore.Save(w, session); err != nil {
			return fmt.Errorf("save cookie: %w", err)
		}
	}

	session := s.cookieStore.New(appSession)
	session.Set(userIDKey, strconv.FormatInt(user.ID, 10))
	if err := s.cookieStore.Save(w, session); err != nil {
		return fmt.Errorf("save cookie: %w", err)
	}

	return nil
}

func (s *Store) ClearAll(w http.ResponseWriter) {
	for _, platformName := range platforms {
		s.cookieStore.Destroy(w, platformName)
	}

	s.cookieStore.Destroy(w, appSession)
}

func (s *Store) DestroyPlatformSession(w http.ResponseWriter, platform domain.Platform) {
	s.cookieStore.Destroy(w, getPlatformSessionName(platform))
}

func (s *Store) DestroyAppSession(w http.ResponseWriter) {
	s.cookieStore.Destroy(w, appSession)
}

func convertSessionToPlatformInfo(session *sessions.Session[string]) PlatformInfo {
	return PlatformInfo{
		ChannelID:   session.Get(channelIDKey),
		Username:    session.Get(usernameKey),
		AccessToken: session.Get(accessTokenKey),
	}
}

func getPlatformSessionName(platform domain.Platform) string {
	switch platform {
	case domain.YouTube:
		return googleSession
	case domain.Twitch:
		return twitchSession
	}

	return ""
}

func sessionToPlatform(session string) domain.Platform {
	switch session {
	case googleSession:
		return domain.YouTube
	case twitchSession:
		return domain.Twitch
	}
	return ""
}
