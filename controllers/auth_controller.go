package controllers

import (
	"net/http"
	"strings"

	"aura-erp/backend/services"

	"github.com/gin-gonic/gin"
)

type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Login godoc
// @Summary Login user
// @Description Authenticate user and return a JWT token
// @Tags auth
// @Accept  json
// @Produce  json
// @Param login body LoginRequest true "Login credentials"
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]string
// @Router /auth/login [post]
func Login(c *gin.Context) {
	var req LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email and password are required"})
		return
	}

	result, err := services.Login(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	c.JSON(http.StatusOK, result)
}

// Verify godoc
// @Summary Verify token
// @Description Verify the JWT token and return user info
// @Tags auth
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Bearer {token}"
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]string
// @Router /auth/verify [get]
func Verify(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No token provided"})
		return
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")

	user, err := services.VerifyTokenAndGetUser(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}

// Logout godoc
// @Summary Logout user
// @Description Logout user (client-side usually handles this by deleting the token)
// @Tags auth
// @Accept  json
// @Produce  json
// @Success 200 {object} map[string]string
// @Router /auth/logout [post]
func Logout(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}
