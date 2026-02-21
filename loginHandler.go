package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Averagejoestudent/Chirpy/internal/auth"
	"github.com/google/uuid"
)

type loginRequest struct {
	Email            string `json:"email"`
	Password         string `json:"password"`
	ExpiresInSeconds *int   `json:"expires_in_seconds"`
}

func (cfg *Config) loginHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	params := loginRequest{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 500, "Something went wrong")
		return
	}
	user, err := cfg.db.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, 401, "Incorrect email or password")
		return
	}
	checking, err := auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil {
		respondWithError(w, 401, "Incorrect email or password")
		return
	}
	if !checking {
		respondWithError(w, 401, "Incorrect email or password")
		return
	}
	expiresIn := time.Hour
	if params.ExpiresInSeconds != nil {
		d := time.Duration(*params.ExpiresInSeconds) * time.Second
		if d > time.Hour {
			d = time.Hour
		}
		if d > 0 {
			expiresIn = d
		}
	}
	token, err := auth.MakeJWT(user.ID, cfg.jwtSecret, expiresIn)
	if err != nil {
		respondWithError(w, 500, "Couldn't create token")
		return
	}

	respondWithJSON(w, 200, struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Email     string    `json:"email"`
		Token     string    `json:"token"`
	}{
		ID: user.ID, CreatedAt: user.CreatedAt, UpdatedAt: user.UpdatedAt, Email: user.Email, Token: token,
	})
}
