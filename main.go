package main

import (
	"github.com/gin-gonic/gin"
	"product_commerce/cmd/product/handler"
	"product_commerce/cmd/product/repository"
	"product_commerce/cmd/product/resource"
	"product_commerce/cmd/product/service"
	"product_commerce/cmd/product/usecase"
	"product_commerce/config"
	"product_commerce/infra/log"
	"product_commerce/routes"
)

func main() {
	cfg := config.LoadConfig()
	redis := resource.InitRedis(&cfg)
	db := resource.InitDB(&cfg)

	// logger
	log.SetupLogger()

	// prepare each layer
	productRepository := repository.NewProductRepo(db, redis)
	productService := service.NewProductService(*productRepository)
	productUseCase := usecase.NewProductUseCase(*productService)
	productHandler := handler.NewProductHandler(*productUseCase)

	port := cfg.App.Port
	router := gin.Default()

	// routes
	routes.SetupRoutes(router, *productHandler)

	_ = router.Run(":" + port)
	log.Logger.Info("Server Running on Port: %s", port)
}
