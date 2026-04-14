package entity

import "github.com/google/uuid"

type Permission struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey"`

	Name     string `json:"name"`     // contoh: "CreateProduct"
	Endpoint string `json:"endpoint"` // contoh: "/products"
	Method   string `json:"method"`   // GET, POST, PUT, DELETE
	Module   string `json:"module"`   // contoh: "user, vendor"

	RolePermissions []*RolePermission `gorm:"foreignKey:PermissionID"`

	TimeStamp
}
