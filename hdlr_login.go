package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/henryEto/go-chirpy/internal/auth"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	var reqUser User

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&reqUser)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Failed to decode login request", err)
		return
	}

	userDB, err := cfg.queries.GetUserByEmail(r.Context(), reqUser.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	err = auth.CheckPasswordHash(reqUser.Password, userDB.HashedPassword)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	tokenString, err := auth.MakeJWT(userDB.ID, cfg.secret, time.Hour)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to generate token", err)
		return
	}

	refreshToken, err := cfg.queries.CreateRefreshToken(r.Context(), userDB.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to generate refresh token", err)
		return
	}

	reqUser.Token = tokenString
	reqUser.RefreshToken = refreshToken.Token
	reqUser.ID = userDB.ID
	reqUser.CreatedAt = userDB.CreatedAt
	reqUser.UpdatedAt = userDB.UpdatedAt
	reqUser.IsRed = userDB.IsChirpyRed
	reqUser.Password = ""
	respondWithJSON(w, http.StatusOK, reqUser)
}
