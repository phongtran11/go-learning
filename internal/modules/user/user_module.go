package user

import (
	"go.uber.org/fx"
)

// Module exports the user module dependencies
var Module = fx.Options(
	fx.Provide(
		NewUserRoutes,
		NewHandlers,
		NewService,
		NewRepository,
	),
	fx.Invoke(Register),
)
