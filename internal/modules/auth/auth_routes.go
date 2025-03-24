package auth

import (
	"modular-fx-fiber/internal/core/server"
	"modular-fx-fiber/internal/shared/middleware"
)

type (
	Routes interface{}

	routes struct {
		handlers   Handlers
		middleware middleware.Middleware
	}
)

// NewRoutes creates new auth routes
func NewRoutes(h Handlers, m middleware.Middleware) Routes {
	return &routes{
		handlers:   h,
		middleware: m,
	}
}

// Register registers auth routes
func Register(s server.Server, h Handlers, m middleware.Middleware) {
	a := s.GetApp()
	group := a.Group("api/auth")

	// Public routes
	group.Post("/login", h.Login)
	group.Post("/register", h.Register)
	group.Post("/refresh-token", h.RefreshToken)
	// Protected routes
	group.Post("/register/verify-email", m.JWT(), h.VerifyEmail)
	group.Post("logout", m.JWT(), h.Logout)
}
