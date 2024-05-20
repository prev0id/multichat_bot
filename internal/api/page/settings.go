package page

import (
	"log/slog"
	"net/http"

	"multichat_bot/internal/domain"
)

type settingsTemplateData struct {
	Platforms []*domain.TemplateSettingsData
}

func (s *Service) HandleSetting(w http.ResponseWriter, r *http.Request) {
	user, ok := s.auth.IsLoggedIn(r)
	if !ok || len(user.Platforms) < len(domain.Platforms) {
		http.Redirect(w, r, "/", http.StatusUnauthorized)
		return
	}

	data := s.getSettingsData(user)

	tmpl, err := s.getTemplate(templateNameSettings)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, data); err != nil {
		slog.Error("handling settings: " + err.Error())
	}
}

func (s *Service) getSettingsData(user domain.User) settingsTemplateData {

	platforms := make([]*domain.TemplateSettingsData, 0, len(user.Platforms))
	for _, platform := range domain.Platforms {
		config := user.Platforms[platform]
		platforms = append(platforms, domain.NewTemplateSettingsData(platform, config, nil))
	}

	return settingsTemplateData{Platforms: platforms}
}
