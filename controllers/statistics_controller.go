package controllers

import (
	"net/http"

	"aura-erp/backend/services"

	"github.com/gin-gonic/gin"
)

// GetStatistics godoc
// @Summary Get system statistics
// @Description Calculate and return key business metrics (counts, totals)
// @Tags statistics
// @Accept  json
// @Produce  json
// @Success 200 {object} models.StatisticsResponse
// @Failure 500 {object} map[string]string
// @Router /statistics [get]
func GetStatistics(c *gin.Context) {
	stats, err := services.GetStatistics()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to compute statistics"})
		return
	}

	c.JSON(http.StatusOK, stats)
}
