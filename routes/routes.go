package routes

import (
	"github.com/gin-gonic/gin"
	"product_commerce/cmd/product/handler"
	"product_commerce/middleware"
)

func SetupRoutes(router *gin.Engine, productHandler handler.ProductHandler) {
	router.Use(middleware.RequestLogger())

	router.POST("v1/product_category", productHandler.ProductCategoryManagement)
	router.GET("v1/product_category/:id", productHandler.GetProductCategoryById)

	router.POST("v1/product", productHandler.ProductManagement)
	router.GET("v1/product/:id", productHandler.GetProductById)

	router.GET("v1/products/search", productHandler.SearchProduct)

}
