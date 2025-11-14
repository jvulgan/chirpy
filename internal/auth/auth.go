package auth

import (
	"errors"
	"net/http"
	"strings"

	"github.com/alexedwards/argon2id"
)

func HashPassword(password string) (string, error) {
    hash, err := argon2id.CreateHash(password, argon2id.DefaultParams)
    if err != nil {
        return "", err
    }
    return hash, nil
}

func CheckPasswordHash(password, hash string) (bool, error) {
    match, err := argon2id.ComparePasswordAndHash(password, hash)
    if err != nil {
        return false, err
    }
    return match, nil
}

func GetBearerToken(headers http.Header) (string, error) {
    authHeader := headers.Get("Authorization")
    if authHeader == "" {
        return authHeader, errors.New("no authorization header found")
    }
    token, ok := strings.CutPrefix(authHeader, "Bearer ")
    if !ok {
        return "", errors.New("authorization header does not match expected format 'Bearer TOKEN_STRING'")
    }
    return token, nil
}
