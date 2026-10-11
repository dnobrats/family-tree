// Command createadmin creates or updates an admin account.
//
// Usage:
//
//	go run ./cmd/createadmin -username admin -password "mat-khau"
//	make createadmin USER=admin PASS=mat-khau
package main

import (
	"context"
	"flag"
	"log"

	"familytree/internal/config"
	"familytree/internal/database"
	"familytree/internal/repository"
	"familytree/internal/service"
)

func main() {
	username := flag.String("username", "", "admin username")
	password := flag.String("password", "", "admin password (at least 6 characters)")
	flag.Parse()

	if *username == "" || *password == "" {
		log.Fatal("both -username and -password are required")
	}

	ctx := context.Background()

	dsn := config.DatabaseURL()
	if dsn == "" {
		log.Fatal("database not configured: set DATABASE_URL or DB_HOST/DB_USER/DB_NAME")
	}

	pool, err := database.Connect(ctx, dsn)
	if err != nil {
		log.Fatalf("database: %v", err)
	}
	defer pool.Close()

	repo := repository.New(pool)
	svc := service.New(repo.Person, repo.Clan, repo.Spouse, repo.Admin)

	if err := svc.SetAdminPassword(ctx, *username, *password); err != nil {
		log.Fatalf("%v", err)
	}
	log.Printf("admin %q is ready", *username)
}
