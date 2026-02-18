package main

import (
	"log"
	"net/http"
)

func (cfg *Config) handlerReset(w http.ResponseWriter, r *http.Request) {
	if cfg.platform != "dev" {
		respondWithError(w, 403, "Forbidden")
		return
	}
	err := cfg.db.DeleteAllUsers(r.Context())
	if err != nil {
		log.Printf("Error Cannot Delete/Reset Database: %s", err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(http.StatusOK)
}
