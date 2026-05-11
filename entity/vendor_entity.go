package entity

import "github.com/google/uuid"

type Vendor struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`

	StoreID uuid.UUID `gorm:"type:uuid;index" json:"store_id"`
	Store   Store     `gorm:"foreignKey:StoreID;references:ID;constraint:OnDelete:CASCADE;"`

	Name        string `json:"name"`
	Email       string `json:"email"`
	PhoneNumber string `gorm:"not null" json:"phone_number"`
	Address     string `json:"address"`
	ImageURL    string `json:"image_url"`
	Description string `json:"description,omitempty"`

	ProductVendors []*ProductVendor `gorm:"foreignKey:VendorID" json:"product_vendors,omitempty"`

	TimeStamp
}
