package processor

import (
	"context"
	"log/slog"

	"multichat_bot/internal/twitch/domain"
)

func (p *Processor) defaultProcessor(_ context.Context, msg *domain.Message) {
	slog.Warn("[message_processor] unsupported message command",
		slog.String("command", msg.Command.Name),
		slog.String("message", msg.RawMessage),
	)
}
