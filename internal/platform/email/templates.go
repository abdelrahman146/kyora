package email

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"strings"
)

//go:embed templates/*.html
var templatesFS embed.FS

const (
	TemplateForgotPassword TemplateID = "forgot_password"
)

// registry maps TemplateID to file paths within the embedded FS
var templateFiles = map[TemplateID]string{
	TemplateForgotPassword: "templates/forgot_password.html",
}

// subjects maps TemplateID to a default subject line
var subjects = map[TemplateID]string{
	TemplateForgotPassword: "Reset your password",
}

// RenderTemplate renders the embedded HTML template with provided data.
// Missing keys render as empty strings.
func RenderTemplate(id TemplateID, data map[string]any) (string, error) {
	path, ok := templateFiles[id]
	if !ok {
		return "", ErrTemplateNotFound(string(id))
	}
	contentBytes, err := templatesFS.ReadFile(path)
	if err != nil {
		return "", err
	}
	// prepare funcs: default returns the first non-empty string
	funcMap := template.FuncMap{
		"default": func(def string, v any) string {
			// best-effort convert to string
			switch vv := v.(type) {
			case string:
				if strings.TrimSpace(vv) != "" {
					return vv
				}
			case []byte:
				s := strings.TrimSpace(string(vv))
				if s != "" {
					return s
				}
			}
			return def
		},
	}
	t, err := template.New(string(id)).Funcs(funcMap).Option("missingkey=zero").Parse(string(contentBytes))
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// SubjectFor returns a sensible default subject for a template id.
func SubjectFor(id TemplateID) string {
	if s, ok := subjects[id]; ok {
		return s
	}
	return string(id)
}

// ErrTemplateNotFound returns an error when a template ID is not registered.
func ErrTemplateNotFound(id string) error { return fmt.Errorf("email template not found: %s", id) }
