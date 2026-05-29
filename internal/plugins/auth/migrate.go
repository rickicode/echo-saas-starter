package auth

import (
	"context"
	"echo-saas-starter/internal/core"
	"embed"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

// RunMigrations executes all embedded SQL migration files.
func RunMigrations(db core.DB) error {
	entries, err := migrationsFS.ReadDir("migrations")
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		data, err := migrationsFS.ReadFile("migrations/" + entry.Name())
		if err != nil {
			return err
		}
		if err := db.Exec(context.Background(), string(data)); err != nil {
			return err
		}
	}
	return nil
}
