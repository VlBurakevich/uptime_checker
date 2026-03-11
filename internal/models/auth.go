package models

import (
	"github.com/google/uuid"
)

type Role struct {
	ID   uuid.UUID `gorm:"type:uuid;primaryKey"`
	Name string    `gorm:"unique;not null"`
}

type User struct {
	Base
	Username string `gorm:"unique;not null" json:"username"`
	Email    string `gorm:"unique;not null" json:"email"`

	Credential Credential `gorm:"foreignKey:UserID;" json:"-"`
	Roles      []Role     `gorm:"many2many:user_roles;" json:"roles"`
	Sites      []Site     `json:"sites"`
}

type Credential struct {
	UserID       uuid.UUID `gorm:"type:uuid;primaryKey"`
	PasswordHash string    `gorm:"not null"`
}
