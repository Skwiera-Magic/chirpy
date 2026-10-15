package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) webhookHandler(writer http.ResponseWriter, req *http.Request) {
	type parametres struct {
		Event string `json:"event"`
		Data struct {
			UserID uuid.UUID `json:"user_id"`
		}
	}

	decoder := json.NewDecoder(req.Body)
	params := parametres{}
	err := decoder.Decode(&params)
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