package main

import (
	"encoding/json"
	"net/http"

	"github.com/henryEto/go-chirpy/internal/auth"
	"github.com/henryEto/go-chirpy/internal/database"
)

func (cfg *apiConfig) handlerUsersPut(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Could not find token", err)
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.secret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Failed to validate token", err)
		return
	}

	type rBody struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var reqBody rBody
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&reqBody)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Failed to read user info from request", err)
		return
	}

	userDB, err := cfg.queries.GetUserByID(r.Context(), userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not find user", err)
		return
	}

	if userID != userDB.ID {
		respondWithError(w, http.StatusUnauthorized, "Can't edit other users", err)
		return
	}

	hashedPsswd, err := auth.HashPassword(reqBody.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to hash password", err)
		return
	}

	updatedUser, err := cfg.queries.UpdateUser(
		r.Context(),
		database.UpdateUserParams{
			Email:          reqBody.Email,
			Hashedpassword: hashedPsswd,
			Userid:         userID,
		},
	)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to update user", err)
		return
	}

	respondWithJSON(
		w,
		http.StatusOK,
		User{
			ID:        updatedUser.ID,
			CreatedAt: updatedUser.CreatedAt,
			UpdatedAt: updatedUser.UpdatedAt,
			Email:     updatedUser.Email,
			IsRed:     updatedUser.IsChirpyRed,
		},
	)
}
