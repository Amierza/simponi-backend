package entity

import "github.com/google/uuid"

type OrderDetail struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey"`

	OrderID *uuid.UUID `gorm:"type:uuid"`
	Order   *Order     `gorm:"foreignKey:OrderID;references:ID;constraint:OnDelete:CASCADE;"`

	ProductID *uuid.UUID `gorm:"type:uuid"`
	Product   *Product   `gorm:"foreignKey:ProductID;references:ID;constraint:OnDelete:SET NULL;"`

	ExternalProductID string
	ProductName       string // snapshot
	SKU               string // snapshot

	Quantity int
	Price    int64
	Total    int64

	TimeStamp
}
