package database_test

import (
	"testing"
	"time"

	"order-system/pkg/infra/config"
	"order-system/pkg/infra/database"
)

func TestStats(t *testing.T) {
	// Create minimal config for testing
	cfg := &config.Config{}
	cfg.Database.MaxOpenConns = 10
	cfg.Database.MaxIdleConns = 5
	cfg.Database.MaxLifetime = time.Minute
	cfg.Database.Host = "localhost"
	cfg.Database.Port = 3306
	cfg.Database.User = "test"
	cfg.Database.Password = "test"
	cfg.Database.Database = "test"

	// Create database instance
	db, err := database.New(cfg)
	if err != nil {
		t.Skip("Skipping test - could not connect to database")
	}
	defer db.Close()

	// Test Stats
	stats := db.Stats()

	// Basic validation that Stats returns expected type
	if stats.OpenConnections < 0 {
		t.Error("OpenConnections should not be negative")
	}
	if stats.InUse < 0 {
		t.Error("InUse should not be negative")
	}
	if stats.Idle < 0 {
		t.Error("Idle should not be negative")
	}
	if stats.WaitCount < 0 {
		t.Error("WaitCount should not be negative")
	}
}
