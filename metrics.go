package main

import (
	"fmt"
	"net/http"
)

func (cfg *apiConfig) hadlerResetMetrics(w http.ResponseWriter, r *http.Request) {
	hits := cfg.fileserverHits.Swap(0)
	w.Header().Add("Content-Type", "text/plain charset=utf-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Hits reseted %v -> %v", hits, cfg.fileserverHits.Load())
}

func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, r *http.Request) {
	htmlPage := `<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>`

	w.Header().Add("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, htmlPage, cfg.fileserverHits.Load())
}

func (cfg *apiConfig) midlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}
