package user

import (
	"go.uber.org/fx"
)

// Module exports the user module dependencies
var Module = fx.Options(
	fx.Provide(
		NewRoutes,
		NewHandlers,
		NewService,
	),
	fx.Invoke(Register),
)
