package main

import (
	"encoding/json"
	"net/http"
)

func validationHandler(writer http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}
	type returns struct {
		Valid bool `json:"valid"`
	}

	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(writer, http.StatusBadRequest, "could not decode parameters", err)
		return
	}

	const maxLength = 140
	if len(params.Body) > maxLength {
		respondWithError(writer, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	respondWithJson(writer, http.StatusOK, returns{
		Valid: true,
	})
}