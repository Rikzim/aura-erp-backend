package controllers

import (
	"net/http"
	"strconv"

	"aura-erp/backend/models"
	"aura-erp/backend/services"

	"github.com/gin-gonic/gin"
)

// GetAllAuditLogs godoc
// @Summary List audit logs
// @Description Get a paginated list of system audit logs
// @Tags audit
// @Accept  json
// @Produce  json
// @Param limit query int false "Limit results" default(50)
// @Param offset query int false "Offset results" default(0)
// @Success 200 {array} models.AuditLog
// @Failure 500 {object} map[string]string
// @Router /audit-log [get]
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

// GetAuditLogsByEntity godoc
// @Summary Get audit logs for an entity
// @Description Get all audit logs associated with a specific entity type and ID
// @Tags audit
// @Accept  json
// @Produce  json
// @Param entityType path string true "Entity Type"
// @Param entityId path int true "Entity ID"
// @Success 200 {array} models.AuditLog
// @Failure 500 {object} map[string]string
// @Router /audit-log/{entityType}/{entityId} [get]
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

// CreateAuditLog godoc
// @Summary Create an audit log entry
// @Description Manually create an audit log record
// @Tags audit
// @Accept  json
// @Produce  json
// @Param auditLog body models.AuditLogCreate true "Audit log details"
// @Success 201 {object} models.AuditLog
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /audit-log [post]
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
