package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// Config holds application configuration
type Config struct {
	App AppConfig `yaml:"app"`
	DB  DBConfig  `yaml:"db"`
}

type AppConfig struct {
	Port string `yaml:"port"`
	Name string `yaml:"name"`
	Env  string `yaml:"env"`
}

type DBConfig struct {
	HOST     string `yaml:"host"`
	PORT     string `yaml:"port"`
	USER     string `yaml:"user"`
	PASSWORD string `yaml:"password"`
	NAME     string `yaml:"name"`
	SSL      string `yaml:"ssl"`
}

// NewConfig creates a new configuration instance
func NewConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("internal/core/config")
	err := viper.ReadInConfig()
	if err != nil {
		return nil, fmt.Errorf("fatal error config file: %w", err)
	}

	var config Config
	err = viper.Unmarshal(&config)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &config, nil
}
