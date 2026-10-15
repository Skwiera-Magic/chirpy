package main

import (
	"encoding/json"
	"net/http"

	"github.com/Skwiera-Magic/chirpy/internal/auth"
	"github.com/Skwiera-Magic/chirpy/internal/database"
)

func (cfg *apiConfig) usersUpdateHandler(writer http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email string `json:"email"`
	}
	type response struct {
		User
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
		respondWithError(writer, http.StatusInternalServerError, "could not decode parameters", err)
		return
	}

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(writer, http.StatusInternalServerError, "could not hash password", err)
		return
	}

	user, err := cfg.db.UpdateUser(req.Context(), database.UpdateUserParams{
		ID: userID,
		Email: params.Email,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		respondWithError(writer, http.StatusInternalServerError, "could not update user", err)
		return
	}

	respondWithJson(writer, http.StatusOK, response{
		User: User{
			ID: user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email: user.Email,
		},
	})
}