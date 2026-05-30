package database

import (
	"strings"
	"testing"
)

func TestNewDB_PostgresScheme(t *testing.T) {
	tests := []struct {
		name string
		url  string
	}{
		{"postgres://", "postgres://user:pass@localhost:5432/db"},
		{"postgresql://", "postgresql://user:pass@localhost:5432/db"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// NewDB will fail to connect since there's no real postgres,
			// but it should attempt to use the postgres driver (not return "unsupported scheme")
			_, err := NewDB(tt.url)
			if err == nil {
				t.Skip("unexpectedly connected to postgres")
			}
			if strings.Contains(err.Error(), "unsupported database URL scheme") {
				t.Errorf("expected postgres driver to be selected, got unsupported scheme error")
			}
		})
	}
}

func TestNewDB_SQLiteScheme(t *testing.T) {
	tests := []struct {
		name string
		url  string
	}{
		{"sqlite://", "sqlite://:memory:"},
		{"file:", "file::memory:"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, err := NewDB(tt.url)
			if err != nil {
				t.Fatalf("expected sqlite to connect with URL %q, got: %v", tt.url, err)
			}
			if db == nil {
				t.Fatal("expected non-nil DB")
			}
			_ = db.Close()
		})
	}
}

func TestNewDB_UnsupportedScheme(t *testing.T) {
	_, err := NewDB("mysql://user:pass@localhost:3306/db")
	if err == nil {
		t.Fatal("expected error for unsupported scheme")
	}
	if !strings.Contains(err.Error(), "unsupported database URL scheme") {
		t.Errorf("expected 'unsupported database URL scheme' error, got: %v", err)
	}
}
