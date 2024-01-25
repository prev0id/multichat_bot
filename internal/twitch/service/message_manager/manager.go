package message_manager

import (
	"errors"

	"multichat_bot/internal/twitch/domain"
	"multichat_bot/internal/twitch/service/message_manager/rate_limit"
)

const (
	rateLimitAuthAttempts = 20
	rateLimitJoinAttempts = 20
	rateLimitChatDefault  = 20
	rateLimitChatMod      = 100
)

var (
	ErrorRateLimitExited = errors.New("twitch's rate limit exited")
)

type ircClient interface {
	Send(msg string) error
}

type Manager struct {
	ircClient ircClient

	chatsRL *rate_limit.Map
	authRL  *rate_limit.Checker
	joinRL  *rate_limit.Checker
}

func New(client ircClient) *Manager {
	return &Manager{
		ircClient: client,

		chatsRL: rate_limit.NewMapChecker(rateLimitChatDefault),
		joinRL:  rate_limit.NewChecker(rateLimitJoinAttempts),
		authRL:  rate_limit.NewChecker(rateLimitAuthAttempts),
	}
}

func (m *Manager) SendChatMessage(chat string, msg domain.IRCMessage) error {
	if !m.chatsRL.IsLimitExceeded(chat) {
		return ErrorRateLimitExited
	}

	return m.ircClient.Send(msg.ToString())
}

func (m *Manager) SendAuthMessage(msg domain.IRCMessage) error {
	if !m.authRL.IsLimitExceeded() {
		return ErrorRateLimitExited
	}

	return m.ircClient.Send(msg.ToString())
}

func (m *Manager) SendJoinMessage(msg domain.IRCMessage) error {
	if !m.joinRL.IsLimitExceeded() {
		return ErrorRateLimitExited
	}

	return m.ircClient.Send(msg.ToString())
}
