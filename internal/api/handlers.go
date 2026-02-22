package api

import (
	"boilerplate/internal/auth"
	"boilerplate/internal/database"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type ApiConfig struct {
	Database *database.Queries
	Secret   string
}

type User struct {
	ID            uuid.UUID `json:"id"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	Username      string    `json:"username"`
	Token         string    `json:"token"`
	Refresh_token string    `json:"refresh_token"`
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte("Service is healthy!"))
}

func homePageHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./static/index.html")
}

func (a *ApiConfig) addUserHandler(w http.ResponseWriter, r *http.Request) {
	type params struct {
		Username string  `json:"username"`
		Password string  `json:"password"`
		Bio      *string `json:"bio"`
	}
	param := params{}
	if err := json.NewDecoder(r.Body).Decode(&param); err != nil {
		// 400 because client send bad json
		respondWithError(w, 400, "Invalid request payload", err)
		return
	}

	hashed_password, err := auth.HashPassword(param.Password)
	if err != nil {
		respondWithError(w, 500, "error with hash password", err)
		return
	}

	dbParams := database.CreateUserParams{
		Username:       param.Username,
		HashedPassword: hashed_password,
	}

	if param.Bio != nil {
		dbParams.Bio = sql.NullString{String: *param.Bio, Valid: true}
	} else {
		dbParams.Bio = sql.NullString{Valid: false}
	}

	user, err := a.Database.CreateUser(r.Context(), dbParams)

	if err != nil {
		if isUniqueViolation(err) {
			respondWithError(w, 400, "Username already exists", err)
			return
		}
		respondWithError(w, 500, "Couldn't create user", err)
		return
	}
	respondWithJSON(w, 201, user)
}

func (a *ApiConfig) loginHandler(w http.ResponseWriter, r *http.Request) {
	type params struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	param := params{}
	if err := json.NewDecoder(r.Body).Decode(&param); err != nil {
		// 400 because client send bad json
		respondWithError(w, 400, "Invalid request payload", err)
		return
	}
	user, err := a.Database.GetUserByName(r.Context(), param.Username)
	if err != nil {
		if isUniqueViolation(err) {
			respondWithError(w, 400, "username already exists", err)
			return
		} else {
			respondWithError(w, 500, "error retrieving user", err)
			return
		}
	}
	authenticated, err := auth.CheckPasswordHash(param.Password, user.HashedPassword)
	if !authenticated {
		respondWithError(w, 400, "Wrong password", err)
		return
	}
	if err != nil {
		respondWithError(w, 500, "Error validating password", err)
		return
	}
	refresh_token_string, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError(w, 500, "error making refresh tokens", err)
		return
	}
	refresh_token_param := database.CreateRefreshTokenParams{
		Token:  refresh_token_string,
		UserID: user.ID,
	}
	_, err = a.Database.CreateRefreshToken(r.Context(), refresh_token_param)

	access_token, err := auth.MakeJWT(
		user.ID,
		a.Secret,
		time.Duration(1)*time.Hour,
	)
	if err != nil {
		respondWithError(w, 400, "error generating jwt token", err)
	}
	user_json := User{
		ID:            user.ID,
		CreatedAt:     user.CreatedAt,
		UpdatedAt:     user.UpdatedAt,
		Username:      user.Username,
		Token:         access_token,
		Refresh_token: refresh_token_string,
	}
	respondWithJSON(w, 200, user_json)

}

func (a *ApiConfig) addPostHandler(w http.ResponseWriter, r *http.Request) {
	type params struct {
		Title   string    `json:"title"`
		Content string    `json:"content"`
		UserID  uuid.UUID `json:"user_id"`
	}
	param := params{}
	if err := json.NewDecoder(r.Body).Decode(&param); err != nil {
		// 400 because client send bad json
		respondWithError(w, 400, "Invalid request payload", err)
		return
	}

	user_id, err := GetUserID(r.Context())
	if err != nil {
		respondWithError(w, 500, "wrong user_id format", err)
		return
	}
	if param.UserID != user_id {
		respondWithError(w, http.StatusUnauthorized, "method not allowed", errors.New("method not allowed!"))
		return
	}

	slug := cleanTitle(param.Title)

	for i := 0; i < 5; i++ {
		dbParams := database.CreatePostParams{
			Content: param.Content,
			Title:   param.Title,
			UserID:  param.UserID,
			Slug:    slug,
		}
		post, err := a.Database.CreatePost(r.Context(), dbParams)

		if err == nil {
			respondWithJSON(w, 201, post)
			return
		}

		if isUniqueViolation(err) {
			slug = fmt.Sprintf("%s-%s", slug, generateRandomString(5))
			continue
		}

		respondWithError(w, 500, "DB Error", err)
		return
	}
}

func (a *ApiConfig) getAllUsersHandler(w http.ResponseWriter, r *http.Request) {
	users, err := a.Database.GetAllUsers(r.Context())
	if err != nil {
		respondWithError(w, 500, "Couldn't retrieve all users", err)
		return
	}
	respondWithJSON(w, 201, users)
}

func (a *ApiConfig) getAllPostsHandler(w http.ResponseWriter, r *http.Request) {
	posts, err := a.Database.GetAllPosts(r.Context())
	if err != nil {
		respondWithError(w, 500, "Couldn't retrieve all posts", err)
		return
	}
	respondWithJSON(w, 201, posts)
}

// func (a *ApiConfig) getPostByIdHandler(w http.ResponseWriter, r *http.Request) {
// 	post_id := r.PathValue("post_id")
// 	parsed_post_id, err := uuid.Parse(post_id)
// 	if err != nil {
// 		respondWithError(w, 500, "Error converting post ID", err)
// 		return
// 	}

// 	post, err := a.Database.GetPostById(r.Context(), parsed_post_id)
// 	if err != nil {
// 		respondWithError(w, 500, "Error while finding post", err)
// 		return
// 	}

//		respondWithJSON(w, 201, post)
//	}
func (a *ApiConfig) getPostBySlugHandler(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")

	post, err := a.Database.GetPostBySlug(r.Context(), slug)
	if err != nil {
		respondWithError(w, 500, "Error while finding post", err)
		return
	}

	respondWithJSON(w, 201, post)
}
func (a *ApiConfig) getUserByIdHandler(w http.ResponseWriter, r *http.Request) {
	user_id := r.PathValue("user_id")
	parsed_user_id, err := uuid.Parse(user_id)
	if err != nil {
		respondWithError(w, 500, "Error converting user ID", err)
		return
	}

	user, err := a.Database.GetUserById(r.Context(), parsed_user_id)
	if err != nil {
		respondWithError(w, 500, "Error while finding user", err)
		return
	}

	respondWithJSON(w, 201, user)
}

func (a *ApiConfig) getProfile(w http.ResponseWriter, r *http.Request) {
	user_id, err := GetUserID(r.Context())
	if err != nil {
		respondWithError(w, 500, "error getting user ID from context", err)
		return
	}
	user, err := a.Database.GetUserById(r.Context(), user_id)
	if err != nil {
		respondWithError(w, 500, "error retrieving user", err)
		return
	}
	posts, err := a.Database.GetPostsByUser(r.Context(), user_id)
	if err != nil {
		respondWithError(w, 500, "error retrieving posts", err)
		return
	}

	type profileVal struct {
		Username      string
		PostsMetadata []database.GetPostsByUserRow
		PostCount     int
	}

	returnVal := profileVal{
		Username:      user.Username,
		PostsMetadata: posts,
		PostCount:     len(posts),
	}
	respondWithJSON(w, 200, returnVal)
}

func (a *ApiConfig) updatePostHandler(w http.ResponseWriter, r *http.Request) {
	post_id := r.PathValue("post_id")
	parsed_post_id, err := uuid.Parse(post_id)
	if err != nil {
		respondWithError(w, 400, "Unable to parse ID", err)
		return
	}
	owner_id, err := GetUserID(r.Context())
	if err != nil {
		respondWithError(w, 500, "Error retrieving user_id", err)
		return
	}
	post, err := a.Database.GetPostById(r.Context(), parsed_post_id)
	if err != nil {
		respondWithError(w, 500, "Error retrieving post by id", err)
		return
	}
	// authorization check
	if post.UserID != owner_id {
		respondWithError(w, http.StatusUnauthorized, "method not allowed", err)
		return
	}
	// input params
	type params struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}
	param := params{}
	if err := json.NewDecoder(r.Body).Decode(&param); err != nil {
		// 400 because client send bad json
		respondWithError(w, 400, "Invalid request payload", err)
		return
	}

	// update logic
	update_params := database.UpdatePostInfoParams{
		Title:   param.Title,
		Content: param.Content,
		ID:      parsed_post_id,
	}
	err = a.Database.UpdatePostInfo(r.Context(), update_params)
	if err != nil {
		respondWithError(w, 500, "error updating post info", err)
		return
	}
	respondWithJSON(w, 200, nil)

}

func (a *ApiConfig) resetHandler(w http.ResponseWriter, r *http.Request) {
	err := a.Database.RefreshDB(r.Context())
	if err != nil {
		respondWithError(w, 500, "Error clearing database", err)
		return
	}
	respondWithJSON(w, 201, nil)
}

func (a *ApiConfig) refreshHandler(w http.ResponseWriter, r *http.Request) {
	type params struct {
		Refresh_token string `json:"refresh_token"`
	}
	param := params{}
	if err := json.NewDecoder(r.Body).Decode(&param); err != nil {
		respondWithError(w, 400, "Invalid request payload", err)
		return
	}

	token_info, err := a.Database.GetUserFromRefreshToken(r.Context(), param.Refresh_token)
	if err != nil {
		respondWithError(w, 500, "error retrieving token", err)
		return
	}
	if token_info.ExpiresAt.Before(time.Now().UTC()) || token_info.RevokedAt.Valid {
		respondWithError(w, http.StatusUnauthorized, "refresh token expired or revoked", nil)
		return
	}
	user_id := token_info.UserID

	new_token, err := auth.MakeJWT(user_id, a.Secret, time.Duration(1)*time.Hour)
	if err != nil {
		respondWithError(w, 500, "error making new jwt token", err)
		return
	}

	err = a.Database.RevokeToken(r.Context(), param.Refresh_token)
	if err != nil {
		respondWithError(w, 401, "something is wrong with revoking", err)
		return
	}
	new_refresh_token, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError(w, 500, "error making new refresh token", err)
		return
	}
	_, err = a.Database.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token:  new_refresh_token,
		UserID: user_id,
	})

	type returnVal struct {
		Token         string `json:"token"`
		Refresh_token string `json:"refresh_token"`
	}

	respondWithJSON(w, 201, returnVal{Token: new_token, Refresh_token: new_refresh_token})
}

func (a *ApiConfig) RevokeHandler(w http.ResponseWriter, r *http.Request) {
	refresh_token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 401, "no refresh token found", err)
		return
	}
	err = a.Database.RevokeToken(r.Context(), refresh_token)
	if err != nil {
		respondWithError(w, 401, "something is wrong with revoking", err)
		return
	}
	respondWithJSON(w, 204, nil)

}

func (a *ApiConfig) deletePostHandler(w http.ResponseWriter, r *http.Request) {
	post_id := r.PathValue("post_id")
	parsed_post_id, err := uuid.Parse(post_id)
	if err != nil {
		respondWithError(w, 400, "Unable to parse ID", err)
		return
	}
	owner_id, err := GetUserID(r.Context())
	if err != nil {
		respondWithError(w, 500, "Error retrieving user_id", err)
		return
	}
	post, err := a.Database.GetPostById(r.Context(), parsed_post_id)
	if err != nil {
		respondWithError(w, 500, "Error retrieving post by id", err)
		return
	}
	// authorization check
	if post.UserID != owner_id {
		respondWithError(w, http.StatusUnauthorized, "method not allowed", err)
		return
	}
	err = a.Database.DeletePost(r.Context(), parsed_post_id)
	if err != nil {
		respondWithError(w, 500, "Couldn't delete post", err)
		return
	}
	respondWithJSON(w, 201, nil)
}
