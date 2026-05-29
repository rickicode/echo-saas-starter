package plugins

import (
	"echo-saas-starter/internal/core"
	"echo-saas-starter/internal/plugins/auth"
	"echo-saas-starter/internal/plugins/plans"
	"echo-saas-starter/internal/plugins/roles"
	"echo-saas-starter/internal/plugins/users"
)

// RegisterAll registers all available plugins with the plugin manager.
func RegisterAll(pm *core.PluginManager) {
	authPlugin := auth.NewPlugin()
	rolesPlugin := roles.NewPlugin(authPlugin)

	pm.Register(authPlugin)
	pm.Register(rolesPlugin)
	pm.Register(users.NewPlugin(authPlugin, rolesPlugin))
	pm.Register(plans.NewPlugin(authPlugin, rolesPlugin))
}
