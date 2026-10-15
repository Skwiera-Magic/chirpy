package main

import (
	"net/http"

	"github.com/Skwiera-Magic/chirpy/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) chirpDeleteHandler(writer http.ResponseWriter, req * http.Request) {
	chirpString := req.PathValue("chirpID")
	chirpID, err := uuid.Parse(chirpString)
	if err != nil {
		respondWithError(writer, http.StatusBadRequest, "invalid chirp ID", err)
		return
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

	chirp, err := cfg.db.GetChirpByID(req.Context(), chirpID)
	if err != nil {
		respondWithError(writer, http.StatusNotFound, "could not get chirp", err)
		return
	}
	if chirp.UserID != userID {
		respondWithError(writer, http.StatusForbidden, "you are not authorized to delete this chirp", err)
		return
	}

	err = cfg.db.DeleteChirp(req.Context(),chirpID)
	if err != nil {
		respondWithError(writer, http.StatusInternalServerError, "could not delete chirp", err)
	}

	writer.WriteHeader(http.StatusNoContent)
}