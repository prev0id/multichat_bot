package auth

import (
	"errors"
	"net/http"

	"github.com/dghubble/gologin/v2"
	"github.com/go-chi/chi/v5"

	"multichat_bot/internal/api/auth/google"
	"multichat_bot/internal/api/auth/twitch"
	"multichat_bot/internal/common/cookie"
	"multichat_bot/internal/config"
	"multichat_bot/internal/database"
	"multichat_bot/internal/database/async_cache"
	"multichat_bot/internal/domain"
)

type Service struct {
	cookieStore *cookie.Store

	db *database.Manager

	callBack map[string]http.Handler
	login    map[string]http.Handler
}

var (
	errForbiddenError = errors.New("This account already linked with other profile")
)

func NewService(cfg config.Auth, store *cookie.Store, db *database.Manager) *Service {
	s := &Service{
		cookieStore: store,
		db:          db,
		callBack:    make(map[string]http.Handler),
		login:       make(map[string]http.Handler),
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
	platformName := domain.StringToPlatform[platform]

	s.cookieStore.DestroyPlatformSession(w, platformName)

	http.Redirect(w, r, "/", http.StatusFound)
}

func (s *Service) DeleteAccount(w http.ResponseWriter, r *http.Request) {
	platform := chi.URLParam(r, domain.URLParamPlatform)

	user, ok := s.cookieStore.GetUser(r)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	platformName := domain.StringToPlatform[platform]

	if err := s.db.DeleteUserPlatform(user.ID, platformName); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	s.cookieStore.DestroyPlatformSession(w, platformName)
	if len(user.Platforms) <= 1 {
		s.cookieStore.DestroyPlatformSession(w, platformName)
	}
	http.Redirect(w, r, "/", http.StatusFound)
}

func (s *Service) issueNewSession(platform domain.Platform) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		var (
			info cookie.PlatformInfo
			err  error
		)

		switch platform {
		case domain.Twitch:
			info, err = twitch.PlatformInfoFromContext(ctx)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		case domain.YouTube:
			info, err = google.PlatformInfoFromContext(ctx)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		cookieUser, ok := s.cookieStore.GetUser(r)
		if !ok {
			s.cookieStore.ClearAll(w)
		}

		dbUser, err := s.getUserFromDB(platform, cookieUser.ID, info.ChannelID)
		if err != nil {
			if errors.Is(err, errForbiddenError) {
				http.Error(w, err.Error(), http.StatusForbidden)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if cookieUser.ID != 0 && cookieUser.ID != dbUser.ID {
			http.Error(w, "This account already linked with other profile", http.StatusForbidden)
			return
		}

		dbUser.Platforms[platform] = mergeConfigs(dbUser.Platforms[platform], info)

		if err := s.db.UpdatePlatform(dbUser.ID, platform, dbUser.Platforms[platform]); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := s.cookieStore.SaveDomainUser(w, dbUser); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/", http.StatusFound)
	})
}

func (s *Service) getUserFromDB(platform domain.Platform, userID int64, channelID string) (domain.User, error) {
	byChannel, err := s.db.GetUserByChannel(platform, channelID)
	if !errors.Is(err, async_cache.ErrNotFound) {
		return domain.User{}, err
	}
	if err != nil {
		byChannel, err = s.db.NewUser()
		if err != nil {
			return domain.User{}, err
		}
	}

	if userID == 0 {
		return byChannel, nil
	}

	byID, err := s.db.GetUserByID(userID)
	if err != nil {
		return domain.User{}, err
	}

	if byID.ID != byChannel.ID {
		return domain.User{}, errForbiddenError
	}

	return byID, nil
}

func mergeConfigs(
	platformConfig *domain.PlatformConfig,
	info cookie.PlatformInfo,
) *domain.PlatformConfig {
	if platformConfig == nil {
		platformConfig = &domain.PlatformConfig{}
	}

	platformConfig.ID = info.ChannelID
	platformConfig.Channel = info.Username
	platformConfig.AccessToken = info.AccessToken

	return platformConfig
}
