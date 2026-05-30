package users

import (
	"echo-saas-starter/internal/config"
	"echo-saas-starter/internal/core"
	"echo-saas-starter/internal/plugins/auth"
	"echo-saas-starter/internal/plugins/roles"

	"github.com/labstack/echo/v4"
)

// Plugin implements the users plugin.
type Plugin struct {
	cfg        *config.Config
	db         core.DB
	repo       *Repository
	handlers   *Handlers
	authPlugin *auth.Plugin
	rolesPlugin *roles.Plugin
}

// NewPlugin creates a new users plugin.
func NewPlugin(authPlugin *auth.Plugin, rolesPlugin *roles.Plugin) *Plugin {
	return &Plugin{authPlugin: authPlugin, rolesPlugin: rolesPlugin}
}

func (p *Plugin) Name() string    { return "users" }
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

	// User self-service routes (must be before :id routes)
	g.GET("/me", p.handlers.GetMyProfile, authMiddleware)
	g.PUT("/me", p.handlers.UpdateMyProfile, authMiddleware)
	g.POST("/me/avatar", p.handlers.UploadAvatar, authMiddleware)
	g.PUT("/me/password", p.handlers.ChangePassword, authMiddleware)

	// Admin routes
	g.GET("", p.handlers.ListUsers, authMiddleware, adminMiddleware)
	g.GET("/stats", p.handlers.GetStats, authMiddleware, adminMiddleware)
	g.GET("/:id", p.handlers.GetUser, authMiddleware, adminMiddleware)
	g.PUT("/:id", p.handlers.UpdateUser, authMiddleware, adminMiddleware)
	g.DELETE("/:id", p.handlers.DeleteUser, authMiddleware, adminMiddleware)
}

func (p *Plugin) Migrate(db core.DB) error {
	return RunMigrations(db)
}

func (p *Plugin) Shutdown() error {
	return nil
}
