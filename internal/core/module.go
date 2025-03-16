package core

import (
	"modular-fx-fiber/internal/core/config"
	"modular-fx-fiber/internal/core/server"

	"go.uber.org/fx"
)

// Module exports the core module dependencies
var Module = fx.Options(
	fx.Provide(
		config.NewConfig,
		server.NewServer,
	),
	fx.Invoke(server.Start),
)
