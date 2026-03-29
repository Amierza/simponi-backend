package entity

import "github.com/google/uuid"

type Product struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	SKU         string    `gorm:"uniqueIndex" json:"sku"`

	// CENTRAL STOCK
	Stock int `json:"stock"`

	CategoryID *uuid.UUID      `gorm:"type:uuid" json:"category_id"`
	Category   ProductCategory `gorm:"foreignKey:CategoryID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"category,omitempty"`

	Images           []*ProductImage    `gorm:"foreignKey:ProductID" json:"images,omitempty"`
	ExternalProducts []*ExternalProduct `gorm:"foreignKey:ProductID" json:"external_products,omitempty"`
	Logs             []*InventoryLog    `gorm:"foreignKey:ProductID" json:"logs,omitempty"`

	TimeStamp
}
