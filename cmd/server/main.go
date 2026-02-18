package main

import (
	"log"
	"net/http"
	"os"

	"genealogy-be/internal/api"
	"genealogy-be/internal/db"
	"genealogy-be/internal/middleware"
)

func main() {
	cfg := map[string]string{
		"host":     os.Getenv("DB_HOST"),
		"port":     os.Getenv("DB_PORT"),
		"user":     os.Getenv("DB_USER"),
		"password": os.Getenv("DB_PASSWORD"),
		"name":     os.Getenv("DB_NAME"),
		"schema":   os.Getenv("DB_SCHEMA"),
	}

	pool, err := db.NewPostgres(cfg)
	if err != nil {
		log.Fatal(err)
	}

	router := api.NewRouter(pool)
	handler := middleware.RateLimit(5)(router)

	log.Println("Server listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}

