package migrations

import (
	"github.com/Amierza/simponi-backend/entity"
	"gorm.io/gorm"
)

func Seed(db *gorm.DB) error {
	if err := SeedFromJSON[entity.Platform](db, "./migrations/json/platforms.json", entity.Platform{}, "Name"); err != nil {
		return err
	}

	if err := SeedFromJSON[entity.User](db, "./migrations/json/users.json", entity.User{}, "Email"); err != nil {
		return err
	}

	if err := SeedFromJSON[entity.Store](db, "./migrations/json/stores.json", entity.Store{}, "ExternalShopID"); err != nil {
		return err
	}

	if err := SeedFromJSON[entity.ProductCategory](db, "./migrations/json/product_categories.json", entity.ProductCategory{}, "Name"); err != nil {
		return err
	}

	if err := SeedFromJSON[entity.Product](db, "./migrations/json/products.json", entity.Product{}, "SKU"); err != nil {
		return err
	}

	if err := SeedFromJSON[entity.ProductImage](db, "./migrations/json/product_images.json", entity.ProductImage{}, "ImageURL"); err != nil {
		return err
	}

	if err := SeedFromJSON[entity.ExternalProduct](db, "./migrations/json/external_products.json", entity.ExternalProduct{}, "ExternalProductID"); err != nil {
		return err
	}

	if err := SeedFromJSON[entity.InventoryLog](db, "./migrations/json/inventory_logs.json", entity.InventoryLog{}, ""); err != nil {
		return err
	}

	return nil
}