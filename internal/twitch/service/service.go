package service

import (
	"context"
	"log/slog"
	"strings"

	"golang.org/x/oauth2"

	"multichat_bot/internal/config"

	"multichat_bot/internal/twitch/domain"
)

type messageManager interface {
	SendChatMessage(chat string, msg domain.IRCMessage) error
	SendAuthMessage(msg domain.IRCMessage) error
	SendJoinMessage(msg domain.IRCMessage) error
}

type Service struct {
	messageManager messageManager
	token          oauth2.TokenSource
	chats          chats
}

func New(manager messageManager, token oauth2.TokenSource) *Service {
	return &Service{
		messageManager: manager,
		token:          token,
		chats:          newChats(),
	}
}

func (s *Service) Connect(cfg config.Twitch) error {
	token, err := s.token.Token()
	if err != nil {
		return err
	}

	if err = s.messageManager.SendAuthMessage(domain.PassMessage(token.AccessToken)); err != nil {
		return err
	}

	if err = s.messageManager.SendAuthMessage(domain.NickMessage(cfg.Username)); err != nil {
		return err
	}

	if err = s.messageManager.SendAuthMessage(domain.CapReqMessage(cfg.Capabilities)); err != nil {
		return err
	}

	return nil
}

func (s *Service) JoinChat(ctx context.Context, chat string) error {
	ctx, cancel := context.WithCancelCause(ctx)
	if err := s.chats.processJoinRequest(chat, cancel); err != nil {
		return err
	}

	msg := &domain.JoinMessage{chat}
	if err := s.messageManager.SendJoinMessage(msg); err != nil {
		slog.Error(
			"[message_processor] unable to join chat",
			slog.String("chat", chat),
			slog.String("error", err.Error()),
			slog.String("type", domain.IRCCommandJoin),
			slog.String("message", msg.ToString()),
		)

		return err
	}

	<-ctx.Done()

	return context.Cause(ctx)
}

func (s *Service) LeaveChat(chat string) error {
	msg := domain.PartMessage(chat)

	err := s.messageManager.SendJoinMessage(msg)
	if err != nil {
		slog.Error(
			"[message_processor] unable to leave chat",
			slog.String("chat", chat),
			slog.String("error", err.Error()),
			slog.String("type", domain.IRCCommandPart),
			slog.String("message", msg.ToString()),
		)
	}

	return err
}

func (s *Service) SendPongMessage(rawPingMessage string) {
	pong := strings.Replace(rawPingMessage, "PING", "PONG", 1)
	msg := domain.PongMessage(pong)

	err := s.messageManager.SendAuthMessage(msg)
	if err != nil {
		slog.Error(
			"[message_processor] unable to send Pong message",
			slog.String("error", err.Error()),
			slog.String("type", domain.IRCCommandPing),
			slog.String("message", msg.ToString()),
		)
	}
}

func (s *Service) SendTextMessage(chat, text string) {
	msg := &domain.PrivMessage{
		Text:    text,
		Channel: chat,
	}

	err := s.messageManager.SendChatMessage(chat, msg)
	if err != nil {
		slog.Error(
			"[message_processor] unable to send PrivMessage",
			slog.String("error", err.Error()),
			slog.String("type", domain.IRCCommandPrivmsg),
			slog.String("message", msg.ToString()),
		)
	}
}
