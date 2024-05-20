package page

import (
	"html/template"
	"net/http"

	"multichat_bot/internal/common/auth"
)

type Service struct {
	templates map[string]*template.Template
	auth      *auth.Auth
	isProd    bool
}

func NewService(isProd bool, authService *auth.Auth) (*Service, error) {
	s := &Service{
		isProd: isProd,
		auth:   authService,
	}

	if !s.isProd {
		return s, nil
	}

	if err := s.initTemplates(); err != nil {
		return nil, err
	}

	return s, nil
}

func (s *Service) HandleCSS(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "website/src/css/main.min.css")
}

func (s *Service) HandleJS(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "website/src/js/htmx.min.js")
}

func (s *Service) HandleIcon(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "website/src/assets/icon.png")
}
