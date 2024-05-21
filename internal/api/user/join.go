package user

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"

	"multichat_bot/internal/domain"
)

var (
	errorTryLater    = errors.New("something went wrong, try again later")
	errorNotLoggedIn = errors.New("you are not logged in")
)

func (s *Service) HandleJoin(w http.ResponseWriter, r *http.Request) {
	rawPlatform := chi.URLParam(r, domain.URLParamPlatform)
	platform, ok := domain.StringToPlatform[rawPlatform]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user, ok := s.auth.IsLoggedIn(r)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		data := domain.NewTemplateSettingsData(platform, nil, errorNotLoggedIn)
		s.executeToggleJointTemplate(w, data, http.StatusUnauthorized)
		return
	}

	cfg, ok := user.Platforms[platform]
	if !ok {
		data := domain.NewTemplateSettingsData(platform, cfg, errorNotLoggedIn)
		s.executeToggleJointTemplate(w, data, http.StatusForbidden)
		return
	}

	service, ok := s.platforms[platform]
	if !ok {
		data := domain.NewTemplateSettingsData(platform, cfg, errorTryLater)
		s.executeToggleJointTemplate(w, data, http.StatusNotImplemented)
		return
	}

	fmt.Printf("%+v\n",cfg)

	to := cfg.ID
	if platform == domain.Twitch {
		to = cfg.Channel
	}

	if err := service.Join(to); err != nil {
		slog.Error(fmt.Sprintf("service.Join[%s]: failed to leave: %v", platform, err))

		data := domain.NewTemplateSettingsData(platform, cfg, err)
		s.executeToggleJointTemplate(w, data, http.StatusInternalServerError)
		return
	}

	if err := s.db.JoinChannel(user.ID, platform); err != nil {
		slog.Error(fmt.Sprintf("s.db.JoinChannel error: %v", err))

		data := domain.NewTemplateSettingsData(platform, cfg, errorTryLater)
		s.executeToggleJointTemplate(w, data, http.StatusInternalServerError)
		return
	}

	cfg.IsJoined = true
	data := domain.NewTemplateSettingsData(platform, cfg, nil)
	s.executeToggleJointTemplate(w, data, http.StatusOK)
}

func (s *Service) HandleLeave(w http.ResponseWriter, r *http.Request) {
	rawPlatform := chi.URLParam(r, domain.URLParamPlatform)
	platform, ok := domain.StringToPlatform[rawPlatform]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user, ok := s.auth.IsLoggedIn(r)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		data := domain.NewTemplateSettingsData(platform, nil, errorNotLoggedIn)
		s.executeToggleJointTemplate(w, data, http.StatusUnauthorized)
		return
	}

	cfg, ok := user.Platforms[platform]
	if !ok {
		data := domain.NewTemplateSettingsData(platform, cfg, errorNotLoggedIn)
		s.executeToggleJointTemplate(w, data, http.StatusForbidden)
		return
	}

	service, ok := s.platforms[platform]
	if !ok {
		data := domain.NewTemplateSettingsData(platform, cfg, errorTryLater)
		s.executeToggleJointTemplate(w, data, http.StatusNotImplemented)
		return
	}

	if err := service.Leave(cfg.Channel); err != nil {
		slog.Error(fmt.Sprintf("service.Leave[%s]: failed to leave: %v", platform, err))

		data := domain.NewTemplateSettingsData(platform, cfg, err)
		s.executeToggleJointTemplate(w, data, http.StatusInternalServerError)
		return
	}

	if err := s.db.LeaveChannel(user.ID, platform); err != nil {
		slog.Error(fmt.Sprintf("s.db.LeaveChannel error: %v", err))

		data := domain.NewTemplateSettingsData(platform, cfg, errorTryLater)
		s.executeToggleJointTemplate(w, data, http.StatusInternalServerError)
		return
	}

	cfg.IsJoined = false
	data := domain.NewTemplateSettingsData(platform, cfg, nil)
	s.executeToggleJointTemplate(w, data, http.StatusOK)
}

func (s *Service) executeToggleJointTemplate(w http.ResponseWriter, data *domain.TemplateSettingsData, status int) {
	if err := s.toggleJoinTmpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(status)
}
