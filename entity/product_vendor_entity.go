package entity

import "github.com/google/uuid"

type ProductVendor struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey"`

	ProductID *uuid.UUID `gorm:"type:uuid" json:"product_id"`
	Product   Product    `gorm:"foreignKey:ProductID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"product,omitempty"`
	VendorID  *uuid.UUID `gorm:"type:uuid" json:"vendor_id"`
	Vendor    Vendor     `gorm:"foreignKey:VendorID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"vendor,omitempty"`

	Price        int64
	MinimumOrder int
	LeadTimeDay  int  // delivery time (day)
	IsSelected   bool // choosen vendor

	TimeStamp
}
