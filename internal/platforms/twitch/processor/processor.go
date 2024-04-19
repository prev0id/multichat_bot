package processor

import (
	"context"
	"log/slog"

	global_domain "multichat_bot/internal/domain"
	"multichat_bot/internal/platforms/twitch/domain"
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
	broadcast  chan<- *global_domain.Message
}

func New(service twitchService, broadcast chan<- *global_domain.Message) *Processor {
	processor := &Processor{
		service:   service,
		broadcast: broadcast,
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
