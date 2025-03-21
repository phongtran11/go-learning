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
func Register(server server.Server, h Handlers, m middleware.Middleware) {
	a := server.GetApp()
	group := a.Group("api/auth")

	group.Post("/login", h.Login)
	group.Post("/register", h.Register)
	group.Post("/register/verify-email", m.JWT(), h.VerifyEmail)
	group.Post("/refresh-token", h.RefreshToken)
}
