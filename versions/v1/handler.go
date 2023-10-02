package v1

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/mail"
)

// Ping provides a simple health check endpoint.
func Ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("pong"))
}

// RegisterUserHandler handles requests to register a new user.
func RegisterUserHandler(w http.ResponseWriter, r *http.Request) {
	var requestData map[string]string
	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	email, emailOk := requestData["email"]
	password, passOk := requestData["password"]
	if !emailOk || !passOk {
		http.Error(w, "Email and password are required", http.StatusBadRequest)
		return
	}

	_, err = mail.ParseAddress(email)
	if err != nil {
		http.Error(w, "Invalid email format", http.StatusBadRequest)
		return
	}

	userID, userEmail, err := RegisterUser(email, password)
	if err != nil {
		// Handle validation errors
		var validationErr *ValidationError
		if errors.As(err, &validationErr) {
			http.Error(w, validationErr.Error(), http.StatusBadRequest)
			return
		}
		// Handle other errors
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(
		map[string]interface{}{
			"id":    userID,
			"email": userEmail,
		},
	)
}
