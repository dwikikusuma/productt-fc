package repository

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"product_commerce/models"
)

func (r *ProductRepository) FindByProductId(ctx context.Context, productId int64) (*models.Product, error) {
	var product models.Product
	err := r.Database.WithContext(ctx).Table("product").Where("id = ?", productId).Last(&product).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &models.Product{}, nil
		}
		return nil, err
	}

	return &product, nil
}

func (r *ProductRepository) FindProductCatById(ctx context.Context, productCatId int64) (*models.ProductCategory, error) {
	var prodCat models.ProductCategory
	err := r.Database.WithContext(ctx).Table("product_category").Where("id = ?", productCatId).Last(&prodCat).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &models.ProductCategory{}, nil
		}
		return nil, err
	}
	return &prodCat, nil
}

func (r *ProductRepository) InsertNewProduct(ctx context.Context, product *models.Product) (int, error) {
	err := r.Database.WithContext(ctx).Table("product").Create(product).Error
	if err != nil {
		return 0, err
	}

	return product.ID, nil
}

func (r *ProductRepository) InsertNewProductCat(ctx context.Context, productCat *models.ProductCategory) (int, error) {
	err := r.Database.WithContext(ctx).Table("product_category").Create(productCat).Error
	if err != nil {
		return 0, err
	}
	return productCat.ID, nil
}

func (r *ProductRepository) UpdateProduct(ctx context.Context, product *models.Product) (*models.Product, error) {
	err := r.Database.WithContext(ctx).Table("product").Save(&product).Error
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (r *ProductRepository) UpdateProductCat(ctx context.Context, product *models.ProductCategory) (*models.ProductCategory, error) {
	err := r.Database.WithContext(ctx).Table("product_category").Save(&product).Error
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (r *ProductRepository) DeleteProduct(ctx context.Context, id int) error {
	err := r.Database.WithContext(ctx).Table("product").Delete(&models.Product{}, id).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *ProductRepository) DeleteProductCat(ctx context.Context, id int) error {
	err := r.Database.WithContext(ctx).Table("product_category").Delete(&models.ProductCategory{}, id).Error
	if err != nil {
		return err
	}
	return nil
}
