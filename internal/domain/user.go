package domain

import (
	"github.com/google/uuid"
)

type User struct {
	Platforms map[Platform]string
	UUID      uuid.UUID
}

type UserSettings struct {
	IsBotJoined map[Platform]bool
	BannedUsers map[Platform][]string
	BannedWords []string
}
