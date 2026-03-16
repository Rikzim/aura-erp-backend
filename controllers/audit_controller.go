package controllers

import (
	"net/http"
	"strconv"

	"aura-erp/backend/models"
	"aura-erp/backend/services"

	"github.com/gin-gonic/gin"
)

func GetAllAuditLogs(c *gin.Context) {
	limitStr := c.Query("limit")
	offsetStr := c.Query("offset")

	limit := 50
	offset := 0

	if l, err := strconv.Atoi(limitStr); err == nil {
		limit = l
	}
	if o, err := strconv.Atoi(offsetStr); err == nil {
		offset = o
	}

	logs, err := services.GetAllAuditLogs(limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	c.JSON(http.StatusOK, logs)
}

func GetAuditLogsByEntity(c *gin.Context) {
	entityType := c.Param("entityType")
	entityID := c.Param("entityId")

	logs, err := services.GetAuditLogsByEntity(entityType, entityID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	c.JSON(http.StatusOK, logs)
}

func CreateAuditLog(c *gin.Context) {
	var data models.AuditLogCreate
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if data.EntityType == "" || data.EntityID == 0 || data.Action == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Fields entity_type, entity_id, and action are required"})
		return
	}

	log, err := services.CreateAuditLog(data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	c.JSON(http.StatusCreated, log)
}
