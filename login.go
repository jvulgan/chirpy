package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/jvulgan/chirpy/internal/auth"
	"github.com/jvulgan/chirpy/internal/database"
)

const jwtExpires = time.Hour
const refreshTokenExpires = time.Hour * 24 * 60

func (cfg *apiConfig) handleLogin(w http.ResponseWriter, req *http.Request) {
	type params struct {
		Email    string `json:"email"`
		Password string `json:"password"`
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

	jwt, err := auth.MakeJWT(usr.ID, cfg.jwtSecret, jwtExpires)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error creating JWT", err)
		return
	}
	refreshToken := auth.MakeRefreshToken()
	refreshTokenParams := database.CreateRefreshTokenParams{
		Token:     refreshToken,
		UserID:    usr.ID,
		ExpiresAt: time.Now().UTC().Add(refreshTokenExpires),
	}
	if err := cfg.dbQueries.CreateRefreshToken(req.Context(), refreshTokenParams); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error creating refresh token", err)
		return
	}

	resp := User{
		ID:           usr.ID,
		CreatedAt:    usr.CreatedAt,
		UpdatedAt:    usr.UpdatedAt,
		Email:        usr.Email,
		Token:        jwt,
		RefreshToken: refreshToken,
		IsChirpyRed:  usr.IsChirpyRed,
	}
	respondWithJSON(w, http.StatusOK, resp)

}

func (cfg *apiConfig) handleRefresh(w http.ResponseWriter, req *http.Request) {
	token, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "No Bearer token provided", err)
		return
	}
	userID, err := cfg.dbQueries.GetUserFromRefreshToken(req.Context(), token)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized", err)
		return
	}
	jwt, err := auth.MakeJWT(userID, cfg.jwtSecret, jwtExpires)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error creating JWT", err)
		return
	}
	type resp struct {
		Token string `json:"token"`
	}
	r := resp{Token: jwt}
	respondWithJSON(w, http.StatusOK, r)
}

func (cfg *apiConfig) handleRevoke(w http.ResponseWriter, req *http.Request) {
	token, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "No Bearer token provided", err)
		return
	}
	if err := cfg.dbQueries.RevokeRefreshToken(req.Context(), token); err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized", err)
		return
	}
	respondWithJSON(w, http.StatusNoContent, struct{}{})
}
