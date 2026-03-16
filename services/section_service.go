package services

import (
	"fmt"

	"aura-erp/backend/config"
	"aura-erp/backend/models"
)

func SearchSections(q string, limit int) ([]models.SectionSearchResult, error) {
	if limit <= 0 || limit > 100 {
		limit = 25
	}
	query := `
		SELECT id, name
		FROM sections
		WHERE name ILIKE $1
		ORDER BY name ASC
		LIMIT $2
	`
	rows, err := config.DB.Query(query, fmt.Sprintf("%%%s%%", q), limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []models.SectionSearchResult
	for rows.Next() {
		var r models.SectionSearchResult
		if err := rows.Scan(&r.ID, &r.Name); err != nil {
			return nil, err
		}
		results = append(results, r)
	}
	if results == nil {
		results = []models.SectionSearchResult{}
	}
	return results, nil
}

func GetAllSections() ([]models.Section, error) {
	query := `SELECT id, name, description, created_at FROM sections ORDER BY id ASC`
	rows, err := config.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sections []models.Section
	for rows.Next() {
		var section models.Section
		if err := rows.Scan(&section.ID, &section.Name, &section.Description, &section.CreatedAt); err != nil {
			return nil, err
		}
		sections = append(sections, section)
	}

	return sections, nil
}

func GetSectionByID(id string) (*models.Section, error) {
	query := `SELECT id, name, description, created_at FROM sections WHERE id = $1`

	var section models.Section
	err := config.DB.QueryRow(query, id).Scan(
		&section.ID, &section.Name, &section.Description, &section.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &section, nil
}

func CreateSection(data models.SectionCreate) (*models.Section, error) {
	query := `INSERT INTO sections (name, description)
	          VALUES ($1, $2)
	          RETURNING id, name, description, created_at`

	var section models.Section
	err := config.DB.QueryRow(query, data.Name, data.Description).Scan(
		&section.ID, &section.Name, &section.Description, &section.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &section, nil
}

func UpdateSection(id string, data models.SectionUpdate) (*models.Section, error) {
	query := `UPDATE sections SET name = $1, description = $2
	          WHERE id = $3
	          RETURNING id, name, description, created_at`

	var section models.Section
	err := config.DB.QueryRow(query, data.Name, data.Description, id).Scan(
		&section.ID, &section.Name, &section.Description, &section.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &section, nil
}

func DeleteSection(id string) error {
	query := `DELETE FROM sections WHERE id = $1`
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
