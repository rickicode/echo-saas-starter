package plugins

import (
	"echo-saas-starter/internal/core"
	"echo-saas-starter/internal/plugins/auth"
	"echo-saas-starter/internal/plugins/roles"
)

// RegisterAll registers all available plugins with the plugin manager.
func RegisterAll(pm *core.PluginManager) {
	authPlugin := auth.NewPlugin()
	pm.Register(authPlugin)
	pm.Register(roles.NewPlugin(authPlugin))
}
