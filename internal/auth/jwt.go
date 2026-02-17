package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "haanhtuandevblog",
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)),
		Subject:   userID.String(),
	})
	signed, err := token.SignedString([]byte(tokenSecret))
	if err != nil {
		return "", err
	}
	return signed, nil

}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	claims := jwt.RegisteredClaims{}
	_, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(tokenSecret), nil
	})
	if err != nil {
		return uuid.UUID{}, errors.New("error parsing with claims")
	}
	user_id, err := uuid.Parse(claims.Subject)
	if err != nil {
		return uuid.UUID{}, errors.New("error parsing user id")
	}
	return user_id, nil

}

func GetBearerToken(headers http.Header) (string, error) {
	bearer_string := headers.Get("Authorization")
	reg := regexp.MustCompile("\\s+")
	bearer_string = reg.ReplaceAllString(bearer_string, "")
	bearer_string = strings.TrimPrefix(bearer_string, "Bearer")
	if bearer_string == "" {
		return "", errors.New("no string bearer found")
	}
	return bearer_string, nil
}

func MakeRefreshToken() (string, error) {
	data := make([]byte, 32) // 16 bytes = 32 hex characters

	_, err := rand.Read(data)
	if err != nil {
		return "", err
	}
	token := hex.EncodeToString(data)
	return token, nil
}
