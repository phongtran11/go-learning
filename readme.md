# API

A highly scalable, modular REST API template built with Go Fiber and Uber FX for dependency injection.

## 🌟 Features

- **Modular Architecture**: Clean separation of concerns with module-based structure
- **Dependency Injection**: Using Uber FX for robust and testable dependency management
- **JWT Authentication**: Complete authentication system with login, register, and refresh token
- **Database Integration**: PostgreSQL with GORM and migrations
- **API Documentation**: Integrated Swagger documentation
- **Hot Reload**: Automatic server restarts during development with Air
- **Logging**: Structured logging with Zap logger
- **Validation**: Request validation using validator
- **Docker Support**: Ready to run in containers for any environment

## 🚀 Getting Started

### Prerequisites

- Go 1.20 or later
- PostgreSQL
- Air for hot reload (optional)

### Setup

1. Clone the repository
2. Install dependencies
3. Set up your database
4. Run the application

## 🛠️ Development

### Creating a Migration

```
make migrate-gen
```

### Running Migrations

```
make migrate-up
```

### Generating Swagger Documentation

```
make swagger
```

### Start Dev Server

```
make dev
```

## 🔐 JWT Authentication

This project uses JWT for authentication:

- Access tokens expire after 60 minutes (configurable)
- Refresh tokens expire after 7 days (configurable)
- Refresh token rotation is implemented for security

## 📚 Used Libraries

- Go Fiber - Web framework
- Uber FX - Dependency injection
- GORM - ORM library
- golang-migrate - Database migrations
- JWT - JWT implementation
- Viper - Configuration management
- Zap - Structured logging
- Validator - Request validation
- Swagger - API documentation
