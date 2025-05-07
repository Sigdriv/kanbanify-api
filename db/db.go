package db

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Connect() (*pgxpool.Pool, error) {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		return nil, fmt.Errorf("DATABASE_URL is not set")
	}

	conn, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database << %v", err)
	}

	err = conn.Ping(context.Background())
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("unable to ping database << %v", err)
	}

	return conn, nil
}
