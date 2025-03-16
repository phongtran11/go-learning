package user

import (
	"modular-fx-fiber/internal/core/server"
	"modular-fx-fiber/internal/shared/middleware"
)

type UserRoutes struct {
	handlers   *Handlers
	middleware *middleware.Middleware
}

// NewUserRoutes creates new user routes
func NewUserRoutes(handlers *Handlers, middleware *middleware.Middleware) *UserRoutes {

	return &UserRoutes{
		handlers:   handlers,
		middleware: middleware,
	}
}

func Register(server *server.Server, handlers *Handlers, middleware *middleware.Middleware) {
	group := server.App.Group("/users")
	group.Post("/", handlers.Create)
	group.Get("/", handlers.List)
	group.Get("/:id", handlers.GetByID)
	group.Put("/:id", handlers.Update)
	group.Delete("/:id", handlers.Delete)

}
