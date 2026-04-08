package entity

import "github.com/google/uuid"

type RolePermission struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey"`

	RoleID *uuid.UUID `gorm:"type:uuid" json:"role_id"`
	Role   Role       `gorm:"foreignKey:RoleID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"role,omitempty"`

	PermissionID *uuid.UUID `gorm:"type:uuid" json:"permission_id"`
	Permission   Permission `gorm:"foreignKey:PermissionID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"permission,omitempty"`

	TimeStamp
}
