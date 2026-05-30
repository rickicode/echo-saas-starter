package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// Config holds application configuration.
type Config struct {
	DatabaseURL  string `mapstructure:"DATABASE_URL"`
	PasetoSecret string `mapstructure:"PASETO_SECRET"`
	ServerPort   string `mapstructure:"SERVER_PORT"`
	ServerHost   string `mapstructure:"SERVER_HOST"`
	LogLevel     string `mapstructure:"LOG_LEVEL"`
	CorsOrigins  string `mapstructure:"CORS_ORIGINS"`
}

// Load reads configuration from environment variables and optional config file.
func Load() (*Config, error) {
	v := viper.New()

	// Defaults
	v.SetDefault("SERVER_PORT", "8080")
	v.SetDefault("SERVER_HOST", "0.0.0.0")
	v.SetDefault("LOG_LEVEL", "info")
	v.SetDefault("CORS_ORIGINS", "*")
	v.SetDefault("DATABASE_URL", "")
	v.SetDefault("PASETO_SECRET", "")

	// Environment variables
	v.AutomaticEnv()

	// Optional config file
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")
	v.AddConfigPath("./config")
	_ = v.ReadInConfig() // ignore error if file not found

	cfg := &Config{}
	if err := v.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	if cfg.DatabaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL is required")
	}

	if cfg.PasetoSecret == "" {
		return nil, fmt.Errorf("PASETO_SECRET is required")
	}
	if len(cfg.PasetoSecret) < 32 {
		return nil, fmt.Errorf("PASETO_SECRET must be at least 32 bytes long")
	}

	return cfg, nil
}
