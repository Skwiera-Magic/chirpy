package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func respondWithError(writer http.ResponseWriter, code int, msg string, err error) {
	if err != nil {
		log.Println(err)
	}
	if code > 499 {
		log.Printf("Responding with 5xxerror: %v", msg)
	}
	type errorResponse struct {
		Error string `json:"error"`
	}
	respondWithJson(writer, code, errorResponse{
		Error: msg,
	})
}

func respondWithJson(writer http.ResponseWriter, code int, payload interface{})  {
	writer.Header().Set("Content-Type", "application/json")
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %v", err)
		writer.WriteHeader(500)
		return
	}
	writer.WriteHeader(code)
	writer.Write(data)
}