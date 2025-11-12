package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func (cfg *apiConfig) addUser(w http.ResponseWriter, req *http.Request) {
	type params struct {
		Email string `json:"email"`
	}
	var p params

	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(&p); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not decode parameters", err)
		return
	}

	u, err := cfg.dbQueries.CreateUser(req.Context(), p.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not create user", err)
		return
	}
	user := User{
		ID:        u.ID,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
		Email:     u.Email,
	}
	respondWithJSON(w, http.StatusCreated, user)
}
