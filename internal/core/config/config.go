package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

// Config holds application configuration
type Config struct {
	App AppConfig `mapstructure:"app"`
	DB  DBConfig  `mapstructure:"db"`
	JWT JWTConfig `mapstructure:"jwt"`
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

// NewConfig creates a new configuration instance
func NewConfig() (*Config, error) {
	// Get environment
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development" // Default to development
	}

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

	// Try to merge with environment config (ignore errors if file doesn't exist)
	_ = viper.MergeInConfig()

	// Override with environment variables
	viper.SetEnvPrefix("APP") // will convert APP_DB_HOST to db.host
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	var config Config
	err = viper.Unmarshal(&config)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	fmt.Printf("Config loaded: %+v\n", config)

	// Debug info for initialization (consider removing in production)
	fmt.Printf("Loaded configuration for environment: %s\n", env)

	return &config, nil
}
