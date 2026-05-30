package docs

import (
	"echo-saas-starter/internal/config"
	"echo-saas-starter/internal/core"
	"echo-saas-starter/internal/plugins/auth"
	"echo-saas-starter/internal/plugins/roles"

	"github.com/labstack/echo/v4"
)

// Plugin implements the docs plugin.
type Plugin struct {
	cfg         *config.Config
	db          core.DB
	repo        *Repository
	handlers    *Handlers
	authPlugin  *auth.Plugin
	rolesPlugin *roles.Plugin
}

// NewPlugin creates a new docs plugin.
func NewPlugin(authPlugin *auth.Plugin, rolesPlugin *roles.Plugin) *Plugin {
	return &Plugin{authPlugin: authPlugin, rolesPlugin: rolesPlugin}
}

func (p *Plugin) Name() string    { return "docs" }
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

	// Admin post routes
	g.GET("/posts", p.handlers.ListPostsAdmin, authMiddleware, adminMiddleware)
	g.GET("/posts/:id", p.handlers.GetPostAdmin, authMiddleware, adminMiddleware)
	g.POST("/posts", p.handlers.CreatePost, authMiddleware, adminMiddleware)
	g.PUT("/posts/:id", p.handlers.UpdatePost, authMiddleware, adminMiddleware)
	g.DELETE("/posts/:id", p.handlers.DeletePost, authMiddleware, adminMiddleware)
	g.POST("/posts/:id/publish", p.handlers.PublishPost, authMiddleware, adminMiddleware)
	g.POST("/posts/:id/unpublish", p.handlers.UnpublishPost, authMiddleware, adminMiddleware)

	// Admin category routes
	g.GET("/categories", p.handlers.ListCategories, authMiddleware, adminMiddleware)
	g.POST("/categories", p.handlers.CreateCategory, authMiddleware, adminMiddleware)
	g.PUT("/categories/:id", p.handlers.UpdateCategory, authMiddleware, adminMiddleware)
	g.DELETE("/categories/:id", p.handlers.DeleteCategory, authMiddleware, adminMiddleware)

	// Admin tag routes
	g.GET("/tags", p.handlers.ListTags, authMiddleware, adminMiddleware)
	g.POST("/tags", p.handlers.CreateTag, authMiddleware, adminMiddleware)
	g.DELETE("/tags/:id", p.handlers.DeleteTag, authMiddleware, adminMiddleware)

	// Public routes (no auth)
	g.GET("/public", p.handlers.ListPublishedPosts)
	g.GET("/public/search", p.handlers.SearchPosts)
	g.GET("/public/:slug", p.handlers.GetPublishedPost)
}

func (p *Plugin) Migrate(db core.DB) error {
	return RunMigrations(db)
}

func (p *Plugin) Shutdown() error {
	return nil
}
