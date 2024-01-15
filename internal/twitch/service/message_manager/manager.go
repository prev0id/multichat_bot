package message_manager

import (
	"errors"

	"multichat_bot/internal/twitch/domain"
)

var (
	ErrorRateLimitExited = errors.New("twitch's rate limit exited")
)

type ircClient interface {
	Send(msg string) error
}

type Manager struct {
	ircClient ircClient

	chatsRL map[string]windowRateLimit
	authRL  windowRateLimit
	joinRL  windowRateLimit
}

func New(client ircClient) *Manager {
	return &Manager{
		ircClient: client,
	}
}

func (m *Manager) SendChatMessage(msg domain.IRCMessage) error {
	return m.ircClient.Send(msg.ToString())
}

func (m *Manager) SendAuthMessage(msg domain.IRCMessage) error {
	if !m.authRL.isSendAllowed() {
		return ErrorRateLimitExited
	}

	return m.ircClient.Send(msg.ToString())
}

func (m *Manager) SendJoinMessage(msg domain.IRCMessage) error {
	if !m.joinRL.isSendAllowed() {
		return ErrorRateLimitExited
	}

	return m.ircClient.Send(msg.ToString())
}
