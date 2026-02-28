package api

import (
	"boilerplate/internal/auth"
	"boilerplate/internal/database"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
)

type ApiConfig struct {
	Database *database.Queries
	Secret   string
}

const (
	StatusOK                  = http.StatusOK                  // 200
	StatusCreated             = http.StatusCreated             // 201
	StatusNoContent           = http.StatusNoContent           // 204
	StatusBadRequest          = http.StatusBadRequest          // 400
	StatusUnauthorized        = http.StatusUnauthorized        // 401
	StatusForbidden           = http.StatusForbidden           // 403
	StatusNotFound            = http.StatusNotFound            // 404
	StatusConflict            = http.StatusConflict            // 409
	StatusInternalServerError = http.StatusInternalServerError // 500
)

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte("Service is healthy!"))
}

func homePageHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./static/index.html")
}

func (a *ApiConfig) addUserHandler(w http.ResponseWriter, r *http.Request) {

	sign_up_request := SignupRequest{}
	if err := json.NewDecoder(r.Body).Decode(&sign_up_request); err != nil {
		// 400 because client send bad json
		respondWithError(w, StatusBadRequest, "Invalid request payload", err)
		return
	}
	err := sign_up_request.Validate()
	if err != nil {
		respondWithError(w, StatusBadRequest, "invalid format", err)
		return
	}

	hashed_password, err := auth.HashPassword(sign_up_request.Password)
	if err != nil {
		respondWithError(w, StatusInternalServerError, "error with hash password", err)
		return
	}

	dbParams := database.CreateUserParams{
		Username:       sign_up_request.Username,
		HashedPassword: hashed_password,
	}

	if sign_up_request.Bio != nil {
		dbParams.Bio = sql.NullString{String: *sign_up_request.Bio, Valid: true}
	} else {
		dbParams.Bio = sql.NullString{Valid: false}
	}

	user, err := a.Database.CreateUser(r.Context(), dbParams)

	if err != nil {
		if isUniqueViolation(err) {
			respondWithError(w, StatusConflict, "Username already exists", err)
			return
		}
		respondWithError(w, StatusInternalServerError, "Couldn't create user", err)
		return
	}
	resp := SafeUserStruct{
		ID:         user.ID,
		Username:   user.Username,
		Created_at: user.CreatedAt,
		Updated_at: user.UpdatedAt,
	}
	respondWithJSON(w, StatusOK, resp)
}

func (a *ApiConfig) loginHandler(w http.ResponseWriter, r *http.Request) {
	log_in_request := LoginRequest{}
	if err := json.NewDecoder(r.Body).Decode(&log_in_request); err != nil {
		// 400 because client send bad json
		respondWithError(w, StatusBadRequest, "Invalid request payload", err)
		return
	}
	err := log_in_request.Validate()
	if err != nil {
		respondWithError(w, StatusBadRequest, "Invalid request format", err)
		return
	}
	user, err := a.Database.GetUserByName(r.Context(), log_in_request.Username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, StatusBadRequest, "invalid credentials", err)
			return
		} else {
			respondWithError(w, StatusInternalServerError, "error retrieving user", err)
			return
		}
	}
	authenticated, err := auth.CheckPasswordHash(log_in_request.Password, user.HashedPassword)
	if !authenticated {
		respondWithError(w, StatusBadRequest, "Wrong password", err)
		return
	}
	if err != nil {
		respondWithError(w, StatusInternalServerError, "Error validating password", err)
		return
	}
	refresh_token_string, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError(w, StatusInternalServerError, "error making refresh tokens", err)
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
		respondWithError(w, StatusBadRequest, "error generating jwt token", err)
	}
	resp := LoginStruct{
		ID:            user.ID,
		Created_at:    user.CreatedAt,
		Updated_at:    user.UpdatedAt,
		Username:      user.Username,
		Token:         access_token,
		Refresh_token: refresh_token_string,
	}
	respondWithJSON(w, StatusOK, resp)

}

