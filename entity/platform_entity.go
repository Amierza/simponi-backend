package entity

import "github.com/google/uuid"

type Platform struct {
	ID   uuid.UUID `gorm:"type:uuid;primaryKey"`
	Name string

	Stores           []*Store           `gorm:"foreignKey:PlatformID"`
	ExternalProducts []*ExternalProduct `gorm:"foreignKey:PlatformID"`
	ExternalOrders   []*ExternalOrder   `gorm:"foreignKey:PlatformID"`

	TimeStamp
}
