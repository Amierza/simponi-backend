package migrations

import (
	"github.com/Amierza/simponi-backend/entity"
	"gorm.io/gorm"
)

func Rollback(db *gorm.DB) error {
	tables := []interface{}{
		// LOG
		&entity.Log{},
		&entity.InventoryLog{},

		// EXTERNAL DOMAIN
		&entity.ExternalProduct{},

		// PAYMENT DOMAIN
		&entity.Payment{},

		// ORDER DOMAIN
		&entity.OrderDetail{},
		&entity.Order{},

		// PRODUCT DOMAIN
		&entity.ProductImage{},
		&entity.Product{},

		// STORE DOMAIN
		&entity.StoreCredential{},
		&entity.Store{},

		// MASTER
		&entity.ProductCategory{},
		&entity.Platform{},
		&entity.User{},
	}

	for _, table := range tables {
		if err := db.Migrator().DropTable(table); err != nil {
			return err
		}
	}

	return nil
}
