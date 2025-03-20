package auth

import (
	"go.uber.org/fx"
)

// Module exports the auth module dependencies
var Module = fx.Options(
	fx.Provide(
		NewAuthRoutes,
		NewHandlers,
		NewService,
		NewRepository,
	),
)
