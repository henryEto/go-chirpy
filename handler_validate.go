package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

func handlerValidateChirp(w http.ResponseWriter, r *http.Request) {
	type resp struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	body := resp{}
	err := decoder.Decode(&body)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong", err)
		return
	}

	if len(body.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp is to long", nil)
		return
	}

	respondWithJSON(w, http.StatusOK, struct {
		CleanBody string `json:"cleaned_body"`
	}{CleanBody: cleanMessage(body.Body)})
}

func cleanMessage(msg string) string {
	badWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}
	msgWords := strings.Fields(msg)
	for i, word := range msgWords {
		if _, ok := badWords[strings.ToLower(word)]; ok {
			msgWords[i] = "****"
		}
	}
	return strings.Join(msgWords, " ")
}
