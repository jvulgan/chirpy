package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/jvulgan/chirpy/internal/database"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	dbQueries      database.Queries
	platform       string
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, req)
	})
}

func (cfg *apiConfig) middlewareMetricsGet(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	data := fmt.Sprintf(`<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>`, cfg.fileserverHits.Load())
	w.Write([]byte(data))
}

func (cfg *apiConfig) middlewareMetricsReset(w http.ResponseWriter, req *http.Request) {
	cfg.fileserverHits = atomic.Int32{}
	if cfg.platform == "dev" {
		if err := cfg.dbQueries.DeleteUsers(req.Context()); err != nil {
			log.Printf("failed to call DeleteUsers: %s", err)
		}
	}
}

func validateChirp(w http.ResponseWriter, req *http.Request) {
	type params struct {
		Body string `json:"body"`
	}
	var p params

	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(&p); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not decode parameters", err)
		return
	}

	if len(p.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	type response struct {
		CleanedBody string `json:"cleaned_body"`
	}
	respBody := response{CleanedBody: replaceBadWords(p.Body)}
	respondWithJSON(w, http.StatusOK, respBody)
}

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

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	platform := os.Getenv("PLATFORM")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Printf("err connecting to postgres: %s", err)
		return
	}
	dbQueries := database.New(db)
	apiCfg := apiConfig{fileserverHits: atomic.Int32{}, dbQueries: *dbQueries, platform: platform}
	mux := http.NewServeMux()
	// app handler
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))
	// api calls handlers
	mux.HandleFunc("GET /api/healthz", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
	mux.HandleFunc("POST /api/validate_chirp", validateChirp)
	mux.HandleFunc("POST /api/users", apiCfg.addUser)
	// admin api handlers
	mux.HandleFunc("GET /admin/metrics", apiCfg.middlewareMetricsGet)
	mux.HandleFunc("POST /admin/reset", apiCfg.middlewareMetricsReset)
	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	server.ListenAndServe()
}
