package page

import (
	"log/slog"
	"net/http"
)

func (s *Service) HandleSetting(w http.ResponseWriter, _ *http.Request) {
	tmpl, err := s.getTemplate(templateNameSettings)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, nil); err != nil {
		slog.Error("handling settings: " + err.Error())
	}

}

func (s *Service) getSettingsData(r *http.Request) {

}
