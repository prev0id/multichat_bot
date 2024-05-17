package page

import (
	"html/template"
	"log/slog"
	"net/http"

	"multichat_bot/internal/common/cookie"
	"multichat_bot/internal/database"
)

const (
	templateNameIndex    = "website/src/index.gohtml"
	templateName404      = "website/src/html/404.gohtml"
	templateNameAccount  = "website/src/html/account.gohtml"
	templateNameSettings = "website/src/html/settings.gohtml"
)

type Service struct {
	templates   map[string]*template.Template
	cookieStore *cookie.Store
	db          *database.Manager
	isProd      bool
}

func NewService(isProd bool, cookieStore *cookie.Store, db *database.Manager) (*Service, error) {
	s := &Service{
		isProd:      isProd,
		cookieStore: cookieStore,
		db:          db,
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

func (s *Service) Handle404(w http.ResponseWriter, _ *http.Request) {
	tmpl, err := s.getTemplate(templateName404)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, nil); err != nil {
		slog.Error("handling 404: " + err.Error())
	}
}
