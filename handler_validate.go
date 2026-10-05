package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

func validationHandler(writer http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}
	type returns struct {
		Cleaned string `json:"cleaned_body"`
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

	profanities := map[string]struct{}{
		"kerfuffle": {},
		"sharbert": {},
		"fornax": {},
	}
	cleaned := cleanBody(params.Body, profanities)

	respondWithJson(writer, http.StatusOK, returns{
		Cleaned: cleaned,
	})
}

func cleanBody(body string, profanities map[string]struct{}) string {
	words := strings.Split(body, " ")
	for i, word := range words {
		lowercase := strings.ToLower(word)
		if _, ok := profanities[lowercase]; ok {
			words[i] = "****"
		}
	}
	cleaned := strings.Join(words, " ")
	return cleaned
}