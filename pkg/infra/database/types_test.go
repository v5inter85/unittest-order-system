package database

import (
	"database/sql"
	"errors"
	"testing"
)

func TestIsNoRows(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "sql.ErrNoRows error",
			err:      sql.ErrNoRows,
			expected: true,
		},
		{
			name:     "nil error",
			err:      nil,
			expected: false,
		},
		{
			name:     "other error",
			err:      errors.New("some other error"),
			expected: false,
		},
		{
			name: "database error with sql.ErrNoRows",
			err: &Error{
				Operation: "query",
				Err:       sql.ErrNoRows,
			},
			expected: true,
		},
		{
			name: "database error with other error",
			err: &Error{
				Operation: "query",
				Err:       errors.New("some other error"),
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsNoRows(tt.err)
			if result != tt.expected {
				t.Errorf("IsNoRows() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestIsDuplicate(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "nil error",
			err:      nil,
			expected: false,
		},
		{
			name:     "non-database error",
			err:      errors.New("some other error"),
			expected: false,
		},
		{
			name: "duplicate key error",
			err: &Error{
				Operation: "insert",
				Err:      errors.New("1062: Duplicate entry"),
			},
			expected: true,
		},
		{
			name: "other database error",
			err: &Error{
				Operation: "insert",
				Err:      errors.New("1064: Syntax error"),
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsDuplicate(tt.err)
			if result != tt.expected {
				t.Errorf("IsDuplicate() = %v, want %v", result, tt.expected)
			}
		})
	}
}
