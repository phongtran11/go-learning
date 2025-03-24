package auth

import (
	"go.uber.org/fx"
)

// Module exports the auth module dependencies
var Module = fx.Options(
	fx.Provide(
		NewRoutes,
		NewHandlers,
		NewService,
	),
	fx.Invoke(Register),
)