func (a *ApiConfig) addPostHandler(w http.ResponseWriter, r *http.Request) {
	create_post_request := CreatePostRequest{}

	if err := json.NewDecoder(r.Body).Decode(&create_post_request); err != nil {
		// 400 because client send bad json
		respondWithError(w, StatusBadRequest, "Invalid request payload", err)
		return
	}

	err := create_post_request.Validate()
	if err != nil {
		respondWithError(w, StatusBadRequest, "bad request", err)
		return
	}

	user_id, err := GetUserID(r.Context())
	if err != nil {
		respondWithError(w, StatusInternalServerError, "wrong user_id format", err)
		return
	}
	dbParams := database.CreatePostParams{
		Content: create_post_request.Content,
		Title:   create_post_request.Title,
		UserID:  user_id,
	}
	post, err := a.Database.CreatePost(r.Context(), dbParams)

	if err != nil {
		respondWithError(w, StatusInternalServerError, "error retrieving post", err)
		return
	}
	respondWithJSON(w, StatusCreated, post)
}

// func (a *ApiConfig) getAllUsersHandler(w http.ResponseWriter, r *http.Request) {

// 	users, err := a.Database.GetAllUsers(r.Context())
// 	if err != nil {
// 		respondWithError(w, 500, "Couldn't retrieve all users", err)
// 		return
// 	}

// 	respondWithJSON(w, 201, users)
// }

// func (a *ApiConfig) getAllPostsHandler(w http.ResponseWriter, r *http.Request) {
// 	posts, err := a.Database.GetAllPosts(r.Context())
// 	if err != nil {
// 		respondWithError(w, 500, "Couldn't retrieve all posts", err)
// 		return
// 	}
// 	respondWithJSON(w, 201, posts)
// }

func (a *ApiConfig) getProfileByIdHandler(w http.ResponseWriter, r *http.Request) {
	// retrieve user
	user_id := r.PathValue("user_id")
	parsed_user_id, err := uuid.Parse(user_id)
	user, err := a.Database.GetUserById(r.Context(), parsed_user_id)
	if err != nil {
		respondWithError(w, StatusInternalServerError, "error retrieving user", err)
		return
	}

	// retrieve posts with pagination
	queryParams := r.URL.Query()
	limitStr := queryParams.Get("limit")
	pageStr := queryParams.Get("page")

	page := 1
	limit := 20

	if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
		page = p
	}
	if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
		limit = l
	}

	offset := (page - 1) * limit

	posts, err := a.Database.GetPostsByUser(r.Context(), database.GetPostsByUserParams{
		UserID: parsed_user_id,
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		respondWithError(w, StatusInternalServerError, "error retrieving paginated posts from database", err)
		return
	}

	post_count, err := a.Database.GetPostsCount(r.Context(), parsed_user_id)

	resp := ProfileStruct{
		Username:      user.Username,
		PostsMetadata: posts,
		PostCount:     int(post_count),
	}
	respondWithJSON(w, StatusOK, resp)
}

func (a *ApiConfig) getProfileHandler(w http.ResponseWriter, r *http.Request) {
	user_id, err := GetUserID(r.Context())
	if err != nil {
		respondWithError(w, StatusInternalServerError, "error getting user ID from context", err)
		return
	}
	user, err := a.Database.GetUserById(r.Context(), user_id)
	if err != nil {
		respondWithError(w, StatusInternalServerError, "error retrieving user", err)
		return
	}

	// retrieve posts with pagination
	queryParams := r.URL.Query()
	limitStr := queryParams.Get("limit")
	pageStr := queryParams.Get("page")

	page := 1
	limit := 20

	if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
		page = p
	}
	if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
		limit = l
	}

	offset := (page - 1) * limit

	posts, err := a.Database.GetPostsByUser(r.Context(), database.GetPostsByUserParams{
		UserID: user_id,
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		respondWithError(w, StatusInternalServerError, "error retrieving paginated posts from database", err)
		return
	}

	post_count, err := a.Database.GetPostsCount(r.Context(), user_id)

	resp := ProfileStruct{
		Username:      user.Username,
		PostsMetadata: posts,
		PostCount:     int(post_count),
	}
	respondWithJSON(w, StatusOK, resp)
}

