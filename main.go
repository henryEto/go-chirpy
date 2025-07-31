package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/henryEto/go-chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var (
	filepathRoot string
	port         string
	dbURL        string
)

type apiConfig struct {
	fileserverHits atomic.Int32
	queries        *database.Queries
	secret         string
	polkaKey       string
}

func loadEnv() error {
	err := godotenv.Load()
	filepathRoot = os.Getenv("FILE_PATH_ROOT")
	port = os.Getenv("PORT")
	dbURL = os.Getenv("DB_URL")
	return err
}

func main() {
	err := loadEnv()
	if err != nil {
		log.Fatalf("Failed to load .env file: %v", err)
		os.Exit(1)
	}

	conn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
		os.Exit(1)
	}
	queries := database.New(conn)

	// Necesary secrets
	serverSecret := os.Getenv("TOKEN_SECRET")
	if serverSecret == "" {
		log.Fatal("Failed to load TOKEN_SECRET")
		os.Exit(1)
	}
	polkaKey := os.Getenv("POLKA_KEY")
	if polkaKey == "" {
		log.Fatal("Failed to load POLKA_KEY")
		os.Exit(1)
	}

	// Create HTTP request multiplexer
	mux := http.NewServeMux()
	cfg := apiConfig{
		fileserverHits: atomic.Int32{},
		queries:        queries,
		secret:         serverSecret,
		polkaKey:       polkaKey,
	}

	// Simple file server a root dir
	mux.Handle("/app/", http.StripPrefix("/app", cfg.middlewareMetricsInc(
		http.FileServer(http.Dir(filepathRoot)))),
	)

	// Readiness check
	mux.HandleFunc("GET /api/healthz", handlerReadiness)

	// Metrics
	mux.HandleFunc("GET /admin/metrics", cfg.handlerMetrics)
	mux.HandleFunc("POST /admin/reset", cfg.handlerReset)

	// Chirps
	mux.HandleFunc("POST /api/validate_chirp", handlerValidateChirp)
	mux.HandleFunc("POST /api/chirps", cfg.handlerChirpsPost)
	mux.HandleFunc("GET /api/chirps", cfg.handlerChirpsGet)
	mux.HandleFunc("GET /api/chirps/{chirpID}", cfg.handlerChirpsGetByID)
	mux.HandleFunc("DELETE /api/chirps/{chirpID}", cfg.handlerChirpsDelete)

	// Users
	mux.HandleFunc("POST /api/users", cfg.handlerUsersPost)
	mux.HandleFunc("PUT /api/users", cfg.handlerUsersPut)

	// Login
	mux.HandleFunc("POST /api/login", cfg.handlerLogin)
	mux.HandleFunc("POST /api/refresh", cfg.handlerRefresh)
	mux.HandleFunc("POST /api/revoke", cfg.handlerRevoke)

	// Webhooks
	mux.HandleFunc("POST /api/polka/webhooks", cfg.handlerPolkaWebhooks)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: mux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())
}
