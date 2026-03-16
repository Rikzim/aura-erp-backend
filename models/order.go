package models

import "time"

type Order struct {
	ID                int         `json:"id"`
	Reference         string      `json:"reference"`
	ProposalID        NullInt64   `json:"proposal_id"`
	ClientID          int         `json:"client_id"`
	SectionID         NullInt64   `json:"section_id"`
	Status            string      `json:"status"`
	DueDate           NullTime    `json:"due_date"`
	CreatedAt         time.Time   `json:"created_at"`
	UpdatedAt         time.Time   `json:"updated_at"`
	ClientName        string      `json:"client_name,omitempty"`
	SectionName       NullString  `json:"section_name,omitempty"`
	ProposalReference NullString  `json:"proposal_reference,omitempty"`
	Items             []OrderItem `json:"items,omitempty"`
}

type OrderCreate struct {
	Reference  string    `json:"reference" binding:"required"`
	ProposalID NullInt64 `json:"proposal_id"`
	ClientID   FlexInt   `json:"client_id" binding:"required"`
	SectionID  NullInt64 `json:"section_id"`
	Status     string    `json:"status"`
	DueDate    NullTime  `json:"due_date"`
}

type OrderUpdate struct {
	Reference  string    `json:"reference" binding:"required"`
	ProposalID NullInt64 `json:"proposal_id"`
	ClientID   FlexInt   `json:"client_id" binding:"required"`
	SectionID  NullInt64 `json:"section_id"`
	Status     string    `json:"status" binding:"required"`
	DueDate    NullTime  `json:"due_date"`
}

type OrderItem struct {
	ID          int     `json:"id"`
	OrderID     int     `json:"order_id"`
	ProductID   int     `json:"product_id"`
	Quantity    float64 `json:"quantity"`
	UnitPrice   float64 `json:"unit_price"`
	Notes       *string `json:"notes"`
	ProductName string  `json:"product_name,omitempty"`
	Unit        string  `json:"unit,omitempty"`
}

type OrderItemCreate struct {
	ProductID FlexInt `json:"product_id" binding:"required"`
	Quantity  float64 `json:"quantity" binding:"required"`
	UnitPrice float64 `json:"unit_price" binding:"required"`
	Notes     *string `json:"notes"`
}

type OrderItemUpdate struct {
	ProductID FlexInt `json:"product_id" binding:"required"`
	Quantity  float64 `json:"quantity" binding:"required"`
	UnitPrice float64 `json:"unit_price" binding:"required"`
	Notes     *string `json:"notes"`
}
