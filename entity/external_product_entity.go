package entity

import "github.com/google/uuid"

type ExternalProduct struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey"`

	ProductID *uuid.UUID `gorm:"type:uuid"`
	Product   *Product   `gorm:"foreignKey:ProductID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`

	PlatformID *uuid.UUID `gorm:"type:uuid"`
	Platform   *Platform  `gorm:"foreignKey:PlatformID;references:ID"`

	ExternalProductID string
	ExternalModelID   string

	TimeStamp
}
