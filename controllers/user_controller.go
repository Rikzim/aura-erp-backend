package controllers

import (
	"net/http"

	"aura-erp/backend/models"
	"aura-erp/backend/services"

	"github.com/gin-gonic/gin"
)

// GetAllUsers godoc
// @Summary List users
// @Description Get all registered users
// @Tags users
// @Accept  json
// @Produce  json
// @Success 200 {array} models.User
// @Failure 500 {object} map[string]string
// @Router /users [get]
func GetAllUsers(c *gin.Context) {
	users, err := services.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	c.JSON(http.StatusOK, users)
}

// GetUserByID godoc
// @Summary Get user by ID
// @Description Get details of a specific user
// @Tags users
// @Accept  json
// @Produce  json
// @Param id path int true "User ID"
// @Success 200 {object} models.User
// @Failure 404 {object} map[string]string
// @Router /users/{id} [get]
func GetUserByID(c *gin.Context) {
	id := c.Param("id")

	user, err := services.GetUserByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// CreateUser godoc
// @Summary Create a user
// @Description Register a new user in the system
// @Tags users
// @Accept  json
// @Produce  json
// @Param user body models.UserCreate true "User details"
// @Success 201 {object} models.User
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /users [post]
func CreateUser(c *gin.Context) {
	var data models.UserCreate
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := services.CreateUser(data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	c.JSON(http.StatusCreated, user)
}

// UpdateUser godoc
// @Summary Update a user
// @Description Update an existing user's details
// @Tags users
// @Accept  json
// @Produce  json
// @Param id path int true "User ID"
// @Param user body models.UserUpdate true "Updated user details"
// @Success 200 {object} models.User
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /users/{id} [put]
func UpdateUser(c *gin.Context) {
	id := c.Param("id")

	var data models.UserUpdate
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := services.UpdateUser(id, data)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// DeleteUser godoc
// @Summary Delete a user
// @Description Remove a user from the system
// @Tags users
// @Accept  json
// @Produce  json
// @Param id path int true "User ID"
// @Success 200 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /users/{id} [delete]
func DeleteUser(c *gin.Context) {
	id := c.Param("id")

	err := services.DeleteUser(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
