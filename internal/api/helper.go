package api

import (
	"context"
	"errors"
	"math/rand/v2"
	"regexp"
	"strings"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

// Remove all HTML tags from string to prevent XSS attacks
func sanitizeInput(input string) string {
	// Remove script tags and their content
	scriptRegex := regexp.MustCompile(`(?i)<script[^>]*>.*?</script>`)
	input = scriptRegex.ReplaceAllString(input, "")

	// Remove all remaining HTML tags
	tagRegex := regexp.MustCompile(`<[^>]*>`)
	input = tagRegex.ReplaceAllString(input, "")

	// Decode common HTML entities to prevent encoding attacks
	input = strings.ReplaceAll(input, "&lt;", "<")
	input = strings.ReplaceAll(input, "&gt;", ">")
	input = strings.ReplaceAll(input, "&amp;", "&")
	input = strings.ReplaceAll(input, "&quot;", "\"")
	input = strings.ReplaceAll(input, "&#39;", "'")

	return strings.TrimSpace(input)
}

func cleanTitle(title string) string {
	stripped := strings.TrimSpace(title)
	lower := strings.ToLower(stripped)
	cleaned := strings.Replace(lower, " ", "-", -1)
	return cleaned
}

func generateRandomString(n int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = charset[rand.IntN(len(charset))]
	}
	return string(b)
}

func isUniqueViolation(err error) bool {
	if pqErr, ok := err.(*pq.Error); ok {
		return pqErr.Code == "23505"
	}
	return false
}

func GetUserID(ctx context.Context) (uuid.UUID, error) {
	val := ctx.Value(userIDKey)
	id, ok := val.(uuid.UUID)
	if !ok {
		return uuid.Nil, errors.New("user_id not found in context")
	}
	return id, nil
}
