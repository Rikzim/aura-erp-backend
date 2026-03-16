package services

import (
	"aura-erp/backend/config"
	"aura-erp/backend/models"
)

func GetAllUsers() ([]models.User, error) {
	query := `SELECT id, name, email, role, created_at FROM users ORDER BY id ASC`
	rows, err := config.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Role, &user.CreatedAt); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func GetUserByID(id string) (*models.User, error) {
	query := `SELECT id, name, email, role, created_at FROM users WHERE id = $1`
	
	var user models.User
	err := config.DB.QueryRow(query, id).Scan(
		&user.ID, &user.Name, &user.Email, &user.Role, &user.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func CreateUser(data models.UserCreate) (*models.User, error) {
	query := `INSERT INTO users (name, email, password_hash, role) 
	          VALUES ($1, $2, $3, $4) 
	          RETURNING id, name, email, role, created_at`
	
	var user models.User
	err := config.DB.QueryRow(query, data.Name, data.Email, data.PasswordHash, data.Role).Scan(
		&user.ID, &user.Name, &user.Email, &user.Role, &user.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func UpdateUser(id string, data models.UserUpdate) (*models.User, error) {
	query := `UPDATE users SET name = $1, email = $2, role = $3 
	          WHERE id = $4 
	          RETURNING id, name, email, role, created_at`
	
	var user models.User
	err := config.DB.QueryRow(query, data.Name, data.Email, data.Role, id).Scan(
		&user.ID, &user.Name, &user.Email, &user.Role, &user.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func DeleteUser(id string) error {
	query := `DELETE FROM users WHERE id = $1`
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
