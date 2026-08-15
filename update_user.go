package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Hir0-shi/chirpy/internal/auth"
	"github.com/Hir0-shi/chirpy/internal/database"
)

type updateUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (cfg *apiConfig) handlerUpdateUser(w http.ResponseWriter, r *http.Request) {
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

	var req updateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Something went wrong")
		return
	}

	hashed, err := auth.HashPassword(req.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	userRow, err := cfg.dbQueries.UpdateUser(r.Context(), database.UpdateUserParams{
		ID:             userID,
		Email:          req.Email,
		HashedPassword: hashed,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	resp := User{
		ID:          userRow.ID.String(),
		CreatedAt:   userRow.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt:   userRow.UpdatedAt.UTC().Format(time.RFC3339),
		Email:       userRow.Email,
		IsChirpyRed: userRow.IsChirpyRed,
	}

	respondWithJSON(w, http.StatusOK, resp)
}
