package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"net/http"
	"strings"
)

func GetBearerToken(headers http.Header) (string, error) {
	const prefix = "Bearer "
	myval := headers.Get("Authorization")
	if myval == "" {
		return "", errors.New("Header is empty")
	}
	if !strings.HasPrefix(myval, prefix) {
		return "", errors.New("Header is Prefix is not as intended")
	}
	token_string := strings.TrimSpace(strings.TrimPrefix(myval, prefix))
	if token_string == "" {
		return "", errors.New("missing token")
	}
	if check := strings.Fields(token_string); len(check) != 1 {
		return "", errors.New("too many argument")
	}

	return token_string, nil
}

func MakeRefreshToken() string {
	key := make([]byte, 32)
	rand.Read(key)
	mystring := hex.EncodeToString(key)
	return mystring
}

func GetAPIKey(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("authorization header missing")
	}

	parts := strings.Fields(authHeader)
	if len(parts) == 2 && parts[0] == "ApiKey" {
		return parts[1], nil
	}
	return "", errors.New("invalid authorization header format")
}
