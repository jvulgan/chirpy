package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
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

func getTokenFromHeader(token_prefix string, headers http.Header) (string, error) {
    authHeader := headers.Get("Authorization")
    if authHeader == "" {
        return authHeader, errors.New("no authorization header found")
    }
    token, ok := strings.CutPrefix(authHeader, fmt.Sprintf("%s ", token_prefix))
    if !ok {
        return "", errors.New("authorization header does not match expected format 'Bearer TOKEN_STRING'")
    }
    return token, nil
}
func GetBearerToken(headers http.Header) (string, error) {
    token, err := getTokenFromHeader("Bearer", headers)
    return token, err
}

func GetAPIKey(headers http.Header) (string, error) {
    token, err := getTokenFromHeader("ApiKey", headers)
    return token, err
}

func MakeRefreshToken() string {
    key := make([]byte, 32)
    rand.Read(key)
    return hex.EncodeToString(key)
}
