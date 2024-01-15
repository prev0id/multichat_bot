package irc

import (
	"context"
	"log/slog"

	"multichat_bot/internal/config"

	"golang.org/x/net/websocket"
)

func (c *Client) Connect(ctx context.Context, cfg config.IRCServer) error {
	ws, err := websocket.Dial(cfg.Address, cfg.Protocol, cfg.Origin)
	if err != nil {
		return err
	}

	c.conn = ws

	go c.startReceiving(ctx)

	return nil
}

func (c *Client) startReceiving(ctx context.Context) {
	slog.Info("[twitch_irc_client] starting receiving messages")

	for {
		select {

		case <-ctx.Done():
			slog.Info("[twitch_irc_client] stopped", slog.String("reason", ctx.Err().Error()))
			return

		default:
			err := c.receive(ctx)
			if err != nil {
				slog.Error("[twitch_irc_client] recieved errror",
					slog.String("error", err.Error()),
				)
				return
			}
		}
	}
}

func (c *Client) receive(ctx context.Context) error {
	var rawMessage string
	err := websocket.Message.Receive(c.conn, &rawMessage)
	if err != nil {
		return err
	}

	messages := parse(rawMessage)

	for _, msg := range messages {
		if msg == nil {
			continue
		}

		c.processor.Process(ctx, msg)
	}

	return nil
}
