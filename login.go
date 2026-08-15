package main

import (
	"encoding/json"
	"github.com/Hir0-shi/chirpy/internal/auth"
	"github.com/Hir0-shi/chirpy/internal/database"
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

	token, err := auth.MakeJWT(userRow.ID, cfg.tokenSecret, time.Hour)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	refreshToken := auth.MakeRefreshToken()

	if err := cfg.dbQueries.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token:     refreshToken,
		UserID:    userRow.ID,
		ExpiresAt: time.Now().Add(60 * 24 * time.Hour),
	}); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	resp := struct {
		ID           string `json:"id"`
		CreatedAt    string `json:"created_at"`
		UpdatedAt    string `json:"updated_at"`
		Email        string `json:"email"`
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
		IsChirpyRed  bool   `json:"is_chirpy_red"`
	}{
		ID:           userRow.ID.String(),
		CreatedAt:    userRow.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt:    userRow.UpdatedAt.UTC().Format(time.RFC3339),
		Email:        userRow.Email,
		Token:        token,
		RefreshToken: refreshToken,
		IsChirpyRed:  userRow.IsChirpyRed,
	}

	respondWithJSON(w, http.StatusOK, resp)
}
