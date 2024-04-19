package processor

import (
	"context"

	"multichat_bot/internal/platforms/twitch/domain"
)

func (p *Processor) ping(_ context.Context, msg *domain.Message) {
	p.service.SendPongMessage(msg.RawMessage)
}
