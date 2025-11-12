package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

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
