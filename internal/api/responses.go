package api

import (
	"boilerplate/internal/database"
	"time"

	"github.com/google/uuid"
)

// for create user handler
type SafeUserStruct struct {
	ID         uuid.UUID
	Username   string
	Created_at time.Time
	Updated_at time.Time
}

type LoginStruct struct {
	ID            uuid.UUID
	Created_at    time.Time
	Updated_at    time.Time
	Username      string
	Token         string
	Refresh_token string
}

type ProfileStruct struct {
	Username      string
	PostsMetadata []database.Post
	PostCount     int
}

type TokenRefreshStruct struct {
	Token         string `json:"token"`
	Refresh_token string `json:"refresh_token"`
}
