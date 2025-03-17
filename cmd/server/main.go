package main

import (
	"modular-fx-fiber/internal/core"
	"modular-fx-fiber/internal/modules/user"
	"modular-fx-fiber/internal/shared"

	"go.uber.org/fx"
)

// @title Your API
// @version 1.0
// @description API Documentation
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email phongtran11.tt@gmail.com
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8080
// @BasePath /api
func main() {
	fx.New(
		// Core module
		core.Module,

		// Shared services
		shared.Module,

		// Feature modules
		user.Module,
	).Run()
}
