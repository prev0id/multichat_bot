package irc

import (
	"log/slog"

	"golang.org/x/net/websocket"
)

func (c *Client) Send(msg string) error {
	slog.Info("[twitch_irc_client] senging msg",
		slog.String("msg", msg),
	)

	return websocket.Message.Send(c.conn, msg)
}
