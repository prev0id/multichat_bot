package page

import (
	"log/slog"
	"net/http"

	"github.com/dghubble/sessions"

	"multichat_bot/internal/common/cookie"
)

type accountTemplateData struct {
	Platforms []platformData
}

type platformData struct {
	PlatformName string
	ProviderName string
	Username     string
	Email        string
	ID           string
	Token        string
	IsLoggedIn   bool
}

func (s *Service) HandleRoot(w http.ResponseWriter, r *http.Request) {
	data := s.getRootData(r)

	tmpl, err := s.getTemplate(templateNameAccount)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, data); err != nil {
		slog.Error("handling root: " + err.Error())
	}
}

func (s *Service) getRootData(r *http.Request) *accountTemplateData {
	twitchSession, err := s.cookieStore.Get(r, cookie.TwitchSession)

	twitchData := platformData{IsLoggedIn: false}
	if err == nil {
		twitchData = getDataFromSession(twitchSession)
	}
	twitchData.ProviderName = "Twitch"
	twitchData.PlatformName = "twitch"

	googleSession, err := s.cookieStore.Get(r, cookie.GoogleSession)

	googleData := platformData{IsLoggedIn: false}
	if err == nil {
		googleData = getDataFromSession(googleSession)
	}
	googleData.ProviderName = "Google"
	googleData.PlatformName = "youtube"

	return &accountTemplateData{
		Platforms: []platformData{googleData, twitchData},
	}
}

func getDataFromSession(session *sessions.Session[string]) platformData {
	return platformData{
		Username:   session.Get(cookie.UsernameKey),
		Email:      session.Get(cookie.EmailKey),
		ID:         session.Get(cookie.IDKey),
		Token:      session.Get(cookie.AccessTokenKey),
		IsLoggedIn: true,
	}
}
