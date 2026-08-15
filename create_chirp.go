package main

import (
	"encoding/json"
	"github.com/Hir0-shi/chirpy/internal/auth"
	"github.com/Hir0-shi/chirpy/internal/database"
	"github.com/google/uuid"
	"net/http"
	"time"
)

type createChirpRequest struct {
	Body string `json:"body"`
}

type Chirp struct {
	ID        string `json:"id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	Body      string `json:"body"`
	UserID    string `json:"user_id"`
}

func (cfg *apiConfig) handlerCreateChirp(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	bearerToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	userID, err := auth.ValidateJWT(bearerToken, cfg.tokenSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req createChirpRequest
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Something went wrong")
		return
	}

	if len(req.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	cleaned := cleanUpProfaneWords(req.Body)

	chirpRow, err := cfg.dbQueries.CreateChirp(r.Context(), database.CreateChirpParams{
		ID:     uuid.New(),
		Body:   cleaned,
		UserID: userID,
	})
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Something went wrong")
		return
	}

	resp := Chirp{
		ID:        chirpRow.ID.String(),
		CreatedAt: chirpRow.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt: chirpRow.UpdatedAt.UTC().Format(time.RFC3339),
		Body:      chirpRow.Body,
		UserID:    chirpRow.UserID.String(),
	}

	respondWithJSON(w, http.StatusCreated, resp)
}
