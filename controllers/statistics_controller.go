package controllers

import (
	"net/http"

	"aura-erp/backend/services"

	"github.com/gin-gonic/gin"
)

func GetStatistics(c *gin.Context) {
	stats, err := services.GetStatistics()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to compute statistics"})
		return
	}

	c.JSON(http.StatusOK, stats)
}
