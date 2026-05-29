package server

import (
	"context"
	"echo-saas-starter/internal/config"
	"echo-saas-starter/internal/core"
	"echo-saas-starter/internal/embed"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

// Server holds the Echo instance and dependencies.
type Server struct {
	echo          *echo.Echo
	cfg           *config.Config
	db            core.DB
	pluginManager *core.PluginManager
}

// NewServer creates a new Server instance.
func NewServer(cfg *config.Config, db core.DB, pm *core.PluginManager) *Server {
	return &Server{
		echo:          echo.New(),
		cfg:           cfg,
		db:            db,
		pluginManager: pm,
	}
}

// Start initializes the Echo server with middleware, routes, and SPA handler.
func (s *Server) Start() error {
	e := s.echo
	e.HideBanner = true

	// Apply core middleware
	e.Use(core.RequestIDMiddleware())
	e.Use(core.RecoveryMiddleware())
	e.Use(core.LoggingMiddleware())
	e.Use(core.CORSMiddleware(s.cfg.CorsOrigins))

	// Let plugins add middleware
	s.pluginManager.OnMiddleware(e)

	// Mount plugin routes
	s.pluginManager.RegisterRoutes(e)

	// Health check
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	})

	// Serve embedded SPA for non-API routes
	e.GET("/*", embed.SPAHandler())

	// Run startup hooks
	if err := s.pluginManager.OnStartup(); err != nil {
		return fmt.Errorf("startup hooks failed: %w", err)
	}

	addr := fmt.Sprintf("%s:%s", s.cfg.ServerHost, s.cfg.ServerPort)
	return e.Start(addr)
}

// Shutdown gracefully shuts down the server.
func (s *Server) Shutdown(ctx context.Context) error {
	return s.echo.Shutdown(ctx)
}
