package main

import (
	"github.com/Hir0-shi/chirpy/internal/auth"
	"net/http"
)

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {
	bearerToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	if err := cfg.dbQueries.RevokeRefreshToken(r.Context(), bearerToken); err != nil {
		respondWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
