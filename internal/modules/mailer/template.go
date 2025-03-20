package mailer

import (
	"bytes"
	"fmt"
	"html/template"
	"io/fs"
	"path/filepath"
	"strings"
)

// TemplateManager handles email templates
type TemplateManager struct {
	templates   map[string]*template.Template
	templateDir string
}

// NewTemplateManager creates a new template manager
func NewTemplateManager() (*TemplateManager, error) {
	templateDir := filepath.Join("./", "templates")
	tm := &TemplateManager{
		templates:   make(map[string]*template.Template),
		templateDir: templateDir,
	}

	// Load all templates
	err := tm.loadTemplates()
	if err != nil {
		return nil, err
	}

	return tm, nil
}

// loadTemplates loads all HTML templates from the template directory
func (tm *TemplateManager) loadTemplates() error {
	pattern := filepath.Join(tm.templateDir, "*.html")
	files, err := filepath.Glob(pattern)
	if err != nil {
		return fmt.Errorf("failed to find templates: %w", err)
	}

	for _, file := range files {
		name := filepath.Base(file)
		// Remove extension to get template name
		name = strings.TrimSuffix(name, filepath.Ext(name))

		tmpl, err := template.ParseFiles(file)
		if err != nil {
			return fmt.Errorf("failed to parse template %s: %w", name, err)
		}

		tm.templates[name] = tmpl
	}

	return nil
}

// LoadTemplatesFromFS loads templates from an embedded filesystem
func (tm *TemplateManager) LoadTemplatesFromFS(fsys fs.FS, dir string) error {
	entries, err := fs.ReadDir(fsys, dir)
	if err != nil {
		return fmt.Errorf("failed to read templates directory: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".html") {
			continue
		}

		name := strings.TrimSuffix(entry.Name(), ".html")
		path := filepath.Join(dir, entry.Name())

		content, err := fs.ReadFile(fsys, path)
		if err != nil {
			return fmt.Errorf("failed to read template %s: %w", name, err)
		}

		tmpl, err := template.New(name).Parse(string(content))
		if err != nil {
			return fmt.Errorf("failed to parse template %s: %w", name, err)
		}

		tm.templates[name] = tmpl
	}

	return nil
}

// RenderTemplate renders a template with the given data
func (tm *TemplateManager) RenderTemplate(name string, data interface{}) (string, error) {
	tmpl, exists := tm.templates[name]
	if !exists {
		return "", fmt.Errorf("template %s not found", name)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to render template %s: %w", name, err)
	}

	return buf.String(), nil
}

// HasTemplate checks if a template exists
func (tm *TemplateManager) HasTemplate(name string) bool {
	_, exists := tm.templates[name]
	return exists
}

// GetTemplateNames returns the names of all loaded templates
func (tm *TemplateManager) GetTemplateNames() []string {
	names := make([]string, 0, len(tm.templates))
	for name := range tm.templates {
		names = append(names, name)
	}
	return names
}
