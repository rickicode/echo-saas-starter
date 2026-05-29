package database

import (
	"echo-saas-starter/internal/core"
	"fmt"
	"strings"
)

// NewDB auto-detects the database driver from the URL scheme and returns
// the appropriate DB implementation.
func NewDB(databaseURL string) (core.DB, error) {
	switch {
	case strings.HasPrefix(databaseURL, "postgres://"),
		strings.HasPrefix(databaseURL, "postgresql://"):
		return NewPostgresDB(databaseURL)
	case strings.HasPrefix(databaseURL, "sqlite://"),
		strings.HasPrefix(databaseURL, "file:"):
		return NewSQLiteDB(databaseURL)
	default:
		return nil, fmt.Errorf("unsupported database URL scheme: %s", databaseURL)
	}
}
