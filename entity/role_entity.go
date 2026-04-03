package entity

import "github.com/google/uuid"

type Role struct {
	ID   uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Name string    `json:"name"`

	Users           []*User           `gorm:"foreignKey:RoleID"`
	RolePermissions []*RolePermission `gorm:"foreignKey:RoleID"`

	TimeStamp
}
