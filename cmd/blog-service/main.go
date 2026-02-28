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
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	dbURL := os.Getenv("DB_URL")
	secret_key := os.Getenv("SECRET")
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
	mux := api.NewHandler(a)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	log.Printf("Serving files on port 8080")
	log.Fatal(http.ListenAndServe(server.Addr, server.Handler))
}
