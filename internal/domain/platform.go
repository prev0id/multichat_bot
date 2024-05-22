package domain

import (
	"slices"
	"strings"
	"time"

	"golang.org/x/oauth2"
)

type Platform string

func (p Platform) String() string {
	return string(p)
}

const (
	Twitch  Platform = "twitch"
	YouTube Platform = "youtube"
)

var (
	StringToPlatform = map[string]Platform{
		Twitch.String():  Twitch,
		YouTube.String(): YouTube,
	}
	Platforms = []Platform{Twitch, YouTube}
)

type PlatformConfig struct {
	ExpiresIn     time.Time
	ID            string
	Channel       string
	AccessToken   string
	RefreshToken  string
	DisabledUsers BannedList
	BannedWords   BannedList
	IsJoined      bool
}

func (c *PlatformConfig) Token() *oauth2.Token {
	return &oauth2.Token{
		AccessToken:  c.AccessToken,
		TokenType:    "Bearer",
		RefreshToken: c.RefreshToken,
		Expiry:       c.ExpiresIn,
	}
}

type BannedList []string

func (w BannedList) Add(word string) BannedList {
	if word == "" {
		return w
	}

	word = strings.TrimSpace(strings.ToLower(word))

	if slices.Contains(w, word) {
		return w
	}

	return append(w, word)
}

func (w BannedList) Remove(word string) BannedList {
	if w == nil {
		return nil
	}

	word = strings.TrimSpace(strings.ToLower(word))

	idx := slices.Index(w, word)
	if idx < 0 {
		return w
	}

	return append(w[:idx], w[idx+1:]...)
}
