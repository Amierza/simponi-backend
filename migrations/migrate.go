package migrations

import (
	"github.com/Amierza/simponi-backend/entity"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	if err := db.AutoMigrate(
		&entity.Role{},
		&entity.User{},
		&entity.Permission{},
		&entity.RolePermission{},

		&entity.Platform{},
		&entity.Store{},
		&entity.StoreUser{},
		&entity.StoreCredential{},
		&entity.Log{},
		&entity.StorePlatform{},

		&entity.ProductCategory{},
		&entity.Product{},
		&entity.ProductImage{},
		&entity.Vendor{},
		&entity.ProductVendor{},
		&entity.ExternalProduct{},
		&entity.InventoryLog{},

		&entity.Order{},
		&entity.OrderDetail{},
	); err != nil {
		return err
	}

	return nil
}
