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
		Username string         `json:"username"`
		Bio      sql.NullString `json:"bio"`
	}
	param := params{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&param)
	if err != nil {
		respondWithError(w, 500, "error decoding json", err)
	}
	if param.Bio.Valid {
		user, err := a.Database.CreateUser(r.Context(), database.CreateUserParams{Username: param.Username, Bio: param.Bio})
		if err != nil {
			respondWithError(w, 500, "error in database operation", err)
		}
		respondWithJSON(w, 201, user)
	} else {
		user, err := a.Database.CreateUser(r.Context(), database.CreateUserParams{Username: param.Username})
		if err != nil {
			respondWithError(w, 500, "error creating user", err)
		}
		respondWithJSON(w, 201, user)
	}

}
