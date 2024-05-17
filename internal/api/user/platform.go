package user

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"multichat_bot/internal/domain"
)

func (s *Service) HandleJoin(w http.ResponseWriter, r *http.Request) {
	user, ok := s.cookies.GetUser(r)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	dbUser, err := s.db.GetUserByID(user.ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	rawPlatform := chi.URLParam(r, domain.URLParamPlatform)
	platform, ok := domain.StringToPlatform[rawPlatform]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	cfg, ok := dbUser.Platforms[platform]
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	service, ok := s.platforms[platform]
	if !ok {
		w.WriteHeader(http.StatusNotImplemented)
		return
	}

	if err := s.db.JoinChannel(user.ID, platform); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := service.Join(cfg.Channel); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *Service) HandleLeave(w http.ResponseWriter, r *http.Request) {
	user, ok := s.cookies.GetUser(r)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	dbUser, err := s.db.GetUserByID(user.ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	rawPlatform := chi.URLParam(r, domain.URLParamPlatform)
	platform, ok := domain.StringToPlatform[rawPlatform]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	cfg, ok := dbUser.Platforms[platform]
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	service, ok := s.platforms[platform]
	if !ok {
		w.WriteHeader(http.StatusNotImplemented)
		return
	}

	if err := s.db.LeaveChannel(user.ID, platform); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := service.Leave(cfg.Channel); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
