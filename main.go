package main

import (
	"encoding/json"
	"fmt"
	"log"
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
}

func validateChirp(w http.ResponseWriter, req *http.Request) {
    type params struct {
        Body string `json:"body"`
    }
    var p params

    type response struct {
        Error string `json:"error"`
        Valid bool `json:"valid"`
    }

    // decode request body
    decoder := json.NewDecoder(req.Body)
    if err := decoder.Decode(&p); err != nil {
        respBody := response{Valid: false, Error: fmt.Sprintf("%v", err)}
        dat, err := json.Marshal(respBody)
        if err != nil {
            log.Printf("Error while marshalling response %s", err)
            w.WriteHeader(500)
            return
        }
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(200)
        w.Write(dat)
        return
    }

    if len(p.Body) > 140 {
        respBody := response{Valid: false, Error: "Chirp is too long"}
        dat, err := json.Marshal(respBody)
        if err != nil {
            log.Printf("Error while marshalling response %s", err)
            w.WriteHeader(500)
            return
        }
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(400)
        w.Write(dat)
        return
    }

    respBody := response{Valid: true, Error: ""}
    dat, err := json.Marshal(respBody)
    if err != nil {
        log.Printf("Error while marshalling response %s", err)
        w.WriteHeader(500)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(200)
    w.Write(dat)
}

func main() {
	apiCfg := apiConfig{fileserverHits: atomic.Int32{}}
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
    // admin api handlers
	mux.HandleFunc("GET /admin/metrics", apiCfg.middlewareMetricsGet)
	mux.HandleFunc("POST /admin/reset", apiCfg.middlewareMetricsReset)
	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	server.ListenAndServe()
}
