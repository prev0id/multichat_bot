package processor

import (
	"context"
	"log/slog"

	"multichat_bot/internal/twitch/domain"
)

func (p *Processor) privmsg(_ context.Context, msg *domain.Message) {
	slog.Info("[message_processor] received PrivMsg msg:",
		slog.Any("msg", *msg),
	)
}
