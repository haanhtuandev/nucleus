package api

import (
	"boilerplate/internal/database"
	"database/sql"
	"encoding/json"
	"net/http"
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
