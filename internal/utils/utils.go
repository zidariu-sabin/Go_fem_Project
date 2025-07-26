package utils

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// we use this type in order to parse input and send as a json reponse regardless of type(valid json object, error, etc)
type Envelope map[string]interface{}

// formatting json sent as response
func WriteJson(w http.ResponseWriter, status int, data Envelope) error {
	js, err := json.MarshalIndent(data, "", " ")

	if err != nil {
		return err
	}

	js = append(js, '\n')
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)
	return nil
}

// read id param from route
func ReadIDParam(r *http.Request) (int64, error) {
	idParam := chi.URLParam(r, "id")

	if idParam == "" {
		return 0, errors.New("invalid id paramater")
	}

	id, err := strconv.ParseInt(idParam, 10, 64)

	if err != nil {
		return 0, errors.New("invalid id paramater")
	}

	return id, nil
}
