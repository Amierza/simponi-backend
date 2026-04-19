package entity

import (
	"time"

	"github.com/google/uuid"
)

type StoreCredential struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey"`

	StoreID uuid.UUID `gorm:"type:uuid;not null;index"`
	Store   *Store    `gorm:"foreignKey:StoreID;references:ID;constraint:OnDelete:CASCADE;"`

	Provider   string  `gorm:"type:varchar(50);not null;index"` // shopee, tokopedia
	ExternalID *string `gorm:"type:varchar(100)"`               // shop_id dari provider

	AccessToken  string     `gorm:"type:text"`
	RefreshToken string     `gorm:"type:text"`
	ExpiresAt    *time.Time `gorm:"column:expires_at"`

	TimeStamp

	// constraint unique: 1 store 1 provider
	// gorm tag
	// ini penting banget biar gak duplicate connect
	// composite unique index
	_ struct{} `gorm:"uniqueIndex:idx_store_provider,unique"`
}
