package utils

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator"
	"github.com/sirupsen/logrus"
)

// respondWithJSON- write json response format
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	response, err := json.Marshal(payload)
	if err != nil {
		logrus.Warn("failed to marshall response", err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, err = w.Write(response)
	if err != nil {
		logrus.Warn("failed to write response", err)
	}
}

// GetBodyParams- get the parameters from the request body
func GetBodyParams(r *http.Request, data interface{}) error {
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		return err
	}
	validator := validator.New()
	return validator.Struct(data)
}

type Response struct {
	Error      *Error      `json:"error,omitempty"`
	Data       interface{} `json:"data,omitempty"`
	StatusCode string      `json:"status"`
}

type Error struct {
	Description string `json:"description"`
}

func WriteHttpSuccess(w http.ResponseWriter, code int, payload interface{}) {
	response := Response{
		Data:       payload,
		StatusCode: fmt.Sprint(code),
	}
	respondWithJSON(w, code, response)
}

func WriteHttpFailure(w http.ResponseWriter, code int, err error) {
	response := Response{
		StatusCode: fmt.Sprint(code),
		Error:      &Error{Description: err.Error()},
	}
	respondWithJSON(w, code, response)
}
