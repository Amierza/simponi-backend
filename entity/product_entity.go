package entity

import "github.com/google/uuid"

type Product struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`

	Name        string `json:"name"`
	Description string `json:"description"`
	SKU         string `gorm:"uniqueIndex" json:"sku"`
	Stock       int    `json:"stock"` // central stok

	CategoryID *uuid.UUID      `gorm:"type:uuid" json:"category_id"`
	Category   ProductCategory `gorm:"foreignKey:CategoryID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`

	Images           []*ProductImage    `gorm:"foreignKey:ProductID"`
	ProductVendors   []*ProductVendor   `gorm:"foreignKey:ProductID"`
	ExternalProducts []*ExternalProduct `gorm:"foreignKey:ProductID"`
	Logs             []*InventoryLog    `gorm:"foreignKey:ProductID"`

	TimeStamp
}
