package entity

import (
	"time"

	"github.com/google/uuid"
)

type Store struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`

	Name        string `json:"name"`
	Description string `json:"description"`
	ImageURL    string `json:"image_url"`

	// Marketplace identity
	ExternalShopID string `json:"external_shop_id"` // shop_id dari Shopee
	ExternalName   string `json:"external_name"`    // nama toko di marketplace

	// Ownership
	UserID *uuid.UUID `gorm:"type:uuid" json:"user_id"`
	User   *User      `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"user,omitempty"`

	PlatformID *uuid.UUID `gorm:"type:uuid" json:"platform_id"`
	Platform   *Platform  `gorm:"foreignKey:PlatformID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"platform,omitempty"`

	// Credential
	Credentials []*StoreCredential `gorm:"foreignKey:StoreID" json:"credentials,omitempty"`
	Logs        []*Log             `gorm:"foreignKey:StoreID" json:"logs,omitempty"`

	// Status
	IsActive    bool       `json:"is_active"`
	IsConnected bool       `json:"is_connected"`
	LastSyncAt  *time.Time `json:"last_sync_at"`

	TimeStamp
}
