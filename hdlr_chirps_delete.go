package main

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/henryEto/go-chirpy/internal/auth"
)

func (cfg *apiConfig) handlerChirpsDelete(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Could not find token", err)
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.secret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Could not validate JWT", err)
		return
	}

	chirpID, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Failed to parse ID", err)
		return
	}

	chirpDB, err := cfg.queries.GetChirpByID(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to find chirp", err)
		return
	}

	if chirpDB.UserID != userID {
		respondWithError(w, http.StatusForbidden, "Chirp does not belong to user", err)
		return
	}

	err = cfg.queries.DeleteChirp(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Failed to delete chirp", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
