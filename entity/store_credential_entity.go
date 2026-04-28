package entity

import (
	"time"

	"github.com/google/uuid"
)

type StoreCredential struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey"`

	StorePlatformID uuid.UUID      `gorm:"type:uuid;not null;uniqueIndex"`
	StorePlatform   *StorePlatform `gorm:"foreignKey:StorePlatformID;constraint:OnDelete:CASCADE;"`

	AccessToken  string     `gorm:"type:text"`
	RefreshToken string     `gorm:"type:text"`
	ExpiresAt    *time.Time `gorm:"column:expires_at"`

	TimeStamp
}
