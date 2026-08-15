package main

import (
	"github.com/Hir0-shi/chirpy/internal/auth"
	"net/http"
	"time"
)

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	bearerToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	rt, err := cfg.dbQueries.GetRefreshToken(r.Context(), bearerToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	if rt.RevokedAt.Valid {
		respondWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	if time.Now().After(rt.ExpiresAt) {
		respondWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	token, err := auth.MakeJWT(rt.UserID, cfg.tokenSecret, time.Hour)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	respondWithJSON(w, http.StatusOK, struct {
		Token string `json:"token"`
	}{
		Token: token,
	})
}
