package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Identity struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
}

type Base struct {
	Identity
	CreatedAt time.Time      `gorm:"type:timestamptz" json:"created_at"`
	UpdatedAt time.Time      `gorm:"type:timestamptz" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index; type:timestamptz" json:"-"`
}

func (b *Identity) BeforeCreate(tx *gorm.DB) (err error) {
	if b.ID == uuid.Nil {
		b.ID = uuid.New()
	}
	return
}
