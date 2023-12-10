package bot

import (
	"golang.org/x/net/websocket"
)

type Implementation struct {
	conn *websocket.Conn
}
