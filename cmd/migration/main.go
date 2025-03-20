package main

import (
	"flag"
	"fmt"
	"log"
	"modular-fx-fiber/internal/core/config"
	"modular-fx-fiber/internal/shared/logger"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

// Command line flags
var (
	migrationPath string
	command       string
	name          string
)

func init() {
	flag.StringVar(&migrationPath, "dir", "internal/shared/database/migrations", "Directory with migration files")
	flag.StringVar(&command, "cmd", "help", "Migration command (up, down, status, create, help)")
	flag.StringVar(&name, "name", "", "Name for new migration (for create command)")
}

func Run() {
	flag.Parse()

	l := logger.NewZapLogger()
	// Load configuration
	cfg, err := config.NewConfig(l)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Build database connection string
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.DB.USER, cfg.DB.PASSWORD, cfg.DB.HOST, cfg.DB.PORT, cfg.DB.NAME, cfg.DB.SSL)

	// Create migration directory if it doesn't exist
	if _, err := os.Stat(migrationPath); os.IsNotExist(err) {
		if err := os.MkdirAll(migrationPath, 0755); err != nil {
			log.Fatalf("Failed to create migrations directory: %v", err)
		}
	}

	// Set goose logger
	goose.SetLogger(log.New(os.Stdout, "", log.LstdFlags))

	// Initialize goose
	if err := goose.SetDialect("pgx"); err != nil {
		log.Fatalf("Failed to set dialect: %v", err)
	}

	// Connect to database
	db, err := goose.OpenDBWithDriver("pgx", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Execute command
	switch command {
	case "up":
		err = goose.Up(db, migrationPath)

	case "down":
		err = goose.Down(db, migrationPath)

	case "status":
		err = goose.Status(db, migrationPath)

	case "create":
		if name == "" {
			log.Fatal("Migration name is required for create command")
		}
		err = goose.Create(db, migrationPath, name, "sql")

	case "reset":
		err = goose.Reset(db, migrationPath)

	case "version":
		var version int64
		version, err = goose.GetDBVersion(db)
		if err == nil {
			fmt.Printf("Current version: %d\n", version)
		}

	default:
		fmt.Println("Usage: migration -cmd=COMMAND [options]")
		fmt.Println("\nCommands:")
		fmt.Println("  up      Migrate the DB to the most recent version")
		fmt.Println("  down    Roll back the version by 1")
		fmt.Println("  status  Display migration status")
		fmt.Println("  create  Create a new migration file (requires -name)")
		fmt.Println("  reset   Roll back all migrations")
		fmt.Println("  version Display current migration version")
		fmt.Println("  help    Show this help")
		fmt.Println("\nOptions:")
		flag.PrintDefaults()
		os.Exit(0)
	}

	if err != nil {
		log.Fatalf("Migration failed: %v", err)
	} else {
		fmt.Println("Migration command completed successfully")
	}
}

// Main function for migration command
func main() {
	Run()
}
