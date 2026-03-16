package services

import (
	"fmt"

	"aura-erp/backend/config"
	"aura-erp/backend/models"
)

func SearchClients(q string, limit int) ([]models.ClientSearchResult, error) {
	if limit <= 0 || limit > 100 {
		limit = 25
	}
	query := `
		SELECT id, name
		FROM clients
		WHERE name ILIKE $1
		ORDER BY name ASC
		LIMIT $2
	`
	rows, err := config.DB.Query(query, fmt.Sprintf("%%%s%%", q), limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []models.ClientSearchResult
	for rows.Next() {
		var r models.ClientSearchResult
		if err := rows.Scan(&r.ID, &r.Name); err != nil {
			return nil, err
		}
		results = append(results, r)
	}
	if results == nil {
		results = []models.ClientSearchResult{}
	}
	return results, nil
}

func GetAllClients() ([]models.Client, error) {
	query := `SELECT id, name, email, phone, address, vat_number, notes, created_at, updated_at FROM clients ORDER BY id ASC`
	rows, err := config.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var clients []models.Client
	for rows.Next() {
		var client models.Client
		if err := rows.Scan(&client.ID, &client.Name, &client.Email, &client.Phone, &client.Address, &client.VatNumber, &client.Notes, &client.CreatedAt, &client.UpdatedAt); err != nil {
			return nil, err
		}
		clients = append(clients, client)
	}

	return clients, nil
}

func GetClientByID(id string) (*models.Client, error) {
	query := `SELECT id, name, email, phone, address, vat_number, notes, created_at, updated_at FROM clients WHERE id = $1`

	var client models.Client
	err := config.DB.QueryRow(query, id).Scan(
		&client.ID, &client.Name, &client.Email, &client.Phone, &client.Address, &client.VatNumber, &client.Notes, &client.CreatedAt, &client.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &client, nil
}

func CreateClient(data models.ClientCreate) (*models.Client, error) {
	query := `INSERT INTO clients (name, email, phone, address, vat_number, notes)
	          VALUES ($1, $2, $3, $4, $5, $6)
	          RETURNING id, name, email, phone, address, vat_number, notes, created_at, updated_at`

	var client models.Client
	err := config.DB.QueryRow(query, data.Name, data.Email, data.Phone, data.Address, data.VatNumber, data.Notes).Scan(
		&client.ID, &client.Name, &client.Email, &client.Phone, &client.Address, &client.VatNumber, &client.Notes, &client.CreatedAt, &client.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &client, nil
}

func UpdateClient(id string, data models.ClientUpdate) (*models.Client, error) {
	query := `UPDATE clients SET name = $1, email = $2, phone = $3, address = $4, vat_number = $5, notes = $6, updated_at = NOW()
	          WHERE id = $7
	          RETURNING id, name, email, phone, address, vat_number, notes, created_at, updated_at`

	var client models.Client
	err := config.DB.QueryRow(query, data.Name, data.Email, data.Phone, data.Address, data.VatNumber, data.Notes, id).Scan(
		&client.ID, &client.Name, &client.Email, &client.Phone, &client.Address, &client.VatNumber, &client.Notes, &client.CreatedAt, &client.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &client, nil
}

func DeleteClient(id string) error {
	query := `DELETE FROM clients WHERE id = $1`
	result, err := config.DB.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected == 0 {
		return err
	}

	return nil
}
