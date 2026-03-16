package controllers

import (
	"net/http"
	"strconv"

	"aura-erp/backend/models"
	"aura-erp/backend/services"

	"github.com/gin-gonic/gin"
)

// SearchSections godoc
// @Summary Search sections
// @Description Search sections by name
// @Tags sections
// @Accept  json
// @Produce  json
// @Param q query string false "Search query"
// @Param limit query int false "Limit results" default(25)
// @Success 200 {array} models.SectionSearchResult
// @Failure 500 {object} map[string]string
// @Router /sections/search [get]
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

// GetAllSections godoc
// @Summary List sections
// @Description Get all sections in the system
// @Tags sections
// @Accept  json
// @Produce  json
// @Success 200 {array} models.Section
// @Failure 500 {object} map[string]string
// @Router /sections [get]
func GetAllSections(c *gin.Context) {
	sections, err := services.GetAllSections()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	c.JSON(http.StatusOK, sections)
}

// GetSectionByID godoc
// @Summary Get section by ID
// @Description Get a single section by its ID
// @Tags sections
// @Accept  json
// @Produce  json
// @Param id path int true "Section ID"
// @Success 200 {object} models.Section
// @Failure 404 {object} map[string]string
// @Router /sections/{id} [get]
func GetSectionByID(c *gin.Context) {
	id := c.Param("id")

	section, err := services.GetSectionByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Section not found"})
		return
	}

	c.JSON(http.StatusOK, section)
}

// CreateSection godoc
// @Summary Create a section
// @Description Create a new section
// @Tags sections
// @Accept  json
// @Produce  json
// @Param section body models.SectionCreate true "Section details"
// @Success 201 {object} models.Section
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /sections [post]
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

// UpdateSection godoc
// @Summary Update a section
// @Description Update an existing section's details
// @Tags sections
// @Accept  json
// @Produce  json
// @Param id path int true "Section ID"
// @Param section body models.SectionUpdate true "Updated section details"
// @Success 200 {object} models.Section
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /sections/{id} [put]
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

// DeleteSection godoc
// @Summary Delete a section
// @Description Remove a section from the system
// @Tags sections
// @Accept  json
// @Produce  json
// @Param id path int true "Section ID"
// @Success 200 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /sections/{id} [delete]
func DeleteSection(c *gin.Context) {
	id := c.Param("id")

	err := services.DeleteSection(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Section not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Section deleted successfully"})
}
