package repository

import (
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