func (a *ApiConfig) updatePostHandler(w http.ResponseWriter, r *http.Request) {
	post_id := r.PathValue("post_id")
	parsed_post_id, err := uuid.Parse(post_id)
	if err != nil {
		respondWithError(w, StatusBadRequest, "Unable to parse ID", err)
		return
	}
	owner_id, err := GetUserID(r.Context())
	if err != nil {
		respondWithError(w, StatusInternalServerError, "Error retrieving user_id", err)
		return
	}
	post, err := a.Database.GetPostById(r.Context(), parsed_post_id)
	if err != nil {
		respondWithError(w, StatusInternalServerError, "Error retrieving post by id", err)
		return
	}
	// authorization check
	if post.UserID != owner_id {
		respondWithError(w, StatusUnauthorized, "method not allowed", err)
		return
	}
	// input params

	update_post_request := UpdatePostRequest{}
	if err := json.NewDecoder(r.Body).Decode(&update_post_request); err != nil {
		// 400 because client send bad json
		respondWithError(w, StatusBadRequest, "Invalid request payload", err)
		return
	}

	err = update_post_request.Validate()
	if err != nil {
		respondWithError(w, StatusInternalServerError, "input failed validation", err)
		return
	}

	update_params := database.UpdatePostInfoParams{
		Title:   update_post_request.Title,
		Content: update_post_request.Content,
		ID:      parsed_post_id,
	}
	err = a.Database.UpdatePostInfo(r.Context(), update_params)
	if err != nil {
		respondWithError(w, StatusInternalServerError, "cannot update post", err)
		return
	}
	respondWithJSON(w, StatusOK, nil)

}

// func (a *ApiConfig) resetHandler(w http.ResponseWriter, r *http.Request) {
// 	err := a.Database.RefreshDB(r.Context())
// 	if err != nil {
// 		respondWithError(w, 500, "Error clearing database", err)
// 		return
// 	}
// 	respondWithJSON(w, 201, nil)
// }

func (a *ApiConfig) refreshHandler(w http.ResponseWriter, r *http.Request) {
	refresh_token_request := RefreshTokenRequest{}
	if err := json.NewDecoder(r.Body).Decode(&refresh_token_request); err != nil {
		respondWithError(w, StatusBadRequest, "Invalid request payload", err)
		return
	}

	err := refresh_token_request.Validate()
	if err != nil {
		respondWithError(w, StatusInternalServerError, "input failed validation", err)
		return
	}

	token_info, err := a.Database.GetUserFromRefreshToken(r.Context(), refresh_token_request.RefreshToken)
	if err != nil {
		respondWithError(w, StatusInternalServerError, "error retrieving token", err)
		return
	}
	if token_info.ExpiresAt.Before(time.Now().UTC()) || token_info.RevokedAt.Valid {
		respondWithError(w, StatusUnauthorized, "refresh token expired or revoked", nil)
		return
	}
	user_id := token_info.UserID

	new_token, err := auth.MakeJWT(user_id, a.Secret, time.Duration(1)*time.Hour)
	if err != nil {
		respondWithError(w, StatusInternalServerError, "error making new jwt token", err)
		return
	}

	err = a.Database.RevokeToken(r.Context(), refresh_token_request.RefreshToken)
	if err != nil {
		respondWithError(w, StatusUnauthorized, "something is wrong with revoking", err)
		return
	}
	new_refresh_token, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError(w, StatusInternalServerError, "error making new refresh token", err)
		return
	}
	_, err = a.Database.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token:  new_refresh_token,
		UserID: user_id,
	})

	resp := TokenRefreshStruct{Token: new_token, Refresh_token: new_refresh_token}

	respondWithJSON(w, StatusCreated, resp)
}

func (a *ApiConfig) RevokeHandler(w http.ResponseWriter, r *http.Request) {
	refresh_token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, StatusUnauthorized, "no refresh token found", err)
		return
	}
	err = a.Database.RevokeToken(r.Context(), refresh_token)
	if err != nil {
		respondWithError(w, StatusUnauthorized, "something is wrong with revoking", err)
		return
	}
	respondWithJSON(w, StatusNoContent, nil)

}

