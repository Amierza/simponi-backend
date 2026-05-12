package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type ExternalProduct struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`

	StorePlatformID *uuid.UUID    `gorm:"type:uuid" json:"store_platform_id"`
	StorePlatform   StorePlatform `gorm:"foreignKey:StorePlatformID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`

	ProductID *uuid.UUID `gorm:"type:uuid" json:"product_id"`
	Product   *Product   `gorm:"foreignKey:ProductID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`

	ExternalID   string `gorm:"type:varchar(255);not null;uniqueIndex:idx_store_platform_external" json:"external_id"`
	ExternalName string `gorm:"type:text" json:"external_name"`
	Price        int64  `gorm:"default:0" json:"price"`

	LastSyncAt *time.Time `json:"last_sync_at"`

	RawPayload datatypes.JSON `gorm:"type:jsonb" json:"raw_payload"`

	OrderDetails []*OrderDetail `gorm:"foreignKey:ExternalProductID"`

	TimeStamp
}
