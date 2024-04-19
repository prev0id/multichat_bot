package domain

import (
	"fmt"
	"strings"
)

type IRCMessage interface {
	ToString() string
}

type PrivMessage struct {
	Text    string
	Channel string
}

func (m *PrivMessage) ToString() string {
	return fmt.Sprintf("PRIVMSG #%s: %s", m.Channel, m.Text)
}

type PongMessage string

func (m PongMessage) ToString() string {
	return string(m)
}

type JoinMessage []string

func (m JoinMessage) ToString() string {
	return "JOIN #" + strings.Join(m, "#,")
}

type PartMessage string

func (m PartMessage) ToString() string {
	return "PART #" + string(m)
}

type CapReqMessage []string

func (m CapReqMessage) ToString() string {
	return "CAP REQ :" + strings.Join(m, " ")
}

type PassMessage string

func (m PassMessage) ToString() string {
	return "PASS oauth:" + string(m)
}

type NickMessage string

func (m NickMessage) ToString() string {
	return "NICK " + string(m)
}
