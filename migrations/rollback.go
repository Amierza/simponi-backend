package migrations

import (
	"github.com/Amierza/simponi-backend/entity"
	"gorm.io/gorm"
)

func Rollback(db *gorm.DB) error {
	tables := []interface{}{
		&entity.OrderDetail{},
		&entity.Order{},

		&entity.InventoryLog{},
		&entity.ExternalProduct{},
		&entity.ProductVendor{},
		&entity.Vendor{},
		&entity.ProductImage{},
		&entity.Product{},
		&entity.ProductCategory{},

		&entity.StorePlatform{},
		&entity.Log{},
		&entity.StoreCredential{},
		&entity.StoreUser{},
		&entity.Store{},
		&entity.Platform{},

		&entity.RolePermission{},
		&entity.Permission{},
		&entity.User{},
		&entity.Role{},
	}

	for _, table := range tables {
		if err := db.Migrator().DropTable(table); err != nil {
			return err
		}
	}

	return nil
}
