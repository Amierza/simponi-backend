package entity

import "github.com/google/uuid"

type Vendor struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`

	Name     string `json:"name"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Address  string `json:"address"`
	ImageURL string `json:"image_url"`

	ProductVendors []*ProductVendor `gorm:"foreignKey:VendorID" json:"product_vendors,omitempty"`

	TimeStamp
}
