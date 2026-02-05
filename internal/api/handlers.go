package api

import (
	"boilerplate/internal/database"
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

type ApiConfig struct {
	Database *database.Queries
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
		Bio      *string `json:"bio"`
	}
	param := params{}
	if err := json.NewDecoder(r.Body).Decode(&param); err != nil {
		// 400 because client send bad json
		respondWithError(w, 400, "Invalid request payload", err)
		return
	}

	dbParams := database.CreateUserParams{
		Username: param.Username,
	}

	if param.Bio != nil {
		dbParams.Bio = sql.NullString{String: *param.Bio, Valid: true}
	} else {
		dbParams.Bio = sql.NullString{Valid: false}
	}

	user, err := a.Database.CreateUser(r.Context(), dbParams)
	if err != nil {
		respondWithError(w, 500, "Couldn't create user", err)
		return
	}
	respondWithJSON(w, 201, user)
}

func (a *ApiConfig) addPostHandler(w http.ResponseWriter, r *http.Request) {
	type params struct {
		UserID  uuid.UUID `json:"user_id"`
		Content string    `json:"content"`
	}
	param := params{}
	if err := json.NewDecoder(r.Body).Decode(&param); err != nil {
		// 400 because client send bad json
		respondWithError(w, 400, "Invalid request payload", err)
		return
	}

	dbParams := database.CreatePostParams{
		Content: param.Content,
		UserID:  param.UserID,
	}

	post, err := a.Database.CreatePost(r.Context(), dbParams)
	if err != nil {
		respondWithError(w, 500, "Couldn't create post", err)
		return
	}
	respondWithJSON(w, 201, post)
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

func (a *ApiConfig) getPostByIdHandler(w http.ResponseWriter, r *http.Request) {
	post_id := r.PathValue("post_id")
	parsed_post_id, err := uuid.Parse(post_id)
	if err != nil {
		respondWithError(w, 500, "Error converting post ID", err)
		return
	}

	post, err := a.Database.GetPostById(r.Context(), parsed_post_id)
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

func (a *ApiConfig) refreshHandler(w http.ResponseWriter, r *http.Request) {
	err := a.Database.RefreshDB(r.Context())
	if err != nil {
		respondWithError(w, 500, "Error clearing database", err)
		return
	}
	respondWithJSON(w, 201, nil)
}
