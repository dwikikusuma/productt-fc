package usecase

import (
	"context"
	"github.com/sirupsen/logrus"
	"product_commerce/cmd/product/service"
	"product_commerce/infra/log"
	"product_commerce/models"
)

type ProductUseCase struct {
	ProductService service.ProductService
}

func NewProductUseCase(orderService service.ProductService) *ProductUseCase {
	return &ProductUseCase{
		ProductService: orderService,
	}
}

func (uc *ProductUseCase) GetProductById(ctx context.Context, productID int64) (*models.Product, error) {
	product, err := uc.ProductService.GetProductById(ctx, productID)
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (uc *ProductUseCase) GetProductCatById(ctx context.Context, productCategoryID int64) (*models.ProductCategory, error) {
	productCategory, err := uc.ProductService.GetProductCatById(ctx, productCategoryID)
	if err != nil {
		return nil, err
	}
	return productCategory, nil
}

func (uc *ProductUseCase) CreateProduct(ctx context.Context, param *models.Product) (int, error) {
	productID, err := uc.ProductService.CreateProduct(ctx, param)
	if err != nil {
		log.Logger.WithFields(logrus.Fields{
			"name":     param.Name,
			"category": param.CategoryID,
		}).Errorf("uc.ProductService.CreateNewProduct got error %v", err)
		return 0, err
	}

	return productID, nil
}

func (uc *ProductUseCase) CreateNewProductCategory(ctx context.Context, param *models.ProductCategory) (int, error) {
	productCategoryID, err := uc.ProductService.CreateProductCat(ctx, param)
	if err != nil {
		log.Logger.WithFields(logrus.Fields{
			"name": param.Name,
		}).Errorf("uc.ProductService.CreateNewProductCategory got error %v", err)
		return 0, err
	}

	return productCategoryID, nil
}

func (uc *ProductUseCase) UpdateProduct(ctx context.Context, param *models.Product) (*models.Product, error) {
	product, err := uc.ProductService.UpdateProduct(ctx, param)
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (uc *ProductUseCase) UpdateProductCat(ctx context.Context, param *models.ProductCategory) (*models.ProductCategory, error) {
	productCategory, err := uc.ProductService.UpdateProductCat(ctx, param)
	if err != nil {
		return nil, err
	}

	return productCategory, nil
}

func (uc *ProductUseCase) DeleteProduct(ctx context.Context, productID int) error {
	err := uc.ProductService.DeleteProduct(ctx, productID)
	if err != nil {
		return err
	}

	return nil
}

func (uc *ProductUseCase) DeleteProductCat(ctx context.Context, productCategoryID int) error {
	err := uc.ProductService.DeleteProductCat(ctx, productCategoryID)
	if err != nil {
		return err
	}

	return nil
}
