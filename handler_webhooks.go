package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Skwiera-Magic/chirpy/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) webhookHandler(writer http.ResponseWriter, req *http.Request) {
	type parametres struct {
		Event string `json:"event"`
		Data struct {
			UserID uuid.UUID `json:"user_id"`
		}
	}	

	apiKey, err := auth.GetAPIKey(req.Header)
	if err != nil {
		respondWithError(writer, http.StatusUnauthorized, "could not get api key", err)
		return
	}
	if apiKey != cfg.polka {
		respondWithError(writer, http.StatusUnauthorized, "api key is not valid", nil)
		return
	}

	decoder := json.NewDecoder(req.Body)
	params := parametres{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(writer, http.StatusInternalServerError, "could not decode parameters", err)
		return
	}

	if params.Event != "user.upgraded" {
		writer.WriteHeader(http.StatusNoContent)
		return
	}

	_, err = cfg.db.Upgrade(req.Context(), params.Data.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(writer, http.StatusNotFound, "could not find user", err)
			return
		}
		respondWithError(writer, http.StatusInternalServerError, "could not upgrade user", err)
		return
	}

	writer.WriteHeader(http.StatusNoContent)
}