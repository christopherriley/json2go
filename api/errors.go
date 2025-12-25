package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type HttpErrorWriter interface {
	write(http.ResponseWriter)
}

type HttpError struct {
	HttpErrorWriter `json:"-"`
	Err             string `json:"error,omitempty"`
}

func writeResponseWithStatus(w http.ResponseWriter, b []byte, statusCode int) {
	w.WriteHeader(statusCode)
	_, err := w.Write(b)
	if err != nil {
		log.Println("** Error: failed to write response: ", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
}

type RequestError struct {
	HttpError
	Cause string `json:"cause,omitempty"`
}

func NewRequestError(err string, cause error) RequestError {
	re := RequestError{HttpError: HttpError{Err: err}}
	if cause != nil {
		re.Cause = cause.Error()
	}

	return re
}

func (e RequestError) Error() string {
	return fmt.Sprintf("%s: %s", e.Err, e.Cause)
}

func (e RequestError) write(w http.ResponseWriter) {
	b, err := json.Marshal(e)
	if err != nil {
		log.Println("** Error: failed to marshal response: ", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	writeResponseWithStatus(w, b, http.StatusBadRequest)
}

type InternalError struct {
	HttpError
}

func NewInternalError(s string) InternalError {
	return InternalError{HttpError: HttpError{Err: s}}
}

func (e InternalError) Error() string {
	return fmt.Sprintf("internal error: %s", e.Err)
}

func (e InternalError) write(w http.ResponseWriter) {
	b, err := json.Marshal(e)
	if err != nil {
		log.Println("** Error: failed to marshal response: ", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	writeResponseWithStatus(w, b, http.StatusInternalServerError)
}
