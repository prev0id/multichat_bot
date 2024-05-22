package youtube

import (
	"fmt"
	"log/slog"

	"multichat_bot/internal/domain"
	"multichat_bot/internal/platform/youtube/client"
)

type Service struct {
	client *client.Adapter
}

func NewService(adapter *client.Adapter) *Service {
	return &Service{
		client: adapter,
	}
}

func (s *Service) SendMessage(message *domain.Message, config *domain.PlatformConfig) error {
	slog.Info(fmt.Sprintf("youtube: send a message to channel %s", config.ID))
	return s.client.SendMessage(message, config)
}

func (s *Service) Join(cfg *domain.PlatformConfig) error {
	return s.client.Join(cfg)
}

func (s *Service) Leave(cfg *domain.PlatformConfig) error {
	s.client.Leave(cfg.ID)
	return nil
}
