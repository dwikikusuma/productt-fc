package repository

import (
	"context"
	"errors"
	"fmt"
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

func (r *ProductRepository) SearchProducts(ctx context.Context, searchParam *models.SearchProductParameter) ([]models.Product, int64, error) {
	var products []models.Product
	var totalCount int64

	query := r.Database.WithContext(ctx).Table("product").
		Select("product.id, product.name, product.description, product.price, product.stock, product.category_id, product_category.name as category").
		Joins("JOIN product_category ON product.category_id = product_category.id")

	if searchParam.Name != "" {
		query = query.Where("product.name LIKE ?", "%"+searchParam.Name+"%")
	}

	if searchParam.Category != "" {
		query = query.Where("product_category.name LIKE ?", "%"+searchParam.Category+"%")
	}

	if searchParam.MinPrice > 0 {
		query = query.Where("product.price >= ?", searchParam.MinPrice)
	}

	if searchParam.MaxPrice > 0 {
		query = query.Where("product.price <= ?", searchParam.MaxPrice)
	}

	// total count
	query.Model(&models.Product{}).Count(&totalCount)

	//default order by
	if searchParam.SortBy == "" {
		searchParam.SortBy = "product.name"

	}

	if searchParam.OrderBy == "" || (searchParam.OrderBy != "asc" && searchParam.OrderBy != "desc") {
		searchParam.OrderBy = "asc"
	}

	query = query.Order(fmt.Sprintf("%s %s", searchParam.SortBy, searchParam.OrderBy))

	//pagination

	offset := (searchParam.Page - 1) * searchParam.Limit
	query = query.Offset(int(offset)).Limit(int(searchParam.Limit))

	err := query.Scan(&products).Error
	if err != nil {
		return []models.Product{}, 0, err
	}

	return products, totalCount, nil
}
