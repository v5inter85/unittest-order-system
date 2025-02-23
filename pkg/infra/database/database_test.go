package database_test

import (
	"context"
	"testing"
	"time"

	"order-system/pkg/infra/config"
	"order-system/pkg/infra/database"
)

func TestStats(t *testing.T) {
	cfg := &config.Config{}
	cfg.Database.MaxOpenConns = 10
	cfg.Database.MaxIdleConns = 5
	cfg.Database.MaxLifetime = time.Minute
	cfg.Database.Host = "localhost"
	cfg.Database.Port = 3306
	cfg.Database.User = "test"
	cfg.Database.Password = "test"
	cfg.Database.Database = "test"

	db, err := database.New(cfg)
	if err != nil {
		t.Skip("Skipping test - could not connect to database")
	}
	defer db.Close()

	stats := db.Stats()

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

func TestQueryRow(t *testing.T) {
	cfg := &config.Config{}
	cfg.Database.MaxOpenConns = 10
	cfg.Database.MaxIdleConns = 5
	cfg.Database.MaxLifetime = time.Minute
	cfg.Database.Host = "localhost"
	cfg.Database.Port = 3306
	cfg.Database.User = "test"
	cfg.Database.Password = "test"
	cfg.Database.Database = "test"

	db, err := database.New(cfg)
	if err != nil {
		t.Skip("Skipping test - could not connect to database")
	}
	defer db.Close()

	ctx := context.Background()
	row := db.QueryRow(ctx, "SELECT 1")
	if row == nil {
		t.Error("QueryRow should not return nil")
	}
}
