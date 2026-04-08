package entity

import "github.com/google/uuid"

type Platform struct {
	ID   uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Name string    `gorm:"uniqueIndex" json:"name"` // shopee, tiktok

	StorePlatforms []*StorePlatform `gorm:"foreignKey:PlatformID" json:"store_platforms,omitempty"`

	TimeStamp
}
