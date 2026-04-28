package entity

import (
	"github.com/google/uuid"
)

type Store struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`

	UserID *uuid.UUID `gorm:"type:uuid" json:"user_id"`
	User   *User      `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`

	Name        string `json:"name"`
	Description string `json:"description"`
	ImageURL    string `json:"image_url"`
	IsActive    bool   `json:"is_active"`

	StorePlatforms []*StorePlatform `gorm:"foreignKey:StoreID"`
	Orders         []*Order         `gorm:"foreignKey:StoreID"`
	Logs           []*Log           `gorm:"foreignKey:StoreID"`

	TimeStamp
}
