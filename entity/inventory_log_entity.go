package entity

import (
	"github.com/google/uuid"
)

type InventoryLog struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`

	ProductID *uuid.UUID `gorm:"type:uuid" json:"product_id"`
	Product   *Product   `gorm:"foreignKey:ProductID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"product,omitempty"`

	Change int `json:"change"` // -1, -2, +10

	Source string `json:"source"` // shopee, tiktok, manual
	Note   string `json:"note"`

	TimeStamp
}
