package core

import (
	"echo-saas-starter/internal/config"

	"github.com/labstack/echo/v4"
)

// Plugin defines the interface all plugins must implement.
type Plugin interface {
	Name() string
	Version() string
	Init(cfg *config.Config, db DB) error
	RegisterRoutes(g *echo.Group)
	Migrate(db DB) error
	Shutdown() error
}

// StartupHook is an optional interface for plugins that need startup logic.
type StartupHook interface {
	OnStartup() error
}

// MiddlewareHook is an optional interface for plugins that register middleware.
type MiddlewareHook interface {
	OnMiddleware(e *echo.Echo)
}
