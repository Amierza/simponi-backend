package entity

import "github.com/google/uuid"

type Product struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey"`
	Name        string
	Description string
	Stock       int
	Status      bool
	Price       int

	StoreID *uuid.UUID `gorm:"type:uuid"`
	Store   Store      `gorm:"foreignKey:StoreID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`

	CategoryID *uuid.UUID      `gorm:"type:uuid"`
	Category   ProductCategory `gorm:"foreignKey:CategoryID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`

	Images           []*ProductImage    `gorm:"foreignKey:ProductID"`
	ExternalProducts []*ExternalProduct `gorm:"foreignKey:ProductID"`

	TimeStamp
}
