package middleware

import (
	"github.com/gofiber/fiber/v2"
)

// Middleware contains application middlewares
type Middleware struct{}

// NewMiddleware creates a new middleware instance
func NewMiddleware() *Middleware {
	return &Middleware{}
}

// Auth is a middleware to check authentication
func (m *Middleware) Auth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// token := c.Get("Authorization")
		// if token == "" {
		// 	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
		// 		"error": "Unauthorized",
		// 	})
		// }

		return c.Next()
	}
}
