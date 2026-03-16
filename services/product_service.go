package services

import (
	"fmt"

	"aura-erp/backend/config"
	"aura-erp/backend/models"
)

func SearchProducts(q string, limit int) ([]models.ProductSearchResult, error) {
	if limit <= 0 || limit > 100 {
		limit = 25
	}
	query := `
		SELECT id, name
		FROM products
		WHERE name ILIKE $1
		ORDER BY name ASC
		LIMIT $2
	`
	rows, err := config.DB.Query(query, fmt.Sprintf("%%%s%%", q), limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []models.ProductSearchResult
	for rows.Next() {
		var r models.ProductSearchResult
		if err := rows.Scan(&r.ID, &r.Name); err != nil {
			return nil, err
		}
		results = append(results, r)
	}
	if results == nil {
		results = []models.ProductSearchResult{}
	}
	return results, nil
}

func GetAllProducts() ([]models.Product, error) {
	query := `SELECT id, name, description, unit_price, unit, created_at FROM products ORDER BY id ASC`
	rows, err := config.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var product models.Product
		if err := rows.Scan(&product.ID, &product.Name, &product.Description, &product.UnitPrice, &product.Unit, &product.CreatedAt); err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	return products, nil
}

func GetProductByID(id string) (*models.Product, error) {
	query := `SELECT id, name, description, unit_price, unit, created_at FROM products WHERE id = $1`

	var product models.Product
	err := config.DB.QueryRow(query, id).Scan(
		&product.ID, &product.Name, &product.Description, &product.UnitPrice, &product.Unit, &product.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &product, nil
}

func CreateProduct(data models.ProductCreate) (*models.Product, error) {
	query := `INSERT INTO products (name, description, unit_price, unit)
	          VALUES ($1, $2, $3, $4)
	          RETURNING id, name, description, unit_price, unit, created_at`

	var product models.Product
	err := config.DB.QueryRow(query, data.Name, data.Description, data.UnitPrice, data.Unit).Scan(
		&product.ID, &product.Name, &product.Description, &product.UnitPrice, &product.Unit, &product.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &product, nil
}

func UpdateProduct(id string, data models.ProductUpdate) (*models.Product, error) {
	query := `UPDATE products SET name = $1, description = $2, unit_price = $3, unit = $4
	          WHERE id = $5
	          RETURNING id, name, description, unit_price, unit, created_at`

	var product models.Product
	err := config.DB.QueryRow(query, data.Name, data.Description, data.UnitPrice, data.Unit, id).Scan(
		&product.ID, &product.Name, &product.Description, &product.UnitPrice, &product.Unit, &product.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &product, nil
}

func DeleteProduct(id string) error {
	query := `DELETE FROM products WHERE id = $1`
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
