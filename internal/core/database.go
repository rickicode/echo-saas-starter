package core

import "context"

// Row represents a single database row.
type Row interface {
	Scan(dest ...interface{}) error
}

// Rows represents multiple database rows.
type Rows interface {
	Next() bool
	Scan(dest ...interface{}) error
	Close() error
	Err() error
}

// DB defines the database abstraction interface.
type DB interface {
	Exec(ctx context.Context, query string, args ...interface{}) error
	Query(ctx context.Context, query string, args ...interface{}) (Rows, error)
	QueryRow(ctx context.Context, query string, args ...interface{}) Row
	Close() error
	GetDB() interface{}
}
