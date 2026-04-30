package entity

import (
	"time"

	"github.com/google/uuid"
)

type StorePlatform struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`

	StoreID *uuid.UUID `gorm:"type:uuid" json:"store_id"`
	Store   *Store     `gorm:"foreignKey:StoreID;references:ID;constraint:OnDelete:CASCADE;"`

	PlatformID *uuid.UUID `gorm:"type:uuid" json:"platform_id"`
	Platform   *Platform  `gorm:"foreignKey:PlatformID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"platform,omitempty"`

	ExternalShopID string `json:"external_shop_id"` // shop_id dari Shopee
	ExternalName   string `json:"external_name"`    // nama toko di marketplace

	IsConnected bool       `json:"is_connected"`
	LastSyncAt  *time.Time `json:"last_sync_at"`

	Credential *StoreCredential `gorm:"foreignKey:StorePlatformID"`

	TimeStamp

	// 1 store 1 platform
	_ struct{} `gorm:"uniqueIndex:idx_store_platform,unique"`
}
