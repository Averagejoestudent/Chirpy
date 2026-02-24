package main

import (
	"encoding/json"
	"net/http"

	"github.com/Averagejoestudent/Chirpy/internal/auth"
	"github.com/google/uuid"
)

type Chirpyredparam struct {
	Event string `json:"event"`
	Data  struct {
		UserID string `json:"user_id"`
	} `json:"data"`
}

func (cfg *Config) ChirpyredWebhookhandler(w http.ResponseWriter, r *http.Request) {
	key, err := auth.GetAPIKey(r.Header)
	if key != cfg.polkaSecret {
		respondWithError(w, 401, "API failed")
		return
	}
	if err != nil {
		respondWithError(w, 401, "API failed")
		return
	}
	decoder := json.NewDecoder(r.Body)
	params := Chirpyredparam{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 500, "Something went wrong")
		return
	}
	if params.Event != "user.upgraded" {
		respondWithError(w, 204, "event not recognized")
		return
	}
	id, err := uuid.Parse(params.Data.UserID)
	if err != nil {
		respondWithError(w, 500, "Issue converting ID")
		return
	}
	if params.Event == "user.upgraded" {
		err := cfg.db.SubChirpyRed(r.Context(), id)
		if err != nil {
			respondWithError(w, 404, "ID not found")
			return
		}
		w.WriteHeader(204)
	}

}
