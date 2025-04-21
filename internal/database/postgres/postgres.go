package postgres

import (
	"database/sql"
	"fmt"
	"go-books-api/internal/config"

	_ "github.com/lib/pq"
)

type DB struct {
	*sql.DB
}

func New(cfg *config.PostgresConfig) (DB, error) {
	dsn := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=%s host=%s port=%d", cfg.User, cfg.Password, cfg.Name, cfg.SSLMode, cfg.Host, cfg.Port)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return DB{}, fmt.Errorf("can`t connect to database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return DB{}, fmt.Errorf("failed connecting to database: %w", err)
	}

	return DB{db}, nil
}
