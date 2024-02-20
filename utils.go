package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type e struct {
	Err string `json:"error"`
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	dat, err := json.Marshal(payload)
	if err != nil {
		log.Fatal("Problem marshalling response")
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(dat)
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	e := e{Err: msg}
	dat, err := json.Marshal(e)
	if err != nil {
		log.Fatal("Problem marshalling error")
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(dat)
}
