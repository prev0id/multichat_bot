package user

import (
	"github.com/google/uuid"

	"multichat_bot/internal/domain"
)

type platformService interface {
	Join(user string) error
	Leave(user string) error
}

type db interface {
	UpdateUserPlatform(userUUID uuid.UUID, platform domain.Platform, channel string) error
	RemoveUserPlatform(userUUID uuid.UUID, platform domain.Platform) error
}

type Service struct {
	platforms map[domain.Platform]platformService
	db        db
}

func NewService(db db) *Service {
	return &Service{
		platforms: make(map[domain.Platform]platformService, len(domain.StringToPlatform)),
		db:        db,
	}
}

func (s *Service) WithPlatformService(platform domain.Platform, service platformService) *Service {
	s.platforms[platform] = service
	return s
}
