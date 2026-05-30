package auth

import (
	"echo-saas-starter/internal/config"
	"echo-saas-starter/internal/core"

	"github.com/labstack/echo/v4"
)

// Plugin implements the auth plugin.
type Plugin struct {
	cfg      *config.Config
	db       core.DB
	repo     *Repository
	handlers *Handlers
}

// NewPlugin creates a new auth plugin.
func NewPlugin() *Plugin {
	return &Plugin{}
}

func (p *Plugin) Name() string    { return "auth" }
func (p *Plugin) Version() string { return "1.0.0" }

func (p *Plugin) Init(cfg *config.Config, db core.DB) error {
	p.cfg = cfg
	p.db = db
	p.repo = NewRepository(db)
	p.handlers = NewHandlers(p.repo, cfg.PasetoSecret)
	return nil
}

func (p *Plugin) RegisterRoutes(g *echo.Group) {
	g.POST("/register", p.handlers.Register)
	g.POST("/login", p.handlers.Login)
	g.POST("/logout", p.handlers.Logout)
	g.POST("/refresh", p.handlers.Refresh)
	g.GET("/me", p.handlers.Me, RequireAuth(p.repo, p.cfg.PasetoSecret))
}

func (p *Plugin) Migrate(db core.DB) error {
	return RunMigrations(db)
}

func (p *Plugin) Shutdown() error {
	return nil
}

// GetRepository returns the auth repository for use by other plugins.
func (p *Plugin) GetRepository() *Repository {
	return p.repo
}

// GetSecret returns the PASETO secret for use by other plugins.
func (p *Plugin) GetSecret() string {
	return p.cfg.PasetoSecret
}
