package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type validateChirpRequest struct {
	Body string `json:"body"`
}

type validateChirpResponseValid struct {
	Valid bool `json:"valid"`
}

type validateChirpResponseErr struct {
	Error string `json:"error"`
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(validateChirpResponseErr{Error: msg})
}

// ■ interface{} can be replaced by any
// func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
func respondWithJSON(w http.ResponseWriter, code int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(payload)
}

func handlerValidateChirp(w http.ResponseWriter, r *http.Request) {
	var req validateChirpRequest
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(&req); err != nil {
		log.Printf("Error decoding chirp request: %s", err)
		respondWithError(w, http.StatusBadRequest, "Something went wrong")
		return
	}

	if len(req.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	respondWithJSON(w, http.StatusOK, validateChirpResponseValid{Valid: true})
}
