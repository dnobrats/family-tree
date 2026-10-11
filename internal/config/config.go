// Package config loads and validates runtime configuration from environment
// variables. It fails fast, so a misconfigured server never starts half-alive.
package config

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

// Config is the fully validated application configuration.
type Config struct {
	Port          int
	DatabaseURL   string
	SessionSecret string
	SessionTTL    time.Duration
	CookieSecure  bool
	AdminUsername string
	AdminPassword string
}

// Load reads configuration from the environment and validates it.
func Load() (*Config, error) {
	cfg := &Config{
		Port:          envInt("SERVER_PORT", 8080),
		DatabaseURL:   DatabaseURL(),
		SessionSecret: strings.TrimSpace(os.Getenv("SESSION_SECRET")),
		SessionTTL:    24 * time.Hour,
		CookieSecure:  envBool("COOKIE_SECURE", false),
		AdminUsername: strings.TrimSpace(os.Getenv("ADMIN_USERNAME")),
		AdminPassword: os.Getenv("ADMIN_PASSWORD"),
	}
	if err := cfg.validate(); err != nil {
		return nil, err
	}
	return cfg, nil
}

// DatabaseURL returns DATABASE_URL, or a connection string assembled from the
// DB_* variables when DATABASE_URL is empty.
func DatabaseURL() string {
	if v := strings.TrimSpace(os.Getenv("DATABASE_URL")); v != "" {
		return v
	}
	return buildDSN()
}

func (c *Config) validate() error {
	if c.DatabaseURL == "" {
		return fmt.Errorf("database not configured: set DATABASE_URL or DB_HOST/DB_USER/DB_NAME")
	}
	if len(c.SessionSecret) < 16 {
		return fmt.Errorf("SESSION_SECRET must be at least 16 characters (try: openssl rand -hex 32)")
	}
	if c.Port <= 0 || c.Port > 65535 {
		return fmt.Errorf("SERVER_PORT %d is out of range", c.Port)
	}
	return nil
}

// buildDSN assembles a connection string from the DB_* variables.
func buildDSN() string {
	host := env("DB_HOST", "")
	name := env("DB_NAME", "")
	if host == "" || name == "" {
		return ""
	}
	u := url.URL{
		Scheme: "postgres",
		Host:   host + ":" + env("DB_PORT", "5432"),
		Path:   "/" + name,
	}
	if user := env("DB_USER", ""); user != "" {
		u.User = url.UserPassword(user, os.Getenv("DB_PASSWORD"))
	}
	q := u.Query()
	q.Set("sslmode", env("DB_SSLMODE", "disable"))
	u.RawQuery = q.Encode()
	return u.String()
}

func env(key, def string) string {
	if v := strings.TrimSpace(os.Getenv(key)); v != "" {
		return v
	}
	return def
}

func envInt(key string, def int) int {
	if v, err := strconv.Atoi(env(key, "")); err == nil {
		return v
	}
	return def
}

func envBool(key string, def bool) bool {
	switch strings.ToLower(env(key, "")) {
	case "1", "true", "yes", "on":
		return true
	case "0", "false", "no", "off":
		return false
	}
	return def
}
