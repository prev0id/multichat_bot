package user

import (
	"io"
	"net/http"
	"net/url"

	"github.com/go-chi/chi/v5"

	"multichat_bot/internal/domain"
)

func (s *Service) HandleAddBanUser(w http.ResponseWriter, r *http.Request) {
	rawPlatform := chi.URLParam(r, domain.URLParamPlatform)
	platform, ok := domain.StringToPlatform[rawPlatform]
	if !ok {
		data := domain.NewTemplateSettingsData(platform, nil, errorTryLater)
		s.executeBanUsersTemplate(w, data, http.StatusBadRequest)

		return
	}

	user, ok := s.auth.IsLoggedIn(r)
	if !ok {
		data := domain.NewTemplateSettingsData(platform, nil, errorNotLoggedIn)
		s.executeBanUsersTemplate(w, data, http.StatusBadRequest)

		return
	}

	config, ok := user.Platforms[platform]
	if !ok {
		data := domain.NewTemplateSettingsData(platform, config, errorNotLoggedIn)
		s.executeBanUsersTemplate(w, data, http.StatusUnauthorized)

		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		data := domain.NewTemplateSettingsData(platform, config, errorTryLater)
		s.executeBanUsersTemplate(w, data, http.StatusBadRequest)

		return
	}

	vals, err := url.ParseQuery(string(body))
	if err != nil {
		data := domain.NewTemplateSettingsData(platform, config, errorTryLater)
		s.executeBanUsersTemplate(w, data, http.StatusBadRequest)

		return
	}

	config.DisabledUsers = config.DisabledUsers.Add(vals.Get("input"))

	if err := s.db.UpdateBannedUsers(user.ID, platform, config.DisabledUsers); err != nil {
		data := domain.NewTemplateSettingsData(platform, config, errorTryLater)
		s.executeBanUsersTemplate(w, data, http.StatusInternalServerError)

		return
	}

	data := domain.NewTemplateSettingsData(platform, config, nil)
	s.executeBanUsersTemplate(w, data, http.StatusOK)
}

func (s *Service) HandleRemoveBanUser(w http.ResponseWriter, r *http.Request) {
	rawPlatform := chi.URLParam(r, domain.URLParamPlatform)
	platform, ok := domain.StringToPlatform[rawPlatform]
	if !ok {
		data := domain.NewTemplateSettingsData(platform, nil, errorTryLater)
		s.executeBanUsersTemplate(w, data, http.StatusBadRequest)

		return
	}

	user, ok := s.auth.IsLoggedIn(r)
	if !ok {
		data := domain.NewTemplateSettingsData(platform, nil, errorNotLoggedIn)
		s.executeBanUsersTemplate(w, data, http.StatusBadRequest)

		return
	}

	config, ok := user.Platforms[platform]
	if !ok {
		data := domain.NewTemplateSettingsData(platform, config, errorNotLoggedIn)
		s.executeBanUsersTemplate(w, data, http.StatusUnauthorized)

		return
	}

	userToRemove := chi.URLParam(r, domain.URLParamUser)

	config.DisabledUsers = config.DisabledUsers.Remove(userToRemove)

	if err := s.db.UpdateBannedUsers(user.ID, platform, config.DisabledUsers); err != nil {
		data := domain.NewTemplateSettingsData(platform, config, errorTryLater)
		s.executeBanUsersTemplate(w, data, http.StatusInternalServerError)

		return
	}

	data := domain.NewTemplateSettingsData(platform, config, nil)
	s.executeBanUsersTemplate(w, data, http.StatusOK)
}

func (s *Service) executeBanUsersTemplate(w http.ResponseWriter, data *domain.TemplateSettingsData, status int) {
	if err := s.bannedUsersTmpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(status)
}
