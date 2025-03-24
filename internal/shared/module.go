package shared

import (
	"modular-fx-fiber/internal/shared/database"
	"modular-fx-fiber/internal/shared/logger"
	"modular-fx-fiber/internal/shared/middleware"
	"modular-fx-fiber/internal/shared/repositories"
	"modular-fx-fiber/internal/shared/swagger"
	"modular-fx-fiber/internal/shared/validator"

	"go.uber.org/fx"
)

// Module exports shared dependencies
var Module = fx.Options(
	fx.Provide(
		database.NewDatabase,
		logger.NewZapLogger,
		middleware.NewMiddleware,
		swagger.NewSwagger,
		validator.NewValidator,
		// Repositories
		repositories.NewUserRepository,
		repositories.NewRefreshTokenRepository,
	),
	fx.Invoke(swagger.Register),
)
