package page

import (
	"log/slog"
	"net/http"

	"multichat_bot/internal/domain"
)

var providers = map[domain.Platform]string{
	domain.Twitch:  "Twitch",
	domain.YouTube: "Google",
}

type accountTemplateData struct {
	Platforms []platformData
}

type platformData struct {
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
	data := &accountTemplateData{
		Platforms: make([]platformData, 0, len(domain.Platforms)),
	}

	for platform, provider := range providers {
		info, ok := s.cookieStore.GetPlatformInfo(r, platform)
		if !ok {
			data.Platforms = append(data.Platforms, platformData{
				PlatformName: platform.String(),
				ProviderName: provider,
			})
			continue
		}

		data.Platforms = append(data.Platforms, platformData{
			PlatformName: platform.String(),
			ChannelID:    info.ChannelID,
			Username:     info.Username,
			IsLoggedIn:   true,
		})

	}

	return data
}
