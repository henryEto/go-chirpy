package main

import (
	"net/http"
	"time"

	"github.com/henryEto/go-chirpy/internal/auth"
)

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Token string `json:"token"`
	}
	refToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Could not find token", err)
		return
	}

	user, err := cfg.queries.GetUserFromRefreshToken(r.Context(), refToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Could not find user for token", err)
		return
	}

	accessToken, err := auth.MakeJWT(user.ID, cfg.secret, time.Hour)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Could not validate token", err)
		return
	}

	respondWithJSON(w, http.StatusOK, response{Token: accessToken})
}

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Token string `json:"token"`
	}
	refToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Could not find token", err)
		return
	}

	_, err = cfg.queries.RevokeRefreshToken(r.Context(), refToken)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not find active token", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
