package database

import (
	"context"
	"echo-saas-starter/internal/core"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// PostgresDB implements core.DB using pgx connection pool.
type PostgresDB struct {
	pool *pgxpool.Pool
}

// pgxRow wraps pgx.Row to implement core.Row.
type pgxRow struct {
	row pgx.Row
}

func (r *pgxRow) Scan(dest ...interface{}) error {
	return r.row.Scan(dest...)
}

// pgxRows wraps pgx.Rows to implement core.Rows.
type pgxRows struct {
	rows pgx.Rows
}

func (r *pgxRows) Next() bool {
	return r.rows.Next()
}

func (r *pgxRows) Scan(dest ...interface{}) error {
	return r.rows.Scan(dest...)
}

func (r *pgxRows) Close() error {
	r.rows.Close()
	return nil
}

func (r *pgxRows) Err() error {
	return r.rows.Err()
}

// NewPostgresDB creates a new PostgreSQL database connection pool.
func NewPostgresDB(databaseURL string) (core.DB, error) {
	pool, err := pgxpool.New(context.Background(), databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgres: %w", err)
	}

	if err := pool.Ping(context.Background()); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to ping postgres: %w", err)
	}

	return &PostgresDB{pool: pool}, nil
}

func (db *PostgresDB) Exec(ctx context.Context, query string, args ...interface{}) error {
	_, err := db.pool.Exec(ctx, RewritePlaceholders(query), args...)
	return err
}

func (db *PostgresDB) Query(ctx context.Context, query string, args ...interface{}) (core.Rows, error) {
	rows, err := db.pool.Query(ctx, RewritePlaceholders(query), args...)
	if err != nil {
		return nil, err
	}
	return &pgxRows{rows: rows}, nil
}

func (db *PostgresDB) QueryRow(ctx context.Context, query string, args ...interface{}) core.Row {
	row := db.pool.QueryRow(ctx, RewritePlaceholders(query), args...)
	return &pgxRow{row: row}
}

func (db *PostgresDB) Close() error {
	db.pool.Close()
	return nil
}

func (db *PostgresDB) GetDB() interface{} {
	return db.pool
}
