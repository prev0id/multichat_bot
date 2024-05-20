package page

import (
	"log/slog"
	"net/http"

	"multichat_bot/internal/domain"
)

var (
	providers = map[domain.Platform]string{
		domain.Twitch:  "Twitch",
		domain.YouTube: "Google",
	}

	loggedOutUserData = &accountTemplateData{
		Platforms: []accountData{
			{
				PlatformName: domain.Twitch.String(),
				ProviderName: providers[domain.Twitch],
			},
			{
				PlatformName: domain.YouTube.String(),
				ProviderName: providers[domain.YouTube],
			},
		},
	}
)

type accountTemplateData struct {
	Platforms []accountData
}

type accountData struct {
	PlatformName string
	ProviderName string
	ChannelID    string
	Username     string
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
	user, ok := s.auth.IsLoggedIn(r)
	if !ok {
		return loggedOutUserData
	}

	data := &accountTemplateData{
		Platforms: make([]accountData, 0, len(domain.Platforms)),
	}

	for _, platform := range domain.Platforms {
		account := accountData{
			PlatformName: platform.String(),
			ProviderName: providers[platform],
		}

		if cfg, ok := user.Platforms[platform]; ok {
			account.IsLoggedIn = true
			account.Username = cfg.Channel
			account.ChannelID = cfg.ID
		}

		data.Platforms = append(data.Platforms, account)
	}

	return data
}
