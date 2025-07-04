package service

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"product_commerce/cmd/product/repository"
	"product_commerce/infra/log"
	"product_commerce/models"
)

type ProductService struct {
	ProductRepository repository.ProductRepository
}

func NewProductService(productRepo repository.ProductRepository) *ProductService {
	return &ProductService{
		ProductRepository: productRepo,
	}
}

func (s *ProductService) GetProductById(ctx context.Context, productId int64) (*models.Product, error) {
	product, err := s.ProductRepository.GetProductByIdFromRedis(ctx, productId)
	if err != nil {
		log.Logger.WithFields(logrus.Fields{
			"error": fmt.Sprintf("got error on s.ProductRepository.GetProductByIdFromRedis: %d", err),
			"id":    productId,
		})
	}

	if product.ID != 0 {
		return product, nil
	}

	product, err = s.ProductRepository.FindByProductId(ctx, productId)
	if err != nil {
		return nil, err
	}

	ctxConcurrent := context.WithValue(context.Background(), "request_id", ctx.Value("request_id"))
	go func(ctx context.Context, product *models.Product) {
		errConcurrent := s.ProductRepository.SetProductById(ctx, product)
		if errConcurrent != nil {
			log.Logger.WithFields(logrus.Fields{
				"product": product,
			}).Errorf("s.ProductRepository.SetProductByID() got error %v", errConcurrent)
		}
	}(ctxConcurrent, product)
	return product, nil
}

func (s *ProductService) GetProductCatById(ctx context.Context, productCatId int64) (*models.ProductCategory, error) {
	productCat, err := s.ProductRepository.FindProductCatById(ctx, productCatId)
	if err != nil {
		return nil, err
	}
	return productCat, nil
}

func (s *ProductService) CreateProduct(ctx context.Context, product *models.Product) (int, error) {
	id, err := s.ProductRepository.InsertNewProduct(ctx, product)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (s *ProductService) CreateProductCat(ctx context.Context, productCat *models.ProductCategory) (int, error) {
	id, err := s.ProductRepository.InsertNewProductCat(ctx, productCat)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (s *ProductService) UpdateProduct(ctx context.Context, product *models.Product) (*models.Product, error) {
	var model *models.Product
	err := s.ProductRepository.WithTransaction(ctx, func(tx *gorm.DB) error {
		productDetail, err := s.ProductRepository.UpdateProduct(ctx, product)
		if err != nil {
			return err
		}
		model = productDetail

		err = s.ProductRepository.SetProductById(ctx, product)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return model, nil
}

func (s *ProductService) UpdateProductCat(ctx context.Context, productCat *models.ProductCategory) (*models.ProductCategory, error) {
	var model *models.ProductCategory
	err := s.ProductRepository.WithTransaction(ctx, func(tx *gorm.DB) error {
		productCatDetail, err := s.ProductRepository.UpdateProductCat(ctx, productCat)
		if err != nil {
			return err
		}
		model = productCatDetail

		err = s.ProductRepository.DeleteProductCache(ctx, productCatDetail.ID)
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return model, nil
}

func (s *ProductService) DeleteProduct(ctx context.Context, productId int) error {
	err := s.ProductRepository.WithTransaction(ctx, func(tx *gorm.DB) error {
		err := s.ProductRepository.DeleteProduct(ctx, productId)
		if err != nil {
			return err
		}

		err = s.ProductRepository.DeleteProductCache(ctx, productId)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}
	return nil
}

func (s *ProductService) DeleteProductCat(ctx context.Context, productCatId int) error {
	err := s.ProductRepository.DeleteProductCat(ctx, productCatId)
	if err != nil {
		return err
	}
	return nil
}

func (s *ProductService) GetProductList(ctx context.Context, paramRequest *models.SearchProductParameter) ([]models.Product, int64, error) {
	products, total, err := s.ProductRepository.SearchProducts(ctx, paramRequest)
	if err != nil {
		return []models.Product{}, 0, err
	}
	return products, total, nil
}
