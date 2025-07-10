package util

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

var (
	ErrInternal = errors.New("internal error")
)

// logs an error if it's not nil
func LogError(err error) {
	if err != nil {
		log.Printf("error: %v", err)
	}
}

// writeJSON writes the data as json.
func WriteJSON(w http.ResponseWriter, data interface{}, statusCode int) error {
	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	err := json.NewEncoder(w).Encode(data)
	return err
}

// writeError writes the message as an error
func WriteError(w http.ResponseWriter, message string, statusCode int) error {
	log.Printf("response errored: %s", message)
	return WriteJSON(
		w,
		map[string]string{"error": message},
		statusCode,
	)
}
