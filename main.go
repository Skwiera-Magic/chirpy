package main

import (
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func main() {
	const path = "."
	const port = ":8080"

	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
	}
	serveMux := http.NewServeMux()
	serveMux.Handle("/app/", apiCfg.metricsMiddleware( http.StripPrefix("/app", http.FileServer(http.Dir(path)))))
	serveMux.HandleFunc("GET /api/healthz", readinessHandler)
	serveMux.HandleFunc("GET /admin/metrics", apiCfg.metricsHandler)
	serveMux.HandleFunc("POST /admin/reset", apiCfg.resetHandler)
	serveMux.HandleFunc("POST /api/validate_chirp", validationHandler)

	server := &http.Server{
		Addr: port,
		Handler: serveMux,
	}
	log.Printf("Serving files from %v on port: %v\n", path, port)
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}