func (a *ApiConfig) deletePostHandler(w http.ResponseWriter, r *http.Request) {
	post_id := r.PathValue("post_id")
	parsed_post_id, err := uuid.Parse(post_id)
	if err != nil {
		respondWithError(w, StatusBadRequest, "Unable to parse ID", err)
		return
	}
	owner_id, err := GetUserID(r.Context())
	if err != nil {
		respondWithError(w, StatusInternalServerError, "Error retrieving user_id", err)
		return
	}
	post, err := a.Database.GetPostById(r.Context(), parsed_post_id)
	if err != nil {
		respondWithError(w, StatusInternalServerError, "Error retrieving post by id", err)
		return
	}
	// authorization check
	if post.UserID != owner_id {
		respondWithError(w, StatusUnauthorized, "method not allowed", err)
		return
	}
	err = a.Database.DeletePost(r.Context(), parsed_post_id)
	if err != nil {
		respondWithError(w, StatusInternalServerError, "Couldn't delete post", err)
		return
	}
	respondWithJSON(w, StatusCreated, nil)
}

func (a *ApiConfig) fetchPostHandler(w http.ResponseWriter, r *http.Request) {

	queryParams := r.URL.Query()
	limitStr := queryParams.Get("limit")
	pageStr := queryParams.Get("page")

	page := 1
	limit := 20

	if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
		page = p
	}
	if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
		limit = l
	}

	offset := (page - 1) * limit

	posts, err := a.Database.FetchPost(r.Context(), database.FetchPostParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		respondWithError(w, StatusInternalServerError, "error retrieving paginated posts from database", err)
		return
	}
	respondWithJSON(w, StatusOK, posts)

}

func (a *ApiConfig) searchPostHandler(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	limitStr := queryParams.Get("limit")
	pageStr := queryParams.Get("page")
	queryStr := queryParams.Get("q")

	page := 1
	limit := 20

	if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
		page = p
	}
	if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
		limit = l
	}

	offset := (page - 1) * limit
	// fetch posts
	posts, err := a.Database.SearchPosts(
		r.Context(),
		database.SearchPostsParams{
			PlaintoTsquery: queryStr,
			Limit:          int32(limit),
			Offset:         int32(offset),
		},
	)
	if err != nil {
		respondWithError(w, StatusInternalServerError, "search failed", err)
		return
	}

	// count
	total, err := a.Database.CountSearchPosts(r.Context(), queryStr)
	if err != nil {
		respondWithError(w, StatusInternalServerError, "count failed", err)
		return
	}

	respondWithJSON(w, StatusOK, map[string]any{
		"query": queryStr,
		"total": total,
		"page":  page,
		"limit": limit,
		"posts": posts,
	})

}

// add an entry to follow table
func (a *ApiConfig) followHandler(w http.ResponseWriter, r *http.Request) {
	// get followee_id
	followee_id := r.PathValue("user_id")
	parsed_followee_id, err := uuid.Parse(followee_id)
	if err != nil {
		respondWithError(w, StatusInternalServerError, "error parsing id", err)
		return
	}

	// get follower_id
	follower_id, err := GetUserID(r.Context())
	if err != nil {
		respondWithError(w, StatusInternalServerError, "cannot retrieve id from auth", err)
		return
	}

	if parsed_followee_id == follower_id {
		respondWithError(w, StatusBadRequest, "you cannot follow yourself", err)
		return
	}

	err = a.Database.CreateFollow(r.Context(), database.CreateFollowParams{
		FolloweeID: parsed_followee_id,
		FollowerID: follower_id,
	})
	if err != nil {
		if isUniqueViolation(err) {
			respondWithError(w, StatusNoContent, "no content", err)
			return
		}
		respondWithError(w, StatusBadRequest, "error creating follow entry", err)
		return
	}
	respondWithJSON(w, StatusOK, nil)
}

