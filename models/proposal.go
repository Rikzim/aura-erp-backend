package models

import "time"

type ProposalSearchResult struct {
	ID        int    `json:"id"`
	Reference string `json:"reference"`
}

type Proposal struct {
	ID          int            `json:"id"`
	Reference   string         `json:"reference"`
	ClientID    int            `json:"client_id"`
	SectionID   NullInt64      `json:"section_id"`
	Status      string         `json:"status"`
	Notes       *string        `json:"notes"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	ClientName  string         `json:"client_name,omitempty"`
	SectionName NullString     `json:"section_name,omitempty"`
	Items       []ProposalItem `json:"items,omitempty"`
}

type ProposalCreate struct {
	Reference string    `json:"reference" binding:"required"`
	ClientID  FlexInt   `json:"client_id" binding:"required"`
	SectionID NullInt64 `json:"section_id"`
	Status    string    `json:"status"`
	Notes     *string   `json:"notes"`
}

type ProposalUpdate struct {
	Reference string    `json:"reference" binding:"required"`
	ClientID  FlexInt   `json:"client_id" binding:"required"`
	SectionID NullInt64 `json:"section_id"`
	Status    string    `json:"status" binding:"required"`
	Notes     *string   `json:"notes"`
}

type ProposalItem struct {
	ID          int     `json:"id"`
	ProposalID  int     `json:"proposal_id"`
	ProductID   int     `json:"product_id"`
	Quantity    float64 `json:"quantity"`
	UnitPrice   float64 `json:"unit_price"`
	Notes       *string `json:"notes"`
	ProductName string  `json:"product_name,omitempty"`
	Unit        string  `json:"unit,omitempty"`
}

type ProposalItemCreate struct {
	ProductID FlexInt `json:"product_id" binding:"required"`
	Quantity  float64 `json:"quantity" binding:"required"`
	UnitPrice float64 `json:"unit_price" binding:"required"`
	Notes     *string `json:"notes"`
}

type ProposalItemUpdate struct {
	ProductID FlexInt `json:"product_id" binding:"required"`
	Quantity  float64 `json:"quantity" binding:"required"`
	UnitPrice float64 `json:"unit_price" binding:"required"`
	Notes     *string `json:"notes"`
}
