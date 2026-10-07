package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/Skwiera-Magic/chirpy/internal/database"
	"github.com/google/uuid"
)

type Chirp struct {
	ID uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	UserID uuid.UUID `json:"user_id"`
	Body string `json:"body"`
}

func (cfg *apiConfig) createChirpHandler(writer http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Body string `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}

	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(writer, http.StatusBadRequest, "could not decode parameters", err)
		return
	}

	cleaned, err := validationHandler(params.Body)
	if err != nil {
		respondWithError(writer, http.StatusBadRequest, err.Error(), err)
		return
	}

	chirp, err := cfg.db.CreateChirp(req.Context(), database.CreateChirpParams{
		Body: cleaned,
		UserID: params.UserID,
	})

	if err != nil {
		respondWithError(writer, http.StatusInternalServerError, "could not create chirp", err)
		return
	}

	respondWithJson(writer, http.StatusCreated, Chirp{
		ID: chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body: chirp.Body,
		UserID: chirp.UserID,
	})
}

func validationHandler(body string) (string,error) {
	const maxLength = 140
	if len(body) > maxLength {
		return "", errors.New("Chirp is too long")
	}

	profanities := map[string]struct{}{
		"kerfuffle": {},
		"sharbert": {},
		"fornax": {},
	}
	cleaned := cleanBody(body, profanities)
	return cleaned, nil
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