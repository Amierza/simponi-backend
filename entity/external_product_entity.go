package entity

import "github.com/google/uuid"

type ExternalProduct struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`

	ProductID *uuid.UUID `gorm:"type:uuid" json:"product_id"`
	Product   *Product   `gorm:"foreignKey:ProductID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"product,omitempty"`

	StoreID *uuid.UUID `gorm:"type:uuid" json:"store_id"`
	Store   Store      `gorm:"foreignKey:StoreID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"store,omitempty"`

	ExternalProductID string `json:"external_product_id"`
	ExternalModelID   string `json:"external_model_id"`
	Price             int64  `json:"price"`

	TimeStamp
}
