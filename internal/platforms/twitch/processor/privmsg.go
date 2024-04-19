package processor

import (
	"context"
	"log/slog"

	global_domain "multichat_bot/internal/domain"
	"multichat_bot/internal/domain/logger"
	"multichat_bot/internal/platforms/twitch/domain"
)

func (p *Processor) privmsg(_ context.Context, msg *domain.Message) {
	slog.Info("[message_processor] received PrivMsg msg",
		slog.Any(logger.Message, *msg),
	)

	resp := &global_domain.Message{
		From:     msg.Command.Channel,
		Text:     msg.Command.BotCommand,
		Platform: global_domain.Twitch,
	}

	p.broadcast <- resp
}
