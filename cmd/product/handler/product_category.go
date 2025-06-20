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

func (h *ProductHandler) GetProductCategoryById(c *gin.Context) {
	productCatIdStr := c.Param("id")
	productCatID, err := strconv.ParseInt(productCatIdStr, 10, 64)
	if err != nil {
		log.Logger.WithFields(logrus.Fields{
			"product_id": productCatID,
		}).Errorf("trconv.ParseInt(productIdStr, 10, 64): %v", err)

		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	productCat, err := h.ProductUseCase.GetProductCatById(c.Request.Context(), productCatID)
	if err != nil {
		log.Logger.WithFields(logrus.Fields{
			"product_cat_id": productCatID,
		}).Errorf("h.ProductUseCase.GetProductCatById got an error: %v", err)

		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	if productCat.ID == 0 {
		log.Logger.WithFields(logrus.Fields{
			"product_cat_id": productCatID,
		}).Info("product cat not found")

		c.JSON(http.StatusNotFound, gin.H{
			"product_category": "product category not found",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"product_category": productCat,
	})
}

func (h *ProductHandler) ProductCategoryManagement(c *gin.Context) {
	var param models.ProductCategoryManagementParameter
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
		if param.ProductCategory.ID != 0 {
			log.Logger.Error("invalid request - product category id is not empty")
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "invalid request",
			})
			return
		}
		productCatId, err := h.ProductUseCase.CreateNewProductCategory(c.Request.Context(), &param.ProductCategory)
		if err != nil {
			log.Logger.WithFields(logrus.Fields{
				"param": param,
			}).Errorf("h.ProductUseCase.CreateNewProductCategory got an error: %v", err)
			return
		}
		c.JSON(http.StatusCreated, gin.H{
			"message": fmt.Sprintf("Successfully create new product category %d", productCatId), // Use fmt.Sprintf to format the message
		})
		return

	case "edit":
		if param.ProductCategory.ID == 0 {
			log.Logger.Error("invalid request - product category id is not empty")
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "invalid request",
			})
			return
		}
		productCategory, err := h.ProductUseCase.UpdateProductCat(c.Request.Context(), &param.ProductCategory)
		if err != nil {
			log.Logger.WithFields(logrus.Fields{
				"param": param,
			}).Errorf("h.ProductUseCase.UpdateProductCat got an error: %v", err)
			return
		}
		c.JSON(http.StatusCreated, gin.H{
			"message":         "Successfully update product category",
			"productCategory": productCategory,
		})
		return

	case "delete":
		if param.ProductCategory.ID <= 0 {
			log.Logger.Error("invalid request - product category is not set")
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "invalid request",
			})
			return
		}

		err := h.ProductUseCase.DeleteProductCat(c.Request.Context(), param.ProductCategory.ID)
		if err != nil {
			log.Logger.WithFields(logrus.Fields{
				"param": param,
			}).Errorf("h.ProductUseCase.DeleteProductCat got an error: %v", err)
			return
		}
		c.JSON(http.StatusCreated, gin.H{
			"message": "Successfully delete product category",
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
