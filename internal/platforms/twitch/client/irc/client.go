package irc

import (
	"context"

	"golang.org/x/net/websocket"

	"multichat_bot/internal/platforms/twitch/domain"
)

type messageProcessor interface {
	Process(ctx context.Context, message *domain.Message)
}

type Client struct {
	conn      *websocket.Conn
	processor messageProcessor
}

func NewClient() *Client {
	return &Client{}
}

func (c *Client) WithMessageProcessor(processor messageProcessor) {
	c.processor = processor
}
