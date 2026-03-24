package entity

import (
	"github.com/google/uuid"
)

type Log struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey"`

	StoreID *uuid.UUID `gorm:"type:uuid"`
	Store   *Store     `gorm:"foreignKey:StoreID;references:ID;constraint:OnDelete:CASCADE;"`

	Action  string
	Message string

	TimeStamp
}
