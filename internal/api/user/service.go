package user

import (
	"multichat_bot/internal/common/cookie"
	"multichat_bot/internal/database"
	"multichat_bot/internal/domain"
)

type platformService interface {
	Join(user string) error
	Leave(user string) error
}

type Service struct {
	platforms map[domain.Platform]platformService
	cookies   *cookie.Store
	db        *database.Manager
}

func NewService(db *database.Manager, cookies *cookie.Store) *Service {
	return &Service{
		platforms: make(map[domain.Platform]platformService, len(domain.StringToPlatform)),
		db:        db,
		cookies:   cookies,
	}
}

func (s *Service) WithPlatformService(platform domain.Platform, service platformService) *Service {
	s.platforms[platform] = service
	return s
}
