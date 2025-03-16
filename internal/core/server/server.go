package server

import (
	"context"
	"fmt"
	"modular-fx-fiber/internal/core/config"
	appLogger "modular-fx-fiber/internal/shared/logger"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"go.uber.org/fx"
)

// Server represents the Fiber server
type Server struct {
	App    *fiber.App
	Config *config.Config
}

// NewServer creates a new server instance
func NewServer(config *config.Config) *Server {
	app := fiber.New(fiber.Config{
		AppName: config.App.Name,
	})

	// Add Logger
	app.Use(logger.New())
	// Add middlewares
	app.Use(recover.New())

	return &Server{
		App:    app,
		Config: config,
	}
}

func Start(lc fx.Lifecycle, s *Server, l *appLogger.ZapLogger) {
	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			l.Info(fmt.Sprint("Server starting ", "port: ", s.Config.App.Port))
			go func() {
				if err := s.App.Listen(":" + s.Config.App.Port); err != nil {
					l.Error(fmt.Sprint("Server failed to start", "error", err))
				}
			}()
			return nil
		},
		OnStop: func(context.Context) error {
			l.Info("Server stopping")
			if err := s.App.Shutdown(); err != nil {
				l.Error(fmt.Sprint("Server failed to stop", "error", err))
			}
			return nil
		},
	})
}
