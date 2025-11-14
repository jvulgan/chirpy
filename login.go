package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/jvulgan/chirpy/internal/auth"
)

func (cfg *apiConfig) Login(w http.ResponseWriter, req *http.Request) {
	type params struct {
		Email            string `json:"email"`
		Password         string `json:"password"`
		ExpiresInSeconds int    `json:"expires_in_seconds"`
	}
	var p params

	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(&p); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not decode parameters", err)
		return
	}

	usr, err := cfg.dbQueries.FindUserByEmail(req.Context(), p.Email)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Could not find specified user", err)
		return
	}
	match, err := auth.CheckPasswordHash(p.Password, usr.HashedPassword)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error while checking password hash", err)
		return
	}
	if !match {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", nil)
		return
	}

	tokenExpires := calculateExpiration(p.ExpiresInSeconds)
	jwt, err := auth.MakeJWT(usr.ID, cfg.jwtSecret, tokenExpires)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error creating JWT", err)
		return
	}

	resp := User{
		ID:        usr.ID,
		CreatedAt: usr.CreatedAt,
		UpdatedAt: usr.UpdatedAt,
		Email:     usr.Email,
		Token:     jwt,
	}
	respondWithJSON(w, http.StatusOK, resp)

}

func calculateExpiration(expiresInSeconds int) time.Duration {
	const defaultExpiresInSeconds = 60 * 60 // one hour
	if expiresInSeconds == 0 || expiresInSeconds > defaultExpiresInSeconds {
		t, _ := time.ParseDuration(fmt.Sprintf("%ds", defaultExpiresInSeconds))
		return t
	}
	t, _ := time.ParseDuration(fmt.Sprintf("%ds", expiresInSeconds))
	return t
}
