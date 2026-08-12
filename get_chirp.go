package main

import (
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"net/http"
	"time"
)

func (cfg *apiConfig) handlerGetChirpByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	chirpIDStr := r.PathValue("chirpID")
	chirpUUID, err := uuid.Parse(chirpIDStr)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	chirpRow, err := cfg.dbQueries.GetChirp(r.Context(), chirpUUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp := Chirp{
		ID:        chirpRow.ID.String(),
		CreatedAt: chirpRow.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt: chirpRow.UpdatedAt.UTC().Format(time.RFC3339),
		Body:      chirpRow.Body,
		UserID:    chirpRow.UserID.String(),
	}

	respondWithJSON(w, http.StatusOK, resp)
}
