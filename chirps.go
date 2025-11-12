package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/jvulgan/chirpy/internal/database"
)

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) addChirp(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Body   string    `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}
	var p parameters

	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(&p); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not decode parameters", err)
		return
	}

	if len(p.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	params := database.CreateChirpParams{
		Body:   replaceBadWords(p.Body),
		UserID: p.UserID,
	}
	ch, err := cfg.dbQueries.CreateChirp(req.Context(), params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not create chirp", err)
		return
	}
	chirp := Chirp{
		ID:        ch.ID,
		CreatedAt: ch.CreatedAt,
		UpdatedAt: ch.UpdatedAt,
		Body:      ch.Body,
		UserID:    ch.UserID,
	}
	respondWithJSON(w, http.StatusCreated, chirp)
}
