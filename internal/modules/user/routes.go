package user

import (
	"modular-fx-fiber/internal/core/server"
	"modular-fx-fiber/internal/shared/middleware"
)

type UserRoutes struct {
	handlers *Handlers
}

// NewUserRoutes creates new user routes
func NewUserRoutes(h *Handlers) *UserRoutes {

	return &UserRoutes{
		handlers: h,
	}
}

func Register(s *server.Server, h *Handlers, m *middleware.Middleware) {
	group := s.App.Group("api/users", m.JWT())
	group.Post("/", h.Create)
}
