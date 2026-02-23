package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func validHandler(r *http.Request) (string , error) {
	type datVals struct {
		Body string `json:"body"`
	}
	decoder := json.NewDecoder(r.Body)
	mydatVals := datVals{}
	err := decoder.Decode(&mydatVals)
	if err != nil {
		 return "Something went wrong" , err
	}
	if len(mydatVals.Body) > 140{
		return "Chirp is too long" , fmt.Errorf("too many charaacters")
	}
	return clean_message(mydatVals.Body) , nil
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
