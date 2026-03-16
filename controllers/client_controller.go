package controllers

import (
	"net/http"
	"strconv"

	"aura-erp/backend/models"
	"aura-erp/backend/services"

	"github.com/gin-gonic/gin"
)

// SearchClients godoc
// @Summary Search clients
// @Description Search clients by name
// @Tags clients
// @Accept  json
// @Produce  json
// @Param q query string false "Search query"
// @Param limit query int false "Limit results" default(25)
// @Success 200 {array} models.ClientSearchResult
// @Failure 500 {object} map[string]string
// @Router /clients/search [get]
func SearchClients(c *gin.Context) {
	q := c.Query("q")
	limitStr := c.DefaultQuery("limit", "25")
	limit := 25
	if n, err := strconv.Atoi(limitStr); err == nil {
		limit = n
	}

	results, err := services.SearchClients(q, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	c.JSON(http.StatusOK, results)
}

// GetAllClients godoc
// @Summary List clients
// @Description Get all clients in the system
// @Tags clients
// @Accept  json
// @Produce  json
// @Success 200 {array} models.Client
// @Failure 500 {object} map[string]string
// @Router /clients [get]
func GetAllClients(c *gin.Context) {
	clients, err := services.GetAllClients()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	c.JSON(http.StatusOK, clients)
}

// GetClientByID godoc
// @Summary Get client by ID
// @Description Get a single client by its ID
// @Tags clients
// @Accept  json
// @Produce  json
// @Param id path int true "Client ID"
// @Success 200 {object} models.Client
// @Failure 404 {object} map[string]string
// @Router /clients/{id} [get]
func GetClientByID(c *gin.Context) {
	id := c.Param("id")

	client, err := services.GetClientByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Client not found"})
		return
	}

	c.JSON(http.StatusOK, client)
}

// CreateClient godoc
// @Summary Create a client
// @Description Create a new client
// @Tags clients
// @Accept  json
// @Produce  json
// @Param client body models.ClientCreate true "Client details"
// @Success 201 {object} models.Client
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /clients [post]
func CreateClient(c *gin.Context) {
	var data models.ClientCreate
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if data.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Field \"name\" is required"})
		return
	}

	client, err := services.CreateClient(data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	c.JSON(http.StatusCreated, client)
}

// UpdateClient godoc
// @Summary Update a client
// @Description Update an existing client's details
// @Tags clients
// @Accept  json
// @Produce  json
// @Param id path int true "Client ID"
// @Param client body models.ClientUpdate true "Updated client details"
// @Success 200 {object} models.Client
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /clients/{id} [put]
func UpdateClient(c *gin.Context) {
	id := c.Param("id")

	var data models.ClientUpdate
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if data.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Field \"name\" is required"})
		return
	}

	client, err := services.UpdateClient(id, data)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Client not found"})
		return
	}

	c.JSON(http.StatusOK, client)
}

// DeleteClient godoc
// @Summary Delete a client
// @Description Remove a client from the system
// @Tags clients
// @Accept  json
// @Produce  json
// @Param id path int true "Client ID"
// @Success 200 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /clients/{id} [delete]
func DeleteClient(c *gin.Context) {
	id := c.Param("id")

	err := services.DeleteClient(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Client not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Client deleted successfully"})
}
