package processor

import (
	"context"

	"multichat_bot/internal/platforms/twitch/domain"
)

// endOfNames 366 message processor = successful chat join
func (p *Processor) endOfNames(_ context.Context, msg *domain.Message) {
	p.service.ValidateJoin(msg.Command.Channel)
}
