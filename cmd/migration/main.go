package main

import (
	"flag"
	"fmt"
	"log"
	"modular-fx-fiber/internal/core/config"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// Command line flags
var (
	migrationPath string
	action        string
	step          int
	version       uint
)

func init() {
	flag.StringVar(&migrationPath, "path", "internal/shared/database/migrations", "Migration files path")
	flag.StringVar(&action, "action", "up", "Migration action (up, down, version, force)")
	flag.IntVar(&step, "step", 0, "Number of migrations to apply (for up/down with step)")
	flag.UintVar(&version, "version", 0, "Target migration version (for force, goto)")
}

func Run() {
	flag.Parse()

	// Load configuration
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Build database connection string
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.DB.USER, cfg.DB.PASSWORD, cfg.DB.HOST, cfg.DB.PORT, cfg.DB.NAME, cfg.DB.SSL)

	// Create migration instance
	m, err := migrate.New(fmt.Sprintf("file://%s", migrationPath), dsn)
	if err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	// Set logger
	m.Log = &MigrationLogger{}

	// Execute migration based on action
	switch action {
	case "up":
		if step > 0 {
			err = m.Steps(step)
		} else {
			err = m.Up()
		}
	case "down":
		if step > 0 {
			err = m.Steps(-step)
		} else {
			err = m.Down()
		}
	case "version":
		version, dirty, vErr := m.Version()
		if vErr != nil {
			if vErr == migrate.ErrNilVersion {
				fmt.Println("No migrations applied")
			} else {
				log.Fatalf("Failed to get migration version: %v", vErr)
			}
		} else {
			fmt.Printf("Current migration version: %d (dirty: %t)\n", version, dirty)
		}
		return
	case "force":
		err = m.Force(int(version))
	case "goto":
		err = m.Migrate(version)
	default:
		log.Fatalf("Unknown action: %s", action)
	}

	if err != nil {
		if err == migrate.ErrNoChange {
			fmt.Println("No migration changes required")
		} else {
			log.Fatalf("Migration failed: %v", err)
		}
	} else {
		fmt.Println("Migration completed successfully")
	}
}

// Custom logger for migrations
type MigrationLogger struct{}

func (l *MigrationLogger) Printf(format string, v ...interface{}) {
	log.Printf(format, v...)
}

func (l *MigrationLogger) Verbose() bool {
	return true
}

// Main function for migration command
func main() {
	Run()
}
