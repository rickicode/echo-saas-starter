package database

import "testing"

func TestRewritePlaceholders(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "no placeholders",
			input:    "SELECT * FROM users",
			expected: "SELECT * FROM users",
		},
		{
			name:     "single placeholder",
			input:    "SELECT * FROM users WHERE id = ?",
			expected: "SELECT * FROM users WHERE id = $1",
		},
		{
			name:     "multiple placeholders",
			input:    "INSERT INTO users (name, email) VALUES (?, ?)",
			expected: "INSERT INTO users (name, email) VALUES ($1, $2)",
		},
		{
			name:     "placeholder in WHERE clause",
			input:    "SELECT * FROM users WHERE email = ? AND is_active = ?",
			expected: "SELECT * FROM users WHERE email = $1 AND is_active = $2",
		},
		{
			name:     "question mark in single-quoted string",
			input:    "SELECT * FROM users WHERE name = '?' AND id = ?",
			expected: "SELECT * FROM users WHERE name = '?' AND id = $1",
		},
		{
			name:     "many placeholders",
			input:    "INSERT INTO t (a,b,c,d,e,f,g,h,i,j,k) VALUES (?,?,?,?,?,?,?,?,?,?,?)",
			expected: "INSERT INTO t (a,b,c,d,e,f,g,h,i,j,k) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)",
		},
		{
			name:     "empty query",
			input:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := RewritePlaceholders(tt.input)
			if result != tt.expected {
				t.Errorf("RewritePlaceholders(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}
