package entity

import (
	"time"

	"github.com/google/uuid"
)

type Store struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey"`

	Name        string
	Description string
	ImageURL    string

	// Marketplace identity
	ExternalShopID string // shop_id dari Shopee
	ExternalName   string // nama toko di marketplace

	// Ownership
	UserID *uuid.UUID `gorm:"type:uuid"`
	User   *User      `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`

	PlatformID *uuid.UUID `gorm:"type:uuid"`
	Platform   *Platform  `gorm:"foreignKey:PlatformID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`

	// Credential
	Credentials []*StoreCredential `gorm:"foreignKey:StoreID"`
	Logs        []*Log             `gorm:"foreignKey:StoreID"`

	// Status
	IsActive    bool
	IsConnected bool
	LastSyncAt  *time.Time

	TimeStamp
}
