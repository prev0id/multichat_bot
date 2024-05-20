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

func (s *Service) SendMessage(message *domain.Message, channelID string) error {
	slog.Info(fmt.Sprintf("youtube: send a message to channel %s", channelID))
	return s.client.SendMessage(message, channelID)
}

func (s *Service) Join(channelID string) error {
	return s.client.Join(channelID)
}

func (s *Service) Leave(channelID string) error {
	s.client.Leave(channelID)
	return nil
}
