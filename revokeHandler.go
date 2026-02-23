package main

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/Averagejoestudent/Chirpy/internal/auth"
	"github.com/Averagejoestudent/Chirpy/internal/database"
)

func (cfg *Config) revokeHandler(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 401, "Invalid token")
		return
	}
	err = cfg.db.RevokeRefreshToken(r.Context(),database.RevokeRefreshTokenParams{
		RevokedAt: sql.NullTime{
			Time: time.Now().UTC(),
			Valid: true,
		},
		Token: token,
	})
	if err != nil {
		respondWithError(w, 401, "Invalid token")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
