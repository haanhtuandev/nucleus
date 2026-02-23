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

	mux.Handle("POST /me/posts", a.authorizeMiddleware(a.addPostHandler))
	mux.HandleFunc("GET /users", a.getAllUsersHandler)
	mux.HandleFunc("GET /posts", a.fetchPostHandler)
	mux.HandleFunc("GET /posts/{slug}", a.getPostBySlugHandler)
	mux.Handle("PUT /me/posts/{post_id}", a.authorizeMiddleware(a.updatePostHandler))
	mux.Handle("GET /me/profile", a.authorizeMiddleware(a.getProfileHandler))
	mux.HandleFunc("GET /profile/{user_id}", a.getProfileByIdHandler)
	mux.HandleFunc("GET /api/reset", a.resetHandler)
	mux.Handle("DELETE /me/posts/{post_id}", a.authorizeMiddleware(a.deletePostHandler))
	mux.HandleFunc("POST /auth/refresh", a.refreshHandler)
	mux.HandleFunc("GET /posts/search", a.searchPostHandler)

	return mux
}
