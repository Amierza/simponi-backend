package entity

import (
	"github.com/google/uuid"
)

type Store struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`

	OwnerID *uuid.UUID `gorm:"type:uuid" json:"owner_id"`
	Owner   User       `gorm:"foreignKey:OwnerID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`

	Name        string `json:"name"`
	Description string `json:"description"`
	ImageURL    string `json:"image_url"`
	IsActive    bool   `json:"is_active"`

	Users          []*StoreUser     `gorm:"foreignKey:StoreID"`
	Products       []*Product       `gorm:"foreignKey:StoreID"`
	StorePlatforms []*StorePlatform `gorm:"foreignKey:StoreID"`
	Orders         []*Order         `gorm:"foreignKey:StoreID"`
	Logs           []*Log           `gorm:"foreignKey:StoreID"`
	InventoryLogs  []*InventoryLog  `gorm:"foreignKey:StoreID"`

	TimeStamp
}
