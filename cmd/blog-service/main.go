package main

import (
	"boilerplate/internal/api"
	"log"
	"net/http"
)

func main() {
	router := api.NewHandler()

	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	log.Printf("Serving files on port 8080")
	log.Fatal(http.ListenAndServe(server.Addr, server.Handler))
}
