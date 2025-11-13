package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const TokenIssuer = "chirpy"

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
    claims := jwt.RegisteredClaims{
        Issuer: TokenIssuer,
        IssuedAt: jwt.NewNumericDate(time.Now().UTC()),
        ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)),
        Subject: userID.String(),
    }
    t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    s, err := t.SignedString([]byte(tokenSecret))
    if err != nil {
        return "", err
    }
    return s, nil
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
    claims := jwt.RegisteredClaims{}
    t, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
        return []byte(tokenSecret), nil
    })
    if err != nil {
        return uuid.Nil, err
    }
    issuer, err := t.Claims.GetIssuer()
    if err != nil {
        return uuid.Nil, err
    }
    if issuer != TokenIssuer {
        return uuid.Nil, errors.New("invalid issuer")
	}

    idStr, err := t.Claims.GetSubject()
    if err != nil {
        return uuid.Nil, err
    }
    id, err :=  uuid.Parse(idStr)
    if err != nil {
        return uuid.Nil, fmt.Errorf("invalid user id: %v", err)
    }
    return id, nil
}
