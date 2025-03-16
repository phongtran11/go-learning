package shared

import (
	"modular-fx-fiber/internal/shared/database"
	"modular-fx-fiber/internal/shared/logger"
	"modular-fx-fiber/internal/shared/middleware"
	"modular-fx-fiber/internal/shared/swagger"

	"go.uber.org/fx"
)

// Module exports shared dependencies
var Module = fx.Options(
	fx.Provide(
		database.NewDatabase,
		logger.NewZapLogger,
		middleware.NewMiddleware,
		swagger.NewSwagger,
	),
	fx.Invoke(swagger.Register),
)
