package db

import (
	"context"
	"fmt"
	"net/url"

	"github.com/jackc/pgx/v5/pgxpool"
)

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

	return pgxpool.New(context.Background(), dsn)
}

