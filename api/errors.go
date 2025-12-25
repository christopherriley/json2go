package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type HttpError interface {
	write(http.ResponseWriter)
}

type RequestError struct {
	Err   string `json:"error,omitempty"`
	Cause string `json:"cause,omitempty"`
}

func (e *RequestError) Error() string {
	return fmt.Sprintf("%s: %s", e.Err, e.Cause)
}

func (e RequestError) write(w http.ResponseWriter) {
	b, err := json.Marshal(e)
	if err != nil {
		log.Println("** Error: failed to marshal response: ", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusBadRequest)
	_, err = w.Write(b)
	if err != nil {
		log.Println("** Error: failed to write response: ", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
}

type InternalError struct {
	Err string `json:"error,omitempty"`
}

func (e *InternalError) Error() string {
	return fmt.Sprintf("internal error: %s", e.Err)
}

func (e InternalError) write(w http.ResponseWriter) {
	b, err := json.Marshal(e)
	if err != nil {
		log.Println("** Error: failed to marshal response: ", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusInternalServerError)
	_, err = w.Write(b)
	if err != nil {
		log.Println("** Error: failed to write response: ", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
}
