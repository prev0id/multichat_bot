package user

import (
	"html/template"

	"multichat_bot/internal/common/auth"
	"multichat_bot/internal/database"
	"multichat_bot/internal/domain"
)

type platformService interface {
	Join(config *domain.PlatformConfig) error
	Leave(config *domain.PlatformConfig) error
}

type Service struct {
	platforms map[domain.Platform]platformService
	db        *database.Manager
	auth      *auth.Auth

	toggleJoinTmpl  *template.Template
	bannedUsersTmpl *template.Template
	bannedWordsTmpl *template.Template
}

func NewService(db *database.Manager, authService *auth.Auth) *Service {
	return &Service{
		platforms:       make(map[domain.Platform]platformService, len(domain.StringToPlatform)),
		db:              db,
		auth:            authService,
		toggleJoinTmpl:  template.Must(template.ParseFiles("website/src/html/settings/toggle_join.gohtml")),
		bannedUsersTmpl: template.Must(template.ParseFiles("website/src/html/settings/banned_users.gohtml")),
		bannedWordsTmpl: template.Must(template.ParseFiles("website/src/html/settings/banned_words.gohtml")),
	}
}

func (s *Service) WithPlatformService(platform domain.Platform, service platformService) *Service {
	s.platforms[platform] = service
	return s
}
