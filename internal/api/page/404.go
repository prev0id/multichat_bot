package page

import (
	"log/slog"
	"net/http"
)

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
