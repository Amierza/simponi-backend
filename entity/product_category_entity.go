package entity

import "github.com/google/uuid"

type ProductCategory struct {
	ID   uuid.UUID `gorm:"type:uuid;primaryKey"`
	Name string

	Products []*Product `gorm:"foreignKey:CategoryID"`

	TimeStamp
}
