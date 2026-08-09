package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

type validateChirpRequest struct {
	Body string `json:"body"`
}

type validateChirpResponseValid struct {
	CleanedBody string `json:"cleaned_body"`
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	type errorResponse struct {
		Error string `json:"error"`
	}
	respondWithJSON(w, code, errorResponse{
		Error: msg,
	})
}

func respondWithJSON(w http.ResponseWriter, code int, payload any) {
	b, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, _ = w.Write(b)
}

func cleanUpProfaneWords(body string) string {
	const replacement = "****"
	profaneWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}

	// split by space, because we have words with punctuation attached that should be match
	parts := strings.Split(body, " ")

	for i, j := range parts {
		lower := strings.ToLower(j)
		if _, isBad := profaneWords[lower]; isBad {
			parts[i] = replacement
		}
	}

	return strings.Join(parts, " ")
}

func handlerValidateChirp(w http.ResponseWriter, r *http.Request) {
	var req validateChirpRequest
	dec := json.NewDecoder(r.Body)

	if err := dec.Decode(&req); err != nil {
		log.Printf("Error decoding chirp request: %s", err)
		respondWithError(w, http.StatusBadRequest, "Something went wrong")
		return
	}

	if len(req.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}
	cleaned := cleanUpProfaneWords(req.Body)
	respondWithJSON(w, http.StatusOK, validateChirpResponseValid{CleanedBody: cleaned})
}
