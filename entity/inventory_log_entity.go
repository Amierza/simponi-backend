package entity

import (
	"github.com/google/uuid"
)

type InventoryLog struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey"`

	ProductID *uuid.UUID `gorm:"type:uuid"`
	Product   *Product   `gorm:"foreignKey:ProductID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`

	Change int // -1, -2, +10

	Source string // shopee, tiktok, manual
	Note   string

	TimeStamp
}
