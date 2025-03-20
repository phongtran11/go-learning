package config

import (
	"fmt"
	"modular-fx-fiber/internal/shared/logger"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// Config holds application configuration
type Config struct {
	App  AppConfig  `mapstructure:"app"`
	DB   DBConfig   `mapstructure:"db"`
	JWT  JWTConfig  `mapstructure:"jwt"`
	Mail MailConfig `mapstructure:"mail"`
}

type AppConfig struct {
	Port string `mapstructure:"port"`
	Name string `mapstructure:"name"`
	Env  string `mapstructure:"env"`
}

type DBConfig struct {
	HOST     string `mapstructure:"host"`
	PORT     string `mapstructure:"port"`
	USER     string `mapstructure:"user"`
	PASSWORD string `mapstructure:"password"`
	NAME     string `mapstructure:"name"`
	SSL      string `mapstructure:"ssl"`
}

type JWTConfig struct {
	Secret              string `mapstructure:"secret"`
	AccessExpiryMinutes int    `mapstructure:"access_expiry_minutes"`
	RefreshExpiryDays   int    `mapstructure:"refresh_expiry_days"`
}

type MailConfig struct {
	FromAddr     string `mapstructure:"from_addr"`
	FromName     string `mapstructure:"from_name"`
	SMTPServer   string `mapstructure:"smtp_server"`
	SMTPPort     int    `mapstructure:"smtp_port"`
	SMTPUsername string `mapstructure:"smtp_username"`
	SMTPPassword string `mapstructure:"smtp_password"`
}

// NewConfig creates a new configuration instance
func NewConfig(l *logger.ZapLogger) (*Config, error) {
	// Get environment
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development" // Default to development
	}

	godotenv.Load()

	// Set default values
	viper.SetDefault("app.env", env)

	// Use base config
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	// Look for config in multiple locations
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")
	viper.AddConfigPath("./internal/core/config")

	// Read base config
	err := viper.ReadInConfig()
	if err != nil {
		return nil, fmt.Errorf("fatal error config file: %w", err)
	}

	// Load environment-specific config if it exists
	envConfigName := fmt.Sprintf("config.%s", env)
	viper.SetConfigName(envConfigName)

	// Override with environment variables
	viper.SetEnvPrefix("APP") // will convert APP_DB_HOST to db.host
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	var config Config
	err = viper.Unmarshal(&config)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	if env == "development" {
		l.Info("Loaded configuration for environment: ", zap.Any("env", env))
		l.Info("App: ", zap.Any("config", &config))

	}

	// Debug info for initialization (consider removing in production)
	fmt.Printf("Loaded configuration for environment: %s\n", env)

	return &config, nil
}
