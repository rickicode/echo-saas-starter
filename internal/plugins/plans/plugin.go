package plans

import (
	"echo-saas-starter/internal/config"
	"echo-saas-starter/internal/core"
	"echo-saas-starter/internal/plugins/auth"
	"echo-saas-starter/internal/plugins/roles"

	"github.com/labstack/echo/v4"
)

// Plugin implements the plans plugin.
type Plugin struct {
	cfg         *config.Config
	db          core.DB
	repo        *Repository
	handlers    *Handlers
	authPlugin  *auth.Plugin
	rolesPlugin *roles.Plugin
}

// NewPlugin creates a new plans plugin.
func NewPlugin(authPlugin *auth.Plugin, rolesPlugin *roles.Plugin) *Plugin {
	return &Plugin{authPlugin: authPlugin, rolesPlugin: rolesPlugin}
}

func (p *Plugin) Name() string    { return "plans" }
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
	adminMiddleware := roles.RequireRole(p.rolesPlugin.GetRepository(), "super_admin", "admin")

	// Public routes
	g.GET("/available", p.handlers.GetAvailablePlans)

	// User routes
	g.GET("/current", p.handlers.GetCurrentPlan, authMiddleware)

	// Admin routes
	g.GET("", p.handlers.ListPlans, authMiddleware, adminMiddleware)
	g.GET("/:id", p.handlers.GetPlan, authMiddleware, adminMiddleware)
	g.POST("", p.handlers.CreatePlan, authMiddleware, adminMiddleware)
	g.PUT("/:id", p.handlers.UpdatePlan, authMiddleware, adminMiddleware)
	g.DELETE("/:id", p.handlers.DeletePlan, authMiddleware, adminMiddleware)
	g.POST("/:id/assign", p.handlers.AssignPlan, authMiddleware, adminMiddleware)
}

func (p *Plugin) Migrate(db core.DB) error {
	return RunMigrations(db)
}

func (p *Plugin) Shutdown() error {
	return nil
}

// GetRepository returns the plans repository for use by other plugins.
func (p *Plugin) GetRepository() *Repository {
	return p.repo
}
