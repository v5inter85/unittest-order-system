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
			name:     "database.Error with sql.ErrNoRows",
			err:      &Error{Err: sql.ErrNoRows},
			expected: true,
		},
		{
			name:     "other error",
			err:      errors.New("some error"),
			expected: false,
		},
		{
			name:     "nil error",
			err:      nil,
			expected: false,
		},
		{
			name:     "database.Error with other error",
			err:      &Error{Err: errors.New("other error")},
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
			name:     "duplicate error code 1062",
			err:      &Error{Err: errors.New("1062: Duplicate entry")},
			expected: true,
		},
		{
			name:     "non-duplicate error",
			err:      &Error{Err: errors.New("1064: Some other error")},
			expected: false,
		},
		{
			name:     "non-database error",
			err:      errors.New("some error"),
			expected: false,
		},
		{
			name:     "nil error",
			err:      nil,
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
