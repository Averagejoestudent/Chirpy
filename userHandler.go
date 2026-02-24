package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Averagejoestudent/Chirpy/internal/auth"
	"github.com/Averagejoestudent/Chirpy/internal/database"
	"github.com/google/uuid"
)

type Params struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type User struct {
	ID          uuid.UUID `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Email       string    `json:"email"`
	IsChirpyRed bool      `json:"is_chirpy_red"`
}

func (cfg *Config) userHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	params := Params{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 500, "Something went wrong")
		return
	}
	hash, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, 500, "Cannot create hash Something went wrong")
		return
	}
	user, err := cfg.db.CreateUser(r.Context(), database.CreateUserParams{
		Email:          params.Email,
		HashedPassword: hash,
	})
	if err != nil {
		respondWithError(w, 500, "Cannot create user Something went wrong")
		return
	}
	respondWithJSON(w, 201, User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	})

}

func (cfg *Config) SetEmailPasswordHandler(w http.ResponseWriter, r *http.Request) {

	tokenString, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 401, "Incorrect token")
		return
	}
	user_id, err := auth.ValidateJWT(tokenString, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, 401, "server broke")
		return
	}
	decoder := json.NewDecoder(r.Body)
	params := Params{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 500, "Something went wrong")
		return
	}
	hash, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, 500, "Cannot create hash Something went wrong")
		return
	}
	user, err := cfg.db.SetUsersEmailPassword(r.Context(), database.SetUsersEmailPasswordParams{
		ID:             user_id,
		HashedPassword: hash,
		Email:          params.Email,
	})
	if err != nil {
		respondWithError(w, 500, "couldn’t update user")
		return
	}
	respondWithJSON(w, 200, User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	})
}
