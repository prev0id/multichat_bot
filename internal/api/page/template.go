package page

import (
	"fmt"
	"html/template"
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
	templateAccount, err := template.ParseFiles(templateNameIndex, templateNameAccount)
	if err != nil {
		return nil, fmt.Errorf("error while parsing account template: %w", err)
	}
	return templateAccount, nil
}

func parseSettings() (*template.Template, error) {
	templateSettings, err := template.ParseFiles(templateNameIndex, templateNameSettings)
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
