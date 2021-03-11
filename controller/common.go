package controller

import (
	"encoding/json"
	"log"
	"net/http"
)

func writeOk(rw http.ResponseWriter, content string) {
	body, err := json.Marshal(content)
	if err != nil {
		log.Print(err)
	}
	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	rw.WriteHeader(http.StatusOK)
	_, err = rw.Write(body)
	if err != nil {
		log.Print(err)
	}
}

func writeError(rw http.ResponseWriter, err error, httpCode int) {
	log.Print(err.Error())
	rw.WriteHeader(httpCode)
	if _, errWrite := rw.Write([]byte(err.Error())); errWrite != nil {
		log.Print("error writing the HTTP response: " + errWrite.Error())
	}
}
