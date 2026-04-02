package entity

import "github.com/google/uuid"

type OrderDetail struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey"`

	OrderID *uuid.UUID `gorm:"type:uuid"`
	Order   *Order     `gorm:"foreignKey:OrderID;references:ID;constraint:OnDelete:CASCADE;"`

	ExternalProductID *uuid.UUID       `gorm:"type:uuid" json:"external_product_id"`
	ExternalProduct   *ExternalProduct `gorm:"foreignKey:ExternalProductID;references:ID;constraint:OnDelete:SET NULL;"`

	Quantity int

	TimeStamp
}
