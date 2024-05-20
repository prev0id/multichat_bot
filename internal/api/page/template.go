package page

import (
	"fmt"
	"html/template"
)

const (
	templateNameIndex = "website/src/index.gohtml"

	templateName404 = "website/src/html/404/404.gohtml"

	templateNameAccount          = "website/src/html/account/account.gohtml"
	templateNameAccountLoggedIn  = "website/src/html/account/logged_in.gohtml"
	templateNameAccountLoggedOut = "website/src/html/account/logged_out.gohtml"
	templateNameAccountLogos     = "website/src/html/account/logos.gohtml"

	templateNameSettings      = "website/src/html/settings/settings.gohtml"
	templateNameSettingsJoin  = "website/src/html/settings/toggle_join.gohtml"
	templateNameSettingsUsers = "website/src/html/settings/banned_users.gohtml"
	templateNameSettingsWords = "website/src/html/settings/banned_words.gohtml"
)

var (
	templateNameToParser = map[string]func() (*template.Template, error){
		templateName404:      parse404,
		templateNameAccount:  parseAccount,
		templateNameSettings: parseSettings,
	}
)

func (s *Service) initTemplates() error {
	s.templates = make(map[string]*template.Template, len(templateNameToParser))

	for name, parser := range templateNameToParser {
		tmpl, err := parser()
		if err != nil {
			return err
		}

		s.templates[name] = tmpl
	}

	return nil
}

func (s *Service) getTemplate(name string) (*template.Template, error) {
	if s.isProd {
		tmpl, ok := s.templates[name]
		if !ok {
			return nil, fmt.Errorf("template %s not found", name)
		}
		return tmpl, nil
	}

	tmpl, err := templateNameToParser[name]()
	if err != nil {
		return nil, err
	}

	return tmpl, nil
}

func parseAccount() (*template.Template, error) {
	templateAccount, err := template.ParseFiles(
		templateNameIndex,
		templateNameAccount,
		templateNameAccountLoggedIn,
		templateNameAccountLoggedOut,
		templateNameAccountLogos,
	)

	if err != nil {
		return nil, fmt.Errorf("error while parsing account template: %w", err)
	}

	return templateAccount, nil
}

func parseSettings() (*template.Template, error) {
	templateSettings, err := template.ParseFiles(
		templateNameIndex,
		templateNameSettings,
		templateNameSettingsJoin,
		templateNameSettingsUsers,
		templateNameSettingsWords,
	)

	if err != nil {
		return nil, fmt.Errorf("error while parsing settings template: %w", err)
	}

	return templateSettings, nil
}

func parse404() (*template.Template, error) {
	template404, err := template.ParseFiles(templateNameIndex, templateName404)
	if err != nil {
		return nil, fmt.Errorf("error while parsing 404 template: %w", err)
	}

	return template404, nil
}
