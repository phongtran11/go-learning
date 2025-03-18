package auth

import (
	"modular-fx-fiber/internal/core/server"
	"modular-fx-fiber/internal/shared/middleware"
)

type AuthRoutes struct {
	handlers   *Handlers
	middleware *middleware.Middleware
}

// NewAuthRoutes creates new auth routes
func NewAuthRoutes(handlers *Handlers, middleware *middleware.Middleware) *AuthRoutes {
	return &AuthRoutes{
		handlers:   handlers,
		middleware: middleware,
	}
}

// Register registers auth routes
func Register(server *server.Server, routes *AuthRoutes) {
	group := server.App.Group("api/auth")

	group.Post("/login", routes.handlers.Login)
	group.Post("/register", routes.handlers.Register)
	group.Post("/refresh", routes.handlers.RefreshToken)
}
