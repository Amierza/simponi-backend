package entity

import "github.com/google/uuid"

type Product struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`

	StoreID uuid.UUID `gorm:"type:uuid;index" json:"store_id"`
	Store   Store     `gorm:"foreignKey:StoreID;references:ID;constraint:OnDelete:CASCADE;"`

	CategoryID *uuid.UUID      `gorm:"type:uuid" json:"category_id"`
	Category   ProductCategory `gorm:"foreignKey:CategoryID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`

	Name        string `json:"name"`
	Description string `json:"description"`
	SKU         string `gorm:"uniqueIndex" json:"sku"`
	Stock       int    `json:"stock"` // central stok

	Images           []*ProductImage    `gorm:"foreignKey:ProductID"`
	ProductVendors   []*ProductVendor   `gorm:"foreignKey:ProductID"`
	ExternalProducts []*ExternalProduct `gorm:"foreignKey:ProductID"`
	Logs             []*InventoryLog    `gorm:"foreignKey:ProductID"`

	TimeStamp
}
