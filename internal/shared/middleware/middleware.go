package middleware

import (
	"errors"
	"fmt"
	"modular-fx-fiber/internal/core/config"
	"modular-fx-fiber/internal/shared/logger"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
)

// Define standard error types
var (
	ErrMissingAuthHeader = errors.New("missing authorization header")
	ErrInvalidAuthFormat = errors.New("invalid authorization format")
	ErrInvalidToken      = errors.New("invalid token")
	ErrTokenExpired      = errors.New("token expired")
	ErrInvalidClaims     = errors.New("invalid token claims")
)

type (
	Middleware interface {
		JWT() fiber.Handler
	}

	middleware struct {
		config *config.Config
		logger *logger.ZapLogger
	}

	// UserClaims defines the structure for JWT claims
	UserClaims struct {
		UserID uint64 `json:"user_id"`
		Email  string `json:"email"`
		jwt.RegisteredClaims
	}
)

// NewMiddleware creates a new middleware instance
func NewMiddleware(config *config.Config, logger *logger.ZapLogger) Middleware {
	return &middleware{
		config: config,
		logger: logger,
	}
}

// JWT middleware for protecting routes
// JWT middleware for protecting routes
func (m *middleware) JWT() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get authorization header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return fiber.NewError(fiber.StatusUnauthorized, ErrMissingAuthHeader.Error())
		}

		// Check if the header has the right format
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return fiber.NewError(fiber.StatusUnauthorized, ErrInvalidAuthFormat.Error())
		}

		// Get the token
		tokenString := parts[1]

		// Parse and validate the token
		token, err := jwt.ParseWithClaims(
			tokenString,
			&UserClaims{},
			func(token *jwt.Token) (any, error) {
				// Validate signing method
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}

				// Return the secret key
				return []byte(m.config.JWT.Secret), nil
			},
		)

		// Handle parsing errors
		if err != nil {
			// Check for specific error types
			if errors.Is(err, jwt.ErrTokenExpired) {
				m.logger.Debug("Token expired", zap.Error(err))
				return fiber.NewError(fiber.StatusUnauthorized, ErrTokenExpired.Error())
			}

			m.logger.Error("JWT parsing error", zap.Error(err))
			return fiber.NewError(fiber.StatusUnauthorized, ErrInvalidToken.Error())
		}

		// Check if token is valid
		if !token.Valid {
			return fiber.NewError(fiber.StatusUnauthorized, ErrInvalidToken.Error())
		}

		// Extract claims with proper type assertion
		claims, ok := token.Claims.(*UserClaims)
		if !ok {
			m.logger.Error("Failed to assert token claims",
				zap.String("claims_type", fmt.Sprintf("%T", token.Claims)))
			return fiber.NewError(fiber.StatusUnauthorized, ErrInvalidClaims.Error())
		}

		// Validate required claims
		if claims.UserID == 0 || claims.Email == "" {
			m.logger.Error("Missing required claims",
				zap.Any("user_id", claims.UserID),
				zap.String("email", claims.Email))
			return fiber.NewError(fiber.StatusUnauthorized, "Invalid token claims")
		}

		// Store user info in context with proper types
		c.Locals("user_id", claims.UserID)
		c.Locals("email", claims.Email)

		m.logger.Debug("JWT successfully validated",
			zap.Uint64("user_id", claims.UserID),
			zap.String("email", claims.Email),
			zap.Time("expires", claims.ExpiresAt.Time))

		// Continue
		return c.Next()
	}
}
