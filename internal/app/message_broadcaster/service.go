package message_broadcaster

import (
	"context"
	"log/slog"

	"multichat_bot/internal/domain"
	"multichat_bot/internal/domain/logger"
)

type db interface {
	GetUserByPlatform(platform domain.Platform, user string) (*domain.User, error)
}

type platformService interface {
	SendMessage(message *domain.Message, user *domain.User) error
}

type Service struct {
	db        db
	platforms map[domain.Platform]platformService
	messageCh chan *domain.Message
}

func New(db db) *Service {
	return &Service{
		db:        db,
		platforms: make(map[domain.Platform]platformService),
	}
}

func (s *Service) AddPlatform(platform domain.Platform, service platformService) {
	s.platforms[platform] = service
}

func (s *Service) StartWorker(ctx context.Context) {
	slog.Info("messageManager::worker start listening")

	go func() {
		select {
		case <-ctx.Done():
			slog.Error("messageManager::worker end of listening", slog.String(logger.Error, ctx.Err().Error()))
			return
		case msg := <-s.messageCh:
			s.broadcast(msg)
		}
	}()
}

func (s *Service) GetMessageChannel() chan<- *domain.Message {
	return s.messageCh
}

func (s *Service) broadcast(msg *domain.Message) {
	user, err := s.db.GetUserByPlatform(msg.Platform, msg.From)
	if err != nil {
		slog.Error(
			"messageManager::broadcast unable to get user from db",
			slog.String(logger.Error, err.Error()),
		)
		return
	}

	for platform := range user.Platforms {
		service, ok := s.platforms[platform]
		if !ok {
			slog.Warn("messageManager::broadcast unable to find platform", slog.String(logger.Platform, string(platform)))
			continue
		}

		if err := service.SendMessage(msg, user); err != nil {
			slog.Error("messageManager::broadcast error occurred while sending the message",
				slog.String(logger.Error, err.Error()),
				slog.Any(logger.User, user),
				slog.Any(logger.Message, msg),
			)
		}
	}
}
