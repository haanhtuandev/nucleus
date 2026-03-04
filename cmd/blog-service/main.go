package main

import (
	"boilerplate/internal/api"
	"boilerplate/internal/database"
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	secret_key := os.Getenv("SECRET")
	if secret_key == "" {
		log.Fatal("SECRET environment variable is required")
	}
	if len(secret_key) < 32 {
		log.Fatal("SECRET must be at least 32 characters for security")
	}
	db, err := sql.Open("postgres", dbURL)

	if err != nil {
		log.Fatal(err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Add health check
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	dbQueries := database.New(db)
	a := &api.ApiConfig{Database: dbQueries, Secret: secret_key}
	rl := &api.RateLimit{RateMap: make(map[string](map[int]int)),
		LimitPerHour: 1000}
	rl.StartGlobalNuke()
	mux := api.NewHandler(a, rl)

	// Wrap with CORS middleware
	corsConfig := api.DefaultCORSConfig()
	handler := api.CORS(corsConfig)(mux)

	server := &http.Server{
		Addr:    ":8080",
		Handler: handler,
	}
	log.Printf("Serving files on port 8080")
	log.Fatal(http.ListenAndServe(server.Addr, server.Handler))
}
