package handlers

import (
	"encoding/json"
	"net/http"
)

type Message struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
	IsError    bool   `json:"is_error"`
}

func sendErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(Message{
		StatusCode: statusCode,
		Message:    message,
		IsError:    true,
	})
}
