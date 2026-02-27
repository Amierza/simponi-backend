package entity

import "github.com/google/uuid"

type ProductImage struct {
	ID       uuid.UUID `gorm:"type:uuid;primaryKey"`
	ImageURL string

	ProductID *uuid.UUID `gorm:"type:uuid"`
	Product   Product    `gorm:"foreignKey:ProductID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`

	TimeStamp
}
