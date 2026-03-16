package models

import "time"

type ProductSearchResult struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Product struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	UnitPrice   float64   `json:"unit_price"`
	Unit        string    `json:"unit"`
	CreatedAt   time.Time `json:"created_at"`
}

type ProductCreate struct {
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description"`
	UnitPrice   float64 `json:"unit_price" binding:"required,gt=0"`
	Unit        string  `json:"unit" binding:"required"`
}

type ProductUpdate struct {
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description"`
	UnitPrice   float64 `json:"unit_price" binding:"required,gt=0"`
	Unit        string  `json:"unit" binding:"required"`
}
