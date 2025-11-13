package main

import (
	"encoding/json"
	"net/http"

	"github.com/jvulgan/chirpy/internal/auth"
)

func (cfg *apiConfig) Login(w http.ResponseWriter, req *http.Request) {
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
	resp := User{
		ID:        usr.ID,
		CreatedAt: usr.CreatedAt,
		UpdatedAt: usr.UpdatedAt,
		Email:     usr.Email,
	}
	respondWithJSON(w, http.StatusOK, resp)

}
