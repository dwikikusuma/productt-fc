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

func (h *ProductHandler) SearchProduct(c *gin.Context) {
	name := c.Query("name")
	category := c.Query("category")

	minPrice, _ := strconv.ParseFloat(c.Query("min_price"), 64)
	maxPrice, _ := strconv.ParseFloat(c.Query("max_price"), 64)

	sortBy := c.Query("sort_by")
	orderBy := c.Query("order_by")

	page, _ := strconv.ParseInt(c.DefaultQuery("page", "1"), 10, 64)
	PageSize, _ := strconv.ParseInt(c.DefaultQuery("page_size", "10"), 10, 64)

	searchParam := models.SearchProductParameter{
		Name:     name,
		Category: category,
		MinPrice: minPrice,
		MaxPrice: maxPrice,
		SortBy:   sortBy,
		OrderBy:  orderBy,
		Limit:    PageSize,
		Page:     page,
	}

	products, total, err := h.ProductUseCase.SearchProduct(c.Request.Context(), &searchParam)
	if err != nil {
		log.Logger.WithFields(logrus.Fields{
			"param": searchParam,
		}).Errorf("h.ProductUseCase.SearchProduct got an error: %v", err)

		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	totalPages := (total + PageSize - 1) / PageSize

	var nextPageURL *string
	if page < totalPages {
		nextPage := page + 1
		nextPageURLStr := fmt.Sprintf("/products?name=%s&category=%s&min_price=%f&max_price=%f&sort_by=%s&order_by=%s&page=%d&page_size=%d",
			name, category, minPrice, maxPrice, sortBy, orderBy, nextPage, PageSize)
		nextPageURL = &nextPageURLStr
	}

	c.JSON(http.StatusOK, models.SearchProductResponse{
		Products:    products,
		Page:        int(page),
		PageSize:    int(PageSize),
		TotalCount:  total,
		TotalPages:  int(totalPages),
		NextPage:    nextPageURL != nil,
		NextPageURL: nextPageURL,
	})

}
