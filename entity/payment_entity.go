package entity

import (
	"time"

	"github.com/google/uuid"
)

type Payment struct {
	ID     uuid.UUID `gorm:"type:uuid;primaryKey"`
	Method string
	Status bool
	PaidAt *time.Time

	OrderID *uuid.UUID `gorm:"type:uuid"`
	Order   *Order     `gorm:"foreignKey:OrderID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`

	TimeStamp
}
