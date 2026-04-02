package entity

import (
	"github.com/google/uuid"
)

type StorePlatform struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`

	PlatformID *uuid.UUID `gorm:"type:uuid" json:"platform_id"`
	Platform   *Platform  `gorm:"foreignKey:PlatformID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"platform,omitempty"`

	StoreID *uuid.UUID `gorm:"type:uuid"`
	Store   *Store     `gorm:"foreignKey:StoreID;references:ID;constraint:OnDelete:CASCADE;"`

	TimeStamp
}
