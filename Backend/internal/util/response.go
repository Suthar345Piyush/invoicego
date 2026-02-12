// response from server in JSON format
// best practice to start function name with "Write" , whenever dealing with any kind of JSON things

package util

import (
	"encoding/json"
	"net/http"
)

// using omitempty means , all validations are skipped when data comes in json format

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// write into JSON

func WriteJSON(w http.ResponseWriter, status int, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(data)

}

// function for data received success or message

func WriteSuccess(w http.ResponseWriter, status int, data interface{}, message string) error {

	return WriteJSON(w, status, Response{
		Success: true,
		Message: message,
		Data:    data,
	})

}

// if any error occured in process

func WriteError(w http.ResponseWriter, status int, err error) error {
	return WriteJSON(w, status, Response{
		Success: false,
		Error:   err.Error(),
	})
}
