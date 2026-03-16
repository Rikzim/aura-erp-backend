package controllers

import (
	"net/http"

	"aura-erp/backend/models"
	"aura-erp/backend/services"

	"github.com/gin-gonic/gin"
)

func GetAllOrders(c *gin.Context) {
	orders, err := services.GetAllOrders()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	c.JSON(http.StatusOK, orders)
}

func GetOrderByID(c *gin.Context) {
	id := c.Param("id")

	order, err := services.GetOrderByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	c.JSON(http.StatusOK, order)
}

func CreateOrder(c *gin.Context) {
	var data models.OrderCreate
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if data.Reference == "" || data.ClientID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Fields reference and client_id are required"})
		return
	}

	order, err := services.CreateOrder(data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	c.JSON(http.StatusCreated, order)
}

func UpdateOrder(c *gin.Context) {
	id := c.Param("id")

	var data models.OrderUpdate
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if data.Reference == "" || data.ClientID == 0 || data.Status == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Fields reference, client_id, and status are required"})
		return
	}

	order, err := services.UpdateOrder(id, data)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	c.JSON(http.StatusOK, order)
}

func DeleteOrder(c *gin.Context) {
	id := c.Param("id")

	err := services.DeleteOrder(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order deleted successfully"})
}

// ── Order Items ───────────────────────────────

func GetOrderItems(c *gin.Context) {
	orderID := c.Query("order_id")

	// No order_id → return all items across every order (used by calendar bulk-load)
	if orderID == "" {
		items, err := services.GetAllOrderItems()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}
		c.JSON(http.StatusOK, items)
		return
	}

	items, err := services.GetOrderItems(orderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	c.JSON(http.StatusOK, items)
}

func CreateOrderItem(c *gin.Context) {
	orderID := c.Query("order_id")
	if orderID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "order_id query param is required"})
		return
	}

	var data models.OrderItemCreate
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if data.ProductID == 0 || data.Quantity == 0 || data.UnitPrice == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Fields product_id, quantity, and unit_price are required"})
		return
	}

	item, err := services.CreateOrderItem(orderID, data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	c.JSON(http.StatusCreated, item)
}

func UpdateOrderItem(c *gin.Context) {
	id := c.Param("id")

	var data models.OrderItemUpdate
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if data.ProductID == 0 || data.Quantity == 0 || data.UnitPrice == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Fields product_id, quantity, and unit_price are required"})
		return
	}

	item, err := services.UpdateOrderItem(id, data)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order item not found"})
		return
	}

	c.JSON(http.StatusOK, item)
}

func DeleteOrderItem(c *gin.Context) {
	id := c.Param("id")

	err := services.DeleteOrderItem(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order item not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order item deleted successfully"})
}
