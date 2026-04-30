package repository

import (
	"context"

	"gorm.io/gorm"
)

type ITransaction interface {
	Run(ctx context.Context, fn func(tx *gorm.DB) error) error
}

type transaction struct {
	db *gorm.DB
}

func NewTransaction(db *gorm.DB) ITransaction {
	return &transaction{db: db}
}

func (t *transaction) Run(ctx context.Context, fn func(tx *gorm.DB) error) error {
	return t.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(tx)
	})
}
