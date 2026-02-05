package api

import (
	"math/rand/v2"
	"strings"

	"github.com/lib/pq"
)

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
