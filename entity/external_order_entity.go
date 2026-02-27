package entity

import "github.com/google/uuid"

type ExternalOrder struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey"`

	OrderID *uuid.UUID `gorm:"type:uuid"`
	Order   *Order     `gorm:"foreignKey:OrderID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`

	PlatformID *uuid.UUID `gorm:"type:uuid"`
	Platform   *Platform  `gorm:"foreignKey:PlatformID;references:ID"`

	ExternalOrderID string

	TimeStamp
}
