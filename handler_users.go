package main

import (
	"encoding/json"
	"net/http"
)

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	userEmail := struct {
		Email string `json:"email"`
	}{}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&userEmail)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error parsing user", err)
		return
	}

	dbUser, err := cfg.dbQueries.CreateUser(r.Context(), userEmail.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error creating new user", err)
		return
	}

	user := User{}
	user.ID = dbUser.ID
	user.CreatedAt = dbUser.CreatedAt
	user.UpdatedAt = dbUser.UpdatedAt
	user.Email = dbUser.Email

	respondWithJSON(w, http.StatusCreated, user)
}
