package main

import (
	"github.com/Hir0-shi/chirpy/internal/database"
	"github.com/google/uuid"
	"net/http"
	"sort"
	"time"
)

func (cfg *apiConfig) handlerListChirps(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	authorID := r.URL.Query().Get("author_id")
	sortParam := r.URL.Query().Get("sort")

	var (
		chirpsRows []database.Chirp
		err        error
	)

	if authorID != "" {
		authorUUID, parseErr := uuid.Parse(authorID)
		if parseErr != nil {
			respondWithError(w, http.StatusBadRequest, "invalid author_id")
			return
		}
		chirpsRows, err = cfg.dbQueries.ListChirpsByUserID(r.Context(), authorUUID)
	} else {
		chirpsRows, err = cfg.dbQueries.ListChirps(r.Context())
	}

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	sort.Slice(chirpsRows, func(i, j int) bool {
		if sortParam == "desc" {
			return chirpsRows[i].CreatedAt.After(chirpsRows[j].CreatedAt)
		}
		return chirpsRows[i].CreatedAt.Before(chirpsRows[j].CreatedAt)
	})

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
