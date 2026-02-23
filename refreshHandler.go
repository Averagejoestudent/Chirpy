package main

import (
	"net/http"
	"time"

	"github.com/Averagejoestudent/Chirpy/internal/auth"
)

func (cfg *Config) refreshHandler(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 401, "Invalid token")
		return
	}
	user, err := cfg.db.GetUserFromRefreshToken(r.Context(), token)
	if err != nil {
		respondWithError(w, 401, "Invalid token")
		return
	}
	newtoken, err := auth.MakeJWT(user.ID, cfg.jwtSecret, time.Hour)
	if err != nil {
		respondWithError(w, 500, "Couldn't create token")
		return
	}
	respondWithJSON(w, 200, struct {
		Token string 	`json:"token"`
	}{Token: newtoken})
}
