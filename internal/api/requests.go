package api

import (
	"errors"
	"strings"
)

// SignupRequest represents the request body for POST /signup
type SignupRequest struct {
	Username string  `json:"username"`
	Password string  `json:"password"`
	Bio      *string `json:"bio"`
}

func (r *SignupRequest) Validate() error {
	r.Username = strings.TrimSpace(r.Username)
	r.Password = strings.TrimSpace(r.Password)

	if r.Username == "" {
		return errors.New("username is required")
	}
	if len(r.Username) < 3 || len(r.Username) > 20 {
		return errors.New("username must be between 3 and 20 characters")
	}
	if !isAlphanumeric(r.Username) {
		return errors.New("username can only contain letters and numbers")
	}
	if r.Password == "" {
		return errors.New("password is required")
	}
	if len(r.Password) < 8 {
		return errors.New("password must be at least 8 characters")
	}
	if r.Bio != nil && len(*r.Bio) > 500 {
		return errors.New("bio must be under 500 characters")
	}
	return nil
}

// LoginRequest represents the request body for POST /login
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (r *LoginRequest) Validate() error {
	r.Username = strings.TrimSpace(r.Username)
	r.Password = strings.TrimSpace(r.Password)

	if r.Username == "" {
		return errors.New("username is required")
	}
	if r.Password == "" {
		return errors.New("password is required")
	}
	return nil
}

// CreatePostRequest represents the request body for POST /me/posts
type CreatePostRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

func (r *CreatePostRequest) Validate() error {
	r.Title = strings.TrimSpace(r.Title)
	r.Content = strings.TrimSpace(r.Content)

	if r.Title == "" {
		return errors.New("title is required")
	}
	if len(r.Title) > 200 {
		return errors.New("title must be under 200 characters")
	}
	if r.Content == "" {
		return errors.New("content is required")
	}
	if len(r.Content) > 50000 {
		return errors.New("content must be under 50000 characters")
	}
	return nil
}

// UpdatePostRequest represents the request body for PUT /me/posts/{post_id}
type UpdatePostRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

func (r *UpdatePostRequest) Validate() error {
	r.Title = strings.TrimSpace(r.Title)
	r.Content = strings.TrimSpace(r.Content)

	if r.Title == "" {
		return errors.New("title is required")
	}
	if len(r.Title) > 200 {
		return errors.New("title must be under 200 characters")
	}
	if r.Content == "" {
		return errors.New("content is required")
	}
	if len(r.Content) > 50000 {
		return errors.New("content must be under 50000 characters")
	}
	return nil
}

// RefreshTokenRequest represents the request body for POST /auth/refresh
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

func (r *RefreshTokenRequest) Validate() error {
	r.RefreshToken = strings.TrimSpace(r.RefreshToken)

	if r.RefreshToken == "" {
		return errors.New("refresh_token is required")
	}
	return nil
}

// Helper functions

func isAlphanumeric(s string) bool {
	for _, c := range s {
		if !((c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9')) {
			return false
		}
	}
	return true
}
