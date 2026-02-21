package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Averagejoestudent/Chirpy/internal/auth"
	"github.com/Averagejoestudent/Chirpy/internal/database"
	"github.com/google/uuid"
)

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

type createChirpRequest struct {
	Body   string `json:"body"`
	
}

func (cfg *Config) chripsHandler(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 401, "Invalid token")
	}
	var CleanedBody string
	decoder := json.NewDecoder(r.Body)
	params := createChirpRequest{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, "Something went wrong")
		return
	}
	if len(params.Body) > 140 {
		respondWithError(w, 400, "Chirp is too long")
		return
	}
	CleanedBody = clean_message(params.Body)
	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, 401, "Invalid user_id")
		return
	}
	chirps, err := cfg.db.UserChirp(r.Context(), database.UserChirpParams{Body: CleanedBody, UserID: userID})
	if err != nil {
		respondWithError(w, 400, "Cannot create chirp Something went wrong")
		return
	}
	respondWithJSON(w, 201, Chirp{
		ID:        chirps.ID,
		CreatedAt: chirps.CreatedAt,
		UpdatedAt: chirps.UpdatedAt,
		Body:      chirps.Body,
		UserID:    userID,
	})
}

func (cfg *Config) GetchripsHandler(w http.ResponseWriter, r *http.Request) {
	Allchrips, err := cfg.db.GetAllChirps(r.Context())
	if err != nil {
		respondWithError(w, 400, "Cannot get All Chrips")
	}
	var list_of_chrips []Chirp
	for _, chirp := range Allchrips {
		list_of_chrips = append(list_of_chrips, Chirp(chirp))
	}
	respondWithJSON(w, 200, list_of_chrips)
}

func (cfg *Config) GetOnechripsHandler(w http.ResponseWriter, r *http.Request) {
	chirp_id := r.PathValue("chirpID")
	id, err := uuid.Parse(chirp_id)
	if err != nil {
		respondWithError(w, 404, "Invalid chirp ID format")
		return
	}
	chirps, err := cfg.db.GetChirpsByID(r.Context(), id)
	if err != nil {
		respondWithError(w, 404, "Cannot get Chrips")
		return
	}
	respondWithJSON(w, 200, Chirp{
		ID:        chirps.ID,
		CreatedAt: chirps.CreatedAt,
		UpdatedAt: chirps.UpdatedAt,
		Body:      chirps.Body,
		UserID:    chirps.UserID,
	})
}
