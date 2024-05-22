package domain

import (
	"fmt"
	"slices"
	"strings"
)

type Message struct {
	From     string
	Text     string
	Channel  string
	Platform Platform
}

func (m *Message) FormatedText() string {
	return fmt.Sprintf("%s says: %s", m.From, m.Text)
}

func (m *Message) Validate(config *PlatformConfig) error {
	if !config.IsJoined {
		return fmt.Errorf("bot not joined to chat")
	}

	if slices.Contains(config.DisabledUsers, m.From) {
		return fmt.Errorf("message is banned for user %s", m.From)
	}

	lowerCase := strings.ToLower(m.Text)
	for _, word := range config.BannedWords {
		if strings.Contains(lowerCase, word) {
			return fmt.Errorf("message is banned for a word %s", word)
		}
	}

	if strings.Contains(m.Text, " says: ") {
		return fmt.Errorf("message from bot %s", m.Text)
	}

	return nil
}
