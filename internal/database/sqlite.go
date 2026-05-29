package database

import (
	"context"
	"database/sql"
	"echo-saas-starter/internal/core"
	"fmt"
	"strings"

	_ "modernc.org/sqlite"
)

// SQLiteDB implements core.DB using modernc.org/sqlite.
type SQLiteDB struct {
	db *sql.DB
}

// sqliteRow wraps *sql.Row to implement core.Row.
type sqliteRow struct {
	row *sql.Row
}

func (r *sqliteRow) Scan(dest ...interface{}) error {
	return r.row.Scan(dest...)
}

// sqliteRows wraps *sql.Rows to implement core.Rows.
type sqliteRows struct {
	rows *sql.Rows
}

func (r *sqliteRows) Next() bool {
	return r.rows.Next()
}

func (r *sqliteRows) Scan(dest ...interface{}) error {
	return r.rows.Scan(dest...)
}

func (r *sqliteRows) Close() error {
	return r.rows.Close()
}

func (r *sqliteRows) Err() error {
	return r.rows.Err()
}

// NewSQLiteDB creates a new SQLite database connection with WAL mode.
func NewSQLiteDB(databaseURL string) (core.DB, error) {
	// Parse the file path from the URL
	dsn := databaseURL
	if strings.HasPrefix(dsn, "sqlite://") {
		dsn = strings.TrimPrefix(dsn, "sqlite://")
	} else if strings.HasPrefix(dsn, "file:") {
		dsn = strings.TrimPrefix(dsn, "file:")
	}

	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open sqlite: %w", err)
	}

	// Enable WAL mode
	if _, err := db.Exec("PRAGMA journal_mode=WAL"); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to enable WAL mode: %w", err)
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping sqlite: %w", err)
	}

	return &SQLiteDB{db: db}, nil
}

func (s *SQLiteDB) Exec(ctx context.Context, query string, args ...interface{}) error {
	_, err := s.db.ExecContext(ctx, query, args...)
	return err
}

func (s *SQLiteDB) Query(ctx context.Context, query string, args ...interface{}) (core.Rows, error) {
	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	return &sqliteRows{rows: rows}, nil
}

func (s *SQLiteDB) QueryRow(ctx context.Context, query string, args ...interface{}) core.Row {
	row := s.db.QueryRowContext(ctx, query, args...)
	return &sqliteRow{row: row}
}

func (s *SQLiteDB) Close() error {
	return s.db.Close()
}

func (s *SQLiteDB) GetDB() interface{} {
	return s.db
}
