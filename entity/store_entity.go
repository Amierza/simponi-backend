package entity

import (
	"github.com/google/uuid"
)

type Store struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`

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
