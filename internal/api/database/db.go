package database

import (
	"fmt"
	"log"
	models2 "uptime-checker/internal/api/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %w", err)
	}

	err = db.AutoMigrate(
		&models2.User{},
		&models2.Site{},
		&models2.Credential{},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	if err := seedRoles(db); err != nil {
		return nil, fmt.Errorf("failed to seedRoles: %w", err)
	}

	log.Println("Successfully connected to database and migrated schemas")
	return db, nil
}

func seedRoles(db *gorm.DB) error {
	roles := []string{models2.RoleUser, models2.RoleAdmin}

	for _, roleName := range roles {
		var role models2.Role
		err := db.Where(models2.Role{Name: roleName}).FirstOrCreate(&role).Error
		if err != nil {
			return err
		}
	}
	return nil
}
