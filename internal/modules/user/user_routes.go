package user

import (
	"modular-fx-fiber/internal/core/server"
	"modular-fx-fiber/internal/shared/middleware"
)

type (
	Routes interface{}

	routes struct {
		handlers *Handlers
	}
)

// NewRoutes creates new user routes
func NewRoutes(h *Handlers) Routes {
	return &routes{
		handlers: h,
	}
}

func Register(s server.Server, m middleware.Middleware, h Handlers) {
	group := s.GetApp().Group("api/users", m.JWT())
	group.Get("/", h.ListUsers)
	group.Post("/", h.Create)
	group.Get("/me", h.GetMe)
}
