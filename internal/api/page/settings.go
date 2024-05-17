package page

import (
	"log/slog"
	"net/http"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

const (
	minLoggedInPlatforms = 2
)

type settingsTemplateData struct {
	Platforms    []platformSettings
	unauthorized bool
}

type platformSettings struct {
	Name          string
	DisabledUsers []string
	BannedWords   []string
	IsJoined      bool
}

func (s *Service) HandleSetting(w http.ResponseWriter, r *http.Request) {
	data, err := s.getSettingsData(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if data.unauthorized {
		http.Redirect(w, r, "/", http.StatusUnauthorized)
		return
	}

	tmpl, err := s.getTemplate(templateNameSettings)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, data); err != nil {
		slog.Error("handling settings: " + err.Error())
	}
}

func (s *Service) getSettingsData(r *http.Request) (settingsTemplateData, error) {
	userFromCookie, ok := s.cookieStore.GetUser(r)
	if !ok {
		return settingsTemplateData{unauthorized: true}, nil
	}

	user, err := s.db.GetUserByID(userFromCookie.ID)
	if err != nil {
		return settingsTemplateData{}, err
	}

	if len(user.Platforms) < minLoggedInPlatforms {
		return settingsTemplateData{unauthorized: false}, nil
	}

	platforms := make([]platformSettings, 0, len(user.Platforms))
	for platform, config := range user.Platforms {
		platforms = append(platforms, platformSettings{
			Name:          cases.Title(language.English).String(platform.String()),
			DisabledUsers: config.DisabledUsers,
			BannedWords:   config.BannedWords,
			IsJoined:      config.IsJoined,
		})
	}

	return settingsTemplateData{
		Platforms: platforms,
	}, nil
}
