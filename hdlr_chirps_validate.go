package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"slices"
	"strings"
)

var badWords = []string{
	"kerfuffle",
	"sharbert",
	"fornax",
}

func handlerValidateChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}
	type returnVals struct {
		CleanedBody string `json:"cleaned_body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := &parameters{}
	err := decoder.Decode(params)
	if err != nil {
		respondWithError(
			w,
			http.StatusInternalServerError,
			fmt.Sprintf("Failed to decode request: %v", err),
			err,
		)
		return
	}

	maxChirpLength := 140
	if len(params.Body) > maxChirpLength {
		respondWithError(
			w,
			http.StatusBadRequest,
			"Chirp is too long",
			nil,
		)
		return
	}

	respondWithJSON(w, http.StatusOK, returnVals{CleanedBody: profaneFilter(params.Body)})
}

func profaneFilter(msg string) string {
	words := strings.Split(msg, " ")
	for i, word := range words {
		if slices.Contains(badWords, strings.ToLower(word)) {
			words[i] = "****"
		}
	}

	return strings.Join(words, " ")
}
