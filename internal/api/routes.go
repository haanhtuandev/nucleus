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
	mux.HandleFunc("POST /signup", a.addUserHandler)

	mux.Handle("POST /posts", a.authorizeMiddleware(a.addPostHandler))
	mux.HandleFunc("GET /users", a.getAllUsersHandler)
	mux.HandleFunc("GET /posts", a.getAllPostsHandler)
	mux.HandleFunc("GET /posts/{slug}", a.getPostBySlugHandler)
	mux.Handle("PUT /posts/{post_id}", a.authorizeMiddleware(a.updatePostHandler))
	mux.Handle("GET /profile", a.authorizeMiddleware(a.getProfile))
	mux.HandleFunc("GET /api/refresh", a.refreshHandler)
	mux.Handle("DELETE /posts/{post_id}", a.authorizeMiddleware(a.deletePostHandler))

	return mux
}