// remove an entry to follow table
func (a *ApiConfig) unfollowHandler(w http.ResponseWriter, r *http.Request) {
	followee_id := r.PathValue("user_id")
	parsed_followee_id, err := uuid.Parse(followee_id)
	if err != nil {
		respondWithError(w, StatusInternalServerError, "error parsing id", err)
		return
	}

	// get follower_id
	follower_id, err := GetUserID(r.Context())
	if err != nil {
		respondWithError(w, StatusInternalServerError, "cannot retrieve id from auth", err)
		return
	}

	err = a.Database.DeleteFollow(r.Context(), database.DeleteFollowParams{
		FolloweeID: parsed_followee_id,
		FollowerID: follower_id,
	})
	if err != nil {
		respondWithError(w, StatusBadRequest, "error deleting follow entry", err)
		return
	}
	respondWithJSON(w, StatusOK, nil)
}

func (a *ApiConfig) getFollowersHandler(w http.ResponseWriter, r *http.Request) {
	user_id := r.PathValue("user_id")
	parsed_user_id, err := uuid.Parse(user_id)
	if err != nil {
		respondWithError(w, StatusInternalServerError, "error parsing id", err)
		return
	}

	queryParams := r.URL.Query()
	limitStr := queryParams.Get("limit")
	pageStr := queryParams.Get("page")

	page := 1
	limit := 20

	if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
		page = p
	}
	if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
		limit = l
	}

	offset := (page - 1) * limit

	follower_result, err := a.Database.GetFollowers(r.Context(), database.GetFollowersParams{
		ID:     parsed_user_id,
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		respondWithError(w, StatusInternalServerError, "cannot retrieve follow info", err)
		return
	}
	respondWithJSON(w, StatusOK, follower_result)

}
func (a *ApiConfig) getFolloweesHandler(w http.ResponseWriter, r *http.Request) {
	user_id := r.PathValue("user_id")
	parsed_user_id, err := uuid.Parse(user_id)
	if err != nil {
		respondWithError(w, StatusInternalServerError, "error parsing id", err)
		return
	}

	queryParams := r.URL.Query()
	limitStr := queryParams.Get("limit")
	pageStr := queryParams.Get("page")

	page := 1
	limit := 20

	if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
		page = p
	}
	if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
		limit = l
	}

	offset := (page - 1) * limit

	followee_result, err := a.Database.GetFollowees(r.Context(), database.GetFolloweesParams{
		ID:     parsed_user_id,
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		respondWithError(w, StatusInternalServerError, "cannot retrieve follow info", err)
		return
	}
	respondWithJSON(w, StatusOK, followee_result)

}

func (a *ApiConfig) getSelfFollowersHandler(w http.ResponseWriter, r *http.Request) {
	user_id, err := GetUserID(r.Context())
	if err != nil {
		respondWithError(w, StatusInternalServerError, "error getting id", err)
		return
	}

	queryParams := r.URL.Query()
	limitStr := queryParams.Get("limit")
	pageStr := queryParams.Get("page")

	page := 1
	limit := 20

	if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
		page = p
	}
	if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
		limit = l
	}

	offset := (page - 1) * limit

	follower_result, err := a.Database.GetFollowers(r.Context(), database.GetFollowersParams{
		ID:     user_id,
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		respondWithError(w, StatusInternalServerError, "cannot retrieve follow info", err)
		return
	}
	respondWithJSON(w, StatusOK, follower_result)

}
func (a *ApiConfig) getSelfFolloweesHandler(w http.ResponseWriter, r *http.Request) {
	user_id, err := GetUserID(r.Context())
	if err != nil {
		respondWithError(w, StatusInternalServerError, "error getting id", err)
		return
	}

	queryParams := r.URL.Query()
	limitStr := queryParams.Get("limit")
	pageStr := queryParams.Get("page")

	page := 1
	limit := 20

	if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
		page = p
	}
	if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
		limit = l
	}

	offset := (page - 1) * limit

	followee_result, err := a.Database.GetFollowees(r.Context(), database.GetFolloweesParams{
		ID:     user_id,
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		respondWithError(w, StatusInternalServerError, "cannot retrieve follow info", err)
		return
	}
	respondWithJSON(w, StatusOK, followee_result)

}
