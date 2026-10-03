package main

import (
	"log"
	"net/http"
)

func main() {
	const path = "."
	const port = ":8080"

	serveMux := http.NewServeMux()
	serveMux.Handle("/app/", http.StripPrefix("/app", http.FileServer(http.Dir(path))))
	serveMux.HandleFunc("/healthz", readiness)

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

func readiness(writer http.ResponseWriter, req *http.Request) {
	writer.Header().Add("Content-Type", "text/plain; charset=utf-8")
	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte(http.StatusText(http.StatusOK)))
}