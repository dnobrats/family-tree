package main

import (
	"fmt"
	"log"
	"net/http"

	"genealogy-be/internal/api"
	"genealogy-be/internal/config"
	"genealogy-be/internal/db"
	"genealogy-be/internal/middleware"
)

func main() {
	// Load config từ environment variables
	cfg, err := config.LoadFromEnv()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Kết nối database
	dbConfig := map[string]string{
		"host":     cfg.Database.Host,
		"port":     cfg.Database.Port,
		"user":     cfg.Database.User,
		"password": cfg.Database.Password,
		"name":     cfg.Database.Name,
		"schema":   cfg.Database.Schema,
	}

	pool, err := db.NewPostgres(dbConfig)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer pool.Close()

	log.Println("✅ Connected to database successfully")

	// Khởi tạo router với dependency injection
	router := api.NewRouter(pool)
	
	// Apply rate limiting middleware
	handler := middleware.RateLimit(5)(router)

	// Start server
	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	log.Printf("🚀 Server listening on %s", addr)
	
	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
