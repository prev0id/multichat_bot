package service

import (
	"context"
	"net/http"
)

type twitchService interface {
	JoinChat(ctx context.Context, chat string) error
	LeaveChat(chat string) error
}

type Service struct {
	twitch twitchService
}

func New(twitch twitchService) *Service {
	return &Service{
		twitch: twitch,
	}
}

func (s *Service) Default(w http.ResponseWriter, r *http.Request) error {
	http.ServeFile(w, r, "website/index.html")
	return nil
}

func (s *Service) Static(w http.ResponseWriter, r *http.Request) error {
	fs := http.FileServer(http.Dir("website/static"))
	fs.ServeHTTP(w, r)

	return nil
}
