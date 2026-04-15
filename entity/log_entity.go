package entity

import (
	"github.com/google/uuid"
)

type Log struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`

	StoreID *uuid.UUID `gorm:"type:uuid" json:"store_id"`
	Store   *Store     `gorm:"foreignKey:StoreID;references:ID;constraint:OnDelete:CASCADE;"`

	Action  string `json:"action"`
	Message string `json:"message"`

	TimeStamp
}
