package controllers

import (
	"net/http"
	"strconv"

	"aura-erp/backend/models"
	"aura-erp/backend/services"

	"github.com/gin-gonic/gin"
)

func SearchSections(c *gin.Context) {
	q := c.Query("q")
	limitStr := c.DefaultQuery("limit", "25")
	limit := 25
	if n, err := strconv.Atoi(limitStr); err == nil {
		limit = n
	}

	results, err := services.SearchSections(q, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	c.JSON(http.StatusOK, results)
}

func GetAllSections(c *gin.Context) {
	sections, err := services.GetAllSections()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	c.JSON(http.StatusOK, sections)
}

func GetSectionByID(c *gin.Context) {
	id := c.Param("id")

	section, err := services.GetSectionByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Section not found"})
		return
	}

	c.JSON(http.StatusOK, section)
}

func CreateSection(c *gin.Context) {
	var data models.SectionCreate
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if data.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Field \"name\" is required"})
		return
	}

	section, err := services.CreateSection(data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	c.JSON(http.StatusCreated, section)
}

func UpdateSection(c *gin.Context) {
	id := c.Param("id")

	var data models.SectionUpdate
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if data.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Field \"name\" is required"})
		return
	}

	section, err := services.UpdateSection(id, data)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Section not found"})
		return
	}

	c.JSON(http.StatusOK, section)
}

func DeleteSection(c *gin.Context) {
	id := c.Param("id")

	err := services.DeleteSection(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Section not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Section deleted successfully"})
}
