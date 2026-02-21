package api

import (
	"boilerplate/internal/auth"
	"context"
	"net/http"
)

type ctxKey int

const userIDKey ctxKey = iota

func (a *ApiConfig) authorizeMiddleware(next http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := auth.GetBearerToken(r.Header)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "no jwt token found", err)
			return
		}
		user_id, err := auth.ValidateJWT(token, a.Secret)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "error validating jwt", err)
			return
		}
		ctx := context.WithValue(r.Context(), userIDKey, user_id)
		next(w, r.WithContext(ctx))
	})
}
