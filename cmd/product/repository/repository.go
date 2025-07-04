package repository

import (
	"context"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type ProductRepository struct {
	Database *gorm.DB
	Redis    *redis.Client
}

func NewProductRepo(db *gorm.DB, redis *redis.Client) *ProductRepository {
	return &ProductRepository{
		Database: db,
		Redis:    redis,
	}
}

func (r *ProductRepository) WithTransaction(ctx context.Context, fn func(tx *gorm.DB) error) error {
	tx := r.Database.Begin().WithContext(ctx)

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	err := fn(tx)
	if r != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
