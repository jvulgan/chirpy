package main

import (
	"fmt"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, req)
	})
}

func (cfg *apiConfig) middlewareMetricsGet(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte(fmt.Sprintf("Hits: %d", cfg.fileserverHits.Load())))
}

func (cfg *apiConfig) middlewareMetricsReset(w http.ResponseWriter, req *http.Request) {
	cfg.fileserverHits = atomic.Int32{}
}

func main() {
	apiCfg := apiConfig{fileserverHits: atomic.Int32{}}
	mux := http.NewServeMux()
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))
	mux.HandleFunc("GET /api/healthz", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
	mux.HandleFunc("GET /api/metrics", apiCfg.middlewareMetricsGet)
	mux.HandleFunc("POST /api/reset", apiCfg.middlewareMetricsReset)
	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	server.ListenAndServe()
}
