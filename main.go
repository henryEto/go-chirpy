package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

const (
	port         = 8080
	filepathRoot = "."
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func main() {
	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
	}

	mux := http.NewServeMux()
	srvr := http.Server{
		Handler: mux,
		Addr:    fmt.Sprintf(":%v", port),
	}

	// App endpoints
	mux.Handle("/app/", apiCfg.midlewareMetricsInc(http.StripPrefix("/app/", http.FileServer(http.Dir(filepathRoot)))))
	// Api endpoints
	mux.HandleFunc("GET /api/healthz", hadlerReadiness)
	mux.HandleFunc("POST /api/validate_chirp", handlerValidateChirp)
	// Admin edpoints
	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)
	mux.HandleFunc("POST /admin/reset", apiCfg.hadlerResetMetrics)

	log.Printf("Serving files from %s on port: %v\n", filepathRoot, port)
	if err := srvr.ListenAndServe(); err != nil {
		log.Println(err)
	}
}
