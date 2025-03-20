package auth

import (
	"modular-fx-fiber/internal/core/server"
	"modular-fx-fiber/internal/shared/middleware"
)

type (
	AuthRoutes interface{}

	authRoutes struct {
		handlers   Handlers
		middleware middleware.Middleware
	}
)

// NewAuthRoutes creates new auth routes
func NewAuthRoutes(handlers Handlers, middleware middleware.Middleware) AuthRoutes {
	return &authRoutes{
		handlers:   handlers,
		middleware: middleware,
	}
}

// Register registers auth routes
func Register(server server.Server, routes authRoutes) {
	a := server.GetApp()
	group := a.Group("api/auth")

	group.Post("/login", routes.handlers.Login)
	group.Post("/register", routes.handlers.Register)
	group.Post("/register/verify-email", routes.handlers.VerifyEmail)
	group.Post("/refresh-token", routes.handlers.RefreshToken)
}
