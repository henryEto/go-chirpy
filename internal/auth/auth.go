package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"net/http"
	"strings"
)

func GetAPIKey(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("no autorization header found")
	}

	tokenSplit := strings.Split(authHeader, " ")
	if len(tokenSplit) < 2 || tokenSplit[0] != "ApiKey" {
		return "", errors.New("auth header malformed")
	}

	return tokenSplit[1], nil
}

func MakeRefreshToken() string {
	token := make([]byte, 32)
	rand.Read(token)
	return hex.EncodeToString(token)
}

func GetBearerToken(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("no autorization header found")
	}

	tokenSplit := strings.Split(authHeader, " ")
	if len(tokenSplit) < 2 || tokenSplit[0] != "Bearer" {
		return "", errors.New("auth header malformed")
	}

	return tokenSplit[1], nil
}
