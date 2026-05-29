package config

import (
	"os"
	"testing"
)

func TestLoad_WithEnvVars(t *testing.T) {
	// Use os.Setenv before calling Load to ensure viper picks them up
	os.Setenv("DATABASE_URL", "sqlite://test.db")
	os.Setenv("PASETO_SECRET", "my-secret-key-that-is-at-least-32-bytes-long")
	os.Setenv("SERVER_PORT", "9090")
	os.Setenv("LOG_LEVEL", "debug")
	defer func() {
		os.Unsetenv("DATABASE_URL")
		os.Unsetenv("PASETO_SECRET")
		os.Unsetenv("SERVER_PORT")
		os.Unsetenv("LOG_LEVEL")
	}()

	cfg, err := Load()
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if cfg.DatabaseURL != "sqlite://test.db" {
		t.Errorf("expected DatabaseURL = 'sqlite://test.db', got %q", cfg.DatabaseURL)
	}
	if cfg.PasetoSecret != "my-secret-key-that-is-at-least-32-bytes-long" {
		t.Errorf("expected PasetoSecret to be set, got %q", cfg.PasetoSecret)
	}
	if cfg.ServerPort != "9090" {
		t.Errorf("expected ServerPort = '9090', got %q", cfg.ServerPort)
	}
	if cfg.LogLevel != "debug" {
		t.Errorf("expected LogLevel = 'debug', got %q", cfg.LogLevel)
	}
}

func TestLoad_Defaults(t *testing.T) {
	os.Setenv("DATABASE_URL", "sqlite://test.db")
	defer os.Unsetenv("DATABASE_URL")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if cfg.ServerPort != "8080" {
		t.Errorf("expected default ServerPort = '8080', got %q", cfg.ServerPort)
	}
	if cfg.ServerHost != "0.0.0.0" {
		t.Errorf("expected default ServerHost = '0.0.0.0', got %q", cfg.ServerHost)
	}
	if cfg.CorsOrigins != "*" {
		t.Errorf("expected default CorsOrigins = '*', got %q", cfg.CorsOrigins)
	}
}

func TestLoad_MissingDatabaseURL(t *testing.T) {
	os.Unsetenv("DATABASE_URL")

	_, err := Load()
	if err == nil {
		t.Fatal("expected error when DATABASE_URL is empty")
	}
}

func TestLoad_ShortPasetoSecret(t *testing.T) {
	os.Setenv("DATABASE_URL", "sqlite://test.db")
	os.Setenv("PASETO_SECRET", "short-secret")
	defer func() {
		os.Unsetenv("DATABASE_URL")
		os.Unsetenv("PASETO_SECRET")
	}()

	_, err := Load()
	if err == nil {
		t.Fatal("expected error when PASETO_SECRET is too short")
	}
}
