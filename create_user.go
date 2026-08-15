package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/Hir0-shi/chirpy/internal/auth"
	"github.com/Hir0-shi/chirpy/internal/database"
)

type createUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type User struct {
	ID          string `json:"id"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
	Email       string `json:"email"`
	IsChirpyRed bool   `json:"is_chirpy_red"`
}

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var req createUserRequest
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Something went wrong")
		return
	}

	hashed, err := auth.HashPassword(req.Password)
	if err != nil {
		log.Printf("HashPassword error: %v", err)
		respondWithError(w, http.StatusBadRequest, "Something went wrong")
		return
	}

	userRow, err := cfg.dbQueries.CreateUser(r.Context(), database.CreateUserParams{
		Email:          req.Email,
		HashedPassword: hashed,
	})
	if err != nil {
		log.Printf("CreateUser error: %v", err)
		respondWithError(w, http.StatusBadRequest, "Something went wrong")
		return
	}

	resp := User{
		ID:          userRow.ID.String(),
		CreatedAt:   userRow.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt:   userRow.UpdatedAt.UTC().Format(time.RFC3339),
		Email:       userRow.Email,
		IsChirpyRed: userRow.IsChirpyRed,
	}

	respondWithJSON(w, http.StatusCreated, resp)
}
