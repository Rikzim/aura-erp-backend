package models

import "time"

type ClientSearchResult struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Client struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     *string   `json:"email"`
	Phone     *string   `json:"phone"`
	Address   *string   `json:"address"`
	VatNumber *string   `json:"vat_number"`
	Notes     *string   `json:"notes"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ClientCreate struct {
	Name      string  `json:"name" binding:"required"`
	Email     *string `json:"email"`
	Phone     *string `json:"phone"`
	Address   *string `json:"address"`
	VatNumber *string `json:"vat_number"`
	Notes     *string `json:"notes"`
}

type ClientUpdate struct {
	Name      string  `json:"name" binding:"required"`
	Email     *string `json:"email"`
	Phone     *string `json:"phone"`
	Address   *string `json:"address"`
	VatNumber *string `json:"vat_number"`
	Notes     *string `json:"notes"`
}
