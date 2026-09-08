// Package database contains database connections, migrations, and repositories.
package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"
)

type Database struct {
	DB *sql.DB
}

func New() (*Database, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		getEnv("DB_HOST", "localhost"),
		getEnv("DB_PORT", "5432"),
		getEnv("DB_USER", "appuser"),
		getEnv("DB_PASSWORD", "secret"),
		getEnv("DB_NAME", "appdb"),
		getEnv("DB_SSLMODE", "disable"),
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	// Пул соединений
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err := db.Ping(); err != nil {
		return nil, err
	}

	log.Println("✅ PostgreSQL подключен")
	return &Database{DB: db}, nil
}

func (s *Database) Close() error {
	return s.DB.Close()
}

// Хелпер
func getEnv(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}
