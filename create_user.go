package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type createUserRequest struct {
	Email string `json:"email"`
}

type User struct {
	ID        string `json:"id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	Email     string `json:"email"`
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

	userRow, err := cfg.dbQueries.CreateUser(r.Context(), req.Email)
	if err != nil {
		log.Printf("CreateUser error: %v", err)
		respondWithError(w, http.StatusBadRequest, "Something went wrong")
		return
	}

	resp := User{
		ID:        userRow.ID.String(),
		CreatedAt: userRow.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt: userRow.UpdatedAt.UTC().Format(time.RFC3339),
		Email:     userRow.Email,
	}

	respondWithJSON(w, http.StatusCreated, resp)
}
