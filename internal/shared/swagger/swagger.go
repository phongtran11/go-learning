package swagger

import (
	"modular-fx-fiber/internal/core/server"

	"github.com/gofiber/swagger"
	"github.com/swaggo/swag"
)

// Config holds swagger configuration
type Config struct {
	// BasePath for the swagger UI
	BasePath string
	// Title of the documentation
	Title string
	// Description of the API
	Description string
	// Version of the API
	Version string
	// Host (optional)
	Host string
}

// Swagger contains swagger handlers
type Swagger struct {
	config *Config
}

// NewSwagger creates a new swagger handler
func NewSwagger() *Swagger {
	config := Config{
		BasePath:    "/swagger",
		Title:       "API Documentation",
		Description: "API Documentation with Swagger",
		Version:     "1.0.0",
	}

	return &Swagger{
		config: &config,
	}
}

// Register registers swagger routes
func Register(s *server.Server, h *Swagger) {
	// Set up swagger info
	swag.Register(swag.Name, &swag.Spec{
		InfoInstanceName: swag.Name,
		Title:            h.config.Title,
		Description:      h.config.Description,
		Version:          h.config.Version,
		Host:             h.config.Host,
		BasePath:         "/api",
	})

	// Setup swagger route
	s.App.Get(h.config.BasePath+"/*", swagger.HandlerDefault)
}
