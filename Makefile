.PHONY: dev swagger


dev:
	@echo "Starting development server"
	@air -c .air.toml 

swagger:
	@echo "Generating swagger"
	@swag init -g cmd/server/main.go