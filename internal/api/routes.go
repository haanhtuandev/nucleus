package api

import (
	"net/http"

	_ "github.com/lib/pq"
)

func NewHandler(a *ApiConfig) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", healthCheckHandler)
	mux.HandleFunc("/", homePageHandler)
	mux.HandleFunc("POST /login", a.loginHandler)
	mux.HandleFunc("POST /users", a.addUserHandler)

	mux.HandleFunc("POST /posts", a.addPostHandler)
	mux.HandleFunc("GET /users", a.getAllUsersHandler)
	mux.HandleFunc("GET /posts", a.getAllPostsHandler)
	mux.HandleFunc("GET /posts/{slug}", a.getPostBySlugHandler)
	mux.HandleFunc("GET /users/{user_id}", a.getUserByIdHandler)
	mux.HandleFunc("GET /api/refresh", a.refreshHandler)
	mux.HandleFunc("DELETE /posts/{post_id}", a.deletePostHandler)

	return mux
}
