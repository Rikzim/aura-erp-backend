package models

import (
	"database/sql"
	"encoding/json"
	"time"
)

type AuditLog struct {
	ID         int              `json:"id"`
	UserID     NullInt64        `json:"user_id"`
	EntityType string           `json:"entity_type"`
	EntityID   int              `json:"entity_id"`
	Action     string           `json:"action"`
	OldValue   *json.RawMessage `json:"old_value"`
	NewValue   *json.RawMessage `json:"new_value"`
	CreatedAt  time.Time        `json:"created_at"`
	UserName   *string          `json:"user_name,omitempty"`
}

type AuditLogCreate struct {
	UserID     sql.NullInt64    `json:"user_id"`
	EntityType string           `json:"entity_type" binding:"required"`
	EntityID   int              `json:"entity_id" binding:"required"`
	Action     string           `json:"action" binding:"required"`
	OldValue   *json.RawMessage `json:"old_value"`
	NewValue   *json.RawMessage `json:"new_value"`
}
