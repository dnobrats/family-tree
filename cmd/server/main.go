// Command server runs the family tree web application.
//
// Flow: load config -> connect to Postgres -> apply migrations -> wire the
// service and HTTP layers -> serve (with graceful shutdown).
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"familytree/internal/auth"
	"familytree/internal/config"
	"familytree/internal/database"
	httpx "familytree/internal/http"
	"familytree/internal/repository"
	"familytree/internal/service"
	"familytree/migrations"
	"familytree/web"
)

func main() {
	migrateOnly := flag.Bool("migrate", false, "apply database migrations and exit")
	healthcheck := flag.Bool("healthcheck", false, "check /healthz and exit (used by docker healthcheck)")
	flag.Parse()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	if *healthcheck {
		os.Exit(runHealthcheck(cfg.Port))
	}

	ctx := context.Background()

	pool, err := database.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("database: %v", err)
	}
	defer pool.Close()
	log.Println("connected to database")

	if err := database.Migrate(ctx, pool, migrations.FS); err != nil {
		log.Fatalf("migrate: %v", err)
	}
	log.Println("migrations applied")

	if *migrateOnly {
		return
	}

	repo := repository.New(pool)
	svc := service.New(repo.Person, repo.Clan, repo.Spouse, repo.Admin)

	if err := svc.EnsureBootstrapAdmin(ctx, cfg.AdminUsername, cfg.AdminPassword); err != nil {
		log.Printf("bootstrap admin skipped: %v", err)
	}

	renderer, err := httpx.NewRenderer(web.FS)
	if err != nil {
		log.Fatalf("templates: %v", err)
	}

	authMgr := auth.NewManager(cfg.SessionSecret, cfg.SessionTTL)
	server := httpx.NewServer(svc, authMgr, cfg, renderer)

	httpServer := &http.Server{
		Addr:              fmt.Sprintf(":%d", cfg.Port),
		Handler:           server.Routes(),
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      60 * time.Second,
		IdleTimeout:       120 * time.Second,
	}

	go func() {
		log.Printf("listening on %s", httpServer.Addr)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		log.Printf("shutdown: %v", err)
	}
	log.Println("stopped")
}

// runHealthcheck requests /healthz on localhost and returns a process exit code.
func runHealthcheck(port int) int {
	client := &http.Client{Timeout: 3 * time.Second}
	resp, err := client.Get(fmt.Sprintf("http://127.0.0.1:%d/healthz", port))
	if err != nil {
		return 1
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return 1
	}
	return 0
}
