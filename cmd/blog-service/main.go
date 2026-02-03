package main

import (
	"boilerplate/internal/api"
	"boilerplate/internal/database"
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	dbURL := os.Getenv("DB_URL")

	log.Println(dbURL)
	db, err := sql.Open("postgres", dbURL)

	if err != nil {
		log.Printf("Database connection error %v", err)
	}
	dbQueries := database.New(db)
	a := &api.ApiConfig{Database: dbQueries}
	mux := api.NewHandler(a)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	log.Printf("Serving files on port 8080")
	log.Fatal(http.ListenAndServe(server.Addr, server.Handler))
}
