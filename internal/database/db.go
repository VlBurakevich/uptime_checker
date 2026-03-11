package database

import (
	"fmt"
	"log"
	"uptime-checker/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %w", err)
	}

	err = db.AutoMigrate(
		&models.User{},
		&models.Site{},
		&models.Credential{},
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
	roles := []string{models.RoleUser, models.RoleAdmin}

	for _, roleName := range roles {
		var role models.Role
		err := db.Where(models.Role{Name: roleName}).FirstOrCreate(&role).Error
		if err != nil {
			return err
		}
	}
	return nil
}
