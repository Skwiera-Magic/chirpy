package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/Skwiera-Magic/chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db *database.Queries
	platform string
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("could not load godotenv: %v", err)
	}
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("could not load database: %v", err)
	}
	platform := os.Getenv("PLATFORM")
	if platform == "" {
		log.Fatal("PLATFORM must be set")
	}
	dbQueries := database.New(db)

	const path = "."
	const port = ":8080"

	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
		db: dbQueries,
		platform: platform,
	}
	serveMux := http.NewServeMux()
	serveMux.Handle("/app/", apiCfg.metricsMiddleware( http.StripPrefix("/app", http.FileServer(http.Dir(path)))))
	serveMux.HandleFunc("GET /api/healthz", readinessHandler)
	serveMux.HandleFunc("GET /admin/metrics", apiCfg.metricsHandler)
	serveMux.HandleFunc("GET /api/chirps", apiCfg.getChirps)
	serveMux.HandleFunc("POST /admin/reset", apiCfg.resetHandler)
	serveMux.HandleFunc("POST /api/chirps", apiCfg.createChirpHandler)
	serveMux.HandleFunc("POST /api/users", apiCfg.createUserHandler)

	server := &http.Server{
		Addr: port,
		Handler: serveMux,
	}
	log.Printf("Serving files from %v on port: %v\n", path, port)
	err = server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}