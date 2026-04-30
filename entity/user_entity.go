package entity

import (
	"github.com/Amierza/simponi-backend/helper"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID       uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Email    string    `gorm:"uniqueIndex" json:"email"`
	Password string    `json:"password"`
	Name     string    `json:"name"`
	ImageURL string    `json:"image_url"`
	Status   string    `json:"status"`

	RoleID *uuid.UUID `gorm:"type:uuid" json:"role_id"`
	Role   Role       `gorm:"foreignKey:RoleID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"role,omitempty"`

	Stores []*StoreUser `gorm:"foreignKey:UserID"`

	TimeStamp
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	defer func() {
		if err := recover(); err != nil {
			tx.Rollback()
		}
	}()

	var err error
	u.Password, err = helper.HashPassword(u.Password)
	if err != nil {
		return err
	}

	return nil
}
