package core

import (
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// CORSMiddleware returns CORS middleware with configurable origins.
func CORSMiddleware(origins string) echo.MiddlewareFunc {
	allowOrigins := strings.Split(origins, ",")
	return middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: allowOrigins,
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowHeaders: []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Request-ID"},
	})
}

// LoggingMiddleware returns a logging middleware that logs method, path, status, and latency.
func LoggingMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()
			err := next(c)
			latency := time.Since(start)

			req := c.Request()
			res := c.Response()

			log.Printf("%s %s %d %s", req.Method, req.URL.Path, res.Status, latency)
			return err
		}
	}
}

// RecoveryMiddleware returns panic recovery middleware.
func RecoveryMiddleware() echo.MiddlewareFunc {
	return middleware.Recover()
}

// RequestIDMiddleware returns middleware that adds a UUID to X-Request-ID header.
func RequestIDMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			reqID := c.Request().Header.Get("X-Request-ID")
			if reqID == "" {
				reqID = uuid.New().String()
			}
			c.Request().Header.Set("X-Request-ID", reqID)
			c.Response().Header().Set("X-Request-ID", reqID)
			return next(c)
		}
	}
}
