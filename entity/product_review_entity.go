package entity

import (
	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type ProductReview struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`

	ProductID *uuid.UUID `gorm:"type:uuid;index" json:"product_id"`
	Product   Product    `gorm:"foreignKey:ProductID;references:ID;constraint:OnDelete:CASCADE;"`

	ReviewText string         `gorm:"type:text" json:"review_text"`
	Tags       datatypes.JSON `gorm:"type:jsonb" json:"tags"` // array string mis. ["Masalah_Logistik"]

	TimeStamp
}
