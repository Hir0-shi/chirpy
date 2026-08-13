package main

import (
	"encoding/json"
	"github.com/Hir0-shi/chirpy/internal/auth"
	"net/http"
	"time"
)

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Something went wrong")
		return
	}

	userRow, err := cfg.dbQueries.GetUserByEmail(r.Context(), req.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
		return
	}

	ok, err := auth.CheckPasswordHash(req.Password, userRow.HashedPassword)
	if err != nil || !ok {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
		return
	}

	resp := User{
		ID:        userRow.ID.String(),
		CreatedAt: userRow.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt: userRow.UpdatedAt.UTC().Format(time.RFC3339),
		Email:     userRow.Email,
	}

	respondWithJSON(w, http.StatusOK, resp)
}
