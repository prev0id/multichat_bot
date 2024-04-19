package processor

import (
	"context"
	"log/slog"

	"multichat_bot/internal/domain/logger"
	"multichat_bot/internal/platforms/twitch/domain"
)

func (p *Processor) defaultProcessor(_ context.Context, msg *domain.Message) {
	slog.Warn("[message_processor] unsupported message command",
		slog.String(logger.Command, msg.Command.Name),
		slog.String(logger.Message, msg.RawMessage),
	)
}
