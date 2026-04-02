package entity

import "github.com/google/uuid"

type Permission struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey"`

	Name     string // contoh: "create_product"
	Endpoint string // contoh: "/products"
	Method   string // GET, POST, DELETE

	RolePermissions []*RolePermission `gorm:"foreignKey:PermissionID"`

	TimeStamp
}
