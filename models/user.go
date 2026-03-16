package models

import "time"

type User struct {
	ID            int       `json:"id"`
	Name          string    `json:"name"`
	Email         string    `json:"email"`
	PasswordHash  string    `json:"password_hash"`
	Role          string    `json:"role"`
	CreatedAt     time.Time `json:"created_at"`
}


type UserCreate struct {
	Name         string `json:"name" binding:"required"`
	Email        string `json:"email" binding:"required,email"`
	PasswordHash string `json:"password_hash" binding:"required"`
	Role         string `json:"role" binding:"required"`
}

type UserUpdate struct {
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required,email"`
	Role  string `json:"role" binding:"required"`
}