package migrations

import (
	"github.com/Amierza/simponi-backend/entity"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	if err := db.AutoMigrate(
		// MASTER
		&entity.User{},
		&entity.Platform{},
		&entity.ProductCategory{},

		// STORE DOMAIN
		&entity.Store{},
		&entity.StoreCredential{},

		// PRODUCT DOMAIN
		&entity.Product{},
		&entity.ProductImage{},

		// ORDER DOMAIN
		&entity.Order{},
		&entity.OrderDetail{},

		// PAYMENT DOMAIN
		&entity.Payment{},

		// EXTERNAL DOMAIN
		&entity.ExternalProduct{},
		&entity.ExternalOrder{},
	); err != nil {
		return err
	}

	return nil
}
