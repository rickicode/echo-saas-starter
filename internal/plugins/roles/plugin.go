package roles

import (
	"echo-saas-starter/internal/config"
	"echo-saas-starter/internal/core"
	"echo-saas-starter/internal/plugins/auth"

	"github.com/labstack/echo/v4"
)

// Plugin implements the roles plugin.
type Plugin struct {
	cfg        *config.Config
	db         core.DB
	repo       *Repository
	handlers   *Handlers
	authPlugin *auth.Plugin
}

// NewPlugin creates a new roles plugin.
func NewPlugin(authPlugin *auth.Plugin) *Plugin {
	return &Plugin{authPlugin: authPlugin}
}

func (p *Plugin) Name() string    { return "roles" }
func (p *Plugin) Version() string { return "1.0.0" }

func (p *Plugin) Init(cfg *config.Config, db core.DB) error {
	p.cfg = cfg
	p.db = db
	p.repo = NewRepository(db)
	p.handlers = NewHandlers(p.repo)
	return nil
}

func (p *Plugin) RegisterRoutes(g *echo.Group) {
	authMiddleware := auth.RequireAuth(p.authPlugin.GetRepository(), p.authPlugin.GetSecret())
	adminMiddleware := RequireRole(p.repo, "super_admin", "admin")

	g.GET("/me", p.handlers.GetMyRolesHandler, authMiddleware)
	g.GET("", p.handlers.ListRoles, authMiddleware)
	g.POST("/:id/assign", p.handlers.AssignRoleHandler, authMiddleware, adminMiddleware)
	g.DELETE("/:id/users/:user_id", p.handlers.RemoveRoleHandler, authMiddleware, adminMiddleware)

	// User roles endpoint under the roles group
	g.GET("/users/:id/roles", p.handlers.GetUserRolesHandler, authMiddleware)
}

func (p *Plugin) Migrate(db core.DB) error {
	return RunMigrations(db)
}

func (p *Plugin) Shutdown() error {
	return nil
}

// GetRepository returns the roles repository for use by other plugins.
func (p *Plugin) GetRepository() *Repository {
	return p.repo
}
