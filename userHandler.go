package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type emailVals struct {
	Email string `json:"email"`
}

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func (cfg *Config) userHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	params := emailVals{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 500, "Something went wrong")
		return
	}
	user, err := cfg.db.CreateUser(r.Context(), params.Email)
	respondWithJSON(w, 201, User{
    ID:        user.ID,
    CreatedAt: user.CreatedAt,
    UpdatedAt: user.UpdatedAt,
    Email:     user.Email,
	})
}
