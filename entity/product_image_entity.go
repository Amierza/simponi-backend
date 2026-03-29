package entity

import "github.com/google/uuid"

type ProductImage struct {
	ID       uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	ImageURL string    `json:"image_url"`

	ProductID *uuid.UUID `gorm:"type:uuid" json:"product_id"`
	Product   Product    `gorm:"foreignKey:ProductID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"product,omitempty"`

	TimeStamp
}
