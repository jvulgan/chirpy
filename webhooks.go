package main

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/jvulgan/chirpy/internal/auth"
)

func (cfg *apiConfig) handleWebhook(w http.ResponseWriter, req *http.Request) {
	token, err := auth.GetAPIKey(req.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "No ApiKey provided", err)
		return
	}
	if token != cfg.polkaKey {
		respondWithError(w, http.StatusUnauthorized, "Invalid ApiKey provided", err)
		return
	}
	type params struct {
		Event string `json:"event"`
		Data  struct {
			UserID uuid.UUID `json:"user_id"`
		} `json:"data"`
	}
	var p params

	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(&p); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not decode parameters", err)
		return
	}

	if p.Event != "user.upgraded" {
		respondWithJSON(w, http.StatusNoContent, struct{}{})
		return
	}
	if _, err := cfg.dbQueries.FindUserByID(req.Context(), p.Data.UserID); err != nil {
		respondWithError(w, http.StatusNotFound, "User doesn't exist", err)
		return
	}
	if err := cfg.dbQueries.SetUserIsChirpyRedTrue(req.Context(), p.Data.UserID); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not update user data", err)
		return
	}
	respondWithJSON(w, http.StatusNoContent, struct{}{})
}
