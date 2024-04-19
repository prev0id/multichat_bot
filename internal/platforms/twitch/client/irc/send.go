package irc

import (
	"log/slog"

	"golang.org/x/net/websocket"

	"multichat_bot/internal/domain/logger"
)

func (c *Client) Send(msg string) error {
	slog.Info("[twitch_irc_client] senging msg",
		slog.String(logger.Message, msg),
	)

	return websocket.Message.Send(c.conn, msg)
}
