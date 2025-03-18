.PHONY: dev swagger migration-gen migrate-up migrate-down migrate-version


dev:
	@echo "Starting development server"
	@air -c .air.toml 

swagger:
	@echo "Generating swagger"
	@swag init -g cmd/server/main.go

migrate-gen:
	@echo "Generating migration file"
	@read -p "Enter migration name: " name; \
	migrate create -ext sql -dir internal/shared/database/migrations -seq $$name

migrate-up:
	@go run cmd/migration/main.go -action up

migrate-down:
	@go run cmd/migration/main.go -action down

migrate-version:
	@go run cmd/migration/main.go -action version