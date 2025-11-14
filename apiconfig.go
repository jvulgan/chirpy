package main

import (
	"sync/atomic"

	"github.com/jvulgan/chirpy/internal/database"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	dbQueries      database.Queries
	platform       string
	jwtSecret      string
}
