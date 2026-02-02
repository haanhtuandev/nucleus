package api

import (
	"net/http"
)

func NewHandler() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", healthCheckHandler)
	mux.HandleFunc("/", homePageHandler)

	return mux
}
