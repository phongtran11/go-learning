.PHONY: dev swagger migrate-up migrate-down migrate-status migrate-create migrate-reset migrate-version docker-dev docker-prod

dev:
	@echo "Starting development server"
	@which air > /dev/null || (echo "Installing air..." && go install github.com/air-verse/air)
	@air -c .air.toml

docker-prod:
	@echo "Starting production environment with Docker"
	@cd deployment/docker && docker-compose up --build

swagger:
	@echo "Generating swagger"
	@swag init -g cmd/server/main.go

migrate-up:
	@go run cmd/migration/main.go -cmd up

migrate-down:
	@go run cmd/migration/main.go -cmd down

migrate-status:
	@go run cmd/migration/main.go -cmd status

migrate-create:
	@read -p "Enter migration name: " name; \
    go run cmd/migration/main.go -cmd create -name $$name

migrate-reset:
	@go run cmd/migration/main.go -cmd reset

migrate-version:
	@go run cmd/migration/main.go -cmd version