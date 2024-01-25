package api

import (
	"context"
	"log/slog"
	"time"

	desc "multichat_bot/internal/api/gen"
)

type twitchService interface {
	JoinChat(ctx context.Context, chat string) error
	LeaveChat(chat string) error
}

type Server struct {
	twitch twitchService
}

func NewServer(twitch twitchService) *Server {
	return &Server{
		twitch: twitch,
	}
}

func (s *Server) LeaveTwitchChat(_ context.Context, request desc.LeaveTwitchChatRequestObject) (desc.LeaveTwitchChatResponseObject, error) {
	slog.Info("LeaveTwitchChat", slog.String("chat", request.Chat))

	if err := s.twitch.LeaveChat(request.Chat); err != nil {
		return desc.LeaveTwitchChat500Response{}, err
	}

	return desc.LeaveTwitchChat200Response{}, nil
}

func (s *Server) JoinTwitchChat(ctx context.Context, request desc.JoinTwitchChatRequestObject) (desc.JoinTwitchChatResponseObject, error) {
	slog.Info("JoinTwitchChat", slog.String("chat", request.Chat))

	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	if err := s.twitch.JoinChat(ctx, request.Chat); err != nil {
		return desc.JoinTwitchChat500Response{}, err
	}

	return desc.JoinTwitchChat200Response{}, nil
}
