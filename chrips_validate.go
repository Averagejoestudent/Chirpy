package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

func validHandler(w http.ResponseWriter, r *http.Request) {
	type datVals struct {
		Body string `json:"body"`
	}
	decoder := json.NewDecoder(r.Body)
	mydatVals := datVals{}
	err := decoder.Decode(&mydatVals)
	if err != nil {
		respondWithError(w, 500, "Something went wrong")
	}
	if len(mydatVals.Body) < 140 {
		type myValid struct {
			CleanedBody string `json:"cleaned_body"`
		}
		passvalid := myValid{
			CleanedBody: clean_message(mydatVals.Body),
		}
		respondWithJSON(w, 200, passvalid)
	} else {
		respondWithError(w, 400, "Chirp is too long")
	}
}

func clean_message(msg string) string {
	message := strings.Split(msg, " ")
	clean_msg := []string{}
	for _, word := range message {
		if word_checker := strings.ToLower(word); word_checker == "kerfuffle" || word_checker == "sharbert" || word_checker == "fornax" {
			word = "****"
		}
		clean_msg = append(clean_msg, word)
	}
	return strings.Join(clean_msg, " ")
}
