.PHONY: dev swagger migrate-up migrate-down migrate-status migrate-create migrate-reset migrate-version


dev:
	@echo "Starting development server"
	@air -c .air.toml 

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