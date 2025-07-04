package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"product_commerce/infra/log"
	"product_commerce/models"
	"time"
)

var (
	cacheKeyProductInfo    = "product:%d"
	cacheKeyProductCatInfo = "product-info:%d"
)

func (r *ProductRepository) GetProductByIdFromRedis(ctx context.Context, productId int64) (*models.Product, error) {
	var product models.Product
	productKey := fmt.Sprintf(cacheKeyProductInfo, productId)
	fmt.Println(productKey)
	productStr, err := r.Redis.Get(ctx, productKey).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return &models.Product{}, nil
		}

		return nil, err
	}

	err = json.Unmarshal([]byte(productStr), &product)
	if err != nil {
		return nil, err
	}

	return &product, nil
}

func (r *ProductRepository) GetProductCatByIdFromRedis(ctx context.Context, productCatId int64) (*models.ProductCategory, error) {
	var productCat models.ProductCategory

	cacheKey := fmt.Sprintf(cacheKeyProductCatInfo, productCatId)
	productCatStr, err := r.Redis.Get(ctx, cacheKey).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return &models.ProductCategory{}, nil
		}

		return nil, err
	}

	err = json.Unmarshal([]byte(productCatStr), &productCat)
	if err != nil {
		return nil, err
	}

	return &productCat, nil
}

func (r *ProductRepository) SetProductById(ctx context.Context, product *models.Product) error {
	cacheKey := fmt.Sprintf(cacheKeyProductInfo, product.ID)
	productJson, err := json.Marshal(product)
	if err != nil {
		log.Logger.Error("failed to marshal product on ProductRepository SetProductById")
		return err
	}

	err = r.Redis.SetEX(ctx, cacheKey, productJson, 10*time.Minute).Err()
	if err != nil {
		log.Logger.Error("failed to set product cache on ProductRepository SetProductById")
		return err
	}
	return nil
}

func (r *ProductRepository) SetProductCatById(ctx context.Context, productCat *models.ProductCategory) error {
	cacheKey := fmt.Sprintf(cacheKeyProductInfo, productCat.ID)
	productCatJson, err := json.Marshal(productCat)
	if err != nil {
		return err
	}

	err = r.Redis.SetEX(ctx, cacheKey, productCatJson, 10*time.Minute).Err()
	if err != nil {
		return err
	}
	return nil
}
