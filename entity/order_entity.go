package entity

import (
	"time"

	"github.com/google/uuid"
)

type Order struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey"`

	// External Identity
	ExternalOrderID string // order_sn
	OrderNumber     string // display ID

	StoreID *uuid.UUID `gorm:"type:uuid"`
	Store   *Store     `gorm:"foreignKey:StoreID;references:ID"`

	// Buyer Snapshot
	BuyerName  string
	BuyerPhone string
	BuyerEmail string

	// Shipping Snapshot
	RecipientName    string
	RecipientPhone   string
	ShippingAddress  string
	ShippingCity     string
	ShippingProvince string
	ShippingPostal   string
	ShippingMethod   string
	TrackingNumber   string

	// Monetary Breakdown
	SubtotalAmount int64
	ShippingFee    int64
	MarketplaceFee int64
	DiscountAmount int64
	TaxAmount      int64
	TotalAmount    int64
	NetAmount      int64

	// Status
	OrderStatus       string // PENDING, READY_TO_SHIP, COMPLETED
	PaymentStatus     string // UNPAID, PAID
	FulfillmentStatus string

	// Timestamps from marketplace
	OrderedAt   *time.Time
	PaidAt      *time.Time
	ShippedAt   *time.Time
	CompletedAt *time.Time
	CancelledAt *time.Time

	OrderDetails []*OrderDetail `gorm:"foreignKey:OrderID"`

	TimeStamp
}
