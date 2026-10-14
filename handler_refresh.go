package main

import (
	"net/http"
	"time"

	"github.com/Skwiera-Magic/chirpy/internal/auth"
)

func (cfg *apiConfig) refreshHandler(writer http.ResponseWriter, req *http.Request) {
	type response struct {
		Token string `json:"token"`
	}

	refreshToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(writer, http.StatusBadRequest, "could not find token", err)
		return
	}

	user, err := cfg.db.GetUserFromRefreshToken(req.Context(), refreshToken)
	if err != nil {
		respondWithError(writer, http.StatusUnauthorized, "could not get user for refresh token", err)
		return
	}

	accessToken, err := auth.MakeJWT(
		user.ID,
		cfg.jwtSecret,
		time.Hour,
	)
	if err != nil {
		respondWithError(writer, http.StatusUnauthorized, "could not validate token", err)
		return
	}

	respondWithJson(writer, http.StatusOK, response{
		Token: accessToken,
	})
}

func (cfg *apiConfig) revokeHandler(writer http.ResponseWriter, req *http.Request) {
	refreshToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(writer, http.StatusBadRequest, "could not find token", err)
		return
	}

	_, err = cfg.db.RevokeRefreshToken(req.Context(), refreshToken)
	if err != nil {
		respondWithError(writer, http.StatusInternalServerError, "could not revoke session", err)
		return
	}

	writer.WriteHeader(http.StatusNoContent)
}