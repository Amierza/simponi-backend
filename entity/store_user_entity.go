package entity

import "github.com/google/uuid"

type StoreUser struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey"`

	UserID *uuid.UUID `gorm:"type:uuid" json:"user_id"`
	User   User       `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"user,omitempty"`

	StoreID *uuid.UUID `gorm:"type:uuid" json:"store_id"`
	Store   Store      `gorm:"foreignKey:StoreID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"store,omitempty"`

	TimeStamp
}
