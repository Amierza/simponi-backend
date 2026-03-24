package entity

import "github.com/google/uuid"

type Platform struct {
	ID   uuid.UUID `gorm:"type:uuid;primaryKey"`
	Name string    `gorm:"uniqueIndex"` // shopee, tiktok

	Stores []*Store `gorm:"foreignKey:PlatformID"`

	TimeStamp
}
