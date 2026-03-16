package config

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB() error {
	connStr := os.Getenv("DATABASE_URL")

	if connStr == "" {
		connStr = "postgres://myuser:mypassword@localhost:5433/crud_db?sslmode=disable"
	}

	var err error

	DB, err = sql.Open("postgres", connStr)

	if err != nil {
		return fmt.Errorf("error opening database: %w", err)
	}

	if err = DB.Ping(); err != nil {
		return fmt.Errorf("error pinging database: %w", err)
	}

	// Set connection pool settings
	DB.SetMaxOpenConns(25)
	DB.SetMaxIdleConns(5)

	return nil
}

func CloseDB() {
	if DB != nil {
		DB.Close()
	}
}
