package entity

import "github.com/google/uuid"

type ExternalProduct struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`

	StorePlatformID *uuid.UUID    `gorm:"type:uuid" json:"store_platform_id"`
	StorePlatform   StorePlatform `gorm:"foreignKey:StorePlatformID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`

	ProductID *uuid.UUID `gorm:"type:uuid" json:"product_id"`
	Product   *Product   `gorm:"foreignKey:ProductID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`

	Price int64 `json:"price"`

	OrderDetails []*OrderDetail `gorm:"foreignKey:ExternalProductID"`

	TimeStamp
}
