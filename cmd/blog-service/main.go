package main

import (
	"boilerplate/internal/api"
	"boilerplate/internal/database"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

// ---------------------------------------------------------------------------
// Configuration
// ---------------------------------------------------------------------------

type config struct {
	dbURL             string
	secretKey         string
	serverPort        string
	dbMaxOpenConns    int
	dbMaxIdleConns    int
	dbConnMaxLifetime time.Duration
	rateLimitPerHour  int
}

func loadConfig() (config, error) {
	_ = godotenv.Load()

	secretKey := os.Getenv("SECRET")
	if secretKey == "" {
		return config{}, fmt.Errorf("SECRET environment variable is required")
	}
	if len(secretKey) < 32 {
		return config{}, fmt.Errorf("SECRET must be at least 32 characters for security")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	return config{
		dbURL:             os.Getenv("DB_URL"),
		secretKey:         secretKey,
		serverPort:        port,
		dbMaxOpenConns:    25,
		dbMaxIdleConns:    5,
		dbConnMaxLifetime: 5 * time.Minute,
		rateLimitPerHour:  1000,
	}, nil
}

// ---------------------------------------------------------------------------
// Database
// ---------------------------------------------------------------------------

func openDB(cfg config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.dbURL)
	if err != nil {
		return nil, fmt.Errorf("sql.Open: %w", err)
	}

	db.SetMaxOpenConns(cfg.dbMaxOpenConns)
	db.SetMaxIdleConns(cfg.dbMaxIdleConns)
	db.SetConnMaxLifetime(cfg.dbConnMaxLifetime)

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("db.Ping: %w", err)
	}

	return db, nil
}

func runMigrations(db *sql.DB) error {
	log.Println("Running database migrations...")
	if err := database.RunMigrations(db); err != nil {
		return fmt.Errorf("migration failed: %w", err)
	}
	log.Println("Migrations complete!")
	return nil
}

// ---------------------------------------------------------------------------
// HTTP server
// ---------------------------------------------------------------------------

func newServer(cfg config, db *sql.DB) *http.Server {
	dbQueries := database.New(db)

	apiCfg := &api.ApiConfig{
		Database: dbQueries,
		Secret:   cfg.secretKey,
	}

	rl := &api.RateLimit{
		RateMap:      make(map[string]map[int]int),
		LimitPerHour: cfg.rateLimitPerHour,
	}
	rl.StartGlobalNuke()

	mux := api.NewHandler(apiCfg, rl)
	handler := api.CORS(api.DefaultCORSConfig())(mux)

	return &http.Server{
		Addr:    ":" + cfg.serverPort,
		Handler: handler,
	}
}

// ---------------------------------------------------------------------------
// Entry point
// ---------------------------------------------------------------------------

func main() {
	cfg, err := loadConfig()
	if err != nil {
		log.Fatalf("Configuration error: %v", err)
	}

	db, err := openDB(cfg)
	if err != nil {
		log.Fatalf("Database error: %v", err)
	}
	defer db.Close()

	if err := runMigrations(db); err != nil {
		log.Fatalf("Migration error: %v", err)
	}

	srv := newServer(cfg, db)

	log.Printf("Starting server on port %s", cfg.serverPort)
	log.Fatal(srv.ListenAndServe())
}
