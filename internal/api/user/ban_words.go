package user

import (
	"io"
	"net/http"
	"net/url"

	"github.com/go-chi/chi/v5"

	"multichat_bot/internal/domain"
)

func (s *Service) HandleAddBanWord(w http.ResponseWriter, r *http.Request) {
	rawPlatform := chi.URLParam(r, domain.URLParamPlatform)
	platform, ok := domain.StringToPlatform[rawPlatform]
	if !ok {
		data := domain.NewTemplateSettingsData(platform, nil, errorTryLater)
		s.executeBanWordsTemplate(w, data, http.StatusBadRequest)

		return
	}

	user, ok := s.auth.IsLoggedIn(r)
	if !ok {
		data := domain.NewTemplateSettingsData(platform, nil, errorNotLoggedIn)
		s.executeBanWordsTemplate(w, data, http.StatusBadRequest)

		return
	}

	config, ok := user.Platforms[platform]
	if !ok {
		data := domain.NewTemplateSettingsData(platform, config, errorNotLoggedIn)
		s.executeBanWordsTemplate(w, data, http.StatusUnauthorized)

		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		data := domain.NewTemplateSettingsData(platform, config, errorTryLater)
		s.executeBanWordsTemplate(w, data, http.StatusBadRequest)

		return
	}

	vals, err := url.ParseQuery(string(body))
	if err != nil {
		data := domain.NewTemplateSettingsData(platform, config, errorTryLater)
		s.executeBanWordsTemplate(w, data, http.StatusBadRequest)

		return
	}

	config.BannedWords = config.BannedWords.Add(vals.Get("input"))

	if err := s.db.UpdateBannedWords(user.ID, platform, config.BannedWords); err != nil {
		data := domain.NewTemplateSettingsData(platform, config, errorTryLater)
		s.executeBanWordsTemplate(w, data, http.StatusInternalServerError)

		return
	}

	data := domain.NewTemplateSettingsData(platform, config, nil)
	s.executeBanWordsTemplate(w, data, http.StatusOK)

}

func (s *Service) HandleRemoveBanWord(w http.ResponseWriter, r *http.Request) {
	rawPlatform := chi.URLParam(r, domain.URLParamPlatform)
	platform, ok := domain.StringToPlatform[rawPlatform]
	if !ok {
		data := domain.NewTemplateSettingsData(platform, nil, errorTryLater)
		s.executeBanWordsTemplate(w, data, http.StatusBadRequest)

		return
	}

	user, ok := s.auth.IsLoggedIn(r)
	if !ok {
		data := domain.NewTemplateSettingsData(platform, nil, errorNotLoggedIn)
		s.executeBanWordsTemplate(w, data, http.StatusBadRequest)

		return
	}

	config, ok := user.Platforms[platform]
	if !ok {
		data := domain.NewTemplateSettingsData(platform, config, errorNotLoggedIn)
		s.executeBanWordsTemplate(w, data, http.StatusUnauthorized)

		return
	}

	wordToRemove := chi.URLParam(r, domain.URLParamWord)

	config.BannedWords = config.BannedWords.Remove(wordToRemove)

	if err := s.db.UpdateBannedWords(user.ID, platform, config.BannedWords); err != nil {
		data := domain.NewTemplateSettingsData(platform, config, errorTryLater)
		s.executeBanWordsTemplate(w, data, http.StatusInternalServerError)

		return
	}

	data := domain.NewTemplateSettingsData(platform, config, nil)
	s.executeBanWordsTemplate(w, data, http.StatusOK)
}

func (s *Service) executeBanWordsTemplate(w http.ResponseWriter, data *domain.TemplateSettingsData, status int) {
	if err := s.bannedWordsTmpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(status)
}
