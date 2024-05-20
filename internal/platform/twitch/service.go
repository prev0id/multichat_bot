package twitch

import (
	"fmt"
	"log"
	"log/slog"

	"multichat_bot/internal/config"
	"multichat_bot/internal/domain"

	"github.com/gempir/go-twitch-irc/v4"
)

const (
	oauthPrefix = "oauth:"
)

type Service struct {
	client         *twitch.Client
	messageChannel chan<- *domain.Message
}

func NewService(cfg config.Twitch, messageChannel chan<- *domain.Message) *Service {
	service := &Service{
		messageChannel: messageChannel,
	}

	client := twitch.NewClient(cfg.Username, oauthPrefix+cfg.Token)
	client.Capabilities = []string{twitch.TagsCapability, twitch.CommandsCapability, twitch.MembershipCapability}
	client.OnPrivateMessage(service.chatMessageCallback)

	service.client = client

	return service
}

func (s *Service) Connect() error {
	go func() {
		if err := s.client.Connect(); err != nil {
			log.Fatalf("err while connecting to twitch: %v", err) //nolint:revive
		}
	}()

	return nil
}

func (s *Service) Join(channel string) error {
	slog.Info(fmt.Sprintf("twitch: joining channel %s", channel))
	s.client.Join(channel)
	return nil
}

func (s *Service) Leave(channel string) error {
	s.client.Depart(channel)
	return nil
}

func (s *Service) SendMessage(message *domain.Message, channel string) error {
	s.client.Say(channel, message.Text)
	return nil
}

func (s *Service) chatMessageCallback(msg twitch.PrivateMessage) {
	if msg.GetType() != twitch.PRIVMSG {
		return
	}

	slog.Info(fmt.Sprintf("twitch: received a message from %s: %s", msg.User.Name, msg.Message))

	s.messageChannel <- convertTwitchPrivmsgToDomain(msg)
}

func convertTwitchPrivmsgToDomain(from twitch.PrivateMessage) *domain.Message {
	return &domain.Message{
		From:     from.User.DisplayName,
		Channel:  from.RoomID,
		Text:     from.Message,
		Platform: domain.Twitch,
	}
}
