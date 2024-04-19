package domain

import (
	"github.com/google/uuid"
)

type User struct {
	UUID      uuid.UUID
	Platforms map[Platform]string
}
