// helpers.go
package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

func respondWithError(w http.ResponseWriter, code int, msg string) {
	type errorResponse struct {
		Error string `json:"error"`
	}
	respondWithJSON(w, code, errorResponse{Error: msg})
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

	parts := strings.Split(body, " ")

	for i := range parts {
		lower := strings.ToLower(parts[i])
		if _, isBad := profaneWords[lower]; isBad {
			parts[i] = replacement
		}
	}

	return strings.Join(parts, " ")
}
