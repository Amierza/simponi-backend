package entity

import (
	"time"

	"github.com/google/uuid"
)

type StoreCredential struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey"`

	StoreID *uuid.UUID `gorm:"type:uuid"`
	Store   *Store     `gorm:"foreignKey:StoreID;references:ID;constraint:OnDelete:CASCADE;"`

	AccessToken  string
	RefreshToken string
	TokenExpiry  *time.Time

	TimeStamp
}
