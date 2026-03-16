package controllers

import (
	"net/http"
	"strconv"

	"aura-erp/backend/models"
	"aura-erp/backend/services"

	"github.com/gin-gonic/gin"
)

// SearchProducts godoc
// @Summary Search products
// @Description Search products by name
// @Tags products
// @Accept  json
// @Produce  json
// @Param q query string false "Search query"
// @Param limit query int false "Limit results" default(25)
// @Success 200 {array} models.ProductSearchResult
// @Failure 500 {object} map[string]string
// @Router /products/search [get]
func SearchProducts(c *gin.Context) {
	q := c.Query("q")
	limitStr := c.DefaultQuery("limit", "25")
	limit := 25
	if n, err := strconv.Atoi(limitStr); err == nil {
		limit = n
	}

	results, err := services.SearchProducts(q, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	c.JSON(http.StatusOK, results)
}

// GetAllProducts godoc
// @Summary List products
// @Description Get all products in the system
// @Tags products
// @Accept  json
// @Produce  json
// @Success 200 {array} models.Product
// @Failure 500 {object} map[string]string
// @Router /products [get]
func GetAllProducts(c *gin.Context) {
	products, err := services.GetAllProducts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	c.JSON(http.StatusOK, products)
}

func GetProductByID(c *gin.Context) {
	id := c.Param("id")

	product, err := services.GetProductByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	c.JSON(http.StatusOK, product)
}

// CreateProduct godoc
// @Summary Create a product
// @Description Create a new product
// @Tags products
// @Accept  json
// @Produce  json
// @Param product body models.ProductCreate true "Product details"
// @Success 201 {object} models.Product
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /products [post]
func CreateProduct(c *gin.Context) {
	var data models.ProductCreate
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if data.Name == "" || data.UnitPrice == 0 || data.Unit == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Fields name, unit_price, and unit are required"})
		return
	}

	product, err := services.CreateProduct(data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	c.JSON(http.StatusCreated, product)
}

func UpdateProduct(c *gin.Context) {
	id := c.Param("id")

	var data models.ProductUpdate
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if data.Name == "" || data.UnitPrice == 0 || data.Unit == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Fields name, unit_price, and unit are required"})
		return
	}

	product, err := services.UpdateProduct(id, data)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	c.JSON(http.StatusOK, product)
}

func DeleteProduct(c *gin.Context) {
	id := c.Param("id")

	err := services.DeleteProduct(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully"})
}
