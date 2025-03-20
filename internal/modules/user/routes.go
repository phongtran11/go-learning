package user

import (
	"modular-fx-fiber/internal/core/server"
	"modular-fx-fiber/internal/shared/middleware"
)

type (
	userRoutes struct {
		handlers *Handlers
	}

	UserRoutes interface{}
)

// NewUserRoutes creates new user routes
func NewUserRoutes(h *Handlers) UserRoutes {
	return &userRoutes{
		handlers: h,
	}
}

func Register(s server.Server, m middleware.Middleware, h Handlers) {
	group := s.GetApp().Group("api/users", m.JWT())
	group.Post("/", h.Create)
}
