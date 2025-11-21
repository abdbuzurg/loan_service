package database

import (
	"context"
	"fmt"
	"loan_service/configs"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPostgresConnection(cfg configs.DatabaseConfig) (*pgxpool.Pool, error) {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.Password,
		cfg.DBName,
		cfg.SSLMode,
	)

	pool, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %v", err)
	}

	log.Printf("Successfully connected to PostgreSQL.")
	return pool, nil
}
