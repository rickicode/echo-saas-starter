package embed

import (
	"embed"
	"io/fs"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

//go:embed all:dist
var staticFS embed.FS

// SPAHandler returns an echo handler that serves the embedded SPA.
// It serves static files from dist/ and falls back to index.html
// for client-side routing.
func SPAHandler() echo.HandlerFunc {
	distFS, _ := fs.Sub(staticFS, "dist")
	fileServer := http.FileServer(http.FS(distFS))

	return func(c echo.Context) error {
		path := c.Request().URL.Path

		// Try to serve the file directly
		if path != "/" && !strings.HasSuffix(path, "/") {
			if file, err := distFS.Open(strings.TrimPrefix(path, "/")); err == nil {
				file.Close()
				fileServer.ServeHTTP(c.Response(), c.Request())
				return nil
			}
		}

		// Fallback to index.html for SPA routing
		index, err := distFS.Open("index.html")
		if err != nil {
			return echo.NewHTTPError(http.StatusNotFound, "index.html not found")
		}
		index.Close()

		c.Request().URL.Path = "/"
		fileServer.ServeHTTP(c.Response(), c.Request())
		return nil
	}
}
