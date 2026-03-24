package entity

import "github.com/google/uuid"

type ExternalProduct struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey"`

	ProductID *uuid.UUID `gorm:"type:uuid"`
	Product   *Product   `gorm:"foreignKey:ProductID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`

	StoreID *uuid.UUID `gorm:"type:uuid"`
	Store   Store      `gorm:"foreignKey:StoreID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`

	ExternalProductID string
	ExternalModelID   string
	Price             int64

	TimeStamp
}
