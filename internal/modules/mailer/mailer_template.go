package mailer

import (
	"bytes"
	"fmt"
	"html/template"
	"io/fs"
	"os"
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
	// Try multiple possible locations for the templates directory
	var templateDir string

	absPath, _ := filepath.Abs("./internal/modules/mailer/templates")
	if _, err := os.Stat(absPath); err == nil {
		templateDir = absPath
		fmt.Printf("Found templates directory: %s\n", templateDir)
	}

	// If no directory was found, default to a sensible location
	if templateDir == "" {
		// Get executable directory as fallback
		execPath, err := os.Executable()
		if err == nil {
			execDir := filepath.Dir(execPath)
			templateDir = filepath.Join(execDir, "templates")
			fmt.Printf("Using executable directory for templates: %s\n", templateDir)
		} else {
			// Last resort
			templateDir = "./templates"
			fmt.Printf("Defaulting to: %s\n", templateDir)
		}
	}

	// Create the directory if it doesn't exist
	if _, err := os.Stat(templateDir); os.IsNotExist(err) {
		fmt.Printf("Creating templates directory: %s\n", templateDir)
		if err := os.MkdirAll(templateDir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create templates directory: %w", err)
		}
	}

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
	// Log the template directory being searched
	fmt.Printf("Looking for templates in directory: %s\n", tm.templateDir)

	// Check if the directory exists
	if _, err := os.Stat(tm.templateDir); os.IsNotExist(err) {
		fmt.Printf("Warning: Template directory does not exist: %s\n", tm.templateDir)
		// You might want to return an error here, or create the directory
	}

	pattern := filepath.Join(tm.templateDir, "*.html")
	fmt.Printf("Using glob pattern: %s\n", pattern)

	files, err := filepath.Glob(pattern)
	if err != nil {
		return fmt.Errorf("failed to find templates: %w", err)
	}

	fmt.Printf("Found %d template files: %v\n", len(files), files)

	if len(files) == 0 {
		fmt.Println("Warning: No template files found!")
		// Optionally return early or continue
	}

	for _, file := range files {
		name := filepath.Base(file)
		// Remove extension to get template name
		name = strings.TrimSuffix(name, filepath.Ext(name))

		fmt.Printf("Loading template: %s from file %s\n", name, file)

		tmpl, err := template.ParseFiles(file)
		if err != nil {
			fmt.Printf("Error parsing template %s: %s\n", name, err)
			return fmt.Errorf("failed to parse template %s: %w", name, err)
		}

		tm.templates[name] = tmpl
		fmt.Printf("Successfully loaded template: %s\n", name)
	}

	fmt.Printf("All templates loaded. Total templates: %d\n", len(tm.templates))
	fmt.Printf("Template names: %v\n", tm.GetTemplateNames())

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
func (tm *TemplateManager) RenderTemplate(name string, data any) (string, error) {
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
