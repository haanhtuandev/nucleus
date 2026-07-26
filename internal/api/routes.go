package api

import (
	"net/http"

	_ "github.com/lib/pq"
)

func NewHandler(a *ApiConfig, rl *RateLimit) http.Handler {
	mux := http.NewServeMux()
	mux.Handle("GET /health", rl.rateLimitMiddleware(healthCheckHandler))
	mux.Handle("POST /login", rl.rateLimitMiddleware(a.loginHandler))
	mux.Handle("POST /signup", rl.rateLimitMiddleware(a.addUserHandler))
	mux.Handle("POST /me/posts", rl.rateLimitMiddleware(a.authorizeMiddleware(a.addPostHandler)))
	mux.Handle("GET /posts", rl.rateLimitMiddleware(a.fetchPostHandler))
	mux.Handle("PUT /me/posts/{post_id}", rl.rateLimitMiddleware(a.authorizeMiddleware(a.updatePostHandler)))
	mux.Handle("GET /me/profile", rl.rateLimitMiddleware(a.authorizeMiddleware(a.getProfileHandler)))
	mux.Handle("GET /me/profile/followers", rl.rateLimitMiddleware(a.authorizeMiddleware(a.getSelfFollowersHandler)))
	mux.Handle("GET /me/profile/followees", rl.rateLimitMiddleware(a.authorizeMiddleware(a.getSelfFolloweesHandler)))
	mux.Handle("GET /profile/{user_id}", rl.rateLimitMiddleware(a.getProfileByIdHandler))
	mux.Handle("GET /profile/{user_id}/followers", rl.rateLimitMiddleware(a.authorizeMiddleware(a.getFollowersHandler)))
	mux.Handle("GET /profile/{user_id}/followees", rl.rateLimitMiddleware(a.authorizeMiddleware(a.getFolloweesHandler)))
	mux.Handle("DELETE /me/posts/{post_id}", rl.rateLimitMiddleware(a.authorizeMiddleware(a.deletePostHandler)))
	mux.Handle("POST /auth/refresh", rl.rateLimitMiddleware(a.refreshHandler))
	mux.Handle("GET /posts/search", rl.rateLimitMiddleware(a.searchPostHandler))
	mux.Handle("GET /follow/{user_id}", rl.rateLimitMiddleware(a.authorizeMiddleware(a.followHandler)))
	mux.Handle("GET /unfollow/{user_id}", rl.rateLimitMiddleware(a.authorizeMiddleware(a.unfollowHandler)))

	return mux
}
