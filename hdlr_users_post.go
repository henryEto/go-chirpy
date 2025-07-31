package main

import (
	"encoding/json"
	"net/http"

	"github.com/henryEto/go-chirpy/internal/auth"
	"github.com/henryEto/go-chirpy/internal/database"
)

func (cfg *apiConfig) handlerUsersPost(w http.ResponseWriter, r *http.Request) {
	var reqUser User

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&reqUser)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Failed to decode user", err)
		return
	}

	hashedPsswd, err := auth.HashPassword(reqUser.Password)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Failed to hash password", err)
		return
	}

	user, err := cfg.queries.CreateUser(r.Context(), database.CreateUserParams{
		Email:          reqUser.Email,
		Hashedpassword: hashedPsswd,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create user", err)
		return
	}

	reqUser.ID = user.ID
	reqUser.CreatedAt = user.CreatedAt
	reqUser.UpdatedAt = user.UpdatedAt
	reqUser.IsRed = user.IsChirpyRed
	reqUser.Password = ""
	respondWithJSON(w, http.StatusCreated, reqUser)
}
