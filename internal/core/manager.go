package core

import (
	"echo-saas-starter/internal/config"
	"fmt"

	"github.com/labstack/echo/v4"
)

// PluginManager handles plugin registration, initialization, and lifecycle.
type PluginManager struct {
	plugins []Plugin
}

// NewPluginManager creates a new plugin manager.
func NewPluginManager() *PluginManager {
	return &PluginManager{
		plugins: make([]Plugin, 0),
	}
}

// Register adds a plugin to the manager.
func (pm *PluginManager) Register(p Plugin) {
	pm.plugins = append(pm.plugins, p)
}

// InitAll initializes all registered plugins.
func (pm *PluginManager) InitAll(cfg *config.Config, db DB) error {
	for _, p := range pm.plugins {
		if err := p.Init(cfg, db); err != nil {
			return fmt.Errorf("failed to init plugin %s: %w", p.Name(), err)
		}
	}
	return nil
}

// MigrateAll runs migrations for all registered plugins.
func (pm *PluginManager) MigrateAll(db DB) error {
	for _, p := range pm.plugins {
		if err := p.Migrate(db); err != nil {
			return fmt.Errorf("failed to migrate plugin %s: %w", p.Name(), err)
		}
	}
	return nil
}

// RegisterRoutes mounts all plugin routes under /api/v1/<plugin-name>.
func (pm *PluginManager) RegisterRoutes(e *echo.Echo) {
	api := e.Group("/api/v1")
	for _, p := range pm.plugins {
		group := api.Group("/" + p.Name())
		p.RegisterRoutes(group)
	}
}

// OnMiddleware calls MiddlewareHook for plugins that implement it.
func (pm *PluginManager) OnMiddleware(e *echo.Echo) {
	for _, p := range pm.plugins {
		if mh, ok := p.(MiddlewareHook); ok {
			mh.OnMiddleware(e)
		}
	}
}

// OnStartup calls StartupHook for plugins that implement it.
func (pm *PluginManager) OnStartup() error {
	for _, p := range pm.plugins {
		if sh, ok := p.(StartupHook); ok {
			if err := sh.OnStartup(); err != nil {
				return fmt.Errorf("startup hook failed for plugin %s: %w", p.Name(), err)
			}
		}
	}
	return nil
}

// Shutdown shuts down all plugins in reverse order.
func (pm *PluginManager) Shutdown() error {
	for i := len(pm.plugins) - 1; i >= 0; i-- {
		if err := pm.plugins[i].Shutdown(); err != nil {
			return fmt.Errorf("failed to shutdown plugin %s: %w", pm.plugins[i].Name(), err)
		}
	}
	return nil
}
