package ws

import (
	"context"
	"fmt"
	"log/slog"

	"golang.org/x/net/websocket"

	"multichat_bot/internal/config"
)

const (
	wsAddress = "ws://irc-ws.chat.twitch.tv:80"
)

type Client struct {
	name      string
	conn      *websocket.Conn
	messageCh chan string
}

func New(ctx context.Context, name string, cfg config.Client) (*Client, error) {
	ws, err := websocket.Dial(cfg.Address, cfg.Protocol, cfg.Origin)
	if err != nil {
		return nil, err
	}

	client := &Client{
		name:      name,
		conn:      ws,
		messageCh: make(chan string, 10),
	}

	go client.startReceiving(ctx)

	return client, nil
}

func (c *Client) startReceiving(ctx context.Context) {
	slog.Info(fmt.Sprintf("[%s] starting receiving messages", c.name))

	for {
		select {
		case <-ctx.Done():
			slog.Info(fmt.Sprintf("[%s] stopped", c.name))
			return

		default:
			err := c.receive()
			if err != nil {
				slog.Error(fmt.Sprintf("[%s] recieved err: %s", c.name, err.Error()))
				return
			}
		}
	}
}

func (c *Client) receive() error {
	var message string
	err := websocket.Message.Receive(c.conn, &message)
	if err != nil {
		return err
	}

	slog.Info(fmt.Sprintf("[%s] recieved message: \"\n%s\"", c.name, message))

	c.messageCh <- message
	return nil
}

func (c *Client) Send(message string) error {
	slog.Info(fmt.Sprintf("[%s] senging message \"%s\"", c.name, message))
	return websocket.Message.Send(c.conn, message)
}

func (c *Client) GetMessageChannel() <-chan string {
	return c.messageCh
}
