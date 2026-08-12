package main

import (
	"net/http"
	"time"
)

func (cfg *apiConfig) handlerListChirps(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	chirpsRows, err := cfg.dbQueries.ListChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	resp := make([]Chirp, 0, len(chirpsRows))
	for _, row := range chirpsRows {
		resp = append(resp, Chirp{
			ID:        row.ID.String(),
			CreatedAt: row.CreatedAt.UTC().Format(time.RFC3339),
			UpdatedAt: row.UpdatedAt.UTC().Format(time.RFC3339),
			Body:      row.Body,
			UserID:    row.UserID.String(),
		})
	}

	respondWithJSON(w, http.StatusOK, resp)
}
