package domain

import (
	"time"
)

type User struct {
	Platforms map[Platform]*PlatformConfig
	ID        int64
}

type PlatformConfig struct {
	ExpiresIn     time.Time
	ID            string
	Channel       string
	AccessToken   string
	RefreshToken  string
	DisabledUsers []string
	BannedWords   []string
	IsJoined      bool
}
