package middleware

import (
	"errors"
	"modular-fx-fiber/internal/core/config"
	"modular-fx-fiber/internal/shared/logger"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
)

// Middleware contains application middlewares
type Middleware struct {
	config *config.Config
	logger *logger.ZapLogger
}

// NewMiddleware creates a new middleware instance
func NewMiddleware(config *config.Config,
	logger *logger.ZapLogger) *Middleware {
	return &Middleware{
		config: config,
		logger: logger,
	}
}

// JWT middleware for protecting routes
func (m *Middleware) JWT() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get authorization header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"error":   "Missing authorization header",
			})
		}

		// Check if the header has the right format
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"error":   "Invalid authorization format",
			})
		}

		// Get the token
		tokenString := parts[1]

		// Parse the token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
			// Validate signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("invalid signing method")
			}

			// Return the secret key
			return []byte(m.config.JWT.Secret), nil
		})

		// Handle parsing errors
		if err != nil {
			m.logger.Error("JWT Error: " + err.Error())
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"error":   "Invalid or expired token",
			})
		}

		// Check if token is valid
		if !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"error":   "Invalid token",
			})
		}

		// Extract claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"error":   "Invalid token claims",
			})
		}

		// Check token expiration
		if exp, ok := claims["exp"].(float64); ok {
			if time.Now().Unix() > int64(exp) {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"success": false,
					"error":   "Token expired",
				})
			}
		}

		// Store user info in context
		c.Locals("user_id", claims["user_id"])
		c.Locals("email", claims["email"])

		// Continue
		return c.Next()
	}
}
