package processor

import (
	"context"
	"log/slog"

	"multichat_bot/internal/twitch/domain"
)

type twitchService interface {
	SendTextMessage(chat, text string)
	SendPongMessage(rawPingMessage string)
	ValidateJoin(chat string)
}

type processFn func(ctx context.Context, msg *domain.Message)

type Processor struct {
	service    twitchService
	processors map[string]processFn
}

func New(service twitchService) *Processor {
	processor := &Processor{
		service: service,
	}

	processor.processors = map[string]processFn{
		domain.IRCCommandPrivmsg: processor.privmsg,
		domain.IRCCommandPing:    processor.ping,
		domain.IRCCommand366:     processor.endOfNames,
	}

	return processor
}

func (p *Processor) Process(ctx context.Context, msg *domain.Message) {
	if msg == nil || msg.Command == nil {
		slog.Error("empty message")
	}

	processor, isExist := p.processors[msg.Command.Name]
	if !isExist {
		processor = p.defaultProcessor
	}

	processor(ctx, msg)
}
