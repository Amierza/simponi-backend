package entity

import "github.com/google/uuid"

type Product struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey"`
	Name        string
	Description string
	SKU         string `gorm:"uniqueIndex"`

	// CENTRAL STOCK
	Stock int

	CategoryID *uuid.UUID      `gorm:"type:uuid"`
	Category   ProductCategory `gorm:"foreignKey:CategoryID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`

	Images           []*ProductImage    `gorm:"foreignKey:ProductID"`
	ExternalProducts []*ExternalProduct `gorm:"foreignKey:ProductID"`
	Logs             []*InventoryLog    `gorm:"foreignKey:ProductID"`

	TimeStamp
}
