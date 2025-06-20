package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"product_commerce/infra/log"
	"product_commerce/models"
	"strconv"
)

func (h *ProductHandler) ProductManagement(c *gin.Context) {
	var param models.ProductManagementParameter
	if err := c.ShouldBindJSON(&param); err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid input",
		})
		return
	}

	if param.Action == "" {
		log.Logger.Error("missing parameter")
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "missing parameter",
		})
		return
	}

	switch param.Action {
	case "add":
		if param.Product.ID != 0 {
			log.Logger.Error("invalid request - product category id is not empty")
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "invalid request",
			})
			return
		}
		productId, err := h.ProductUseCase.CreateProduct(c.Request.Context(), &param.Product)
		if err != nil {
			log.Logger.WithFields(logrus.Fields{
				"param": param,
			}).Errorf("h.ProductUseCase.CreateProduct got an error: %v", err)

			c.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
			return
		}
		c.JSON(http.StatusCreated, gin.H{
			"message": fmt.Sprintf("Successfully create new product %d", productId), // Use fmt.Sprintf to format the message
		})
		return

	case "edit":
		if param.Product.ID == 0 {
			log.Logger.Error("invalid request - product category id is not empty")
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "invalid request",
			})
			return
		}
		product, err := h.ProductUseCase.UpdateProduct(c.Request.Context(), &param.Product)
		if err != nil {
			log.Logger.WithFields(logrus.Fields{
				"param": param,
			}).Errorf("h.ProductUseCase.UpdateProduct got an error: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
			return
		}
		c.JSON(http.StatusCreated, gin.H{
			"message":         "Successfully update product",
			"productCategory": product,
		})
		return
	case "delete":
		if param.Product.ID <= 0 {
			log.Logger.Error("invalid request - product category is not set")
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "invalid request",
			})
			return
		}

		err := h.ProductUseCase.DeleteProduct(c.Request.Context(), param.Product.ID)
		if err != nil {
			log.Logger.WithFields(logrus.Fields{
				"param": param,
			}).Errorf("h.ProductUseCase.DeleteProduct got an error: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
			return
		}
		c.JSON(http.StatusCreated, gin.H{
			"message": "Successfully delete product",
		})
		return
	default:
		log.Logger.Error("invalid action: %s", param.Action)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid input",
		})
		return
	}

}

func (h *ProductHandler) GetProductById(c *gin.Context) {
	productIdStr := c.Param("id")
	productID, err := strconv.ParseInt(productIdStr, 10, 64)
	if err != nil {
		log.Logger.WithFields(logrus.Fields{
			"product_id": productIdStr,
		}).Errorf("trconv.ParseInt(productIdStr, 10, 64): %v", err)

		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	product, err := h.ProductUseCase.GetProductById(c.Request.Context(), productID)

	if product.ID == 0 {
		log.Logger.WithFields(logrus.Fields{
			"product_id": productID,
		}).Info("product not found")

		c.JSON(http.StatusNotFound, gin.H{
			"product": "not found",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"product": product,
	})
}
