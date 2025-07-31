package main

import (
	"encoding/json"
	"net/http"

	"github.com/henryEto/go-chirpy/internal/auth"
	"github.com/henryEto/go-chirpy/internal/database"
)

func (cfg *apiConfig) handlerChirpsPost(w http.ResponseWriter, r *http.Request) {
	var reqChirp Chirp

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Failed to find JWT", err)
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.secret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "JWT validation failed", err)
		return
	}

	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&reqChirp)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Failed to decode chirp", err)
		return
	}

	postParms := database.PostChirpParams{
		Body:   reqChirp.Body,
		Userid: userID,
	}
	chirp, err := cfg.queries.PostChirp(r.Context(), postParms)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to post chirp", err)
		return
	}

	reqChirp.ID = chirp.ID
	reqChirp.CreatedAt = chirp.CreatedAt
	reqChirp.UpdatedAt = chirp.UpdatedAt
	reqChirp.UserID = userID
	respondWithJSON(w, http.StatusCreated, reqChirp)
}
