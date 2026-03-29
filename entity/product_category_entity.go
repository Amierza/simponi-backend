package entity

import "github.com/google/uuid"

type ProductCategory struct {
	ID   uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Name string    `json:"name"`

	Products []*Product `gorm:"foreignKey:CategoryID" json:"products,omitempty"`

	TimeStamp
}
