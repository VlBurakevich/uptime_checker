package database

import (
	"fmt"
	"log/slog"

	"github.com/pressly/goose/v3"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"uptime-checker/migrations/api_db"
)

func InitDB(dsn string) (*gorm.DB, error) {
	if err := runMigrations(dsn); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	slog.Info("Successfully connected to database and migrated schemas")
	return db, nil
}

func runMigrations(dsn string) error {
	sqlDB, err := goose.OpenDBWithDriver("postgres", dsn)
	if err != nil {
		return fmt.Errorf("goose: failed to open database: %w", err)
	}
	defer func() {
		if closeErr := sqlDB.Close(); closeErr != nil {
			slog.Error("Failed to close sqlDB", "error", closeErr)
		}
	}()

	goose.SetBaseFS(migrations.FS)
	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}

	slog.Info("Checking and applying migrations...")

	if err := goose.Up(sqlDB, "."); err != nil {
		return fmt.Errorf("goose up failed: %w", err)
	}

	return nil
}
