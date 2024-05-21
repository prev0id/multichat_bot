package page

import (
	"net/http"

	"multichat_bot/internal/domain"
)

type settingsTemplateData struct {
	Platforms  []*domain.TemplateSettingsData
	IsLoggedIn bool
}

func (s *Service) HandleSetting(w http.ResponseWriter, r *http.Request) {
	tmpl, err := s.getTemplate(templateNameSettings)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	user, ok := s.auth.IsLoggedIn(r)
	if !ok || len(user.Platforms) < len(domain.Platforms) {
		if err := tmpl.Execute(w, settingsTemplateData{IsLoggedIn: false}); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		return
	}

	data := s.getSettingsData(user)

	if err := tmpl.Execute(w, data); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *Service) getSettingsData(user domain.User) settingsTemplateData {
	platforms := make([]*domain.TemplateSettingsData, 0, len(user.Platforms))
	for _, platform := range domain.Platforms {
		config := user.Platforms[platform]
		platforms = append(platforms, domain.NewTemplateSettingsData(platform, config, nil))
	}

	return settingsTemplateData{Platforms: platforms, IsLoggedIn: true}
}
