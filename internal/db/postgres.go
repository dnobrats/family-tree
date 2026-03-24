package db

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// NewPostgres tạo connection pool tới PostgreSQL
func NewPostgres(cfg map[string]string) (*pgxpool.Pool, error) {
	schema := cfg["schema"]
	if schema == "" {
		schema = "public"
	}

	q := url.Values{}
	q.Set("search_path", schema)

	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?%s",
		cfg["user"],
		cfg["password"],
		cfg["host"],
		cfg["port"],
		cfg["name"],
		q.Encode(),
	)

	// Parse config với timeout
	poolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("parse database config: %w", err)
	}

	// Cấu hình connection pool
	poolConfig.MaxConns = 25
	poolConfig.MinConns = 5
	poolConfig.MaxConnLifetime = time.Hour
	poolConfig.MaxConnIdleTime = 30 * time.Minute

	// Tạo pool với context timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("create connection pool: %w", err)
	}

	// Ping để kiểm tra kết nối
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping database: %w", err)
	}

	return pool, nil
}
