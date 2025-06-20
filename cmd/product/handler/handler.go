package handler

import (
	"product_commerce/cmd/product/usecase"
)

type ProductHandler struct {
	ProductUseCase usecase.ProductUseCase
}

func NewProductHandler(productUseCase usecase.ProductUseCase) *ProductHandler {
	return &ProductHandler{
		ProductUseCase: productUseCase,
	}
}
