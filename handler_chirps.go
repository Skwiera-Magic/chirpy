package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/Skwiera-Magic/chirpy/internal/auth"
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
	}

	token, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(writer, http.StatusUnauthorized, "could not find JWT", err)
		return
	}
	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(writer, http.StatusUnauthorized, "could not validate JWT", err)
		return
	}

	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err = decoder.Decode(&params)
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
		UserID: userID,
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

func (cfg *apiConfig) getChirps(writer http.ResponseWriter, req *http.Request) {
	authorID, err := getAuthorID(req)
	if err != nil {
		respondWithError(writer, http.StatusBadRequest, "invalid author ID", err)
		return
	}

	var chirpsDB []database.Chirp

	if authorID != uuid.Nil {
		chirpsDB, err = cfg.db.GetChirpsByAuthor(req.Context(), authorID)
	} else {
		chirpsDB, err = cfg.db.GetChirps(req.Context())
	}
	if err != nil {
		respondWithError(writer, http.StatusInternalServerError, "could not retreive chirps", err)
		return
	}

	chirps := []Chirp{}
	for _, chirp := range chirpsDB {
		chirps = append(chirps, Chirp{
			ID: chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			UserID: chirp.UserID,
			Body: chirp.Body,
		})
	}

	respondWithJson(writer, http.StatusOK, chirps)
}

func (cfg *apiConfig) getChirpByID(writer http.ResponseWriter, req *http.Request) {
	chirpID, err := uuid.Parse(req.PathValue("chirpID"))
	if err != nil {
		respondWithError(writer, http.StatusBadRequest, "no chirp with requested ID", err)
		return
	}


	chirpByID, err := cfg.db.GetChirpByID(req.Context(), chirpID)
	if err != nil {
		respondWithError(writer, http.StatusNotFound, "could not get chirp", err)
		return
	}

	respondWithJson(writer, http.StatusOK, Chirp{
		ID: chirpByID.ID,
		CreatedAt: chirpByID.CreatedAt,
		UpdatedAt: chirpByID.UpdatedAt,
		UserID: chirpByID.UserID,
		Body: chirpByID.Body,
	})
}

func getAuthorID(req *http.Request) (uuid.UUID, error) {
	authorIDString := req.URL.Query().Get("author_id")
	if authorIDString == "" {
		return uuid.Nil, nil
	}
	authorID, err := uuid.Parse(authorIDString)
	if err != nil {
		return uuid.Nil, err
	}
	return authorID, nil
}