package utils

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type Envelope map[string]interface{}

func WriteJSON(w http.ResponseWriter, status int, data Envelope) error {
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

func ReadIDParam(r *http.Request) (int, error) {
	id := chi.URLParam(r, "id")
	if id == "" {
		return 0, errors.New("invalid id parameter")
	}
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return 0, errors.New("invalid id parameter")
	}

	return idInt, nil
}
