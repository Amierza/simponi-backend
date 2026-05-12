package entity

import (
	"time"

	"github.com/google/uuid"
)

type StorePlatform struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`

	StoreID *uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_store_platform" json:"store_id"`
	Store   *Store     `gorm:"foreignKey:StoreID;references:ID;constraint:OnDelete:CASCADE;"`

	PlatformID *uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_store_platform" json:"platform_id"`
	Platform   *Platform  `gorm:"foreignKey:PlatformID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"platform,omitempty"`

	// marketpalce info
	ExternalShopID string `json:"external_shop_id"` // contoh: shop_id dari marketplace
	ExternalName   string `json:"external_name"`    // nama toko di marketplace

	IsConnected bool       `json:"is_connected"`
	LastSyncAt  *time.Time `json:"last_sync_at"`

	Credential *StoreCredential `gorm:"foreignKey:StorePlatformID"`

	TimeStamp
}
