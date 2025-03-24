package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"modular-fx-fiber/internal/core/config"
	appLogger "modular-fx-fiber/internal/shared/logger"
	"modular-fx-fiber/internal/shared/validator"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"go.uber.org/fx"
)

// Server represents the Fiber server
type Server interface {
	GetApp() *fiber.App
	GetConfig() *config.Config
}

type server struct {
	App    *fiber.App
	Config *config.Config
}

// NewServer creates a new server instance
func NewServer(l *appLogger.ZapLogger, config *config.Config) Server {
	app := fiber.New(fiber.Config{
		AppName:      config.App.Name,
		JSONEncoder:  json.Marshal,
		JSONDecoder:  json.Unmarshal,
		ErrorHandler: customErrorHandler,
	})

	// Add CORS
	app.Use(cors.New())

	// Add Logger
	app.Use(logger.New())

	// Add middlewares
	app.Use(recover.New())

	return &server{
		App:    app,
		Config: config,
	}
}

// GetApp returns the Fiber app
func (s *server) GetApp() *fiber.App {
	return s.App
}

// GetConfig returns the server config
func (s *server) GetConfig() *config.Config {
	return s.Config
}

func Start(lc fx.Lifecycle, s Server, l *appLogger.ZapLogger) {
	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			c := s.GetConfig()
			a := s.GetApp()
			l.Info(fmt.Sprint("Server starting ", "port: ", c.App.Port))
			go func() {
				if err := a.Listen(":" + c.App.Port); err != nil {
					l.Error(fmt.Sprint("Server failed to start", "error", err))
				}
			}()
			return nil
		},
		OnStop: func(context.Context) error {
			a := s.GetApp()
			l.Info("Server stopping")
			if err := a.Shutdown(); err != nil {
				l.Error(fmt.Sprint("Server failed to stop", "error", err))
			}
			return nil
		},
	})
}

func customErrorHandler(c *fiber.Ctx, err error) error {

	// Default status code is 500
	code := fiber.StatusInternalServerError

	// Check if it's a Fiber error
	var e *fiber.Error
	if errors.As(err, &e) {
		code = e.Code
	}

	// Return JSON response if the error is internal server error, not exposed details to the client
	if code == fiber.StatusInternalServerError {
		// Send 500 status code when Internal Server Error
		return c.Status(code).JSON(validator.GlobalErrorHandlerResponse{
			Success: false,
			Message: "Internal Server Error",
			Status:  code,
		})
	}

	// Return JSON response with error details
	return c.Status(code).JSON(validator.GlobalErrorHandlerResponse{
		Success: false,
		Message: err.Error(),
		Status:  code,
	})
}
