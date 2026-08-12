package main

import (
	"log"
	"net/http"
	"os"
)

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	platform := os.Getenv("PLATFORM")
	if platform != "dev" {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	cfg.fileserverHits.Store(0)

	if err := cfg.dbQueries.DeleteAllUsers(r.Context()); err != nil {
		log.Printf("DeleteAllUsers error: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits reset to 0"))
}
