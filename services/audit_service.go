package services

import (
	"aura-erp/backend/config"
	"aura-erp/backend/models"
	"database/sql"
	"strconv"
)

func GetAllAuditLogs(limit int, offset int) ([]models.AuditLog, error) {
	query := `
		SELECT a.id, a.user_id, a.entity_type, a.entity_id, a.action, a.old_value, a.new_value, a.created_at,
		       u.name
		FROM audit_log a
		LEFT JOIN users u ON a.user_id = u.id
		ORDER BY a.created_at DESC
		LIMIT $1 OFFSET $2
	`
	rows, err := config.DB.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []models.AuditLog
	for rows.Next() {
		var log models.AuditLog
		if err := rows.Scan(&log.ID, &log.UserID, &log.EntityType, &log.EntityID, &log.Action, &log.OldValue, &log.NewValue, &log.CreatedAt, &log.UserName); err != nil {
			return nil, err
		}
		logs = append(logs, log)
	}

	return logs, nil
}

func GetAuditLogsByEntity(entityType string, entityID string) ([]models.AuditLog, error) {
	query := `
		SELECT a.id, a.user_id, a.entity_type, a.entity_id, a.action, a.old_value, a.new_value, a.created_at,
		       u.name
		FROM audit_log a
		LEFT JOIN users u ON a.user_id = u.id
		WHERE a.entity_type = $1 AND a.entity_id = $2
		ORDER BY a.created_at DESC
	`
	rows, err := config.DB.Query(query, entityType, entityID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []models.AuditLog
	for rows.Next() {
		var log models.AuditLog
		if err := rows.Scan(&log.ID, &log.UserID, &log.EntityType, &log.EntityID, &log.Action, &log.OldValue, &log.NewValue, &log.CreatedAt, &log.UserName); err != nil {
			return nil, err
		}
		logs = append(logs, log)
	}

	return logs, nil
}

func CreateAuditLog(data models.AuditLogCreate) (*models.AuditLog, error) {
	query := `INSERT INTO audit_log (user_id, entity_type, entity_id, action, old_value, new_value)
	          VALUES ($1, $2, $3, $4, $5, $6)
	          RETURNING id, user_id, entity_type, entity_id, action, old_value, new_value, created_at`

	var log models.AuditLog
	var userID sql.NullInt64
	if data.UserID.Valid {
		userID = data.UserID
	}

	err := config.DB.QueryRow(query, userID, data.EntityType, data.EntityID, data.Action, data.OldValue, data.NewValue).Scan(
		&log.ID, &log.UserID, &log.EntityType, &log.EntityID, &log.Action, &log.OldValue, &log.NewValue, &log.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &log, nil
}

// Helper function to convert string to int
func StringToInt(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}
