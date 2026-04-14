package migrations

import (
	"github.com/Amierza/simponi-backend/entity"
	"gorm.io/gorm"
)

func Seed(db *gorm.DB) error {
	if err := SeedFromJSON[entity.Role](db, "./migrations/json/roles.json", entity.Role{}, "ID"); err != nil {
		return err
	}

	if err := SeedFromJSON[entity.Permission](db, "./migrations/json/permissions.json", entity.Permission{}, "ID"); err != nil {
		return err
	}

	if err := SeedFromJSON[entity.RolePermission](db, "./migrations/json/role_permissions.json", entity.RolePermission{}, "ID"); err != nil {
		return err
	}

	if err := SeedFromJSON[entity.Platform](db, "./migrations/json/platforms.json", entity.Platform{}, "ID"); err != nil {
		return err
	}

	if err := SeedFromJSON[entity.User](db, "./migrations/json/users.json", entity.User{}, "ID"); err != nil {
		return err
	}

	if err := SeedFromJSON[entity.Store](db, "./migrations/json/stores.json", entity.Store{}, "ID"); err != nil {
		return err
	}

	if err := SeedFromJSON[entity.StorePlatform](db, "./migrations/json/store_platforms.json", entity.StorePlatform{}, "ID"); err != nil {
		return err
	}

	if err := SeedFromJSON[entity.ProductCategory](db, "./migrations/json/product_categories.json", entity.ProductCategory{}, "ID"); err != nil {
		return err
	}

	if err := SeedFromJSON[entity.Product](db, "./migrations/json/products.json", entity.Product{}, "ID"); err != nil {
		return err
	}

	if err := SeedFromJSON[entity.ProductImage](db, "./migrations/json/product_images.json", entity.ProductImage{}, "ID"); err != nil {
		return err
	}

	if err := SeedFromJSON[entity.ExternalProduct](db, "./migrations/json/external_products.json", entity.ExternalProduct{}, "ID"); err != nil {
		return err
	}

	if err := SeedFromJSON[entity.InventoryLog](db, "./migrations/json/inventory_logs.json", entity.InventoryLog{}, "ID"); err != nil {
		return err
	}

	return nil
}
